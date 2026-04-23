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

type UserStore struct {
	db *pgxpool.Pool
}

func NewUserStore(db *pgxpool.Pool) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, nombre, apellidos, dni, username, password, rol, activo, created_at, updated_at
		FROM usuarios_web
		WHERE username = $1
	`

	var u models.User
	err := s.db.QueryRow(ctx, query, username).Scan(
		&u.ID, &u.Nombre, &u.Apellidos, &u.DNI, &u.Username, &u.Password,
		&u.Rol, &u.Activo, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	return &u, nil
}

func (s *UserStore) Create(ctx context.Context, payload models.CreateUserDTO, passwordHash string) (*models.User, error) {
	query := `
		INSERT INTO usuarios_web (nombre, apellidos, dni, username, password, rol)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, nombre, apellidos, dni, username, password, rol, activo, created_at, updated_at
	`

	var u models.User
	err := s.db.QueryRow(ctx, query,
		payload.Nombre, payload.Apellidos, payload.DNI,
		payload.Username, passwordHash, payload.Rol,
	).Scan(
		&u.ID, &u.Nombre, &u.Apellidos, &u.DNI, &u.Username, &u.Password,
		&u.Rol, &u.Activo, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *UserStore) FindByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, nombre, apellidos, dni, username, password, rol, activo, created_at, updated_at
		FROM usuarios_web
		WHERE id = $1
	`

	var u models.User
	err := s.db.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Nombre, &u.Apellidos, &u.DNI, &u.Username, &u.Password,
		&u.Rol, &u.Activo, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	return &u, nil
}

func (s *UserStore) List(ctx context.Context, page, limit int, search string) ([]models.User, int, error) {
	offset := (page - 1) * limit

	countQuery := `SELECT COUNT(*) FROM usuarios_web WHERE 1=1`
	listQuery := `
		SELECT id, nombre, apellidos, dni, username, password, rol, activo, created_at, updated_at
		FROM usuarios_web WHERE 1=1`
	var args []interface{}
	argIdx := 1

	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		countQuery += fmt.Sprintf(" AND (LOWER(nombre) LIKE $%d OR LOWER(dni) LIKE $%d OR LOWER(username) LIKE $%d)", argIdx, argIdx, argIdx)
		listQuery += fmt.Sprintf(" AND (LOWER(nombre) LIKE $%d OR LOWER(dni) LIKE $%d OR LOWER(username) LIKE $%d)", argIdx, argIdx, argIdx)
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

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID, &u.Nombre, &u.Apellidos, &u.DNI, &u.Username, &u.Password,
			&u.Rol, &u.Activo, &u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	return users, total, nil
}

func (s *UserStore) Update(ctx context.Context, id string, payload models.UpdateUserDTO) (*models.User, error) {
	var sets []string
	var args []interface{}
	argIdx := 1

	if payload.Nombre != nil {
		sets = append(sets, fmt.Sprintf("nombre = $%d", argIdx))
		args = append(args, *payload.Nombre)
		argIdx++
	}
	if payload.Apellidos != nil {
		sets = append(sets, fmt.Sprintf("apellidos = $%d", argIdx))
		args = append(args, *payload.Apellidos)
		argIdx++
	}
	if payload.DNI != nil {
		sets = append(sets, fmt.Sprintf("dni = $%d", argIdx))
		args = append(args, *payload.DNI)
		argIdx++
	}
	if payload.Username != nil {
		sets = append(sets, fmt.Sprintf("username = $%d", argIdx))
		args = append(args, *payload.Username)
		argIdx++
	}
	if payload.Password != nil {
		sets = append(sets, fmt.Sprintf("password = $%d", argIdx))
		args = append(args, *payload.Password)
		argIdx++
	}
	if payload.Rol != nil {
		sets = append(sets, fmt.Sprintf("rol = $%d", argIdx))
		args = append(args, *payload.Rol)
		argIdx++
	}

	if len(sets) == 0 {
		return s.FindByID(ctx, id)
	}

	query := fmt.Sprintf("UPDATE usuarios_web SET %s WHERE id = $%d RETURNING id, nombre, apellidos, dni, username, password, rol, activo, created_at, updated_at", strings.Join(sets, ", "), argIdx)
	args = append(args, id)

	var u models.User
	err := s.db.QueryRow(ctx, query, args...).Scan(
		&u.ID, &u.Nombre, &u.Apellidos, &u.DNI, &u.Username, &u.Password,
		&u.Rol, &u.Activo, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	return &u, nil
}

func (s *UserStore) ToggleActive(ctx context.Context, id string) (*models.User, error) {
	query := `
		UPDATE usuarios_web SET activo = NOT activo WHERE id = $1
		RETURNING id, nombre, apellidos, dni, username, password, rol, activo, created_at, updated_at
	`

	var u models.User
	err := s.db.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Nombre, &u.Apellidos, &u.DNI, &u.Username, &u.Password,
		&u.Rol, &u.Activo, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	return &u, nil
}
