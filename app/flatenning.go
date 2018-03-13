package app

import (
	"math/big"
	"time"

	uuid "github.com/satori/go.uuid"
)

// FlattenedTransaction - merges transaction and account data
type FlattenedTransaction struct {
	Id          uuid.UUID // nolint: golint
	Description string
	Creation    time.Time
	Date        time.Time

	Value1   *big.Rat
	Account1 string
	Ref1     uuid.UUID
	AccType1 string

	Value2   *big.Rat
	Account2 string
	Ref2     uuid.UUID
	AccType2 string
}

func groupAccountsById(accounts []Account) map[string]Account { // nolint: golint
	grouped := make(map[string]Account)
	for _, acc := range accounts {
		grouped[acc.Id.String()] = acc
	}
	return grouped
}

func flatten(tx *Transaction, account1, account2 *Account) *FlattenedTransaction {
	return &FlattenedTransaction{
		Id:          tx.Id,
		Account1:    account1.Name,
		Account2:    account2.Name,
		AccType1:    account1.Type,
		AccType2:    account2.Type,
		Creation:    tx.Creation,
		Date:        tx.Date,
		Description: tx.Description,
		Ref1:        tx.Ref1,
		Ref2:        tx.Ref2,
		Value1:      tx.Value1,
		Value2:      tx.Value2,
	}
}

// Flatten combines all transaction and their respective accounts
func Flatten(book *Book) []FlattenedTransaction {
	accountsById := groupAccountsById(book.Accounts) // nolint: golint
	flattened := []FlattenedTransaction{}
	for _, tx := range book.Transactions {
		account1, account2 := accountsById[tx.Account1.String()], accountsById[tx.Account2.String()]
		flattened = append(flattened, *flatten(&tx, &account1, &account2))
	}
	return flattened
}
