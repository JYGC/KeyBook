package dtos

type ItemDto struct {
	Id          string `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	Picture     string `db:"picture" json:"picture"`
}

type PropertyItemDto struct {
	Id       string `db:"id" json:"id"`
	Item     string `db:"item" json:"item"`
	Property string `db:"property" json:"property"`
}

type PersonItemDto struct {
	Id     string `db:"id" json:"id"`
	Person string `db:"person" json:"person"`
	Item   string `db:"item" json:"item"`
}

type EntryDeviceDto struct {
	Id            string `db:"id" json:"id"`
	Item          string `db:"item" json:"item"`
	DeviceType    string `db:"deviceType" json:"deviceType"`
	Identifier    string `db:"identifier" json:"identifier"`
	DefunctReason string `db:"defunctReason" json:"defunctReason"`
}
