package services

import (
	"errors"
	"fmt"
	"io"
	"llm_inference_web/accessor/dao"
	"net/http"

	"llm_inference_web/accessor/client"
	"llm_inference_web/accessor/dto"
	"llm_inference_web/accessor/resource"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type OnlineInferenceOperator struct {
	userId int
}

func NewOnlineInferenceOperator(userId int) *OnlineInferenceOperator {
	return &OnlineInferenceOperator{userId: userId}
}

func (o *OnlineInferenceOperator) ChatCompletion(ctx *gin.Context, params *dto.ChatCompletionReq) error {
	// 调用分词服务：获取当前prompt的分词结果
	tokenIds, tokenTypeIds, err := client.NewTokenizer().Encode(ctx, params.Prompt)
	if err != nil {
		return err
	}

	// 创建对话历史：如果为第一次对话，会自动创建各种id
	historyID, err := dao.NewChatHistoryDao().Create(
		o.userId, params.ChatSessionId, params.ParentMessageId, params.Prompt, len(tokenIds),
	)
	if err != nil {
		return err
	}

	// 调用模型服务，开启流式传输
	modelSrvClt := client.NewModelServer()
	modelAnswerStream, err := modelSrvClt.ChatCompletion(
		ctx, params.ChatSessionId, tokenIds, tokenTypeIds, params.InferenceParams,
	)
	defer modelSrvClt.CloseChatCompletionStream()
	if err != nil {
		return err
	}

	// 开启sse流
	flusher, ok := ctx.Writer.(http.Flusher)
	if !ok {
		return errors.New("failed to open http flusher")
	}

	// 从下游读取流式数据
	var modelAnswer string
	var answerTokensCnt int
	for {
		// 接收数据
		resp, err := modelAnswerStream.Recv()
		if err != nil {
			if err == io.EOF {
				flusher.Flush()
				return dao.NewChatHistoryDao().UpdateAnswer(historyID, modelAnswer, answerTokensCnt)
			}
			resource.Logger.WithFields(
				logrus.Fields{
					"userID":        o.userId,
					"chatSessionID": params.ChatSessionId,
					"error":         err.Error(),
				},
			).Error("read from model server failed")
			flusher.Flush()
			return dao.NewChatHistoryDao().UpdateStatus2Failed(historyID)
		}
		// 解码
		assistantTokenIds := resp.GetTokenIds()
		data, err := client.NewTokenizer().Decode(ctx, assistantTokenIds)
		if err != nil {
			resource.Logger.WithFields(
				logrus.Fields{
					"userID":        o.userId,
					"chatSessionID": params.ChatSessionId,
					"error":         err.Error(),
				},
			).Error("tokenizer decode failed")
			return err
		}
		answerTokensCnt += len(assistantTokenIds)
		// SSE 格式数据返回
		_, _ = fmt.Fprintf(ctx.Writer, "event: customEvent\n")
		_, _ = fmt.Fprintf(ctx.Writer, "data: %s\n\n", data)
		modelAnswer += data
		flusher.Flush()
	}
}

func (o *OnlineInferenceOperator) GetChatHistory(_ *gin.Context, params dto.ChatHistoryReq) (
	res dto.ChatHistoryResp, err error) {
	res.TotalCount, err = dao.NewChatHistoryDao().CountByChatSessionID(o.userId, params.ChatSessionID)
	if err != nil {
		return res, err
	}

	histories, err := dao.NewChatHistoryDao().GetByChatSessionID(
		o.userId, params.ChatSessionID, params.PageIndex, params.PageSize,
	)
	if err != nil {
		return res, err
	}
	for _, h := range histories {
		res.ChatHistories = append(res.ChatHistories, dto.ChatHistoryItem{
			MessageID:       h.MessageID,
			ParentMessageID: h.ParentMessageID,
			Prompt:          h.Prompt,
			Answer:          h.Answer,
			TokenCnt:        h.PromptTokenCnt + h.AnswerTokenCnt,
			Timestamp:       h.CreatedAt.Unix(),
			Status:          h.Status,
		})
	}
	return
}
