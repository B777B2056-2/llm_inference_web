package dto

type Page struct {
	PageIndex int `json:"page_index" validate:"required, min=1"`
	PageSize  int `json:"page_size" validate:"required, min=1"`
}
