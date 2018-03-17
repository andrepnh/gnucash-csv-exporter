# gnucash-csv-exporter
Exports a subset of GnuCash data to two csv files. One contains transactions and matching accounts merged into single observations; the other tracks account balances over time.

Usage:
```
go build exporter.go
exporter -in ./xml-samples/simple-compressed.gnucash -outdir .
```
