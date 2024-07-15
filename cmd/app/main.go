package main

import (
	"database/sql"
	"fmt"
	"sync"
	"vortex_test/internal/config"
	"vortex_test/internal/database"
	podsmanagement "vortex_test/internal/podsManagement"
	server "vortex_test/internal/web"
	"vortex_test/pkg/logging"
)

func main() {
	cfg := config.GetConfig()

	logger := logging.New()

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	storage := database.New(db, logger)
	if err = storage.CreateTables(); err != nil {
		panic(err)
	}

	server := server.New(storage, logger)

	pm := podsmanagement.New(storage, logger)

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		pm.Start()
	}()

	go func() {
		defer wg.Done()
		server.Start(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port))
	}()
	wg.Wait()
}
