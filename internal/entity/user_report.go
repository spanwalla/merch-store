package entity

type UserReport struct {
	Coins     int `db:"coins" json:"coins"`
	Inventory []struct {
		Type     string `json:"type"`
		Quantity int    `json:"quantity"`
	} `db:"inventory" json:"inventory"`
	CoinHistory struct {
		Received []struct {
			FromUser string `json:"fromUser"`
			Amount   int    `json:"amount"`
		} `json:"received"`
		Sent []struct {
			ToUser string `json:"toUser"`
			Amount int    `json:"amount"`
		} `json:"sent"`
	} `db:"coin_history" json:"coinHistory"`
}
