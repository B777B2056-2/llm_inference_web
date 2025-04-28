package dto

type ChatHistoryItem struct {
	MessageID       string `json:"message_id"`
	ParentMessageID string `json:"parent_message_id"`
	Prompt          string `json:"prompt"`
	Answer          string `json:"answer"`
	TokenCnt        int    `json:"token_cnt"`
	Timestamp       int64  `json:"timestamp"`
	Status          string `json:"status"`
}

type ChatHistoryResp struct {
	TotalCount    int64             `json:"total_count"`
	ChatHistories []ChatHistoryItem `json:"results"`
}

type BatchInferenceResultItem struct {
	Prompt string `json:"prompt"`
	Answer string `json:"answer"`
}

type BatchInferenceTaskResp struct {
	TotalCount         int64                      `json:"total_count"`
	BatchInferenceId   string                     `json:"batch_inference_id"`
	BatchInferenceName string                     `json:"batch_inference_name"`
	Results            []BatchInferenceResultItem `json:"results"`
}
