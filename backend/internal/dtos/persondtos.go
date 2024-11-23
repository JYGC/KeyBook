package dtos

type PersonDto struct {
	Id   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Type string `db:"type" json:"type"`
}

type PersonIdNameDto struct {
	Id   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}
