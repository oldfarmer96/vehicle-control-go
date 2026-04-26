package services

import (
	"context"
	"errors"
	"math"

	"github.com/oldfarmer96/vehicle-control-go/internal/models"
	"github.com/oldfarmer96/vehicle-control-go/internal/store"
)

type VehicleService struct {
	vehicleStore *store.VehicleStore
	personaStore *store.PersonaStore
}

func NewVehicleService(store *store.VehicleStore, personaStore *store.PersonaStore) *VehicleService {
	return &VehicleService{vehicleStore: store, personaStore: personaStore}
}

func (s *VehicleService) CreateVehicle(ctx context.Context, payload models.CreaateVehicleDTO) (*models.Vehicle, error) {
	_, err := s.vehicleStore.FindByPlaca(ctx, payload.Placa)
	if err == nil {
		return nil, errors.New("la placa ya esta registrada")
	}

	return s.vehicleStore.Create(ctx, payload)
}

func (s *VehicleService) GetAllVehicles(ctx context.Context, page, limit int, placa string) (*models.VehicleListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	vehicles, total, err := s.vehicleStore.GetAll(ctx, page, limit, placa)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &models.VehicleListResponse{
		Vehiculos:  vehicles,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *VehicleService) GetVehicleByPlaca(ctx context.Context, placa string) (*models.Vehicle, error) {
	return s.vehicleStore.FindByPlacaWithOwner(ctx, placa)
}

func (s *VehicleService) AssignOwner(ctx context.Context, vehiculoID, personaID string) (*models.Vehicle, error) {
	_, err := s.vehicleStore.FindByID(ctx, vehiculoID)
	if err != nil {
		return nil, err
	}

	_, err = s.personaStore.FindByID(ctx, personaID)
	if err != nil {
		return nil, err
	}

	err = s.vehicleStore.AssignOwner(ctx, vehiculoID, personaID)
	if err != nil {
		return nil, err
	}

	return s.vehicleStore.GetVehicleWithOwner(ctx, vehiculoID)
}
