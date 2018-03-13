package main

import (
	"flag"
	"log"
	"os"

	"github.com/andrepnh/gnucash-es-exporter/app"
)

func toCsv(gnucashFile, csvFile string) {
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
	app.Export(flattened, csvFile)
	log.Println("Done.")
}

func main() {
	gnucashFile := flag.String("i", "", "Input .gnucash file")
	csvFile := flag.String("o", "transactions.csv", "Output csv file")
	flag.Parse()
	if *gnucashFile == "" {
		panic("Input file is mandatory")
	}

	toCsv(*gnucashFile, *csvFile)
}
