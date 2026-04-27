package store

import (
	"context"
	"errors"
	"fmt"
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

func (s *AccessEventStore) GetAll(ctx context.Context, page, limit int, placaFilter string) ([]models.AccessEventWithVehicle, int, error) {
	offset := (page - 1) * limit

	countQuery := `SELECT COUNT(*) FROM eventos_acceso ea
		JOIN vehiculos v ON ea.vehiculo_id = v.id
		WHERE 1=1`
	listQuery := `
		SELECT
			ea.id, ea.vehiculo_id, ea.tipo_evento, ea.punto_control, ea.confianza_ocr, ea.fecha_hora,
			v.placa, v.marca, v.modelo, v.color, v.vin, v.motor,
			p.id as owner_id, p.dni, p.nombre_completo, p.rol, p.tiene_acceso_permitido
		FROM eventos_acceso ea
		JOIN vehiculos v ON ea.vehiculo_id = v.id
		LEFT JOIN vehiculos_personas vp ON v.id = vp.vehiculo_id
		LEFT JOIN personas p ON vp.persona_id = p.id
		WHERE 1=1`

	var args []any
	argIdx := 1

	if placaFilter != "" {
		countQuery += fmt.Sprintf(" AND v.placa ILIKE $%d", argIdx)
		listQuery += fmt.Sprintf(" AND v.placa ILIKE $%d", argIdx)
		args = append(args, "%"+placaFilter+"%")
		argIdx++
	}

	var total int
	err := s.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	listQuery += fmt.Sprintf(" ORDER BY ea.fecha_hora DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := s.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	eventMap := make(map[string]*models.AccessEventWithVehicle)
	for rows.Next() {
		var e models.AccessEventWithVehicle
		var ownerID, ownerDNI, ownerNombreCompleto, ownerRol *string
		var ownerTieneAcceso *bool

		err := rows.Scan(
			&e.ID, &e.VehiculoID, &e.TipoEvento, &e.PuntoControl, &e.ConfianzaOcr, &e.FechaHora,
			&e.Placa, &e.Marca, &e.Modelo, &e.Color, &e.Vin, &e.Motor,
			&ownerID, &ownerDNI, &ownerNombreCompleto, &ownerRol, &ownerTieneAcceso,
		)
		if err != nil {
			return nil, 0, err
		}

		if _, exists := eventMap[e.ID]; !exists {
			if ownerID != nil {
				e.Duenio = &models.VehicleOwner{
					ID:             *ownerID,
					DNI:            *ownerDNI,
					NombreCompleto: *ownerNombreCompleto,
					Rol:            *ownerRol,
					TieneAcceso:    *ownerTieneAcceso,
				}
			}
			eventMap[e.ID] = &e
		}
	}

	events := make([]models.AccessEventWithVehicle, 0, len(eventMap))
	for _, e := range eventMap {
		events = append(events, *e)
	}

	return events, total, nil
}