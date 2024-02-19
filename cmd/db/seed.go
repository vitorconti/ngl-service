package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

const (
	dbConnectionString = "postgresql://postgres:postgres@localhost:5432/your_database_name"
)

func Seed() {
	conn, err := pgx.Connect(context.Background(), dbConnectionString)
	if err != nil {
		log.Fatal("Unable to connect to the database:", err)
	}
	defer conn.Close(context.Background())

	createTables(conn)
	customerIDs := seedCustomers(conn)

	contractIDs := seedContracts(conn, customerIDs)

	seedCharges(conn, contractIDs)
}

func createTables(conn *pgx.Conn) {
	_, err := conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS customer (
			id UUID PRIMARY KEY,
			name VARCHAR(255) NOT NULL
		);

		CREATE TABLE IF NOT EXISTS contract (
			id UUID PRIMARY KEY,
			id_customer UUID REFERENCES customer(id),
			description VARCHAR(255) NOT NULL
		);

		CREATE TABLE IF NOT EXISTS charges (
			id UUID PRIMARY KEY,
			reference TIMESTAMP NOT NULL,
			expiration_date TIMESTAMP NOT NULL,
			payment_date TIMESTAMP,
			status BOOLEAN NOT NULL,
			id_contract UUID REFERENCES contract(id)
		);
	`)
	if err != nil {
		log.Fatal("Error creating tables:", err)
	}
}

func seedCustomers(conn *pgx.Conn) []uuid.UUID {
	customerIDs := make([]uuid.UUID, 0)

	for i := 1; i <= 5; i++ {
		id := uuid.New()
		name := fmt.Sprintf("Customer%d", i)

		_, err := conn.Exec(context.Background(), "INSERT INTO customer (id, name) VALUES ($1, $2)", id, name)
		if err != nil {
			log.Fatal("Error seeding customers:", err)
		}

		customerIDs = append(customerIDs, id)
	}

	return customerIDs
}

func seedContracts(conn *pgx.Conn, customerIDs []uuid.UUID) []uuid.UUID {
	contractIDs := make([]uuid.UUID, 0)

	for i, customerID := range customerIDs {
		id := uuid.New()
		description := fmt.Sprintf("Contract%d", i+1)

		_, err := conn.Exec(context.Background(), "INSERT INTO contract (id, id_customer, description) VALUES ($1, $2, $3)", id, customerID, description)
		if err != nil {
			log.Fatal("Error seeding contracts:", err)
		}

		contractIDs = append(contractIDs, id)
	}

	return contractIDs
}

func seedCharges(conn *pgx.Conn, contractIDs []uuid.UUID) {

	for _, contractID := range contractIDs {
		for i := 1; i <= 3; i++ {
			id := uuid.New()
			reference := time.Now().AddDate(0, i, 0)
			expirationDate := reference.AddDate(0, 1, 0)
			paymentDate := sql.NullTime{Time: reference.AddDate(0, 0, 15), Valid: true}
			status := i%2 == 0

			_, err := conn.Exec(context.Background(), "INSERT INTO charges (id, reference, expiration_date, payment_date, status, id_contract) VALUES ($1, $2, $3, $4, $5, $6)",
				id, reference, expirationDate, paymentDate, status, contractID)
			if err != nil {
				log.Fatal("Error seeding charges:", err)
			}
		}
	}
}
