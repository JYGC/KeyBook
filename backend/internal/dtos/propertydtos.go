package dtos

type PropertyIdAddressDto struct {
	Id      string `db:"id" json:"id"`
	Address string `db:"address" json:"address"`
}

type PropertyAddressOwnersDtos struct {
	Address string `db:"address" json:"address"`
	Owners  string `db:"owners" json:"owners"`
}

type PropertyFullDto struct {
	Id       string   `db:"id" json:"id"`
	Address  string   `db:"address" json:"address"`
	Owners   []string `db:"owners" json:"owners"`
	Managers []string `db:"managers" json:"managers"`
}

type PropertyIdFromJsonDto struct {
	Property string `json:"property"`
}
