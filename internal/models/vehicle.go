package models

import (
	"strings"
	"time"
)

type Vehicle struct {
	ID        string    `json:"id"`
	Placa     string    `json:"placa"`
	Marca     *string   `json:"marca"`
	Modelo    *string   `json:"modelo"`
	Color     *string   `json:"color"`
	Vin       *string   `json:"vin"`
	Motor     *string   `json:"motor"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreaateVehicleDTO struct {
	Placa  string  `json:"placa" validate:"required,min=6,max=15"`
	Marca  *string `json:"marca,omitempty" validate:"omitempty,min=2"`
	Modelo *string `json:"modelo,omitempty" validate:"omitempty,min=2"`
	Color  *string `json:"color,omitempty" validate:"omitempty,min=3"`
	Vin    *string `json:"vin,omitempty" validate:"omitempty,alphanum"`
	Motor  *string `json:"motor,omitempty" validate:"omitempty,min=5"`
}

func (d *CreaateVehicleDTO) Normalize() {
	if d.Marca != nil && strings.TrimSpace(*d.Marca) == "" {
		d.Marca = nil
	}
	if d.Modelo != nil && strings.TrimSpace(*d.Modelo) == "" {
		d.Modelo = nil
	}
	if d.Color != nil && strings.TrimSpace(*d.Color) == "" {
		d.Color = nil
	}
	if d.Vin != nil && strings.TrimSpace(*d.Vin) == "" {
		d.Vin = nil
	}
	if d.Motor != nil && strings.TrimSpace(*d.Motor) == "" {
		d.Motor = nil
	}
}
