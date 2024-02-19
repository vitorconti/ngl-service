package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID   uuid.UUID
	Name string
}

type Contract struct {
	ID          uuid.UUID
	CustomerID  uuid.UUID
	Description string
}

type Charge struct {
	ID             uuid.UUID
	Reference      time.Time
	ExpirationDate time.Time
	PaymentDate    sql.NullTime
	Status         bool
	ContractID     uuid.UUID
}
