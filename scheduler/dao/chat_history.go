package dao

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"llm_online_inference/scheduler/resource"
)

const (
	ChatStatusDoing  = "Doing"
	ChatStatusDone   = "Done"
	ChatStatusFailed = "Failed"
)

// ChatHistory 业务层对话历史
type ChatHistory struct {
	gorm.Model
	ChatSessionID   string `gorm:"column:chat_session_id"`
	MessageID       string `gorm:"column:message_id"`
	ParentMessageID string `gorm:"column:parent_message_id"`
	UserID          int    `gorm:"column:user_id"`
	Prompt          string `gorm:"column:prompt"`
	Answer          string `gorm:"column:answer"`
	PromptTokenCnt  uint32 `gorm:"column:prompt_token_cnt"`
	AnswerTokenCnt  uint32 `gorm:"column:answer_token_cnt"`
	Status          string `gorm:"column:status"`
}

func (c ChatHistory) TableName() string {
	return "chat_history"
}

func (c ChatHistory) BeforeCreate(_ *gorm.DB) error {
	c.Status = ChatStatusDoing
	if c.ChatSessionID == "" {
		id, err := uuid.NewV4()
		if err != nil {
			return err
		}
		c.ChatSessionID = id.String()
	}
	if c.MessageID == "" {
		id, err := uuid.NewV4()
		if err != nil {
			return err
		}
		c.MessageID = id.String()
	}
	return nil
}

type ChatHistoryDao struct {
	*gorm.DB
}

func NewChatHistoryDao() *ChatHistoryDao {
	return &ChatHistoryDao{resource.DB}
}

func NewChatHistoryDaoWithTx(tx *gorm.DB) *ChatHistoryDao {
	return &ChatHistoryDao{tx}
}

func (c *ChatHistoryDao) Create(userID int, chatSessionID,
	parentMessageID, prompt string, curTokenCnt uint32) (pk uint, err error) {
	record := ChatHistory{
		UserID:          userID,
		ChatSessionID:   chatSessionID,
		ParentMessageID: parentMessageID,
		Prompt:          prompt,
		PromptTokenCnt:  curTokenCnt,
	}
	err = c.DB.Create(&record).Error
	if err != nil {
		return 0, err
	}
	return record.ID, nil
}

func (c *ChatHistoryDao) UpdateAnswer(pk uint, answer string, answerTokenCnt uint32) error {
	return c.DB.
		Where("id = ?", pk).
		Updates(&ChatHistory{Answer: answer, AnswerTokenCnt: answerTokenCnt, Status: ChatStatusDone}).Error
}

func (c *ChatHistoryDao) UpdateStatus2Failed(pk uint) error {
	return c.DB.Where("id = ?", pk).Updates(&ChatHistory{Status: ChatStatusFailed}).Error
}

func (c *ChatHistoryDao) GetByChatSessionID(userID int, chatSessionID string) (records []ChatHistory, err error) {
	err = c.DB.
		Where("user_id = ?", userID).
		Where("chat_session_id = ?", chatSessionID).
		Find(&records).Error
	return
}
