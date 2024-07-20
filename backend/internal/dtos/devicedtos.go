package dtos

type NewDeviceDto struct {
	Name          string
	Identifier    string
	Type          string
	DefunctReason string
}

type DeviceDto struct {
	Id            string `db:"id" json:"id"`
	Name          string `db:"name" json:"name"`
	Identifier    string `db:"identifier" json:"identifier"`
	Type          string `db:"type" json:"type"`
	DefunctReason string `db:"defunctreason" json:"defunctReason"`
}
