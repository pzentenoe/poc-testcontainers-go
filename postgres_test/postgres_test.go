package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestPostgresContainer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Postgres Container Suite")
}

var _ = Describe("Postgres Container", func() {
	var (
		ctx               context.Context
		postgresContainer *postgres.PostgresContainer
		db                *sql.DB
		err               error
	)

	BeforeEach(func() {
		ctx = context.Background()

		// Configurar y lanzar el contenedor de PostgreSQL
		postgresContainer, err = postgres.Run(ctx,
			"docker.io/postgres:16-alpine",
			postgres.WithDatabase("testdb"),
			postgres.WithUsername("testuser"),
			postgres.WithPassword("testpass"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).
					WithStartupTimeout(5*time.Second)),
		)
		Expect(err).NotTo(HaveOccurred(), "Failed to start Postgres container")
	})

	AfterEach(func() {
		Expect(postgresContainer.Terminate(ctx)).To(Succeed(), "Failed to terminate Postgres container")
	})

	It("should have the container running", func() {
		state, err := postgresContainer.State(ctx)
		Expect(err).NotTo(HaveOccurred(), "Failed to get container state")
		Expect(state.Running).To(BeTrue(), "Expected the container to be running")
	})

	It("should create and query a table in Postgres", func() {
		dsn, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
		Expect(err).NotTo(HaveOccurred(), "Failed to get connection string")

		// Conectar a la base de datos
		db, err = sql.Open("postgres", dsn)
		Expect(err).NotTo(HaveOccurred(), "Failed to connect to Postgres")
		defer db.Close()

		// Crear una tabla, insertar un dato y luego consultarlo
		_, err = db.Exec(`CREATE TABLE users (id SERIAL PRIMARY KEY, name VARCHAR(50));`)
		Expect(err).NotTo(HaveOccurred(), "Failed to create table")

		_, err = db.Exec(`INSERT INTO users (name) VALUES ('Alice');`)
		Expect(err).NotTo(HaveOccurred(), "Failed to insert data")

		var name string
		err = db.QueryRow(`SELECT name FROM users WHERE name='Alice';`).Scan(&name)
		Expect(err).NotTo(HaveOccurred(), "Failed to query data")
		Expect(name).To(Equal("Alice"), "Expected name to be 'Alice'")
	})
})
