package dao

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"llm_inference_web/accessor/resource"
)

type BatchInferenceResult struct {
	BatchInferenceId   string `bson:"batch_inference_id"`
	BatchInferenceName string `bson:"batch_inference_name"`
	Prompt             string `bson:"prompt"`
	Text               string `bson:"text"`
}

type BatchInferenceResultsDao struct {
	c *mongo.Collection
}

func NewBatchInferenceResultsDao() *BatchInferenceResultsDao {
	return &BatchInferenceResultsDao{
		c: resource.MongoDB.Collection("batch_inference_results"),
	}
}

// GetByBatchInferenceName 根据id查询记录
func (b *BatchInferenceResultsDao) GetByBatchInferenceName(ctx context.Context,
	userId int, batchInferenceName string, pageIdx, pageSize int) (results []BatchInferenceResult, err error) {
	filter := bson.D{
		{"user_id", userId},
		{"batch_inference_name", batchInferenceName},
	}

	// 设置分页参数
	opts := options.Find()
	opts.SetSort(bson.D{{"_id", 1}})              // 按id排序
	opts.SetSkip(int64((pageIdx - 1) * pageSize)) // 跳过记录数
	opts.SetLimit(int64(pageSize))                // 每页数量

	cursor, err := b.c.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// CountByBatchInferenceName 计数
func (b *BatchInferenceResultsDao) CountByBatchInferenceName(ctx context.Context,
	userId int, batchInferenceName string) (int64, error) {
	filter := bson.D{
		{"user_id", userId},
		{"batch_inference_name", batchInferenceName},
	}
	return b.c.CountDocuments(ctx, filter)
}

// CheckTaskNameIsExists 检查任务名称是否已存在（用户粒度）
func (b *BatchInferenceResultsDao) CheckTaskNameIsExists(ctx context.Context,
	userId int, batchInferenceName string) (bool, error) {
	filter := bson.D{
		{"user_id", userId},
		{"batch_inference_name", batchInferenceName},
	}

	var record bson.M
	if err := b.c.FindOne(ctx, filter).Decode(&record); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
