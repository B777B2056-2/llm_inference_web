package services

import (
	"github.com/gin-gonic/gin"
	"llm_online_inference/scheduler/dao"
	"llm_online_inference/scheduler/dto"
)

type BatchInferenceOperator struct{}

func NewBatchInferenceOperator() *BatchInferenceOperator {
	return &BatchInferenceOperator{}
}

// CreateTask 创建批量推理任务
func (b *BatchInferenceOperator) CreateTask(ctx *gin.Context, userID int, param dto.CreateBatchInferenceTaskReq) error {
	// TODO 开启mongodb事务
	// TODO 检查任务名称是否已存在
	// TODO 生成任务id
	// TODO 塞入kafka
	return nil
}

// TaskResults 获取批量推理任务结果
func (b *BatchInferenceOperator) TaskResults(ctx *gin.Context, userID int,
	param dto.GetBatchInferenceTaskResultsReq) (dto.BatchInferenceTaskResp, error) {
	results, err := dao.NewBatchInferenceResultsDao().GetByBatchInferenceName(
		ctx, userID, param.BatchInferenceName, param.PageIndex, param.PageSize,
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
		ctx, userID, param.BatchInferenceName,
	)
	return resp, err
}
