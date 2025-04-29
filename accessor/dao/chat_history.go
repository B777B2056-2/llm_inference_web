package dao

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"llm_online_inference/accessor/resource"
)

const (
	ChatStatusDoing  = "Doing"
	ChatStatusDone   = "Done"
	ChatStatusFailed = "Failed"
)

// ChatHistory 对话历史
type ChatHistory struct {
	gorm.Model
	ChatSessionID   string `gorm:"column:chat_session_id"`
	MessageID       string `gorm:"column:message_id"`
	ParentMessageID string `gorm:"column:parent_message_id"`
	UserID          int    `gorm:"column:user_id"`
	Prompt          string `gorm:"column:prompt"`
	Answer          string `gorm:"column:answer"`
	PromptTokenCnt  int    `gorm:"column:prompt_token_cnt"`
	AnswerTokenCnt  int    `gorm:"column:answer_token_cnt"`
	Status          string `gorm:"column:status"`
}

func (c ChatHistory) TableName() string {
	return "chat_history"
}

func (c *ChatHistory) BeforeCreate(tx *gorm.DB) error {
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

func (c *ChatHistoryDao) Create(userID int, chatSessionID, parentMessageID, prompt string,
	tokenCnt int) (pk uint, err error) {
	record := ChatHistory{
		UserID:          userID,
		ChatSessionID:   chatSessionID,
		ParentMessageID: parentMessageID,
		Prompt:          prompt,
		PromptTokenCnt:  tokenCnt,
	}
	err = c.DB.Create(&record).Error
	if err != nil {
		return 0, err
	}
	return record.ID, nil
}

func (c *ChatHistoryDao) UpdateAnswer(pk uint, answer string, tokenCnt int) error {
	return c.DB.
		Where("id = ?", pk).
		Updates(
			&ChatHistory{
				Answer:         answer,
				AnswerTokenCnt: tokenCnt,
				Status:         ChatStatusDone,
			},
		).Error
}

func (c *ChatHistoryDao) UpdateStatus2Failed(pk uint) error {
	return c.DB.Where("id = ?", pk).Updates(&ChatHistory{Status: ChatStatusFailed}).Error
}

func (c *ChatHistoryDao) CountChatSessionContextTokens(userID int, chatSessionID string) (uint64, error) {
	var totalTokens struct {
		Total uint64 `gorm:"column:total_tokens"`
	}

	err := c.DB.Model(&ChatHistory{}).
		Select("SUM(prompt_token_cnt + answer_token_cnt) AS total_tokens").
		Where("user_id = ?", userID).
		Where("chat_session_id = ?", chatSessionID).
		Scan(&totalTokens).Error
	if err != nil {
		return 0, err
	}
	return totalTokens.Total, nil
}

func (c *ChatHistoryDao) GetByChatSessionID(userID int,
	chatSessionID string, pageIdx, pageSize int) (records []ChatHistory, err error) {
	offset := (pageIdx - 1) * pageSize
	err = c.DB.Model(&ChatHistory{}).
		Where("user_id = ?", userID).
		Where("chat_session_id = ?", chatSessionID).
		Order("id asc").
		Offset(offset).Limit(pageSize).
		Find(&records).Error
	return
}

func (c *ChatHistoryDao) CountByChatSessionID(userID int, chatSessionID string) (count int64, err error) {
	err = c.DB.Model(&ChatHistory{}).
		Where("user_id = ?", userID).
		Where("chat_session_id = ?", chatSessionID).
		Order("id asc").
		Count(&count).Error
	return
}
