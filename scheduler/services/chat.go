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

func ChatCompletion(ctx *gin.Context, userID int, params *dto.ChatCompletionReq) error {
	// 调用分词服务：获取原先上下文+当前prompt的滑动窗口分词结果（窗口大小由分词服务控制）
	inputIDs, attnMasks, promptTokenCnt, err := client.NewTokenizer().Do(ctx, params.ParentMessageId, params.Prompt)
	if err != nil {
		return err
	}

	// 创建对话历史：如果为第一次对话，会自动创建各种id
	historyID, err := dao.NewChatHistoryDao().Create(
		userID, params.ChatSessionId, params.ParentMessageId, params.Prompt, promptTokenCnt,
	)
	if err != nil {
		return err
	}

	// 调用模型服务，开启流式传输
	modelAnswerStream, err := client.NewModelServer().ChatCompletion(ctx, inputIDs, attnMasks, params.ChatSessionId)
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
	var answerTokenCnt uint32
	for {
		// 接收数据
		resp, err := modelAnswerStream.Recv()
		if err != nil {
			if err == io.EOF {
				flusher.Flush()
				_ = dao.NewChatHistoryDao().UpdateAnswer(historyID, modelAnswer, answerTokenCnt)
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
		// SSE 格式数据
		data := resp.GetAnswer()
		answerTokenCnt += resp.GetTokenCnt()
		_, _ = fmt.Fprintf(ctx.Writer, "event: customEvent\n")
		_, _ = fmt.Fprintf(ctx.Writer, "data: %s\n\n", data)
		modelAnswer += data
		flusher.Flush()
	}
}

func GetChatHistory(_ *gin.Context, userID int, chatSessionID string) (res []dto.ChatHistoryResp, err error) {
	histories, err := dao.NewChatHistoryDao().GetByChatSessionID(userID, chatSessionID)
	if err != nil {
		return nil, err
	}
	if len(histories) == 0 {
		return nil, nil
	}

	for _, h := range histories {
		res = append(res, dto.ChatHistoryResp{
			MessageID:       h.MessageID,
			ParentMessageID: h.ParentMessageID,
			Prompt:          h.Prompt,
			Answer:          h.Answer,
			TokenCnt:        h.PromptTokenCnt,
			Timestamp:       h.CreatedAt.Unix(),
			Status:          h.Status,
		})
	}
	return
}
