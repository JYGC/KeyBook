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
	AddCobrandAdmin(userId, cobrandId string) (dtos.CobrandAdminDto, error)
	RemoveCobrandAdmin(id string) error
	AddPropertyManager(cobrandId, propertyId string) (dtos.CobrandPropertyManagerDto, error)
	RemovePropertyManager(id string) error
}

type CobrandApplicationService struct {
	cobrandService             services.ICobrandService
	cobrandRepo                repositories.ICobrandRepository
	cobrandAdminRepo           repositories.ICobrandAdminRepository
	cobrandPropertyManagerRepo repositories.ICobrandPropertyManagerRepository
}

func NewCobrandApplicationService(
	cobrandService services.ICobrandService,
	cobrandRepo repositories.ICobrandRepository,
	cobrandAdminRepo repositories.ICobrandAdminRepository,
	cobrandPropertyManagerRepo repositories.ICobrandPropertyManagerRepository,
) ICobrandApplicationService {
	return &CobrandApplicationService{cobrandService, cobrandRepo, cobrandAdminRepo, cobrandPropertyManagerRepo}
}

func (s *CobrandApplicationService) CreateCobrand(name string) (dtos.CobrandDto, error) {
	if err := s.cobrandService.ValidateCobrand(name); err != nil {
		return dtos.CobrandDto{}, err
	}
	return s.cobrandRepo.CreateCobrand(name)
}

func (s *CobrandApplicationService) UpdateCobrand(id, name string) error {
	if err := s.cobrandService.ValidateCobrand(name); err != nil {
		return err
	}
	return s.cobrandRepo.UpdateCobrand(id, name)
}

func (s *CobrandApplicationService) DeleteCobrand(id string) error {
	return s.cobrandRepo.DeleteCobrand(id)
}

func (s *CobrandApplicationService) AddCobrandAdmin(userId, cobrandId string) (dtos.CobrandAdminDto, error) {
	existing, err := s.cobrandAdminRepo.GetCobrandAdminsByCobrandId(cobrandId)
	if err != nil {
		return dtos.CobrandAdminDto{}, err
	}
	if err := s.cobrandService.EnsureAdminIsUnique(existing, userId); err != nil {
		return dtos.CobrandAdminDto{}, err
	}
	return s.cobrandAdminRepo.CreateCobrandAdmin(userId, cobrandId)
}

func (s *CobrandApplicationService) RemoveCobrandAdmin(id string) error {
	return s.cobrandAdminRepo.DeleteCobrandAdmin(id)
}

func (s *CobrandApplicationService) AddPropertyManager(cobrandId, propertyId string) (dtos.CobrandPropertyManagerDto, error) {
	return s.cobrandPropertyManagerRepo.CreateCobrandPropertyManager(cobrandId, propertyId)
}

func (s *CobrandApplicationService) RemovePropertyManager(id string) error {
	return s.cobrandPropertyManagerRepo.DeleteCobrandPropertyManager(id)
}
