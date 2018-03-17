package app

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func shuffle(txs []Transaction) []Transaction {
	shuffled := make([]Transaction, len(txs))
	for i, v := range rand.Perm(len(txs)) {
		shuffled[v] = txs[i]
	}
	return shuffled
}

func assertBalanceOfAccount(t *testing.T, account *Account, balance *AccountBalance) {
	assert.Equal(t, account.Id, balance.Id)
	assert.Equal(t, account.Name, balance.Name)
	assert.Equal(t, account.Parent, balance.Parent)
	assert.Equal(t, account.Type, balance.Type)
}

func assertBalanceOfTransaction(t *testing.T, tx *Transaction, balance *AccountBalance) {
	assert.Equal(t, tx.Id, balance.Transaction)
	assert.Equal(t, tx.Date, balance.Date)
}

func TestCalcAccountBalancesShouldReturnCorrectBalances(t *testing.T) {
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
	var txValue int64 = 100
	tx1, tx2, tx3 := *newTransaction("tx1", accounts[0].Id, txValue, accounts[1].Id, -txValue),
		*newTransaction("tx2", accounts[1].Id, -txValue, accounts[0].Id, txValue),
		*newTransaction("tx3", accounts[0].Id, -txValue, accounts[1].Id, txValue)
	book := Book{
		Accounts:     accounts,
		Transactions: shuffle([]Transaction{tx3, tx2, tx1}),
	}

	balances := CalcAccountBalances(&book)
	assert.Equal(t, 6, len(balances))
	assert.Equal(t, big.NewRat(txValue, 1).FloatString(4), balances[0].Balance.FloatString(4))
	assert.Equal(t, big.NewRat(-txValue, 1).FloatString(4), balances[1].Balance.FloatString(4))
	assert.Equal(t, big.NewRat(-txValue*2, 1).FloatString(4), balances[2].Balance.FloatString(4))
	assert.Equal(t, big.NewRat(txValue*2, 1).FloatString(4), balances[3].Balance.FloatString(4))
	assert.Equal(t, big.NewRat(txValue, 1).FloatString(4), balances[4].Balance.FloatString(4))
	assert.Equal(t, big.NewRat(-txValue, 1).FloatString(4), balances[5].Balance.FloatString(4))
}

func TestCalcAccountBalancesShouldCopyAccountData(t *testing.T) {
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
			*newTransaction("tx1", accounts[0].Id, 1, accounts[1].Id, -1),
			*newTransaction("tx1", accounts[1].Id, 1, accounts[0].Id, -1),
		},
	}

	balances := CalcAccountBalances(&book)
	assert.Equal(t, 4, len(balances))
	assertBalanceOfAccount(t, &accounts[0], &balances[0])
	assertBalanceOfAccount(t, &accounts[1], &balances[1])
	assertBalanceOfAccount(t, &accounts[1], &balances[2])
	assertBalanceOfAccount(t, &accounts[0], &balances[3])
}

func TestCalcAccountBalancesShouldCopyTransactionData(t *testing.T) {
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
			*newTransaction("tx1", accounts[0].Id, 1, accounts[1].Id, -1),
			*newTransaction("tx1", accounts[1].Id, 1, accounts[0].Id, -1),
		},
	}

	balances := CalcAccountBalances(&book)
	assert.Equal(t, 4, len(balances))
	assertBalanceOfTransaction(t, &book.Transactions[0], &balances[0])
	assertBalanceOfTransaction(t, &book.Transactions[0], &balances[1])
	assertBalanceOfTransaction(t, &book.Transactions[1], &balances[2])
	assertBalanceOfTransaction(t, &book.Transactions[1], &balances[3])
}
