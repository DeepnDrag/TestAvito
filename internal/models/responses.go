package models

type TransactionsFromUser struct {
	ToUser string `json:"to_user"`
	Amount int    `json:"amount"`
}

type TransactionsToUser struct {
	FromUser string `json:"from_user"`
	Amount   int    `json:"amount"`
}
