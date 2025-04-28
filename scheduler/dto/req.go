package dto

type ChatCompletionReq struct {
	ChatSessionId   string `json:"chat_session_id"`            // 对话id
	ParentMessageId string `json:"parent_message_id"`          // 当前prompt的父问答id
	Prompt          string `json:"prompt" validate:"required"` // 当前prompt
}

type ChatHistoryReq struct {
	ChatSessionID string `json:"chat_session_id" validate:"required"`
	Page
}

type CreateBatchInferenceTaskReq struct {
	BatchInferenceName string `json:"batch_inference_name" validate:"required"`
	Prompts            string `json:"prompts" validate:"required"`
}

type GetBatchInferenceTaskResultsReq struct {
	BatchInferenceName string `json:"batch_inference_name" validate:"required"`
	Page
}
