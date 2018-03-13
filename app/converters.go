package app

import (
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/satori/go.uuid"
)

// GcDateTimeFormat - format used in GnuCash XML
const GcDateTimeFormat = "2006-01-02 15:04:05 -0700"

func uuidFromStringOrPanic(str string) uuid.UUID {
	uuid, err := uuid.FromString(str)
	if err != nil {
		panic(err)
	}
	return uuid
}

func parseTimeOrPanic(str string) time.Time {
	time, err := time.Parse(GcDateTimeFormat, str)
	if err != nil {
		panic(err)
	}
	return time
}

func parseIntOrPanic(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func parseMoneyOrPanic(str string) *big.Rat {
	parts := strings.Split(str, "/")
	return big.NewRat(parseIntOrPanic(parts[0]), parseIntOrPanic(parts[1]))
}

func convertAccount(account *GnuCashAccount) *Account {
	converted := Account{
		Id:   uuidFromStringOrPanic(account.Id),
		Name: account.Name,
		Type: account.Type,
	}
	if len(strings.TrimSpace(account.Parent)) > 0 {
		converted.Parent = uuidFromStringOrPanic(account.Parent)
	}
	return &converted
}

func incomplete(split *GnuCashSplit) bool {
	if len(split.Account) == 0 || len(split.Value) == 0 {
		return true
	}
	return false
}

func convertTransaction(gcTx *GnuCashTransaction) (*Transaction, error) {
	if len(gcTx.Splits) != 2 {
		return nil, fmt.Errorf(
			"Ignoring GnuCashTransaction with unexpected splits: %v; tx id: %v",
			len(gcTx.Splits), gcTx.Id)
	}
	if incomplete(&gcTx.Splits[0]) || incomplete(&gcTx.Splits[1]) {
		return nil, fmt.Errorf(
			"Ignoring incomplete GnuCashTransaction missing required fields; tx id: %v",
			gcTx.Id)
	}
	desc := gcTx.Description
	if len(strings.TrimSpace(desc)) == 0 {
		log.Printf("Replacing empty transaction description by \"[UNKNOWN]\"; tx id: %v", gcTx.Id)
		desc = "[UNKNOWN]"
	}
	tx := Transaction{
		Id:          uuidFromStringOrPanic(gcTx.Id),
		Description: desc,
		Creation:    parseTimeOrPanic(gcTx.DateEntered),
		Date:        parseTimeOrPanic(gcTx.DatePosted),
		Value1:      parseMoneyOrPanic(gcTx.Splits[0].Value),
		Account1:    uuidFromStringOrPanic(gcTx.Splits[0].Account),
		Ref1:        uuidFromStringOrPanic(gcTx.Splits[0].Id),
		Value2:      parseMoneyOrPanic(gcTx.Splits[1].Value),
		Account2:    uuidFromStringOrPanic(gcTx.Splits[1].Account),
		Ref2:        uuidFromStringOrPanic(gcTx.Splits[1].Id),
	}
	return &tx, nil
}

func ConvertBook(book *GnuCashBook) *Book { // nolint: golint
	converted := Book{}
	for _, acc := range book.Accounts {
		converted.Accounts = append(converted.Accounts, *convertAccount(&acc))
	}
	for _, tx := range book.Transactions {
		convertedTx, ignored := convertTransaction(&tx)
		if ignored != nil {
			log.Print(ignored)
			continue
		}
		converted.Transactions = append(converted.Transactions, *convertedTx)
	}
	return &converted
}
