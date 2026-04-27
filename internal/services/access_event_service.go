package services

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/oldfarmer96/vehicle-control-go/internal/models"
	"github.com/oldfarmer96/vehicle-control-go/internal/store"
	"github.com/oldfarmer96/vehicle-control-go/pkg/external"
)

type AccessEventService struct {
	store       *store.AccessEventStore
	placaClient *external.PlacaClient
}

func NewAccessEventService(store *store.AccessEventStore, client *external.PlacaClient) *AccessEventService {
	return &AccessEventService{
		store:       store,
		placaClient: client,
	}
}

func (s *AccessEventService) ProcessAccessEvent(ctx context.Context, req *models.AccessEventRequest) (*models.AccessEventResponse, error) {
	var vehicle *models.Vehicle
	var source string
	var found bool

	vehicle, err := s.store.FindByPlaca(ctx, req.Placa)
	if err == nil {
		found = true
		source = "local"
		log.Printf("Vehículo encontrado en BD local: %s", req.Placa)
	} else {
		log.Printf("Vehículo no encontrado en BD local, consultando API externa: %s", req.Placa)
		vehicle, source, err = s.createVehicleFromExternalOrUnknown(ctx, req.Placa)
		if err != nil {
			return nil, err
		}
		found = (source != "unknown")
	}

	err = s.store.CreateAccessEvent(ctx, vehicle.ID, req.Evento, req.PuntoControl, req.ConfianzaOcr)
	if err != nil {
		return nil, err
	}

	return &models.AccessEventResponse{
		Evento:  *req,
		Vehicle: *vehicle,
		Source:  source,
		Found:   found,
	}, nil
}

func (s *AccessEventService) createVehicleFromExternalOrUnknown(ctx context.Context, placa string) (*models.Vehicle, string, error) {
	if s.placaClient != nil {
		externalResp, err := s.placaClient.GetPlacaData(ctx, placa)
		if err == nil && externalResp.Success {
			marca := externalResp.Data.Marca
			modelo := externalResp.Data.Modelo
			color := externalResp.Data.Color
			vin := externalResp.Data.Vin
			motor := externalResp.Data.Motor

			vehicle, err := s.store.CreateVehicle(ctx, placa, &marca, &modelo, &color, &vin, &motor)
			if err == nil {
				return vehicle, "external_api", nil
			}
			log.Printf("Error al crear vehículo desde API externa: %v", err)
		} else {
			log.Printf("API externa no devolvió datos para placa: %s", placa)
		}
	}

	vehicle, err := s.store.CreateVehicle(ctx, placa, nil, nil, nil, nil, nil)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || err != nil {
			return nil, "", errors.New("error al crear vehículo con placa desconocida")
		}
		return nil, "", err
	}

	return vehicle, "unknown", nil
}