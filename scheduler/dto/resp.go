package dto

type ChatHistoryResp struct {
	MessageID       string `json:"message_id"`
	ParentMessageID string `json:"parent_message_id"`
	Prompt          string `json:"prompt"`
	Answer          string `json:"answer"`
	TokenCnt        uint32 `json:"token_cnt"`
	Timestamp       int64  `json:"timestamp"`
	Status          string `json:"status"`
}
