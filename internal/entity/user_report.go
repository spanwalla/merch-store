package entity

type ReceivedTransaction struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type SentTransaction struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type CoinHistory struct {
	Received []ReceivedTransaction `json:"received"`
	Sent     []SentTransaction     `json:"sent"`
}

type Inventory struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type UserReport struct {
	Coins       int         `db:"coins" json:"coins"`
	Inventory   []Inventory `db:"inventory" json:"inventory"`
	CoinHistory CoinHistory `db:"coin_history" json:"coinHistory"`
}
