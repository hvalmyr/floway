// Package certpdf renders the "подарочный сертификат" PDF — one fixed-design
// page (background artwork + brand fonts, both embedded as assets below)
// with a handful of dynamic fields drawn on top: certificate number, issue
// date, what it's redeemable for, and the recipient's name.
//
// All layout constants below are in "design pixels" — the 1181×1772px frame
// the certificate was designed at, sized for a 10×15cm print at 300dpi.
// Only the pt() helper (and the handful of fpdf calls that use it) convert
// into the PDF's real physical points, so the page comes out an actual
// 10×15cm sheet instead of a same-pixel-count-but-huge page — everything
// else in this file reasons in the original design-pixel space.
//
// Vertical positions are computed by a sequential flow (see the flow type)
// using each line's real font metrics (ascent/descent from the font's own
// descriptor) rather than guessed line-heights — this is what keeps the gap
// before the QR row visually identical across certificates even though the
// title, purpose and recipient lines each render at a different font size
// depending on their own text.
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

//go:embed assets/NonBureau-Medium.ttf
var nonBureauMediumTTF []byte

const (
	// dpi is the resolution the 1181×1772 design frame was authored at —
	// 1181px/300dpi = 10cm and 1772px/300dpi ≈ 15cm exactly (a standard
	// 10×15cm photo-print size), which is what fixes pxToPt below.
	dpi    = 300.0
	pxToPt = 72.0 / dpi

	pageWidth  = 1181.0
	pageHeight = 1772.0

	// #82B1CC / #41342A — the certificate's two design colors (background
	// base and text/QR ink). The background artwork PNG already bakes
	// #82B1CC blended with the texture image at "soft light", so this solid
	// fill is only a fallback for any edge the image doesn't quite cover.
	bgR, bgG, bgB    = 0x82, 0xB1, 0xCC
	inkR, inkG, inkB = 0x41, 0x34, 0x2A
	// #B7D4E4 is the page's actual *visible* color once the texture is
	// blended in (lighter than the raw #82B1CC base) — used as the QR
	// codes' flat background so they blend into the page instead of
	// standing out as a visibly different, more saturated square.
	qrBgR, qrBgG, qrBgB = 0xB7, 0xD4, 0xE4

	marginLeft  = 100.0
	marginRight = 1080.0 // right edge text aligns against
	textWidth   = marginRight - marginLeft

	// Display font for the heading, title and purpose block.
	fontDisplay = "SoyuzGrotesk"
	// Body font for the recipient line, certificate number/date, and the
	// contact row (links + phone).
	fontBody = "NonBureau"

	// Vertical rhythm, all as gaps between one element's visual bottom
	// (its descent) and the next element's visual top (its ascent) — not
	// baseline-to-baseline distances, since those would drift whenever two
	// elements end up at different font sizes.
	topMargin          = 100.0 // page top -> "школа фловей"
	headingToTitleGap  = 50.0  // "школа фловей" -> "подарочный"
	titleLineGap       = 20.0  // "подарочный" -> "сертификат"
	titleToNumberGap   = 100.0 // "сертификат" -> number/date row
	numberToPurposeGap = 100.0
	purposeToRecipient = 24.0
	recipientToQRGap   = 100.0 // fixed regardless of content — see drawPurpose's shrink loop
	qrToLinkGap        = 24.0
	bottomMargin       = 100.0 // contact row's ink bottom -> page bottom edge, fixed — see the bottom-anchored cluster in Render
	headingFontSize    = 60.0
	numberDateFontSize = 44.0
	purposeMaxFontSize = 104.0
	purposeMinFontSize = 60.0
	columnGutter       = 100.0
	linkFontSizeInset  = 0.95 // fit the longest contact string to 95% of its column, for breathing room

	telegramContact = "t.me/floway_school"
	websiteContact  = "kursfloristiki.ru"
	phoneContact    = "+7 985 226 1948"
	telegramURL     = "https://t.me/floway_school"
	websiteURL      = "https://kursfloristiki.ru"
)

// pt converts a design-pixel measurement into the PDF's real points.
func pt(px float64) float64 { return px * pxToPt }

// kindPurposeLabel is the fixed lead-in text for each certificate kind — see
// the reference designs: "на сумму" / "5000 рублей", "на любой мастер-класс"
// (no value), "на курс" / "“<title>”", "на мастер-класс" / "“<title>”".
var kindPurposeLabel = map[model.GiftCertificateKind]string{
	model.GiftCertificateKindAmount:         "на сумму",
	model.GiftCertificateKindCourse:         "на курс",
	model.GiftCertificateKindMasterclass:    "на мастер-класс",
	model.GiftCertificateKindAnyMasterclass: "на любой мастер-класс",
}

// quotedKinds get their free-text value wrapped in “...” curly quotes,
// matching the reference designs' 'на курс "основы флористики"' style.
// amount's value ("5000 рублей") and any_masterclass (no value) never quote.
var quotedKinds = map[model.GiftCertificateKind]bool{
	model.GiftCertificateKindCourse:      true,
	model.GiftCertificateKindMasterclass: true,
}

