package dtos

type AddDeviceHistoriesDto struct {
	ActionDescription string
	DateSpecified     string
}

type AddPersonDeviceHistoriesDto struct {
	DeviceHolder      string
	DateSpecified     string
	ActionDescription string
}

type AddDeviceAndHistoriesDto struct {
	Name                  string
	Identifier            string
	Type                  string
	DefunctReason         string
	CurrentHolder         string
	DeviceHistories       []AddDeviceHistoriesDto
	PersonDeviceHistories []AddPersonDeviceHistoriesDto
}

type AddPropertyDeviceAndHistoriesDto struct {
	PropertyAddress                  string
	DevicesPersonDevicesAndHistories []AddDeviceAndHistoriesDto
}
