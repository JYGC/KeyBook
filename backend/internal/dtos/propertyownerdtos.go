package dtos

type PropertyOwnerDto struct {
	Id       string `db:"id" json:"id"`
	Property string `db:"property" json:"property"`
}

type PersonPropertyOwnerDto struct {
	Id            string `db:"id" json:"id"`
	Person        string `db:"person" json:"person"`
	PropertyOwner string `db:"propertyOwner" json:"propertyOwner"`
}