// flow lays out elements top-to-bottom using real font metrics: advance()
// adds blank space, place() reserves room for one line of text at the given
// font/size (ascent above the cursor, descent below it) and returns the
// baseline y to draw at.
type flow struct {
	pdf *fpdf.Fpdf
	y   float64
}

func (fl *flow) advance(gap float64) { fl.y += gap }

func (fl *flow) place(family string, size float64) float64 {
	ascent, descent := fontMetrics(fl.pdf, family, size)
	baseline := fl.y + ascent
	fl.y = baseline + descent
	return baseline
}

// fontMetrics returns a font's ascent/descent (both positive, in design
// pixels) at the given size, read from the font's own descriptor — this is
// what lets gaps between elements stay visually constant regardless of
// which font size an auto-fitted line ended up at.
func fontMetrics(pdf *fpdf.Fpdf, family string, size float64) (ascent, descent float64) {
	desc := pdf.GetFontDesc(family, "")
	return float64(desc.Ascent) / 1000 * size, float64(-desc.Descent) / 1000 * size
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
	pdf.AddUTF8FontFromBytes(fontDisplay, "", soyuzGroteskBoldTTF)
	pdf.AddUTF8FontFromBytes(fontBody, "", nonBureauMediumTTF)
	pdf.SetFont(fontDisplay, "", 12)
	pdf.AddPage()

	pdf.SetFillColor(bgR, bgG, bgB)
	pdf.Rect(0, 0, pt(pageWidth), pt(pageHeight), "F")

	pdf.RegisterImageOptionsReader("background", fpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(backgroundPNG))
	pdf.ImageOptions("background", 0, 0, pt(pageWidth), pt(pageHeight), false, fpdf.ImageOptions{ImageType: "PNG"}, 0, "")

	pdf.SetTextColor(inkR, inkG, inkB)

	fl := &flow{pdf: pdf}

	fl.advance(topMargin)
	headingBaseline := fl.place(fontDisplay, headingFontSize)
	drawRightAligned(pdf, fontDisplay, marginRight, headingBaseline, headingFontSize, "школа фловей")

	// "подарочный" and "сертификат" each scale to fill the full text
	// width (margin to margin) at their own font size, not a shared fixed
	// size — matching the reference design's edge-to-edge title treatment.
	fl.advance(headingToTitleGap)
	title1Size := fitFontSizeToWidth(pdf, fontDisplay, "подарочный", textWidth)
	title1Baseline := fl.place(fontDisplay, title1Size)
	drawText(pdf, fontDisplay, marginLeft, title1Baseline, title1Size, "подарочный")

	fl.advance(titleLineGap)
	title2Size := fitFontSizeToWidth(pdf, fontDisplay, "сертификат", textWidth)
	title2Baseline := fl.place(fontDisplay, title2Size)
	drawText(pdf, fontDisplay, marginLeft, title2Baseline, title2Size, "сертификат")

	fl.advance(titleToNumberGap)
	numberDateBaseline := fl.place(fontBody, numberDateFontSize)
	drawText(pdf, fontBody, marginLeft, numberDateBaseline, numberDateFontSize, cert.Number)
	drawRightAligned(pdf, fontBody, marginRight, numberDateBaseline, numberDateFontSize, dateLabel)

	// The recipient/QR/contact-row cluster is anchored to the *bottom* of
	// the page, working upward at fixed gaps from a link-row baseline that
	// sits exactly bottomMargin above the page's bottom edge — not stacked
	// downward from the purpose block. That way this cluster's position
	// never depends on how tall the purpose block above it ends up: the
	// purpose block is the one thing that shrinks (see drawPurpose) to fit
	// whatever space is left above it.
	recipientSize := fitFontSizeToWidth(pdf, fontBody, cert.Recipient, textWidth)
	rAscent, rDescent := fontMetrics(pdf, fontBody, recipientSize)

	colWidth := (textWidth - 2*columnGutter) / 3
	linkSize := fitFontSizeToWidth(pdf, fontBody, telegramContact, colWidth*linkFontSizeInset)
	lAscent, lDescent := fontMetrics(pdf, fontBody, linkSize)

	linkBaseline := pageHeight - bottomMargin - lDescent
	qrTop := linkBaseline - lAscent - qrToLinkGap - colWidth
	recipientBaseline := qrTop - recipientToQRGap - rDescent
	recipientTop := recipientBaseline - rAscent

	fl.advance(numberToPurposeGap)
	purposeValue := purposeValueText(cert.Kind, cert.Value)
	drawPurpose(fl, pdf, cert.Kind, purposeValue, recipientTop-purposeToRecipient)

	drawText(pdf, fontBody, marginLeft, recipientBaseline, recipientSize, cert.Recipient)

	col1X := marginLeft
	col2X := marginLeft + colWidth + columnGutter
	if err := drawQRCodes(pdf, col1X, col2X, qrTop, colWidth); err != nil {
		return nil, fmt.Errorf("render qr codes: %w", err)
	}

	drawText(pdf, fontBody, col1X, linkBaseline, linkSize, telegramContact)
	drawText(pdf, fontBody, col2X, linkBaseline, linkSize, websiteContact)
	drawRightAligned(pdf, fontBody, marginRight, linkBaseline, linkSize, phoneContact)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("render certificate pdf: %w", err)
	}
	return buf.Bytes(), nil
}

