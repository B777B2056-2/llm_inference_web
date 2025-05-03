package services

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/segmentio/kafka-go"
	"llm_inference_web/accessor/client"
	"llm_inference_web/accessor/dao"
	"llm_inference_web/accessor/dto"
)

type BatchInferenceOperator struct {
	userID int
}

func NewBatchInferenceOperator(userID int) *BatchInferenceOperator {
	return &BatchInferenceOperator{userID: userID}
}

// CreateTask 创建批量推理任务
func (b *BatchInferenceOperator) CreateTask(ctx *gin.Context, param dto.CreateBatchInferenceTaskReq) error {
	// 检查任务名称是否已存在
	nameIsExits, err := dao.NewBatchInferenceResultsDao().CheckTaskNameIsExists(ctx, b.userID, param.BatchInferenceName)
	if err != nil {
		return err
	}
	if nameIsExits {
		return errors.New("task name already exists")
	}
	// 生成任务id
	uid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	taskId := uid.String()
	// 塞入kafka
	productor := client.NewKafkaProductorClient(client.TopicBatchInferenceRequests)
	traceId, _ := ctx.Get("trace_id")
	payload := dto.BatchInferenceKafkaPayload{
		TraceId:            traceId.(string),
		BatchInferenceId:   taskId,
		BatchInferenceName: param.BatchInferenceName,
		UserId:             b.userID,
		SamplingParams:     param.InferenceParams,
		Prompts:            param.Prompts,
	}
	payloadBytes, _ := json.Marshal(payload)
	msgs := []kafka.Message{
		{
			Value: payloadBytes,
		},
	}
	return productor.Send(ctx, msgs)
}

// TaskResults 获取批量推理任务结果
func (b *BatchInferenceOperator) TaskResults(ctx *gin.Context, param dto.GetBatchInferenceTaskResultsReq) (
	dto.BatchInferenceTaskResp, error) {
	results, err := dao.NewBatchInferenceResultsDao().GetByBatchInferenceName(
		ctx, b.userID, param.BatchInferenceName, param.PageIndex, param.PageSize,
	)
	if err != nil {
		return dto.BatchInferenceTaskResp{}, err
	}

	resp := dto.BatchInferenceTaskResp{BatchInferenceName: param.BatchInferenceName}
	for _, result := range results {
		resp.BatchInferenceId = result.BatchInferenceId
		resp.Results = append(resp.Results, dto.BatchInferenceResultItem{
			Prompt: result.Prompt,
			Answer: result.Text,
		})
	}

	resp.TotalCount, err = dao.NewBatchInferenceResultsDao().CountByBatchInferenceName(
		ctx, b.userID, param.BatchInferenceName,
	)
	return resp, err
}
