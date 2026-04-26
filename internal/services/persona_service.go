package services

import (
	"context"
	"errors"
	"math"

	"github.com/oldfarmer96/vehicle-control-go/internal/models"
	"github.com/oldfarmer96/vehicle-control-go/internal/store"
)

type PersonaService struct {
	personaStore *store.PersonaStore
}

func NewPersonaService(store *store.PersonaStore) *PersonaService {
	return &PersonaService{personaStore: store}
}

func (s *PersonaService) CreatePersona(ctx context.Context, payload models.CreatePersonaDTO) (*models.Persona, error) {
	_, err := s.personaStore.FindByDNI(ctx, payload.DNI)
	if err == nil {
		return nil, errors.New("el dni ya esta registrado")
	}

	return s.personaStore.Create(ctx, payload)
}

func (s *PersonaService) GetAllPersonas(ctx context.Context, page, limit int, search string) (*models.PersonaListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	personas, total, err := s.personaStore.GetAll(ctx, page, limit, search)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &models.PersonaListResponse{
		Personas:  personas,
		Total:     total,
		Page:      page,
		Limit:     limit,
		TotalPages: totalPages,
	}, nil
}