package models

import (
	"strings"
	"time"
)

type Persona struct {
	ID                   string    `json:"id"`
	DNI                  string    `json:"dni"`
	NombreCompleto       string    `json:"nombre_completo"`
	Rol                  string    `json:"rol"`
	TieneAccesoPermitido bool      `json:"tiene_acceso_permitido"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type CreatePersonaDTO struct {
	DNI                  string `json:"dni" validate:"required,len=8"`
	NombreCompleto       string `json:"nombre_completo" validate:"required,min=5,max=150"`
	Rol                  string `json:"rol" validate:"required,oneof=DOCENTE ALUMNO ADMINISTRATIVO VISITANTE"`
	TieneAccesoPermitido *bool  `json:"tiene_acceso_permitido" validate:"required"`
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