package app

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func deserialize(t *testing.T, file string) *Root {
	xmlFile, err := os.Open("../xml-samples/" + file)
	if assert.NoError(t, err, "Could not open file") {
		defer xmlFile.Close()
		bytes, err := ioutil.ReadAll(xmlFile)
		if assert.NoError(t, err) {
			var root Root
			xml.Unmarshal(bytes, &root)
			return &root
		}
	}
	t.FailNow()
	return &Root{}
}

func TestAccountAmount(t *testing.T) {
	root := deserialize(t, "simple-uncompressed.xml")
	assert.Equal(t, 19, len(root.Book.Accounts))
}

func TestTransactionAmount(t *testing.T) {
	root := deserialize(t, "simple-uncompressed.xml")
	assert.Equal(t, 12, len(root.Book.Transactions))
}

func TestRootAccountDeserialization(t *testing.T) {
	assert := assert.New(t)
	root := deserialize(t, "simple-uncompressed.xml")
	rootAcc, err := root.Book.findRootAccount()
	assert.NoError(err)
	assert.Equal("ROOT", rootAcc.Type)
	assert.Equal("Root Account", rootAcc.Name)
	assert.Equal("0d81582f04c310d343d7a9270a9dad7b", rootAcc.Id)
	assert.Empty(rootAcc.Parent)
}

func TestAlimentacaoAccountDeserialization(t *testing.T) {
	assert := assert.New(t)
	root := deserialize(t, "simple-uncompressed.xml")
	for _, account := range root.Book.Accounts {
		if account.Name == "Alimentação" {
			assert.Equal("EXPENSE", account.Type)
			assert.Equal("491ef27091f545f74a74325c8cde20dd", account.Id)
			assert.Equal("7346f237a234ba1373a38653e6369191", account.Parent)
			return
		}
	}
	assert.FailNow("Account not found by name")
}

// nolint: golint
func assertSplit(t *testing.T, expectedId string, expectedValue string, expectedAccount string, split *GnuCashSplit) {
	assert := assert.New(t)
	assert.Equal(expectedId, split.Id)
	assert.Equal(expectedValue, split.Value)
	assert.Equal(expectedAccount, split.Account)
}

func TestTransactionDeserialization(t *testing.T) {
	assert := assert.New(t)
	root := deserialize(t, "simple-uncompressed.xml")
	for _, tx := range root.Book.Transactions {
		if tx.Id == "c5d90c78d80428aacfe63f8c3dcb46a0" {
			assert.Equal("2017-01-19 00:00:00 -0200", tx.DatePosted.Data)
			assert.Equal("2018-01-20 16:44:08 -0200", tx.DateEntered.Data)
			assert.Equal("Salário", tx.Description)
			assert.Equal(2, len(tx.Splits))

			assertSplit(t, "1f476aeb6bfbaa1b7de801d5915454f8", "100000/100", "0010af774cf308a29fa01b20ba3228b7", &tx.Splits[0])
			assertSplit(t, "7453cef925b00358da4950421b4bc887", "-100000/100", "7c14ba7c920fe7a430a4f1ef6b61aa2c", &tx.Splits[1])
			return
		}
	}
	assert.FailNow("Transaction not found by id")
}
