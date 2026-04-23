package main

import (
	"log"

	"github.com/oldfarmer96/vehicle-control-go/internal/bootstrap"
	"github.com/oldfarmer96/vehicle-control-go/pkg/database"
	"github.com/oldfarmer96/vehicle-control-go/pkg/env"
)

func main() {

	cfg, err := env.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbPool, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("No se pudo iniciar la BD: %v", err)
	}
	defer dbPool.Close()

	app := bootstrap.NewApp(dbPool)

	log.Printf("Servidor corriendo en el puerto %s...", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
