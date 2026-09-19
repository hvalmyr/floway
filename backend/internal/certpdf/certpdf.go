// Package certpdf renders the "подарочный сертификат" PDF — one fixed-design
// page (background artwork + brand font, both embedded as assets below)
// with a handful of dynamic fields drawn on top: certificate number, issue
// date, what it's redeemable for, and the recipient's name.
//
// All layout constants below are in "design pixels" — the 1181×1772px frame
// the certificate was designed at, sized for a 10×15cm print at 300dpi.
// Only the pt() helper (and the handful of fpdf calls that use it) convert
// into the PDF's real physical points, so the page comes out an actual
// 10×15cm sheet instead of a same-pixel-count-but-huge page — everything
// else in this file reasons in the original design-pixel space.
package certpdf

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/color"

	"github.com/go-pdf/fpdf"
	"github.com/skip2/go-qrcode"

	"floway-backend/internal/model"
)

//go:embed assets/background.png
var backgroundPNG []byte

//go:embed assets/SoyuzGrotesk-Bold.ttf
var soyuzGroteskBoldTTF []byte

const (
	// dpi is the resolution the 1181×1772 design frame was authored at —
	// 1181px/300dpi = 10cm and 1772px/300dpi ≈ 15cm exactly (a standard
	// 10×15cm photo-print size), which is what fixes pxToPt below.
	dpi    = 300.0
	pxToPt = 72.0 / dpi

	pageWidth  = 1181.0
	pageHeight = 1772.0

	// #82B1CC / #41342A — the certificate's two colors (background base and
	// text/QR ink). The background artwork PNG already bakes #82B1CC
	// blended with the texture image at "soft light", so this solid fill is
	// only a fallback for any edge the image doesn't quite cover.
	bgR, bgG, bgB    = 0x82, 0xB1, 0xCC
	inkR, inkG, inkB = 0x41, 0x34, 0x2A

	marginLeft  = 100.0
	marginRight = 1080.0 // right edge text aligns against
	fontFamily  = "SoyuzGrotesk"

	purposeStartY         = 900.0
	purposeMaxFontSize    = 104.0
	purposeMinFontSize    = 60.0
	purposeLineHeight     = 108.0 // one line at purposeMaxFontSize; scales with the chosen size
	purposeToRecipientGap = 24.0
	recipientToQRGap      = 300.0
	qrSize                = 160.0
	qrToContactGap        = 90.0
	bottomMargin          = 60.0
	telegramContact       = "t.me/floway_school"
	websiteContact        = "kursfloristiki.ru"
	phoneContact          = "+7 985 226 1948"
	telegramURL           = "https://t.me/floway_school"
	websiteURL            = "https://kursfloristiki.ru"
)

// pt converts a design-pixel measurement into the PDF's real points.
func pt(px float64) float64 { return px * pxToPt }

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
		Size:    fpdf.SizeType{Wd: pt(pageWidth), Ht: pt(pageHeight)},
	})
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddUTF8FontFromBytes(fontFamily, "", soyuzGroteskBoldTTF)
	pdf.SetFont(fontFamily, "", 12)
	pdf.AddPage()

	pdf.SetFillColor(bgR, bgG, bgB)
	pdf.Rect(0, 0, pt(pageWidth), pt(pageHeight), "F")

	pdf.RegisterImageOptionsReader("background", fpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(backgroundPNG))
	pdf.ImageOptions("background", 0, 0, pt(pageWidth), pt(pageHeight), false, fpdf.ImageOptions{ImageType: "PNG"}, 0, "")

	pdf.SetTextColor(inkR, inkG, inkB)

	drawRightAligned(pdf, marginRight, 140, 60, "школа фловей")

	// "подарочный" and "сертификат" each scale to fill the full text
	// width (margin to margin), not a fixed size — matching the reference
	// design's edge-to-edge title treatment.
	titleWidth := marginRight - marginLeft
	drawFitted(pdf, marginLeft, 300, titleWidth, "подарочный")
	drawFitted(pdf, marginLeft, 478, titleWidth, "сертификат")

	drawText(pdf, marginLeft, 600, 74, cert.Number)
	drawRightAligned(pdf, marginRight, 600, 74, dateLabel)

	// The purpose block wraps to a variable number of lines (a long course
	// or masterclass title, or a long label like "на любой мастер-класс",
	// can each take 1-2 lines). Shrinking its font size when needed keeps
	// everything below it — recipient, QR codes, contact row — on the page
	// regardless of how long the admin's free-text value turns out to be.
	purposeBottom := drawPurpose(pdf, cert.Kind, cert.Value)
	recipientY := purposeBottom + purposeToRecipientGap

	// The recipient's name always fills the full line width too.
	drawFitted(pdf, marginLeft, recipientY, titleWidth, cert.Recipient)

	qrY := recipientY + recipientToQRGap
	if err := drawQRCodes(pdf, qrY); err != nil {
		return nil, fmt.Errorf("render qr codes: %w", err)
	}

	contactY := qrY + qrSize + qrToContactGap
	drawText(pdf, marginLeft, contactY, 28, telegramContact)
	drawText(pdf, 433, contactY, 28, websiteContact)
	drawRightAligned(pdf, marginRight, contactY, 28, phoneContact)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("render certificate pdf: %w", err)
	}
	return buf.Bytes(), nil
}

