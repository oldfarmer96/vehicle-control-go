package models

import "time"

type AccessEventRequest struct {
	Placa        string  `json:"placa"`
	Evento       string  `json:"evento"`
	ConfianzaOcr float64 `json:"confianzaOcr"`
	PuntoControl string  `json:"puntoControl"`
}

type AccessEventResponse struct {
	Evento   AccessEventRequest `json:"evento"`
	Vehicle  Vehicle            `json:"vehiculo"`
	Source   string             `json:"source"`
	Found    bool               `json:"found"`
}

type AccessEvent struct {
	ID           string    `json:"id"`
	VehiculoID   string    `json:"vehiculo_id"`
	TipoEvento   string    `json:"tipo_evento"`
	PuntoControl string    `json:"punto_control"`
	ConfianzaOcr float64   `json:"confianza_ocr"`
	FechaHora    time.Time `json:"fecha_hora"`
}

type AccessEventWithVehicle struct {
	ID              string         `json:"id"`
	VehiculoID      string         `json:"vehiculo_id"`
	TipoEvento      string         `json:"tipo_evento"`
	PuntoControl    string         `json:"punto_control"`
	ConfianzaOcr    float64        `json:"confianza_ocr"`
	FechaHora       time.Time      `json:"fecha_hora"`
	Placa           string         `json:"placa"`
	Marca           *string        `json:"marca"`
	Modelo          *string        `json:"modelo"`
	Color           *string        `json:"color"`
	Vin             *string        `json:"vin"`
	Motor           *string        `json:"motor"`
	Duenio          *VehicleOwner  `json:"duenio,omitempty"`
}

type ListAccessEventsQuery struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Placa string `query:"placa"`
}

type AccessEventListResponse struct {
	Eventos     []AccessEventWithVehicle `json:"eventos"`
	Total       int                     `json:"total"`
	Page        int                     `json:"page"`
	Limit       int                     `json:"limit"`
	TotalPages  int                     `json:"totalPages"`
}

type ExternalPlacaResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Placa  string `json:"placa"`
		Marca  string `json:"marca"`
		Modelo string `json:"modelo"`
		Serie  string `json:"serie"`
		Color  string `json:"color"`
		Motor  string `json:"motor"`
		Vin    string `json:"vin"`
	} `json:"data"`
}