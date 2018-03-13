package app

import (
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func dashlessString(id uuid.UUID) string {
	return strings.Replace(id.String(), "-", "", -1)
}

func truncateClock(dateTime time.Time) time.Time {
	year, month, day := dateTime.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, dateTime.Location())
}

func newGcTransaction(description, value1, value2 string) *GnuCashTransaction {
	return &GnuCashTransaction{
		Id:          dashlessString(NewUUID()),
		Description: description,
		DateEntered: time.Now().Format(GcDateTimeFormat),
		DatePosted:  truncateClock(time.Now()).Format(GcDateTimeFormat),
		Splits: []GnuCashSplit{
			GnuCashSplit{
				Account: dashlessString(NewUUID()),
				Id:      dashlessString(NewUUID()),
				Value:   value1,
			},
			GnuCashSplit{
				Account: dashlessString(NewUUID()),
				Id:      dashlessString(NewUUID()),
				Value:   value2,
			},
		},
	}
}

func TestShouldConvertAccount(t *testing.T) {
	assert := assert.New(t)
	gcAcc := GnuCashAccount{
		Id:     "5401c8035feb761617af83de00f2ec16",
		Type:   "EQUITY",
		Name:   "Líquido",
		Parent: "0d81582f04c310d343d7a9270a9dad7b",
	}
	gcBook := GnuCashBook{
		Accounts: []GnuCashAccount{gcAcc},
	}
	book := ConvertBook(&gcBook)
	assert.Equal(1, len(book.Accounts))
	account := book.Accounts[0]
	assert.Equal(gcAcc.Id, dashlessString(account.Id))
	assert.Equal(gcAcc.Name, account.Name)
	assert.Equal("EQUITY", account.Type)
	parent, err := uuid.FromString(gcAcc.Parent)
	if assert.NoError(err) {
		assert.Equal(parent, account.Parent)
	}
}

func TestShouldConvertGnuCashTransaction(t *testing.T) {
	gcTx := *newGcTransaction("Fádida", "100/100", "-100/100")
	gcBook := GnuCashBook{
		Transactions: []GnuCashTransaction{gcTx},
	}
	book := ConvertBook(&gcBook)
	assert.Equal(t, 1, len(book.Transactions))
	tx := book.Transactions[0]
	assert.Equal(t, gcTx.Id, dashlessString(tx.Id))
	assert.Equal(t, gcTx.Description, tx.Description)
	assert.Equal(t, gcTx.DateEntered, tx.Creation.Format("2006-01-02 15:04:05 -0700"))
	assert.Equal(t, gcTx.DatePosted, tx.Date.Format("2006-01-02 15:04:05 -0700"))
}

func TestShouldConvertTransactionWithEmptyDescriptionAsUknown(t *testing.T) {
	gcBook := GnuCashBook{
		Transactions: []GnuCashTransaction{
			*newGcTransaction("", "100/100", "-100/100"),
		},
	}
	book := ConvertBook(&gcBook)
	assert.Equal(t, "[UNKNOWN]", book.Transactions[0].Description)
}

func TestShouldIgnoreTransactionMissingSplits(t *testing.T) {
	gcTx := *newGcTransaction("Foo", "100/100", "-100/100")
	gcTx.Splits = gcTx.Splits[:len(gcTx.Splits)-1]
	gcBook := GnuCashBook{
		Transactions: []GnuCashTransaction{gcTx},
	}
	book := ConvertBook(&gcBook)
	assert.Equal(t, 0, len(book.Transactions))
}

func TestShouldIgnoreTransactionMissingAccount(t *testing.T) {
	gcTx := *newGcTransaction("Foo", "100/100", "-100/100")
	gcBook := GnuCashBook{
		Transactions: []GnuCashTransaction{gcTx},
	}
	gcTx.Splits[0].Account = ""
	book := ConvertBook(&gcBook)
	assert.Equal(t, 0, len(book.Transactions))
}

func TestShouldIgnoreTransactionMissingValues(t *testing.T) {
	gcBook := GnuCashBook{
		Transactions: []GnuCashTransaction{
			*newGcTransaction("Foo", "", ""),
		},
	}
	book := ConvertBook(&gcBook)
	assert.Equal(t, 0, len(book.Transactions))
}

func TestShouldConvertGnuCashSplitIntoTransaction(t *testing.T) {
	gcTx := *newGcTransaction("ASD", "100000/100", "-100000/100")
	gcBook := GnuCashBook{
		Transactions: []GnuCashTransaction{gcTx},
	}
	book := ConvertBook(&gcBook)
	assert.Equal(t, 1, len(book.Transactions))
	tx := book.Transactions[0]
	assert.Equal(t, gcTx.Splits[0].Id, dashlessString(tx.Ref1))
	assert.Equal(t, gcTx.Splits[0].Account, dashlessString(tx.Account1))
	assert.Equal(t, big.NewRat(100000, 100), tx.Value1)
	assert.Equal(t, gcTx.Splits[1].Id, dashlessString(tx.Ref2))
	assert.Equal(t, gcTx.Splits[1].Account, dashlessString(tx.Account2))
	assert.Equal(t, big.NewRat(-100000, 100), tx.Value2)
}
