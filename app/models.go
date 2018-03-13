package app

import (
	"math/big"
	"time"

	"github.com/satori/go.uuid"
)

// NewUUID generates an UUID V4 or panics
func NewUUID() uuid.UUID {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return id
}

type Book struct { // nolint: golint
	Accounts     []Account
	Transactions []Transaction
}

type Account struct { // nolint: golint
	Id     uuid.UUID // nolint: golint
	Name   string
	Type   string
	Parent uuid.UUID
}

type Transaction struct { // nolint: golint
	Id          uuid.UUID // nolint: golint
	Description string
	Creation    time.Time
	Date        time.Time
	Value1      *big.Rat
	Account1    uuid.UUID
	Ref1        uuid.UUID
	Value2      *big.Rat
	Account2    uuid.UUID
	Ref2        uuid.UUID
}