// purposeValueText applies the reference designs' curly-quote convention
// for kinds that name a specific course/masterclass — amount's value and
// any_masterclass (which has none) are never quoted.
func purposeValueText(kind model.GiftCertificateKind, value string) string {
	if value == "" || !quotedKinds[kind] {
		return value
	}
	return "“" + value + "”"
}

func drawText(pdf *fpdf.Fpdf, family string, x, baseline, size float64, text string) {
	pdf.SetFont(family, "", pt(size))
	pdf.Text(pt(x), pt(baseline), text)
}

func drawRightAligned(pdf *fpdf.Fpdf, family string, rightX, baseline, size float64, text string) {
	pdf.SetFont(family, "", pt(size))
	pdf.Text(pt(rightX)-pdf.GetStringWidth(text), pt(baseline), text)
}

// fitFontSizeToWidth returns the design-pixel font size at which text's
// rendered width equals targetWidth (also in design pixels). GetStringWidth
// scales linearly with font size for a fixed string, so one probe size is
// enough to solve for it directly — no search loop needed.
func fitFontSizeToWidth(pdf *fpdf.Fpdf, family, text string, targetWidth float64) float64 {
	const probe = 100.0
	pdf.SetFont(family, "", pt(probe))
	widthAtProbe := pdf.GetStringWidth(text) / pxToPt
	return probe * targetWidth / widthAtProbe
}

// drawPurpose renders the "на <тип> / “<значение>”" block — one line per
// wrapped segment (label, then value if present), advancing fl past it.
// Either segment can wrap onto more than one line on its own (a long
// course/masterclass title, or the long any_masterclass label "на любой
// мастер-класс"), and the admin's value is free text of unpredictable
// length — so the font size shrinks (down to purposeMinFontSize) until the
// block's bottom is projected to land at or above maxBottom, instead of
// assuming a fixed line count ever fits. maxBottom is the recipient's fixed
// (bottom-anchored) top edge minus purposeToRecipient — everything below
// the purpose block has an unchanging position regardless of certificate
// content, so this is the one element that absorbs the variability.
func drawPurpose(fl *flow, pdf *fpdf.Fpdf, kind model.GiftCertificateKind, value string, maxBottom float64) {
	label := kindPurposeLabel[kind]

	var lines []string
	var ascent, descent float64
	for size := purposeMaxFontSize; ; size -= 4 {
		pdf.SetFont(fontDisplay, "", pt(size))
		ascent, descent = fontMetrics(pdf, fontDisplay, size)

		lines = lines[:0]
		for _, l := range pdf.SplitLines([]byte(label), pt(textWidth)) {
			lines = append(lines, string(l))
		}
		if value != "" {
			for _, l := range pdf.SplitLines([]byte(value), pt(textWidth)) {
				lines = append(lines, string(l))
			}
		}

		lineHeight := ascent + descent
		purposeHeight := float64(len(lines)) * lineHeight
		if fl.y+purposeHeight <= maxBottom || size <= purposeMinFontSize {
			break
		}
	}

	for _, line := range lines {
		baseline := fl.y + ascent
		pdf.Text(pt(marginLeft), pt(baseline), line)
		fl.y = baseline + descent
	}
}

func drawQRCodes(pdf *fpdf.Fpdf, col1X, col2X, y, size float64) error {
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

	pdf.ImageOptions("qr-telegram", pt(col1X), pt(y), pt(size), pt(size), false, fpdf.ImageOptions{ImageType: "PNG"}, 0, "")
	pdf.ImageOptions("qr-website", pt(col2X), pt(y), pt(size), pt(size), false, fpdf.ImageOptions{ImageType: "PNG"}, 0, "")
	return nil
}

// qrPNG renders a QR code in the certificate's own ink color on the page's
// actual visible background color (qrBg) — not white (no box), but also not
// transparent: a fully transparent code lets the page's mottled texture
// show through its light modules, which measurably hurt scan reliability in
// testing (a scanner failed on codes sitting over a busier patch of the
// artwork). A flat fill matching what the page actually looks like keeps
// the "no box" look while giving every code the same reliable contrast
// regardless of where it lands on the texture.
func qrPNG(url string) ([]byte, error) {
	qr, err := qrcode.New(url, qrcode.High)
	if err != nil {
		return nil, err
	}
	qr.BackgroundColor = color.RGBA{R: qrBgR, G: qrBgG, B: qrBgB, A: 0xFF}
	qr.ForegroundColor = color.RGBA{R: inkR, G: inkG, B: inkB, A: 0xFF}
	qr.DisableBorder = true
	// Fixed source resolution regardless of the column's design-pixel
	// size — plenty crisp once scaled down onto the page, and simpler than
	// re-deriving a size from the (already scaling) column width.
	return qr.PNG(600)
}
