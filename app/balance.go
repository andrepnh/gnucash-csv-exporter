package app

import (
	"math/big"
	"sort"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

type transactions []Transaction

func (txs transactions) Len() int {
	return len(txs)
}

func (txs transactions) Swap(i, j int) {
	txs[i], txs[j] = txs[j], txs[i]
}

func (txs transactions) Less(i, j int) bool {
	if txs[i].Date == txs[j].Date {
		if txs[i].Creation == txs[j].Creation {
			return strings.Compare(txs[i].Description, txs[j].Description) < 0
		}
		return txs[i].Creation.Before(txs[j].Creation)
	}
	return txs[i].Date.Before(txs[j].Date)
}

// AccountBalance - keeps track of account balance at a given Date
type AccountBalance struct {
	Id     uuid.UUID // nolint: golint
	Name   string
	Type   string
	Parent uuid.UUID

	Date        time.Time
	Balance     *big.Rat
	Transaction uuid.UUID
}

func CalcAccountBalances(book *Book) []AccountBalance { // nolint: golint
	accountsById := groupAccountsById(book.Accounts) // nolint: golint
	sortedTxs := make([]Transaction, len(book.Transactions))
	copy(sortedTxs, book.Transactions)
	sort.Sort(transactions(sortedTxs))
	balanceByAccount := make(map[string]*big.Rat)
	balances := make([]AccountBalance, len(book.Transactions)*2)
	for i, tx := range sortedTxs {
		acc := accountsById[tx.Account1.String()]
		balances[i*2] = *newAccountBalance(&tx, tx.Value1, &acc, balanceByAccount)
		acc = accountsById[tx.Account2.String()]
		balances[i*2+1] = *newAccountBalance(&tx, tx.Value2, &acc, balanceByAccount)
	}
	return balances
}

func newAccountBalance(tx *Transaction, value *big.Rat, acc *Account, balanceByAccount map[string]*big.Rat) *AccountBalance {
	currentBalance, found := balanceByAccount[acc.Id.String()]
	if !found {
		balanceByAccount[acc.Id.String()], currentBalance = big.NewRat(0, 1), big.NewRat(0, 1)
	}
	newBalance := value.Add(value, currentBalance)
	balanceByAccount[acc.Id.String()] = newBalance
	return &AccountBalance{
		Id:          acc.Id,
		Name:        acc.Name,
		Type:        acc.Type,
		Parent:      acc.Parent,
		Date:        tx.Date,
		Balance:     newBalance,
		Transaction: tx.Id,
	}
}
