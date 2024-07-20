package dtos

import "time"

type NewDeviceDto struct {
	Name          string
	Identifier    string
	Type          string
	DefunctReason string
	Histories     []NewDeviceHistoryDto
}

type NewDeviceHistoryDto struct {
	Description    string
	StatedDateTime time.Time
}

type DeviceDto struct {
	Id            string `db:"id" json:"id"`
	Name          string `db:"name" json:"name"`
	Identifier    string `db:"identifier" json:"identifier"`
	Type          string `db:"type" json:"type"`
	DefunctReason string `db:"defunctreason" json:"defunctReason"`
}
