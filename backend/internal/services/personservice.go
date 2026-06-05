package services

import "errors"

type IPersonService interface {
	ValidatePerson(name, dob string) error
}

type PersonService struct{}

func NewPersonService() IPersonService {
	return &PersonService{}
}

func (s *PersonService) ValidatePerson(name, dob string) error {
	if name == "" {
		return errors.New("name is required")
	}
	if dob == "" {
		return errors.New("DOB is required")
	}
	return nil
}
