package database

import (
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"
	"vortex_test/internal/model"
	"vortex_test/pkg/logging"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

const dbURL = "postgres://postgres@localhost:5432/example"

type clientTest struct {
	testName   string
	testInput  model.Client
	testOutput struct {
		client model.Client
		err    error
	}
}

type algorithmStatusTest struct {
	testName   string
	testInput  model.Algorithm
	testOutput struct {
		algorithmStatus model.Algorithm
		err             error
	}
}

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		panic(err)
	}

	if err = pool.Client.Ping(); err != nil {
		panic(err)
	}

	db, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Name:       "Test",
		Tag:        "16.3-alpine",
		Env: []string{
			"POSTGRES_DB=example",
			"POSTGRES_HOST_AUTH_METHOD=trust",
			"listen_addresses = '*'",
		},
		ExposedPorts: []string{"5432/tcp"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432/tcp": {{HostIP: "localhost", HostPort: "5432/tcp"}},
		},
	},
		func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		})

	if err != nil {
		panic(err)
	}

	db.Expire(10)

	if err := pool.Retry(func() error {
		d, err := sql.Open("pgx", dbURL)

		if err != nil {
			return err
		}
		s := New(d, logging.New())
		if err = s.CreateTables(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		panic("Could not connect to postgres: " + err.Error())
	}

	os.Exit(m.Run())
}

func TestAddClient(t *testing.T) {
	tests := []clientTest{
		{
			testName: "without errors",
			testInput: model.Client{
				ClientName:  "Boba",
				Version:     1,
				Image:       "unknown",
				CPU:         "first",
				Memory:      "Low",
				Priority:    1.25,
				NeedRestart: false,
			},
			testOutput: struct {
				client model.Client
				err    error
			}{
				client: model.Client{},
				err:    nil,
			},
		},
		{
			testName: "without errors",
			testInput: model.Client{
				ClientName:  "Biba",
				Version:     1,
				Image:       "unknown",
				CPU:         "first",
				Memory:      "Low",
				Priority:    1.25,
				NeedRestart: false,
			},
			testOutput: struct {
				client model.Client
				err    error
			}{
				client: model.Client{},
				err:    nil,
			},
		},
	}
	db, err := sql.Open("pgx", dbURL)
	require.NoError(t, err)
	storage := New(db, logging.New())

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			_, err := storage.AddClient(test.testInput)
			require.NoError(t, err)
		})
	}
}

func TestUpdateClient(t *testing.T) {
	tests := []clientTest{
		{
			testName: "without errors",
			testInput: model.Client{
				ID:          1,
				ClientName:  "Boba",
				Version:     1,
				Image:       "unknown",
				CPU:         "first",
				Memory:      "Low",
				Priority:    7,
				NeedRestart: true,
				SpawnedAt:   time.Now(),
				CreatedAt:   time.Now(),
			},
			testOutput: struct {
				client model.Client
				err    error
			}{
				client: model.Client{},
				err:    nil,
			},
		},
		{
			testName: "with errors",
			testInput: model.Client{
				ID:          836,
				ClientName:  "Boba",
				Version:     1,
				Image:       "unknown",
				CPU:         "first",
				Memory:      "Low",
				Priority:    7,
				NeedRestart: true,
			},
			testOutput: struct {
				client model.Client
				err    error
			}{
				client: model.Client{},
				err:    errors.New("sql: no rows in result set"),
			},
		},
	}

	db, err := sql.Open("pgx", dbURL)
	require.NoError(t, err)

	storage := New(db, logging.New())

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			_, err := storage.UpdateClient(test.testInput)
			require.Equal(t, test.testOutput.err, err)
		})
	}
}

func TestDeleteClient(t *testing.T) {
	tests := []clientTest{
		{
			testName: "without errors",
			testInput: model.Client{
				ID:          1,
				ClientName:  "Boba",
				Version:     1,
				Image:       "unknown",
				CPU:         "first",
				Memory:      "Low",
				Priority:    7,
				NeedRestart: true,
			},
			testOutput: struct {
				client model.Client
				err    error
			}{
				client: model.Client{},
				err:    nil,
			},
		},
	}

	db, err := sql.Open("pgx", dbURL)
	require.NoError(t, err)

	storage := New(db, logging.New())

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			err := storage.DeleteClient(test.testInput)
			require.NoError(t, err)
		})
	}
}

func TestUpdateAlgorithmStatus(t *testing.T) {
	tests := []algorithmStatusTest{
		{
			testName: "without errors",
			testInput: model.Algorithm{
				ClientID: 2,
				VWAP:     true,
				TWAP:     false,
				HFT:      true,
			},
			testOutput: struct {
				algorithmStatus model.Algorithm
				err             error
			}{
				algorithmStatus: model.Algorithm{
					ID:       2,
					ClientID: 2,
					VWAP:     true,
					TWAP:     false,
					HFT:      true,
				},
				err: nil,
			},
		},
		{
			testName: "with errors",
			testInput: model.Algorithm{
				ClientID: 4,
				VWAP:     true,
				TWAP:     false,
				HFT:      true,
			},
			testOutput: struct {
				algorithmStatus model.Algorithm
				err             error
			}{
				algorithmStatus: model.Algorithm{},
				err:             errors.New("sql: no rows in result set"),
			},
		},
	}
	db, err := sql.Open("pgx", dbURL)
	require.NoError(t, err)
	storage := New(db, logging.New())

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			res, err := storage.UpdateAlgorithmStatus(test.testInput)
			result := struct {
				algorithmStatus model.Algorithm
				err             error
			}{
				res,
				err,
			}
			require.Equal(t, test.testOutput, result)
		})
	}
}
