package tests

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

// Test PDF Generation
func TestGeneratePDF(t *testing.T) {
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
	if err := invoice.Generate("Development Services", id, "2025-04-20", "2025-06-01", "Test Company\n123 Address Avenue\nMilan, Italy", "Thank you for your business!", "./invoice"+id+".pdf", 0.10, 0.15, lineItems); err != nil {
		t.Fatal(err)
	}
}
