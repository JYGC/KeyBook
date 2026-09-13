package dtos

type PersonDto struct {
	Id           string `db:"id" json:"id"`
	Name         string `db:"name" json:"name"`
	DOB          string `db:"DOB" json:"DOB"`
	User         string `db:"user" json:"user"`
	ProfileImage string `db:"profileImage" json:"profileImage"`
}
