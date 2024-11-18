package user

import (
	"bytes"
	"fmt"

	"admin/models"

	"github.com/signintech/gopdf"
)

func GeneratePDF(invoice models.Invoice) ([]byte, error) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	err := pdf.AddTTFFont("DejaVuSans", "/home/athul/Documents/The Furnish spot/fonts/arial.ttf")
	if err != nil {
		return nil, fmt.Errorf("failed to add font: %w", err)
	}
	err = pdf.SetFont("DejaVuSans", "", 12)
	if err != nil {
		return nil, fmt.Errorf("cannot set font: %w", err)
	}

	pdf.Cell(nil, fmt.Sprintf("Invoice ID: %s", invoice.InvoiceID))
	pdf.Br(17)

	pdf.Cell(nil, fmt.Sprintf("Date: %s", invoice.Date.Format("02-Jan-2006")))
	pdf.Br(17)

	pdf.Cell(nil, fmt.Sprintf("Customer ID: %d", invoice.UserID))
	pdf.Br(17)

	pdf.SetFont("DejaVuSans", "", 10)

	pdf.Cell(nil, "ProductID")
	pdf.SetX(150)
	pdf.Cell(nil, "Qty")
	pdf.SetX(180)
	pdf.Cell(nil, "Unit Price")
	pdf.SetX(230)
	pdf.Cell(nil, "Discount")
	pdf.SetX(280)
	pdf.Cell(nil, "Total")
	pdf.Br(15)

	pdf.Line(5, pdf.GetY(), 400, pdf.GetY())
	pdf.Br(10)

	subtotal := 0.0
	for _, item := range invoice.Items {
		fmt.Printf("Item: %+v\n", item)

		itemTotal := (item.UnitPrice * float64(item.Quantity)) - item.Discount
		subtotal += itemTotal

		pdf.Cell(nil, fmt.Sprintf("%d", item.ProductID))
		pdf.SetX(150)
		pdf.Cell(nil, fmt.Sprintf("%d", item.Quantity))

		pdf.SetX(180)
		pdf.Cell(nil, fmt.Sprintf("%.2f", item.UnitPrice))

		pdf.SetX(230)
		pdf.Cell(nil, fmt.Sprintf("%.2f", item.Discount))

		pdf.SetX(280)
		pdf.Cell(nil, fmt.Sprintf("%.2f", itemTotal))

		pdf.Br(10)
	}

	discountApplied := invoice.Subtotal - subtotal
	total := subtotal - discountApplied

	pdf.Br(15)
	pdf.SetFont("DejaVuSans", "", 12)
	pdf.Cell(nil, fmt.Sprintf("Subtotal: %.2f", subtotal))
	pdf.Br(15)
	pdf.Cell(nil, fmt.Sprintf("Discount Applied: %.2f", discountApplied))
	pdf.Br(15)

	pdf.Cell(nil, fmt.Sprintf("Total: %.2f", total))
	pdf.Br(15)

	var buffer bytes.Buffer
	_, err = pdf.WriteTo(&buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to write PDF: %w", err)
	}

	return buffer.Bytes(), nil
}