func drawText(pdf *fpdf.Fpdf, x, y, size float64, text string) {
	pdf.SetFontSize(pt(size))
	pdf.Text(pt(x), pt(y), text)
}

func drawRightAligned(pdf *fpdf.Fpdf, rightX, y, size float64, text string) {
	pdf.SetFontSize(pt(size))
	pdf.Text(pt(rightX)-pdf.GetStringWidth(text), pt(y), text)
}

// drawFitted draws text at whatever font size makes it exactly width wide —
// used for the two title words and the recipient line, which the design
// always stretches edge to edge rather than showing at a fixed size.
func drawFitted(pdf *fpdf.Fpdf, x, y, width float64, text string) {
	size := fitFontSizeToWidth(pdf, text, width)
	drawText(pdf, x, y, size, text)
}

// fitFontSizeToWidth returns the design-pixel font size at which text's
// rendered width equals targetWidth (also in design pixels). GetStringWidth
// scales linearly with font size for a fixed string, so one probe size is
// enough to solve for it directly — no search loop needed.
func fitFontSizeToWidth(pdf *fpdf.Fpdf, text string, targetWidth float64) float64 {
	const probe = 100.0
	pdf.SetFontSize(pt(probe))
	widthAtProbe := pdf.GetStringWidth(text) / pxToPt
	return probe * targetWidth / widthAtProbe
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
		pdf.SetFontSize(pt(size))
		lineHeight = size * (purposeLineHeight / purposeMaxFontSize)

		lines = lines[:0]
		for _, l := range pdf.SplitLines([]byte(label), pt(boxWidth)) {
			lines = append(lines, string(l))
		}
		if value != "" {
			for _, l := range pdf.SplitLines([]byte(value), pt(boxWidth)) {
				lines = append(lines, string(l))
			}
		}

		purposeBottom := purposeStartY + float64(len(lines))*lineHeight
		contactY := purposeBottom + purposeToRecipientGap + recipientToQRGap + qrSize + qrToContactGap
		if contactY <= pageHeight-bottomMargin || size <= purposeMinFontSize {
			break
		}
	}

	y := purposeStartY
	for _, line := range lines {
		pdf.Text(pt(marginLeft), pt(y), line)
		y += lineHeight
	}
	return y
}

func drawQRCodes(pdf *fpdf.Fpdf, y float64) error {
	telegramPNG, err := qrPNG(telegramURL)
	if err != nil {
		return err
	}
	websitePNG, err := qrPNG(websiteURL)
	if err != nil {
		return err
	}

	pdf.RegisterImageOptionsReader("qr-telegram", fpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(telegramPNG))
	pdf.RegisterImageOptionsReader("qr-website", fpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(websitePNG))

	pdf.ImageOptions("qr-telegram", pt(marginLeft), pt(y), pt(qrSize), pt(qrSize), false, fpdf.ImageOptions{ImageType: "PNG"}, 0, "")
	pdf.ImageOptions("qr-website", pt(433), pt(y), pt(qrSize), pt(qrSize), false, fpdf.ImageOptions{ImageType: "PNG"}, 0, "")
	return nil
}

// qrPNG renders a QR code in the certificate's own ink color on a
// transparent background (no white box), border disabled since the
// transparency already provides the quiet zone against the page artwork.
func qrPNG(url string) ([]byte, error) {
	qr, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		return nil, err
	}
	qr.BackgroundColor = color.Transparent
	qr.ForegroundColor = color.RGBA{R: inkR, G: inkG, B: inkB, A: 0xFF}
	qr.DisableBorder = true
	// Render at 2x the design-pixel size for crispness once embedded and
	// scaled down onto the page — qrSize itself is already a 300dpi design
	// pixel measurement, no further dpi conversion needed here.
	return qr.PNG(int(qrSize * 2))
}
