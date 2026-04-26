package services

import (
	"context"
	"errors"

	"github.com/oldfarmer96/vehicle-control-go/internal/models"
	"github.com/oldfarmer96/vehicle-control-go/internal/store"
)

type VehicleService struct {
	vehicleStore *store.VehicleStore
}

func NewVehicleService(store *store.VehicleStore) *VehicleService {
	return &VehicleService{vehicleStore: store}
}

func (s *VehicleService) CreateVehicle(ctx context.Context, payload models.CreaateVehicleDTO) (*models.Vehicle, error) {
	_, err := s.vehicleStore.FindByPlaca(ctx, payload.Placa)
	if err == nil {
		return nil, errors.New("la placa ya esta registrada")
	}

	return s.vehicleStore.Create(ctx, payload)
}
