package database

import (
	"database/sql"
	"time"
	"vortex_test/internal/model"
	"vortex_test/pkg/logging"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db     *sql.DB
	logger *logging.Logger
}

func New(url *sql.DB, logger *logging.Logger) *Storage {
	return &Storage{
		db:     url,
		logger: logger,
	}
}

func (s *Storage) CreateTables() error {
	query := `CREATE TABLE IF NOT EXISTS clients (
			id SERIAL PRIMARY KEY,
			clientName VARCHAR(100) UNIQUE,
			version INT,
			image VARCHAR(100),
			cpu VARCHAR(100),
			memory VARCHAR(100),
			priority FLOAT,
			needRestart BOOL DEFAULT FALSE,
			spawnedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    		createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    		updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
			CREATE TABLE IF NOT EXISTS algorithmStatus (
			id SERIAL PRIMARY KEY,
			clientID INT UNIQUE REFERENCES clients(id) ON DELETE CASCADE,
			vwap BOOLEAN DEFAULT FALSE,
    		twap BOOLEAN DEFAULT FALSE,
    		hft BOOLEAN DEFAULT FALSE
			)
			`
	_, err := s.db.Exec(query)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	return nil
}

func (s *Storage) AddClient(client model.Client) (model.Client, error) {
	query := `INSERT INTO clients (clientName, version, image, cpu, memory, priority, needRestart)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING *`
	res := s.db.QueryRow(
		query,
		client.ClientName,
		client.Version,
		client.Image,
		client.CPU,
		client.Memory,
		client.Priority,
		client.NeedRestart,
	)

	var cl model.Client
	err := res.Scan(
		&cl.ID,
		&cl.ClientName,
		&cl.Version,
		&cl.Image,
		&cl.CPU,
		&cl.Memory,
		&cl.Priority,
		&cl.NeedRestart,
		&cl.SpawnedAt,
		&cl.CreatedAt,
		&cl.UpdatedAt,
	)
	if err != nil {
		s.logger.Error(err)
		return model.Client{}, err
	}

	_, err = s.db.Exec("INSERT INTO algorithmStatus(clientID) VALUES ($1)", cl.ID)
	if err != nil {
		s.logger.Error(err)
		return model.Client{}, err
	}
	return cl, nil
}

func (s *Storage) UpdateClient(client model.Client) (model.Client, error) {
	query := `UPDATE clients SET clientName=$2,
				version=$3,
				image=$4,
				cpu=$5,
				memory=$6,
				priority=$7,
				needRestart=$8,
				updatedAt=$9
			WHERE id=$1
		RETURNING *`

	res := s.db.QueryRow(
		query,
		client.ID,
		client.ClientName,
		client.Version,
		client.Image,
		client.CPU,
		client.Memory,
		client.Priority,
		client.NeedRestart,
		time.Now(),
	)
	err := res.Scan(
		&client.ID,
		&client.ClientName,
		&client.Version,
		&client.Image,
		&client.CPU,
		&client.Memory,
		&client.Priority,
		&client.NeedRestart,
		&client.SpawnedAt,
		&client.CreatedAt,
		&client.UpdatedAt,
	)
	if err != nil {
		s.logger.Error(err)
		return model.Client{}, err
	}

	return client, nil
}

func (s *Storage) DeleteClient(client model.Client) error {
	query := `DELETE FROM clients WHERE id=$1`
	_, err := s.db.Exec(query, client.ID)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	return nil
}

func (s *Storage) UpdateAlgorithmStatus(algorithm model.Algorithm) (model.Algorithm, error) {
	query := `UPDATE algorithmStatus SET vwap=$2,
				twap=$3,
				hft=$4
			WHERE clientID=$1
			RETURNING *`

	res := s.db.QueryRow(
		query,
		algorithm.ClientID,
		algorithm.VWAP,
		algorithm.TWAP,
		algorithm.HFT,
	)

	err := res.Scan(
		&algorithm.ID,
		&algorithm.ClientID,
		&algorithm.VWAP,
		&algorithm.TWAP,
		&algorithm.HFT,
	)

	if err != nil {
		s.logger.Error(err)
		return model.Algorithm{}, err
	}

	return algorithm, nil
}

func (s *Storage) GetCurrentAlgoritmStatus() ([]model.Algorithm, error) {
	query := `SELECT * FROM algorithmStatus`
	res, err := s.db.Query(query)
	if err != nil {
		s.logger.Error(err)
		return []model.Algorithm{}, err
	}
	var algorithmRow model.Algorithm
	ans := []model.Algorithm{}
	for res.Next() {
		err := res.Scan(
			&algorithmRow.ID,
			&algorithmRow.ClientID,
			&algorithmRow.VWAP,
			&algorithmRow.TWAP,
			&algorithmRow.HFT,
		)

		if err != nil {
			s.logger.Error(err)
			return []model.Algorithm{}, err
		}

		ans = append(ans, algorithmRow)
	}

	return ans, nil
}
