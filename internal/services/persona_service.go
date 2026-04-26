package services

import (
	"context"
	"errors"

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