package entity

type User struct {
	Id       int    `db:"id"`
	Name     string `db:"name"`
	Password string `db:"password"`
	Balance  int    `db:"balance"`
}
