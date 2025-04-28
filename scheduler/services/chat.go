package services

import (
	"errors"
	"fmt"
	"io"
	"llm_online_inference/scheduler/dao"
	"net/http"

	"llm_online_inference/scheduler/client"
	"llm_online_inference/scheduler/dto"
	"llm_online_inference/scheduler/resource"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type OnlineInferenceOperator struct{}

func NewOnlineInferenceOperator() *OnlineInferenceOperator {
	return &OnlineInferenceOperator{}
}

func (o *OnlineInferenceOperator) ChatCompletion(ctx *gin.Context, userID int, params *dto.ChatCompletionReq) error {
	// 调用分词服务：获取当前prompt的分词结果
	tokenIds, tokenTypeIds, err := client.NewTokenizer().Encode(ctx, params.Prompt)
	if err != nil {
		return err
	}

	// 创建对话历史：如果为第一次对话，会自动创建各种id
	historyID, err := dao.NewChatHistoryDao().Create(
		userID, params.ChatSessionId, params.ParentMessageId, params.Prompt, len(tokenIds),
	)
	if err != nil {
		return err
	}

	// 调用模型服务，开启流式传输
	modelAnswerStream, err := client.NewModelServer().ChatCompletion(ctx, params.ChatSessionId, tokenIds, tokenTypeIds)
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
				_ = dao.NewChatHistoryDao().UpdateAnswer(historyID, modelAnswer, answerTokensCnt)
				return nil
			}
			resource.Logger.WithFields(
				logrus.Fields{
					"userID":        userID,
					"chatSessionID": params.ChatSessionId,
					"error":         err.Error(),
				},
			).Error("read from model server failed")
			flusher.Flush()
			_ = dao.NewChatHistoryDao().UpdateStatus2Failed(historyID)
			return err
		}
		// 解码
		assistantTokenIds := resp.GetTokenIds()
		data, err := client.NewTokenizer().Decode(ctx, assistantTokenIds)
		if err != nil {
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

func (o *OnlineInferenceOperator) GetChatHistory(_ *gin.Context, userID int, params dto.ChatHistoryReq) (
	res dto.ChatHistoryResp, err error) {
	res.TotalCount, err = dao.NewChatHistoryDao().CountByChatSessionID(userID, params.ChatSessionID)
	if err != nil {
		return res, err
	}

	histories, err := dao.NewChatHistoryDao().GetByChatSessionID(
		userID, params.ChatSessionID, params.PageIndex, params.PageSize,
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
