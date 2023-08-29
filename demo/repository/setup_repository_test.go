package repository

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Testconteinners-go でPostgreSQLを起動する
func setupDBContainer(t *testing.T) (*sql.DB, error) {
	t.Helper()
	path, _ := filepath.Abs("../db")
	// PostgreSQLのコンテナを起動する
	req := testcontainers.ContainerRequest{
		Image: "postgres:15.4",
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "test",
		},
		ExposedPorts: []string{"5432/tcp"},
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.AutoRemove = true
		},
		// DBの初期化用のSQLをマウントする
		Mounts: testcontainers.ContainerMounts{
			testcontainers.BindMount(path, "/docker-entrypoint-initdb.d"),
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(1 * time.Minute),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	t.Cleanup(func() {
		require.NoError(t, postgresC.Terminate(context.Background()))
	})
	// テスト用のDBに接続する
	host, err := postgresC.Host(ctx)
	if err != nil {
		return nil, err
	}
	port, err := postgresC.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=postgres password=password dbname=test sslmode=disable", host, port.Port()))
	if err != nil {
		return nil, err
	}
	return db, nil
}
