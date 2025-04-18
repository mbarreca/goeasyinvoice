package goeasyinvoice

/*
	Heavily inspired by https://github.com/maaslalani/invoice
*/

import (
	"image"
	"os"
	"strconv"
	"strings"

	"github.com/signintech/gopdf" // MIT
)

const (
	quantityColumnOffset = 360
	rateColumnOffset     = 405
	amountColumnOffset   = 480
	logoOffset           = 460
)

func (i *Invoice) writeLogo(pdf *gopdf.GoPdf, logo string, from string) error {
	if logo != "" {
		err, width, height := getImageDimension(logo)
		if err != nil {
			return err
		}
		scaledWidth := 50.0
		scaledHeight := float64(height) * scaledWidth / float64(width)
		if err := pdf.Image(logo, (logoOffset + pdf.GetX()), pdf.GetY(), &gopdf.Rect{W: scaledWidth, H: scaledHeight}); err != nil {
			return err
		}
		pdf.Br(scaledHeight + 24)
	}
	pdf.SetTextColor(55, 55, 55)

	formattedFrom := strings.ReplaceAll(from, `\n`, "\n")
	fromLines := strings.Split(formattedFrom, "\n")

	for j := range fromLines {
		if j == 0 {
			if err := pdf.SetFont(i.headlineFont.Title, "", 12); err != nil {
				return err
			}
			if err := pdf.Cell(nil, fromLines[j]); err != nil {
				return err
			}
			pdf.Br(18)
		} else {
			if err := pdf.SetFont(i.headlineFont.Title, "", 10); err != nil {
				return err
			}
			if err := pdf.Cell(nil, fromLines[j]); err != nil {
				return err
			}
			pdf.Br(15)
		}
	}
	pdf.Br(21)
	pdf.SetStrokeColor(225, 225, 225)
	pdf.Line(pdf.GetX(), pdf.GetY(), 260, pdf.GetY())
	pdf.Br(36)
	return nil
}

func (i *Invoice) writeTitle(pdf *gopdf.GoPdf, title, id, date string) error {
	if err := pdf.SetFont(i.headlineFont.Title, "", 24); err != nil {
		return err
	}
	pdf.SetTextColor(0, 0, 0)
	if err := pdf.Cell(nil, title); err != nil {
		return err
	}
	pdf.Br(36)
	if err := pdf.SetFont(i.headlineFont.Title, "", 12); err != nil {
		return err
	}
	pdf.SetTextColor(100, 100, 100)
	if err := pdf.Cell(nil, "#"); err != nil {
		return err
	}
	if err := pdf.Cell(nil, id); err != nil {
		return err
	}
	pdf.SetTextColor(150, 150, 150)
	if err := pdf.Cell(nil, "  ·  "); err != nil {
		return err
	}
	pdf.SetTextColor(100, 100, 100)
	if err := pdf.Cell(nil, date); err != nil {
		return err
	}
	pdf.Br(48)
	return nil
}

func (i *Invoice) writeDueDate(pdf *gopdf.GoPdf, due string) error {
	if err := pdf.SetFont(i.regularFont.Title, "", 9); err != nil {
		return err
	}
	pdf.SetTextColor(75, 75, 75)
	pdf.SetX(rateColumnOffset)
	if err := pdf.Cell(nil, "Due Date"); err != nil {
		return err
	}
	pdf.SetTextColor(0, 0, 0)
	if err := pdf.SetFontSize(11); err != nil {
		return err
	}
	pdf.SetX(amountColumnOffset - 0)
	if err := pdf.Cell(nil, due); err != nil {
		return err
	}
	pdf.Br(12)
	return nil
}

func (i *Invoice) writeBillTo(pdf *gopdf.GoPdf, to string) error {
	pdf.SetTextColor(75, 75, 75)
	if err := pdf.SetFont(i.regularFont.Title, "", 9); err != nil {
		return err
	}
	if err := pdf.Cell(nil, "BILL TO"); err != nil {
		return err
	}
	pdf.Br(18)
	pdf.SetTextColor(75, 75, 75)

	formattedTo := strings.ReplaceAll(to, `\n`, "\n")
	toLines := strings.Split(formattedTo, "\n")

	for j := range toLines {
		if j == 0 {
			if err := pdf.SetFont(i.regularFont.Title, "", 15); err != nil {
				return err
			}
			if err := pdf.Cell(nil, toLines[j]); err != nil {
				return err
			}
			pdf.Br(20)
		} else {
			if err := pdf.SetFont(i.regularFont.Title, "", 10); err != nil {
				return err
			}
			if err := pdf.Cell(nil, toLines[j]); err != nil {
				return err
			}
			pdf.Br(15)
		}
	}
	pdf.Br(64)
	return nil
}

