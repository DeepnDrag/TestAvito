package models

type AuthorizeUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SendCoinRequest struct {
	RecipientUsername string `json:"recipient_username" validate:"required"`
	Amount            int    `json:"amount" validate:"required,min=1"`
}

type BuyItemRequest struct {
	Quantity int `json:"quantity" validate:"required, min=1"`
}
