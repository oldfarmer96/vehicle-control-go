package store

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oldfarmer96/vehicle-control-go/internal/models"
)

type PersonaStore struct {
	db *pgxpool.Pool
}

func NewPersonaStore(db *pgxpool.Pool) *PersonaStore {
	return &PersonaStore{db: db}
}

func (s *PersonaStore) FindByDNI(ctx context.Context, dni string) (*models.Persona, error) {
	query := `
		SELECT id, dni, nombre_completo, rol, tiene_acceso_permitido, created_at, updated_at
		FROM personas
		WHERE dni = $1
	`

	var p models.Persona
	err := s.db.QueryRow(ctx, query, dni).Scan(
		&p.ID, &p.DNI, &p.NombreCompleto, &p.Rol, &p.TieneAccesoPermitido, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("persona no encontrada")
		}
		return nil, err
	}

	return &p, nil
}

func (s *PersonaStore) Create(ctx context.Context, payload models.CreatePersonaDTO) (*models.Persona, error) {
	query := `
		INSERT INTO personas (dni, nombre_completo, rol, tiene_acceso_permitido)
		VALUES ($1, $2, $3, $4)
		RETURNING id, dni, nombre_completo, rol, tiene_acceso_permitido, created_at, updated_at
	`

	var p models.Persona
	err := s.db.QueryRow(ctx, query,
		payload.DNI, payload.NombreCompleto, payload.Rol, payload.TieneAccesoPermitido,
	).Scan(
		&p.ID, &p.DNI, &p.NombreCompleto, &p.Rol, &p.TieneAccesoPermitido, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *PersonaStore) ToggleAccess(ctx context.Context, id string) (*models.Persona, error) {
	query := `
		UPDATE personas SET tiene_acceso_permitido = NOT tiene_acceso_permitido WHERE id = $1
		RETURNING id, dni, nombre_completo, rol, tiene_acceso_permitido, created_at, updated_at
	`

	var p models.Persona
	err := s.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.DNI, &p.NombreCompleto, &p.Rol, &p.TieneAccesoPermitido, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("persona no encontrada")
		}
		return nil, err
	}

	return &p, nil
}

func (s *PersonaStore) Update(ctx context.Context, id string, payload models.UpdatePersonaDTO) (*models.Persona, error) {
	var sets []string
	var args []interface{}
	argIdx := 1

	if payload.NombreCompleto != nil {
		sets = append(sets, fmt.Sprintf("nombre_completo = $%d", argIdx))
		args = append(args, *payload.NombreCompleto)
		argIdx++
	}
	if payload.Rol != nil {
		sets = append(sets, fmt.Sprintf("rol = $%d", argIdx))
		args = append(args, *payload.Rol)
		argIdx++
	}
	if payload.TieneAccesoPermitido != nil {
		sets = append(sets, fmt.Sprintf("tiene_acceso_permitido = $%d", argIdx))
		args = append(args, *payload.TieneAccesoPermitido)
		argIdx++
	}

	if len(sets) == 0 {
		return s.findByID(ctx, id)
	}

	query := fmt.Sprintf("UPDATE personas SET %s WHERE id = $%d RETURNING id, dni, nombre_completo, rol, tiene_acceso_permitido, created_at, updated_at", strings.Join(sets, ", "), argIdx)
	args = append(args, id)

	var p models.Persona
	err := s.db.QueryRow(ctx, query, args...).Scan(
		&p.ID, &p.DNI, &p.NombreCompleto, &p.Rol, &p.TieneAccesoPermitido, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("persona no encontrada")
		}
		return nil, err
	}

	return &p, nil
}

func (s *PersonaStore) findByID(ctx context.Context, id string) (*models.Persona, error) {
	query := `
		SELECT id, dni, nombre_completo, rol, tiene_acceso_permitido, created_at, updated_at
		FROM personas
		WHERE id = $1
	`

	var p models.Persona
	err := s.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.DNI, &p.NombreCompleto, &p.Rol, &p.TieneAccesoPermitido, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("persona no encontrada")
		}
		return nil, err
	}

	return &p, nil
}

func (s *PersonaStore) GetAll(ctx context.Context, page, limit int, search string) ([]models.Persona, int, error) {
	offset := (page - 1) * limit

	countQuery := `SELECT COUNT(*) FROM personas WHERE 1=1`
	listQuery := `
		SELECT id, dni, nombre_completo, rol, tiene_acceso_permitido, created_at, updated_at
		FROM personas WHERE 1=1`
	var args []interface{}
	argIdx := 1

	if search != "" {
		searchPattern := "%" + search + "%"
		countQuery += fmt.Sprintf(" AND (dni ILIKE $%d OR nombre_completo ILIKE $%d)", argIdx, argIdx)
		listQuery += fmt.Sprintf(" AND (dni ILIKE $%d OR nombre_completo ILIKE $%d)", argIdx, argIdx)
		args = append(args, searchPattern)
		argIdx++
	}

	var total int
	err := s.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	listQuery += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := s.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var personas []models.Persona
	for rows.Next() {
		var p models.Persona
		err := rows.Scan(
			&p.ID, &p.DNI, &p.NombreCompleto, &p.Rol, &p.TieneAccesoPermitido, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		personas = append(personas, p)
	}

	if len(personas) == 0 {
		return personas, total, nil
	}

	personaIDs := make([]string, len(personas))
	personaMap := make(map[string]*models.Persona)
	for i, p := range personas {
		personaIDs[i] = p.ID
		cp := p
		personaMap[p.ID] = &cp
	}

	vehiclesQuery := `
		SELECT vp.persona_id, v.id, v.placa, v.marca
		FROM vehiculos_personas vp
		JOIN vehiculos v ON vp.vehiculo_id = v.id
		WHERE vp.persona_id = ANY($1)
	`
	vehicleRows, err := s.db.Query(ctx, vehiclesQuery, personaIDs)
	if err != nil {
		return nil, 0, err
	}
	defer vehicleRows.Close()

	for vehicleRows.Next() {
		var personaID string
		var pv models.PersonaVehicle
		err := vehicleRows.Scan(&personaID, &pv.ID, &pv.Placa, &pv.Marca)
		if err != nil {
			return nil, 0, err
		}
		if p, exists := personaMap[personaID]; exists {
			p.Vehiculos = append(p.Vehiculos, pv)
		}
	}

	result := make([]models.Persona, len(personas))
	for i, id := range personaIDs {
		result[i] = *personaMap[id]
	}

	return result, total, nil
}
