package dtos

type CobrandDto struct {
	Id   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type CobrandAdminDto struct {
	Id       string `db:"id" json:"id"`
	User     string `db:"user" json:"user"`
	Cobrand  string `db:"cobrand" json:"cobrand"`
	Approved bool   `db:"approved" json:"approved"`
}

type CobrandPropertyManagerDto struct {
	Id       string `db:"id" json:"id"`
	Cobrand  string `db:"cobrand" json:"cobrand"`
	Property string `db:"property" json:"property"`
}

type CobrandPropertyOwnerDto struct {
	Id            string `db:"id" json:"id"`
	Cobrand       string `db:"cobrand" json:"cobrand"`
	PropertyOwner string `db:"propertyOwner" json:"propertyOwner"`
}
