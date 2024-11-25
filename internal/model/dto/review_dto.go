package dto

type CreateReviewRequest struct {
	Rating  int    `json:"rating" validate:"required,min=1,max=5"`
	Comment string `json:"comment" validate:"required,min=3,max=100"`
	MenuId  string `json:"menu_id" validate:"required"`
}

type UpdateReviewRequest struct {
	Rating  int    `json:"rating,omitempty" validate:"omitempty,min=1,max=5"`
	Comment string `json:"comment,omitempty" validate:"omitempty,min=3,max=100"`
}
