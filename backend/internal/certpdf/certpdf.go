// Package certpdf renders the "подарочный сертификат" PDF — one fixed-design
// page (background artwork + brand font, both embedded as assets below)
// with a handful of dynamic fields drawn on top: certificate number, issue
// date, what it's redeemable for, and the recipient's name.
package certpdf

import (
	"bytes"
	_ "embed"
	"fmt"

	"github.com/go-pdf/fpdf"
	"github.com/skip2/go-qrcode"

	"floway-backend/internal/model"
)

//go:embed assets/background.png
var backgroundPNG []byte

//go:embed assets/SoyuzGrotesk-Bold.ttf
var soyuzGroteskBoldTTF []byte

// Page dimensions match the reference design exactly (extracted from the
// example PDFs via `pdfinfo`) — reusing them verbatim means the background
// artwork's proportions never need rescaling.
const (
	pageWidth  = 1181.0
	pageHeight = 1772.0

	// Sampled from the reference PDFs' rendered composite (pdftoppm -r 150
	// front.pdf, then picked pixel colors) — solid page fill under the
	// background artwork, and the dark-brown ink used for every text
	// element in the design.
	bgR, bgG, bgB    = 182, 212, 228 // #B6D4E4
	inkR, inkG, inkB = 65, 52, 42    // #41342A
	marginLeft       = 100.0
	marginRight      = 1080.0 // right edge text aligns against
	fontFamily       = "SoyuzGrotesk"

	purposeStartY         = 900.0
	purposeMaxFontSize    = 104.0
	purposeMinFontSize    = 60.0
	purposeLineHeight     = 108.0 // one line at purposeMaxFontSize; scales with the chosen size
	purposeToRecipientGap = 24.0
	recipientToQRGap      = 340.0
	qrToContactGap        = 170.0
	bottomMargin          = 60.0
	telegramContact       = "t.me/floway_school"
	websiteContact        = "kursfloristiki.ru"
	phoneContact          = "+7 985 226 1948"
	telegramURL           = "https://t.me/floway_school"
	websiteURL            = "https://kursfloristiki.ru"
)

// kindPurposeLabel is the fixed lead-in text for each certificate kind — see
// the reference designs: "на сумму" / "5000 рублей", "на любой мастер-класс"
// (no value), "на курс" / "<title>", "на мастер-класс" / "<title>".
var kindPurposeLabel = map[model.GiftCertificateKind]string{
	model.GiftCertificateKindAmount:         "на сумму",
	model.GiftCertificateKindCourse:         "на курс",
	model.GiftCertificateKindMasterclass:    "на мастер-класс",
	model.GiftCertificateKindAnyMasterclass: "на любой мастер-класс",
}

// Render builds the certificate PDF for cert and returns its raw bytes.
// dateLabel is the already-formatted Russian issue date (e.g. "1 сентября
// 2026") — formatting the month name is calendar/locale logic that belongs
// to the caller, not this rendering package.
func Render(cert model.GiftCertificate, dateLabel string) ([]byte, error) {
	pdf := fpdf.NewCustom(&fpdf.InitType{
		UnitStr: "pt",
		Size:    fpdf.SizeType{Wd: pageWidth, Ht: pageHeight},
	})
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddUTF8FontFromBytes(fontFamily, "", soyuzGroteskBoldTTF)
	pdf.SetFont(fontFamily, "", 12)
	pdf.AddPage()

	pdf.SetFillColor(bgR, bgG, bgB)
	pdf.Rect(0, 0, pageWidth, pageHeight, "F")

	pdf.RegisterImageOptionsReader("background", fpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(backgroundPNG))
	pdf.ImageOptions("background", 0, 0, pageWidth, pageHeight, false, fpdf.ImageOptions{ImageType: "PNG"}, 0, "")

	pdf.SetTextColor(inkR, inkG, inkB)

	drawRightAligned(pdf, marginRight, 140, 60, "школа фловей")
	pdf.SetFontSize(148)
	pdf.Text(marginLeft, 300, "подарочный")
	pdf.Text(marginLeft, 478, "сертификат")

	pdf.SetFontSize(74)
	pdf.Text(marginLeft, 600, cert.Number)
	drawRightAligned(pdf, marginRight, 600, 74, dateLabel)

	// The purpose block wraps to a variable number of lines (a long course
	// or masterclass title, or a long label like "на любой мастер-класс",
	// can each take 1-2 lines). Shrinking its font size when needed keeps
	// everything below it — recipient, QR codes, contact row — on the page
	// regardless of how long the admin's free-text value turns out to be.
	purposeBottom := drawPurpose(pdf, cert.Kind, cert.Value)
	recipientY := purposeBottom + purposeToRecipientGap

	pdf.SetFontSize(62)
	pdf.Text(marginLeft, recipientY, cert.Recipient)

	qrY := recipientY + recipientToQRGap
	if err := drawQRCodes(pdf, qrY); err != nil {
		return nil, fmt.Errorf("render qr codes: %w", err)
	}

	contactY := qrY + qrToContactGap
	pdf.SetFontSize(28)
	pdf.Text(marginLeft, contactY, telegramContact)
	pdf.Text(433, contactY, websiteContact)
	drawRightAligned(pdf, marginRight, contactY, 28, phoneContact)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("render certificate pdf: %w", err)
	}
	return buf.Bytes(), nil
}

