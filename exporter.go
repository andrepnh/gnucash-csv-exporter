package main

import (
	"flag"
	"log"
	"os"

	"github.com/andrepnh/gnucash-es-exporter/app"
)

func toCsv(gnucashFile, outdir string) {
	gcFile, err := os.Open(gnucashFile)
	if err != nil {
		panic(err)
	}
	xml := app.Decompress(gcFile)
	root := app.Deserialize(xml)
	log.Println("Deserialized", len(root.Book.Accounts), "accounts")
	log.Println("Deserialized", len(root.Book.Transactions), "transactions")
	book := *app.ConvertBook(&root.Book)
	log.Println(len(book.Accounts), "accounts successfully converted")
	log.Println(len(book.Transactions), "transactions successfully converted")
	flattened := app.Flatten(&book)
	log.Println("Merged all accounts into transactions, exporting...")
	app.ExportTransactions(flattened, outdir+"/transactions.csv")
	balances := app.CalcAccountBalances(&book)
	log.Println("Calculated account balances, exporting...")
	app.ExportBalances(balances, outdir+"/accounts.csv")
	log.Println("Done.")
}

func main() {
	gnucashFile := flag.String("in", "", "Input .gnucash file")
	outdir := flag.String("outdir", ".", "Path to write output csvs")
	flag.Parse()
	if *gnucashFile == "" {
		panic("Input file is mandatory")
	}

	toCsv(*gnucashFile, *outdir)
}
