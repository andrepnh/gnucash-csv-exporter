# gnucash-csv-exporter
Exports a subset of GnuCash data to a single csv file. Transactions and matching accounts will be merged into single observations.

Usage:
```
go build exporter.go
exporter -i ./xml-samples/simple-compressed.gnucash -o simple.csv
```
