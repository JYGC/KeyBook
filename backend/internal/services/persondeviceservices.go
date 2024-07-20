package services

import "keybook/backend/internal/repositories"

type IPersonDeviceServices interface {
}

type PersonDeviceServices struct {
	personDeviceRepository repositories.IPersonDeviceRepository
}

func NewPersonDeviceServices(personDeviceRepository repositories.IPersonDeviceRepository) IPersonDeviceServices {
	return &PersonDeviceServices{
		personDeviceRepository,
	}
}
