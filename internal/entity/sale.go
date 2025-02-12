package entity

type Sale struct {
	Id       int `db:"id"`
	UserId   int `db:"user_id"`
	ItemId   int `db:"item_id"`
	Quantity int `db:"quantity"`
}
