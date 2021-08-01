package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/piquette/finance-go/quote"
	"github.com/rodaine/table"
)

type SymbolPrice struct {
	Symbol          string
	Price           float64
	Change24        float64
	Change24Percent float64
	Currency        string
}

type SymbolNotFoundError string

func (e SymbolNotFoundError) Error() string {
	return fmt.Sprintf("Symbol not found: %s", string(e))
}

func getQuote(symbol string) (SymbolPrice, error) {
	q, err := quote.Get(symbol)
	if err != nil {
		return SymbolPrice{}, err
	}
	if q == nil {
		return SymbolPrice{}, SymbolNotFoundError(symbol)
	}
	return SymbolPrice{
		Symbol:          q.Symbol,
		Price:           q.RegularMarketPrice,
		Change24:        q.RegularMarketChange,
		Change24Percent: q.RegularMarketChangePercent,
		Currency:        q.CurrencyID,
	}, err
}

func getQuotes(symbols []string) []SymbolPrice {
	var symbolPrices []SymbolPrice
	for _, symbol := range symbols {
		sp, err := getQuote(symbol)
		if err != nil {
			fmt.Println(err)
			break
		}
		symbolPrices = append(symbolPrices, sp)
	}
	return symbolPrices
}

func displayTable(symbolPrices []SymbolPrice) {
	tbl := table.New("", "", "", "").WithHeaderFormatter(func(s string, a ...interface{}) string { return "" })

	changeColor := func(change float64) *color.Color {
		if change < 0 {
			return color.New(color.FgRed)
		} else {
			return color.New(color.FgGreen)
		}
	}

	for _, sp := range symbolPrices {
		currencySymbol := sp.Currency + " "
		if sp.Currency == "USD" {
			currencySymbol = "$"
		} else if sp.Currency == "CAD" {
			currencySymbol = "CA$"
		}
		tbl.AddRow(
			color.New(color.FgYellow).Sprintf(sp.Symbol),
			fmt.Sprintf("%25s", color.New(color.Faint).Sprintf(currencySymbol)+color.New(color.Bold).Sprintf("%.2f", sp.Price)),
			fmt.Sprintf("%18s", changeColor(sp.Change24).Sprintf("%+.2f", sp.Change24)),
			fmt.Sprintf("%13s", changeColor(sp.Change24).Sprintf("%+.2f%%", sp.Change24Percent)),
		)
	}

	tbl.Print()
}

func main() {
	symbols := os.Args[1:]
	symbolPrices := getQuotes(symbols)
	displayTable(symbolPrices)
}
