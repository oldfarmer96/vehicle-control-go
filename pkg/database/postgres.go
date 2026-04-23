package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)

	if err != nil {
		return nil, err
	}

	// CONFIGURACIÓN ROBUSTA DEL POOL
	// Máximo de conexiones abiertas simultáneamente
	config.MaxConns = 20
	// Mínimo de conexiones siempre abiertas listas para usarse
	config.MinConns = 5
	// Tiempo máximo de vida de una conexión
	config.MaxConnLifetime = time.Hour
	// Tiempo máximo que una conexión puede estar inactiva antes de cerrarse
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	log.Println("🚀 Conectado exitosamente a PostgreSQL (pgxpool configurado)")
	return pool, nil

}
