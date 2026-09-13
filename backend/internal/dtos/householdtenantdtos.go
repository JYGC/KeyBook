package dtos

type HouseholdDto struct {
	Id       string `db:"id" json:"id"`
	Person   string `db:"person" json:"person"`
	Property string `db:"property" json:"property"`
}

type TenantDto struct {
	Id       string `db:"id" json:"id"`
	Person   string `db:"person" json:"person"`
	Property string `db:"property" json:"property"`
}
