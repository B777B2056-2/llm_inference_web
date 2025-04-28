package dto

type InferenceParams struct {
	PresencePenalty   float32 `json:"presence_penalty" validate:"required"`
	FrequencyPenalty  float32 `json:"frequency_penalty" validate:"required"`
	RepetitionPenalty float32 `json:"repetition_penalty" validate:"required"`
	Temperature       float32 `json:"temperature" validate:"required"`
	TopP              float32 `json:"top_p" validate:"required"`
	TopK              int32   `json:"top_k" validate:"required"`
}

type ChatCompletionReq struct {
	ChatSessionId   string `json:"chat_session_id"`            // 对话id
	ParentMessageId string `json:"parent_message_id"`          // 当前prompt的父问答id
	Prompt          string `json:"prompt" validate:"required"` // 当前prompt
	InferenceParams
}

type ChatHistoryReq struct {
	ChatSessionID string `json:"chat_session_id" validate:"required"`
	Page
}

type CreateBatchInferenceTaskReq struct {
	BatchInferenceName string   `json:"batch_inference_name" validate:"required"`
	Prompts            []string `json:"prompts" validate:"required"`
	InferenceParams
}

type GetBatchInferenceTaskResultsReq struct {
	BatchInferenceName string `json:"batch_inference_name" validate:"required"`
	Page
}
