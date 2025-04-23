package dto

type ChatCompletionReq struct {
	ChatSessionId   string `json:"chat_session_id"`            // 对话id
	ParentMessageId string `json:"parent_message_id"`          // 当前prompt的父问答id
	Prompt          string `json:"prompt" validate:"required"` // 当前prompt
}
