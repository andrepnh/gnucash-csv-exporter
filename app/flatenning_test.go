package app

import (
	"math/big"
	"testing"
	"time"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func newTransaction(description string, account1 uuid.UUID, value1 int64, account2 uuid.UUID, value2 int64) *Transaction {
	return &Transaction{
		Id:          NewUUID(),
		Description: description,
		Creation:    time.Now(),
		Date:        time.Now(),
		Account1:    account1,
		Ref1:        NewUUID(),
		Value1:      big.NewRat(value1, 1),
		Account2:    account2,
		Ref2:        NewUUID(),
		Value2:      big.NewRat(value2, 1),
	}
}

/*
func groupAccountsById(accounts []Account) map[string]Account {
	grouped := make(map[string]Account)
	for _, acc := range accounts {
		grouped[acc.Id.String()] = acc
	}
	return grouped
}
*/

func TestShouldFlattenAccountsIntoTransactions(t *testing.T) {
	accounts := []Account{
		Account{
			Id:   NewUUID(),
			Name: "Account1",
			Type: "ASSET",
		},
		Account{
			Id:   NewUUID(),
			Name: "Account2",
			Type: "BANK",
		},
	}
	book := Book{
		Accounts: accounts,
		Transactions: []Transaction{
			*newTransaction("Tx1", accounts[0].Id, -100, accounts[1].Id, 100),
			*newTransaction("Tx2", accounts[1].Id, -500, accounts[0].Id, 500),
		},
	}

	accountsByID := groupAccountsById(accounts)

	flattened := Flatten(&book)
	for _, flattenedTx := range flattened {
		match := false
		for _, tx := range book.Transactions {
			if flattenedTx.Id == tx.Id {
				match = true
				assert.Equal(t, tx.Creation, flattenedTx.Creation)
				assert.Equal(t, tx.Date, flattenedTx.Date)
				assert.Equal(t, tx.Description, flattenedTx.Description)
				assert.Equal(t, tx.Value1, flattenedTx.Value1)
				assert.Equal(t, tx.Value2, flattenedTx.Value2)
				assert.Equal(t, tx.Ref1, flattenedTx.Ref1)
				assert.Equal(t, tx.Ref2, flattenedTx.Ref2)
				account1 := accountsByID[tx.Account1.String()]
				assert.Equal(t, account1.Name, flattenedTx.Account1)
				assert.Equal(t, account1.Type, flattenedTx.AccType1)
				account2 := accountsByID[tx.Account2.String()]
				assert.Equal(t, account2.Name, flattenedTx.Account2)
				assert.Equal(t, account2.Type, flattenedTx.AccType2)
			}
		}
		assert.True(t, match)
	}
}
