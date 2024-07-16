package dtos

type PropertyIdAddressDto struct {
	Id      string `db:"id" json:"id"`
	Address string `db:"address" json:"address"`
}
