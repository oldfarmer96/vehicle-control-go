package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oldfarmer96/vehicle-control-go/internal/models"
)

type AccessEventStore struct {
	db *pgxpool.Pool
}

func NewAccessEventStore(db *pgxpool.Pool) *AccessEventStore {
	return &AccessEventStore{db: db}
}

func (s *AccessEventStore) FindByPlaca(ctx context.Context, placa string) (*models.Vehicle, error) {
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
		return nil, err
	}

	return &v, nil
}

func (s *AccessEventStore) CreateVehicle(ctx context.Context, placa string, marca, modelo, color, vin, motor *string) (*models.Vehicle, error) {
	query := `
	INSERT INTO vehiculos (placa, marca, modelo, color, vin, motor)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, placa, marca, modelo, color, vin, motor, created_at, updated_at
	`

	var v models.Vehicle
	err := s.db.QueryRow(ctx, query, placa, marca, modelo, color, vin, motor).Scan(
		&v.ID, &v.Placa, &v.Marca, &v.Modelo, &v.Color, &v.Vin, &v.Motor, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (s *AccessEventStore) CreateAccessEvent(ctx context.Context, vehiculoID, tipoEvento, puntoControl string, confianzaOcr float64) error {
	query := `
	INSERT INTO eventos_acceso (vehiculo_id, tipo_evento, punto_control, confianza_ocr, fecha_hora)
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err := s.db.Exec(ctx, query, vehiculoID, tipoEvento, puntoControl, confianzaOcr, time.Now())
	return err
}