func drawRightAligned(pdf *fpdf.Fpdf, rightX, y, size float64, text string) {
	pdf.SetFontSize(size)
	pdf.Text(rightX-pdf.GetStringWidth(text), y, text)
}

// drawPurpose renders the "на <тип> / <значение>" block — matching the
// reference designs' "на сумму" / "5000 рублей" style layout, one line per
// wrapped segment (label, then value if present). Either segment can wrap
// onto more than one line on its own (a long course/masterclass title, or
// the long any_masterclass label "на любой мастер-класс"), and the admin's
// value is free text of unpredictable length — so the font size shrinks
// (down to purposeMinFontSize) until everything that follows it on the page
// (recipient, QR codes, contact row) is projected to still fit above
// bottomMargin, instead of assuming a fixed line count ever fits.
// Returns the y just below the last line drawn, so the caller can position
// everything that follows relative to however much space this actually
// took.
func drawPurpose(pdf *fpdf.Fpdf, kind model.GiftCertificateKind, value string) float64 {
	boxWidth := marginRight - marginLeft
	label := kindPurposeLabel[kind]

	var lines []string
	var lineHeight float64
	for size := purposeMaxFontSize; ; size -= 4 {
		pdf.SetFontSize(size)
		lineHeight = size * (purposeLineHeight / purposeMaxFontSize)

		lines = lines[:0]
		for _, l := range pdf.SplitLines([]byte(label), boxWidth) {
			lines = append(lines, string(l))
		}
		if value != "" {
			for _, l := range pdf.SplitLines([]byte(value), boxWidth) {
				lines = append(lines, string(l))
			}
		}

		purposeBottom := purposeStartY + float64(len(lines))*lineHeight
		contactY := purposeBottom + purposeToRecipientGap + recipientToQRGap + qrToContactGap
		if contactY <= pageHeight-bottomMargin || size <= purposeMinFontSize {
			break
		}
	}

	y := purposeStartY
	for _, line := range lines {
		pdf.Text(marginLeft, y, line)
		y += lineHeight
	}
	return y
}

func drawQRCodes(pdf *fpdf.Fpdf, y float64) error {
	telegramPNG, err := qrcode.Encode(telegramURL, qrcode.Medium, 300)
	if err != nil {
		return err
	}
	websitePNG, err := qrcode.Encode(websiteURL, qrcode.Medium, 300)
	if err != nil {
		return err
	}

	pdf.RegisterImageOptionsReader("qr-telegram", fpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(telegramPNG))
	pdf.RegisterImageOptionsReader("qr-website", fpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(websitePNG))

	const qrSize = 100.0
	pdf.ImageOptions("qr-telegram", marginLeft, y, qrSize, qrSize, false, fpdf.ImageOptions{ImageType: "PNG"}, 0, "")
	pdf.ImageOptions("qr-website", 433, y, qrSize, qrSize, false, fpdf.ImageOptions{ImageType: "PNG"}, 0, "")
	return nil
}
