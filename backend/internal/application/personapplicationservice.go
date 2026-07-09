package application

import (
	"keybook/backend/internal/dtos"
	"keybook/backend/internal/repositories"
	"keybook/backend/internal/services"
)

type IPersonApplicationService interface {
	CreatePerson(name, dob, userID string) (dtos.PersonDto, error)
	UpdatePerson(id, name, dob string) error
	DeletePerson(id string) error
}

type PersonApplicationService struct {
	personService services.IPersonService
	personRepo    repositories.IPersonRepository
}

func NewPersonApplicationService(
	personService services.IPersonService,
	personRepo repositories.IPersonRepository,
) IPersonApplicationService {
	return &PersonApplicationService{personService, personRepo}
}

func (s *PersonApplicationService) CreatePerson(name, dob, userID string) (dtos.PersonDto, error) {
	if err := s.personService.ValidatePerson(name, dob); err != nil {
		return dtos.PersonDto{}, err
	}
	return s.personRepo.CreatePerson(name, dob, userID)
}

func (s *PersonApplicationService) UpdatePerson(id, name, dob string) error {
	if err := s.personService.ValidatePerson(name, dob); err != nil {
		return err
	}
	return s.personRepo.UpdatePerson(id, name, dob)
}

func (s *PersonApplicationService) DeletePerson(id string) error {
	return s.personRepo.DeletePerson(id)
}
