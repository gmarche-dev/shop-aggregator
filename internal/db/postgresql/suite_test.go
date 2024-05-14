package postgresql

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"shop-aggregator/internal/config"
	"shop-aggregator/tools/migrations"
	"strconv"
	"strings"
)

type DBTestSuite struct {
	suite.Suite
	DB        *Client
	Container testcontainers.Container
	ctx       context.Context
}

func (s *DBTestSuite) SetupSuite() {
	s.ctx = context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "shopdb",
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "password",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	postgresContainer, err := testcontainers.GenericContainer(s.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err)

	s.Container = postgresContainer
	host, err := postgresContainer.Host(s.ctx)
	s.Require().NoError(err)
	portString, err := postgresContainer.MappedPort(s.ctx, "5432")
	s.Require().NoError(err)

	parts := strings.Split(string(portString), "/")
	if len(parts) == 0 {
		s.T().Fatal("invalid port")
	}
	i, err := strconv.Atoi(parts[0])
	s.Require().NoError(err)

	dbConfig := &config.DatabaseConfig{
		Host:     host,
		Port:     i,
		Username: "user",
		Password: "password",
		DBName:   "shopdb",
	}

	db, err := NewDB(s.ctx, dbConfig)
	s.Require().NoError(err)

	s.DB = db

	// migration
	err = migrations.Run(s.ctx, s.DB, "../../../migrations/deploy")
	s.Require().NoError(err)
}

func (s *DBTestSuite) TearDownSuite() {
	s.DB.Close()
	err := s.Container.Terminate(s.ctx)
	s.Require().NoError(err)
}
