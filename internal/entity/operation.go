package entity

type Operation struct {
	Id         int `db:"id"`
	SenderId   int `db:"sender_id"`
	ReceiverId int `db:"receiver_id"`
	Amount     int `db:"amount"`
}
