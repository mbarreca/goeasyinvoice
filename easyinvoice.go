package goeasyinvoice

/*
	Heavily inspired by https://github.com/maaslalani/invoice
*/

import (
	"errors"
	"os"

	"github.com/signintech/gopdf" // MIT
)

type Invoice struct {
	From     string
	Logo     string
	Title    string
	Subtotal string
	Discount string
	Tax      string
	Total    string

	currency *Currency

	regularFont  *Font
	boldFont     *Font
	headlineFont *Font
}

type Font struct {
	Title string
	File  *[]byte
}
type InvoiceItem struct {
	Item     string
	Quantity int
	Rate     float64
}

// Label - The Label of the Currency
type Currency struct {
	Label  string
	Symbol string
}

// Invoice Instantiator
// regularFont - A goeasyinvoice Font Object that contains a reference to the byte[] of the font and title of the font
// boldFont - A goeasyinvoice Font Object that contains a reference to the byte[] of the font and title of the font
// headlineFont - A goeasyinvoice Font Object that contains a reference to the byte[] of the font and title of the font
// subtotal - The label for your subtitle
// discount - The label for a discount
// tax - The label for taxes
// total - The label for the total sum
// logopath - The path to your logo
// from - Who the invoice is from
func New(regularFont, boldFont, headlineFont *Font, currency *Currency, subtotal, discount, tax, total, logopath, from string) (*Invoice, error) {
	var i Invoice
	i.regularFont = regularFont
	i.boldFont = boldFont
	i.headlineFont = headlineFont
	i.Logo = logopath
	i.From = from
	i.Subtotal = subtotal
	i.Discount = discount
	i.Tax = tax
	i.Total = total
	i.currency = currency
	return &i, nil
}

// Generate a PDF with the following parameters
// title - The title of the Invoice in the PDF
// id - Your Invoice ID
// issueDate - The Invoice's Issue Date
// dueDate - The Invoice's Due Date
// to - Who the Invoice is addressed to
// note - A note at the bottom left corner of the invoice - typically used for remarks
// filePathName - Where you would like the invoice saved. Make sure to clean this up if you're using it for SaaS
// discount - A discount percentage that will be multiplied by the subtotal - i.e 0.25 is a 25% discount - Please note that the subtotal is calculated pre-tax and the tax is applied after this
// tax - The taxable percentage (we will add 1 to it so if you apply a 15% tax then you would put 0.15)
// lineItems - An array of Invoice Items that represent the line items in your invoice
// RETURNS
// err - We will return an error or nil if successful. Check your filePathName for a successful PDF on nil return
func (i *Invoice) Generate(title, id, issueDate, dueDate, to, note, filePathName string, discount, tax float64, lineItems []InvoiceItem) error {
	// Validate inputs
	if len(title) < 2 || len(id) < 2 || len(issueDate) < 2 || len(dueDate) < 2 || len(to) < 2 || len(filePathName) < 2 || len(lineItems) < 1 {
		return errors.New("Goeasyinvoice: All inputs required")
	}
	// Create the PDF and set it up
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{
		PageSize: *gopdf.PageSizeA4,
	})
	// Setup Margins
	pdf.SetMargins(50, 50, 50, 50)
	pdf.AddPage()
	// Add fonts
	if err := pdf.AddTTFFontData(i.regularFont.Title, *i.regularFont.File); err != nil {
		return err
	}
	if err := pdf.AddTTFFontData(i.boldFont.Title, *i.boldFont.File); err != nil {
		return err
	}
	if err := pdf.AddTTFFontData(i.headlineFont.Title, *i.headlineFont.File); err != nil {
		return err
	}
	// Add the Top Structure
	if err := i.writeLogo(&pdf, i.Logo, i.From); err != nil {
		return err
	}
	if err := i.writeTitle(&pdf, title, id, issueDate); err != nil {
		return err
	}
	if err := i.writeBillTo(&pdf, to); err != nil {
		return err
	}
	if err := i.writeHeaderRow(&pdf); err != nil {
		return err
	}
	// Add Line Items
	subtotal := 0.0
	for _, item := range lineItems {
		if err := i.writeRow(&pdf, item.Item, item.Quantity, item.Rate); err != nil {
			return err
		}
		subtotal += float64(item.Quantity) * item.Rate
	}
	// OPTIONAL - Add the Note
	if note != "" {
		if err := i.writeNotes(&pdf, note); err != nil {
			return err
		}
	}
	// Write the final totals
	// Sort out final totals
	taxable := subtotal - (subtotal * discount)
	// Taxes are taken off of a discount
	if err := i.writeTotals(&pdf, subtotal, taxable*tax, tax, subtotal*discount, discount); err != nil {
		return err
	}
	// Write the due date
	if err := i.writeDueDate(&pdf, dueDate); err != nil {
		return err
	}
	// Write the footer
	if err := i.writeFooter(&pdf, id); err != nil {
		return err
	}
	// Create a file in the provided directory
	tempFile, err := os.Create(filePathName)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	// Write the PDF to the temporary file
	if err := pdf.WritePdf(tempFile.Name()); err != nil {
		return err
	}

	return nil
}
