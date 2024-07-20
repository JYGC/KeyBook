package dtos

type PersonIdNameDto struct {
	Id   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}
