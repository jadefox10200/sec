Note: This code was originally writen by github.com/tradyfinance.
However, it was taken down and so I changed it over to my own repo
for dependancy purposes. 

# SEC

SEC is a library for accessing SEC filings.

## Examples

### EDGAR Index Entries

```go
end := time.Now()
start := end.AddDate(0, -1, 0)
if err := sec.GetEDGARIndexEntries(start, end, func(e sec.EDGARIndexEntry) error {
    fmt.Printf("%+v\n", e)
    return nil
}); err != nil {
    log.Fatal(err)
}
```

### Form 4 Filings

```go
end := time.Now()
start := end.AddDate(0, -1, 0)
if err := sec.GetForm4Filings(start, end, func(form sec.Form4) error {
    fmt.Printf("%+v\n", form)
    return nil
}); err != nil {
    log.Fatal(err)
}
```

## Documentation

Documentation is available [here](https://godoc.org/github.com/tradyfinance/sec).

## License

This project is released under the [Apache License, Version 2.0](LICENSE).
