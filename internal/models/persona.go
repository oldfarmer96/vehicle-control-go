package models

import (
	"strings"
	"time"
)

type PersonaVehicle struct {
	ID    string  `json:"id"`
	Placa string  `json:"placa"`
	Marca *string `json:"marca"`
}

type Persona struct {
	ID                   string           `json:"id"`
	DNI                  string           `json:"dni"`
	NombreCompleto       string           `json:"nombreCompleto"`
	Rol                  string           `json:"rol"`
	TieneAccesoPermitido bool             `json:"tieneAccesoPermitido"`
	CreatedAt            time.Time        `json:"createdAt"`
	UpdatedAt            time.Time        `json:"updatedAt"`
	Vehiculos            []PersonaVehicle `json:"vehiculos,omitempty"`
}

type CreatePersonaDTO struct {
	DNI                  string `json:"dni" validate:"required,len=8"`
	NombreCompleto       string `json:"nombreCompleto" validate:"required,min=5,max=150"`
	Rol                  string `json:"rol" validate:"required,oneof=DOCENTE ALUMNO ADMINISTRATIVO VISITANTE"`
	TieneAccesoPermitido *bool  `json:"tieneAccesoPermitido" validate:"required"`
}

func (d *CreatePersonaDTO) Normalize() {
	d.DNI = strings.TrimSpace(d.DNI)
	d.NombreCompleto = strings.TrimSpace(d.NombreCompleto)
	d.Rol = strings.TrimSpace(d.Rol)
	if d.TieneAccesoPermitido == nil {
		t := true
		d.TieneAccesoPermitido = &t
	}
}

type UpdatePersonaDTO struct {
	NombreCompleto       *string `json:"nombreCompleto,omitempty"`
	Rol                  *string `json:"rol,omitempty"`
	TieneAccesoPermitido *bool   `json:"tieneAccesoPermitido,omitempty"`
}

func (d *UpdatePersonaDTO) Normalize() {
	if d.NombreCompleto != nil && strings.TrimSpace(*d.NombreCompleto) == "" {
		d.NombreCompleto = nil
	}
	if d.Rol != nil && strings.TrimSpace(*d.Rol) == "" {
		d.Rol = nil
	}
}

type ListPersonasQuery struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Search string `query:"search"`
}

type PersonaListResponse struct {
	Personas   []Persona `json:"personas"`
	Total      int       `json:"total"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalPages int       `json:"totalPages"`
}