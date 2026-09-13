package application

import (
	"keybook/backend/internal/dtos"
	"keybook/backend/internal/repositories"
	"keybook/backend/internal/services"
)

type ICobrandApplicationService interface {
	CreateCobrand(name string) (dtos.CobrandDto, error)
	UpdateCobrand(id, name string) error
	DeleteCobrand(id string) error
	AddCobrandAdmin(userId, cobrandId, inviterUserId string) (dtos.CobrandAdminDto, error)
	RemoveCobrandAdmin(id string) error
	AddPropertyManager(cobrandId, propertyId string) (dtos.CobrandPropertyManagerDto, error)
	RemovePropertyManager(id string) error
}

type CobrandApplicationService struct {
	cobrandService                   services.ICobrandService
	cobrandRepository                repositories.ICobrandRepository
	cobrandAdminRepository           repositories.ICobrandAdminRepository
	cobrandPropertyManagerRepository repositories.ICobrandPropertyManagerRepository
}

func NewCobrandApplicationService(
	cobrandService services.ICobrandService,
	cobrandRepository repositories.ICobrandRepository,
	cobrandAdminRepository repositories.ICobrandAdminRepository,
	cobrandPropertyManagerRepository repositories.ICobrandPropertyManagerRepository,
) ICobrandApplicationService {
	return &CobrandApplicationService{cobrandService, cobrandRepository, cobrandAdminRepository, cobrandPropertyManagerRepository}
}

func (s *CobrandApplicationService) CreateCobrand(name string) (dtos.CobrandDto, error) {
	if err := s.cobrandService.ValidateCobrand(name); err != nil {
		return dtos.CobrandDto{}, err
	}
	return s.cobrandRepository.CreateCobrand(name)
}

func (s *CobrandApplicationService) UpdateCobrand(id, name string) error {
	if err := s.cobrandService.ValidateCobrand(name); err != nil {
		return err
	}
	return s.cobrandRepository.UpdateCobrand(id, name)
}

func (s *CobrandApplicationService) DeleteCobrand(id string) error {
	return s.cobrandRepository.DeleteCobrand(id)
}

func (s *CobrandApplicationService) AddCobrandAdmin(userId, cobrandId, inviterUserId string) (dtos.CobrandAdminDto, error) {
	existingAdmins, err := s.cobrandAdminRepository.GetCobrandAdminsByCobrandId(cobrandId)
	if err != nil {
		return dtos.CobrandAdminDto{}, err
	}
	if err := s.cobrandService.EnsureAdminIsUnique(existingAdmins, userId); err != nil {
		return dtos.CobrandAdminDto{}, err
	}
	if err := s.cobrandService.EnsureInviterIsApprovedAdmin(existingAdmins, inviterUserId); err != nil {
		return dtos.CobrandAdminDto{}, err
	}
	return s.cobrandAdminRepository.CreateCobrandAdmin(userId, cobrandId)
}

func (s *CobrandApplicationService) RemoveCobrandAdmin(id string) error {
	return s.cobrandAdminRepository.DeleteCobrandAdmin(id)
}

func (s *CobrandApplicationService) AddPropertyManager(cobrandId, propertyId string) (dtos.CobrandPropertyManagerDto, error) {
	return s.cobrandPropertyManagerRepository.CreateCobrandPropertyManager(cobrandId, propertyId)
}

func (s *CobrandApplicationService) RemovePropertyManager(id string) error {
	return s.cobrandPropertyManagerRepository.DeleteCobrandPropertyManager(id)
}