func (i *Invoice) writeHeaderRow(pdf *gopdf.GoPdf) error {
	if err := pdf.SetFont(i.regularFont.Title, "", 9); err != nil {
		return err
	}
	pdf.SetTextColor(55, 55, 55)
	if err := pdf.Cell(nil, "ITEM"); err != nil {
		return err
	}
	pdf.SetX(quantityColumnOffset)
	if err := pdf.Cell(nil, "QTY"); err != nil {
		return err
	}
	pdf.SetX(rateColumnOffset)
	if err := pdf.Cell(nil, "RATE"); err != nil {
		return err
	}
	pdf.SetX(amountColumnOffset)
	if err := pdf.Cell(nil, "AMOUNT"); err != nil {
		return err
	}
	pdf.Br(24)
	return nil
}

func (i *Invoice) writeNotes(pdf *gopdf.GoPdf, notes string) error {
	pdf.SetY(600)

	if err := pdf.SetFont(i.regularFont.Title, "", 9); err != nil {
		return err
	}
	pdf.SetTextColor(55, 55, 55)
	if err := pdf.Cell(nil, "NOTES"); err != nil {
		return err
	}
	pdf.Br(18)
	if err := pdf.SetFont(i.regularFont.Title, "", 9); err != nil {
		return err
	}
	pdf.SetTextColor(0, 0, 0)

	formattedNotes := strings.ReplaceAll(notes, `\n`, "\n")
	notesLines := strings.Split(formattedNotes, "\n")

	for i := range notesLines {
		if err := pdf.Cell(nil, notesLines[i]); err != nil {
			return err
		}
		pdf.Br(15)
	}

	pdf.Br(48)
	return nil
}
func (i *Invoice) writeFooter(pdf *gopdf.GoPdf, id string) error {
	pdf.SetY(800)
	if err := pdf.SetFont(i.regularFont.Title, "", 10); err != nil {
		return err
	}
	pdf.SetTextColor(55, 55, 55)
	if err := pdf.Cell(nil, id); err != nil {
		return err
	}
	pdf.SetStrokeColor(225, 225, 225)
	pdf.Line(pdf.GetX()+10, pdf.GetY()+6, 550, pdf.GetY()+6)
	pdf.Br(48)
	return nil
}

func (i *Invoice) writeRow(pdf *gopdf.GoPdf, item string, quantity int, rate float64) error {
	if err := pdf.SetFont(i.regularFont.Title, "", 11); err != nil {
		return err
	}
	pdf.SetTextColor(0, 0, 0)

	total := float64(quantity) * rate
	amount := strconv.FormatFloat(total, 'f', 2, 64)

	if err := pdf.Cell(nil, item); err != nil {
		return err
	}
	pdf.SetX(quantityColumnOffset)
	if err := pdf.Cell(nil, strconv.Itoa(quantity)); err != nil {
		return err
	}
	pdf.SetX(rateColumnOffset)
	if err := pdf.Cell(nil, i.currency.Symbol+strconv.FormatFloat(rate, 'f', 2, 64)); err != nil {
		return err
	}
	pdf.SetX(amountColumnOffset)
	if err := pdf.Cell(nil, i.currency.Symbol+amount); err != nil {
		return err
	}
	pdf.Br(24)
	return nil
}

func (i *Invoice) writeTotals(pdf *gopdf.GoPdf, subtotal, tax, taxPercent, discount, discountPercent float64) error {
	pdf.SetY(600)

	if err := i.writeTotal(pdf, i.Subtotal, subtotal); err != nil {
		return err
	}
	if tax > 0 {
		if err := i.writeTotal(pdf, i.Tax+" ("+strconv.FormatFloat(taxPercent*100, 'f', -2, 64)+"%)", tax); err != nil {
			return err
		}
	}
	if discount > 0 {
		if err := i.writeTotal(pdf, i.Discount+" ("+strconv.FormatFloat(discountPercent*100, 'f', -2, 64)+"%)", discount); err != nil {
			return err
		}
	}
	if err := i.writeTotal(pdf, i.Total, subtotal+tax-discount); err != nil {
		return err
	}
	return nil
}

func (i *Invoice) writeTotal(pdf *gopdf.GoPdf, label string, total float64) error {
	if err := pdf.SetFont(i.regularFont.Title, "", 9); err != nil {
		return err
	}
	pdf.SetTextColor(75, 75, 75)
	pdf.SetX(rateColumnOffset)
	if err := pdf.Cell(nil, label); err != nil {
		return err
	}
	pdf.SetTextColor(0, 0, 0)
	if err := pdf.SetFontSize(12); err != nil {
		return err
	}
	pdf.SetX(amountColumnOffset)
	if label == i.Total {
		if err := pdf.SetFont(i.boldFont.Title, "", 11.5); err != nil {
			return err
		}
	}
	if err := pdf.Cell(nil, i.currency.Symbol+strconv.FormatFloat(total, 'f', 2, 64)); err != nil {
		return err
	}
	pdf.Br(24)
	return nil
}

func getImageDimension(imagePath string) (error, int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		return err, 0, 0
	}
	defer file.Close()

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		return err, 0, 0
	}
	return nil, image.Width, image.Height
}
