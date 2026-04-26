package store

import (
	"context"
	"errors"

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