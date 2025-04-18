# goeasyinvoice
<div align="center">

Create a PDF Invoice the easy way using Go Easy Invoice.

`Special thanks to https://github.com/maaslalani/invoice and https://github.com/signintech/gopdf`

[![GoDoc][doc-img]][doc]

<div align="left">

## Installation

<div align="center">

`go get github.com/mbarreca/goeasyinvoice`

<div align="left">

Go Easy Invoice was built and tested with Go 1.24, it may still work with prior versions however it has not been tested so use at your own risk.

## Why would I use this library?

In an world thats increasingly going towards vendor lock-in I believe its important to keep yourself as agnostic as possible. Billing is a big culprit. I couldn't find an easy to use library to generate me a modern looking invoice so I decided to take whats out there, improve it and give it back to the open source world. Enjoy.

## Core Concepts

This library is built to be extensible within a particular format and template. It's fairly simple to modify. Everything is error handled, well tested and should work without issue. See the example below to make it work.

### Telemetry

A library like this does not have much in the way of telemetry outside of error handling. That has intentionally been made robust, make sure your errors are being handled well and it will integrate into your telemetry stack easily.

## Example
```
package main

import (
	"embed"
	"testing"

	invoice "github.com/mbarreca/goeasyinvoice" // Apache 2.0
)
/*
Create Embeds for TTF Files
*/
var f embed.FS

//go:embed "assets/fonts/Go-Regular.ttf"
var regularFont []byte

//go:embed "assets/fonts/Go-Medium.ttf"
var boldFont []byte

//go:embed "assets/fonts/Go-Bold.ttf"
var headlineFont []byte

func main() {
	/*
		Setup Invoice Object
	*/
	// Setup Fonts
	var rFont invoice.Font
	var bFont invoice.Font
	var hFont invoice.Font
	rFont.Title = "Go-Regular"
	rFont.File = &regularFont
	bFont.Title = "Go-Medium"
	bFont.File = &boldFont
	hFont.Title = "Go-Bold"
	hFont.File = &headlineFont

	// Setup Currency
	var currency invoice.Currency
	currency.Label = "EUR"
	currency.Symbol = "€"

	// Generate Line Item
	var lineItems []invoice.InvoiceItem
	var lineItem invoice.InvoiceItem
	lineItem.Item = "Software Development Services"
	lineItem.Rate = 100.00
	lineItem.Quantity = 100
	lineItems = append(lineItems, lineItem)

	// Generate Invoice Object
	invoice, err := invoice.New(&rFont, &bFont, &hFont, &currency, "Subtotal", "Discount", "Tax", "Total", "assets/logo.png", "Test Company")
	if err != nil {
		t.Fatal(err)
	}
	id := "DEV-001"

	// Generate Invoice
	if err := invoice.Generate("Development Services", id, "2025-04-20", "2025-06-01", "Test Company\n123 Address Avenue\nMilan, Italy", "Thank you for your business!", "./invoice"+id+".pdf", 0.10, 0.15, lineItems); err != nil {
	panic(err)
	}
}

```

## License

Goeasyinvoice is licensed under Apache 2.0.

[doc]: https://pkg.go.dev/github.com/mbarreca/goeasyinvoice
[doc-img]: https://pkg.go.dev/badge/github.com/mbarreca/goeasyinvoice
