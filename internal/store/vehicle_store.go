package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oldfarmer96/vehicle-control-go/internal/models"
)

type VehicleStore struct {
	db *pgxpool.Pool
}

func NewVehicleStore(db *pgxpool.Pool) *VehicleStore {
	return &VehicleStore{db: db}
}

func (s *VehicleStore) FindByPlaca(ctx context.Context, placa string) (*models.Vehicle, error) {
	query := `
	SELECT id, placa, marca, modelo, color, vin, motor, created_at, updated_at
	FROM vehiculos
	WHERE placa = $1
	`

	var v models.Vehicle
	err := s.db.QueryRow(ctx, query, placa).Scan(

		&v.ID, &v.Placa, &v.Marca, &v.Modelo, &v.Color, &v.Vin, &v.Motor, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("vehiculo no encontrado")
		}
	}

	return &v, nil
}

func (s *VehicleStore) Create(ctx context.Context, payload models.CreaateVehicleDTO) (*models.Vehicle, error) {
	query := `
	INSERT INTO vehiculos (placa, marca, modelo, color, vin, motor) 
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, placa, marca, modelo, color, vin, motor, created_at, updated_at
	`

	var v models.Vehicle

	err := s.db.QueryRow(ctx, query, payload.Placa, payload.Marca, payload.Modelo, payload.Color, payload.Vin, payload.Motor).Scan(
		&v.ID, &v.Placa, &v.Marca, &v.Modelo, &v.Color, &v.Vin, &v.Motor, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &v, nil
}
