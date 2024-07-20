package dtos

import "time"

type NewPersonDeviceAndHistoriesDto struct {
	Person    string
	Device    string
	Histories []NewPersonDeviceHistoryDto
}

type NewPersonDeviceHistoryDto struct {
	Person         string
	Device         string
	StatedDateTime time.Time
	Description    string
}

type PersonDeviceDto struct {
	Id     string `db:"id" json:"id"`
	Person string `db:"person" json:"person"`
	Device string `db:"device" json:"device"`
}
