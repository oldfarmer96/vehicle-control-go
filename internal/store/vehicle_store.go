package store

import (
	"context"
	"errors"
	"fmt"

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

func (s *VehicleStore) GetAll(ctx context.Context, page, limit int, placaFilter string) ([]models.Vehicle, int, error) {
	offset := (page - 1) * limit

	countQuery := `SELECT COUNT(*) FROM vehiculos WHERE 1=1`
	listQuery := `
		SELECT
			v.id, v.placa, v.marca, v.modelo, v.color, v.vin, v.motor, v.created_at, v.updated_at,
			p.id as owner_id, p.dni, p.nombre_completo, p.rol, p.tiene_acceso_permitido
		FROM vehiculos v
		LEFT JOIN vehiculos_personas vp ON v.id = vp.vehiculo_id
		LEFT JOIN personas p ON vp.persona_id = p.id
		WHERE 1=1`
	var args []any
	argIdx := 1

	if placaFilter != "" {
		countQuery += fmt.Sprintf(" AND placa ILIKE $%d", argIdx)
		listQuery += fmt.Sprintf(" AND v.placa ILIKE $%d", argIdx)
		args = append(args, "%"+placaFilter+"%")
		argIdx++
	}

	var total int
	err := s.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	listQuery += fmt.Sprintf(" ORDER BY v.created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := s.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	vehicleMap := make(map[string]*models.Vehicle)
	for rows.Next() {
		var v models.Vehicle
		var ownerID *string
		var ownerDNI, ownerNombreCompleto, ownerRol *string
		var ownerTieneAcceso *bool

		err := rows.Scan(
			&v.ID, &v.Placa, &v.Marca, &v.Modelo, &v.Color, &v.Vin, &v.Motor, &v.CreatedAt, &v.UpdatedAt,
			&ownerID, &ownerDNI, &ownerNombreCompleto, &ownerRol, &ownerTieneAcceso,
		)
		if err != nil {
			return nil, 0, err
		}

		if _, exists := vehicleMap[v.ID]; !exists {
			if ownerID != nil {
				v.Duenio = &models.VehicleOwner{
					ID:             *ownerID,
					DNI:            *ownerDNI,
					NombreCompleto: *ownerNombreCompleto,
					Rol:            *ownerRol,
					TieneAcceso:    *ownerTieneAcceso,
				}
			}
			vehicleMap[v.ID] = &v
		}
	}

	vehicles := make([]models.Vehicle, 0, len(vehicleMap))
	for _, v := range vehicleMap {
		vehicles = append(vehicles, *v)
	}

	return vehicles, total, nil
}
