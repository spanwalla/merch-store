package entity

type Item struct {
	Id    int    `db:"id"`
	Name  string `db:"name"`
	Price int    `db:"price"`
}
