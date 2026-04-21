package report

import (
	"bytes"
	"fmt"

	"github.com/jung-kurt/gofpdf"

	"gitlab.com/Uranury/tunescape/internal/leaderboard"
	"gitlab.com/Uranury/tunescape/internal/track"
)

func buildPDF(name string, tracks []track.Track, rankings *leaderboard.UserRankings) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)

	drawCoverPage(pdf, name)

	pdf.AddPage()
	drawGradientBackground(pdf, [3]int{30, 58, 138}, [3]int{34, 197, 94})
	drawSectionHeader(pdf, "Your Top Tracks")
	drawTracksWrapped(pdf, tracks)

	pdf.AddPage()
	drawGradientBackground(pdf, [3]int{80, 0, 120}, [3]int{34, 197, 94})
	drawSectionHeader(pdf, "Your Music Profile")
	drawRankingsWrapped(pdf, rankings)

	pdf.AddPage()
	drawGradientBackground(pdf, [3]int{20, 20, 20}, [3]int{34, 197, 94})
	drawSectionHeader(pdf, "Your Stats")
	drawSummaryWrapped(pdf, tracks)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func drawGradientBackground(pdf *gofpdf.Fpdf, top [3]int, bottom [3]int) {
	for i := 0; i < 297; i++ {
		r := top[0] + (bottom[0]-top[0])*i/297
		g := top[1] + (bottom[1]-top[1])*i/297
		b := top[2] + (bottom[2]-top[2])*i/297

		pdf.SetFillColor(r, g, b)
		pdf.Rect(0, float64(i), 210, 1, "F")
	}
}

func drawCoverPage(pdf *gofpdf.Fpdf, name string) {
	pdf.AddPage()

	drawGradientBackground(pdf, [3]int{10, 10, 10}, [3]int{34, 197, 94})

	pdf.SetY(100)

	pdf.SetFont("Helvetica", "B", 32)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(0, 12, "Tunescape", "", 1, "C", false, 0, "")

	pdf.SetFont("Helvetica", "", 16)
	pdf.SetTextColor(200, 200, 200)
	pdf.CellFormat(0, 10, "Your Music Journey", "", 1, "C", false, 0, "")

	pdf.Ln(20)

	pdf.SetFont("Helvetica", "B", 28)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(0, 12, name, "", 1, "C", false, 0, "")
}

func drawSectionHeader(pdf *gofpdf.Fpdf, title string) {
	pdf.SetFont("Helvetica", "B", 22)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(0, 12, title, "", 1, "L", false, 0, "")

	pdf.SetDrawColor(255, 255, 255)
	pdf.SetLineWidth(1)
	pdf.Line(15, pdf.GetY(), 100, pdf.GetY())

	pdf.Ln(10)
}

func drawTracksWrapped(pdf *gofpdf.Fpdf, tracks []track.Track) {
	for i, t := range tracks {
		if i >= 5 {
			break
		}

		pdf.SetFont("Helvetica", "B", 48)
		pdf.SetTextColor(255, 255, 255)
		pdf.CellFormat(0, 20, fmt.Sprintf("#%d", i+1), "", 1, "L", false, 0, "")

		pdf.SetFont("Helvetica", "B", 18)
		pdf.CellFormat(0, 10, t.Name, "", 1, "L", false, 0, "")

		pdf.SetFont("Helvetica", "", 12)
		pdf.SetTextColor(220, 220, 220)
		pdf.CellFormat(0, 8, fmt.Sprintf("Popularity: %d%%", t.Popularity), "", 1, "L", false, 0, "")

		pdf.Ln(8)
	}
}

func drawRankingsWrapped(pdf *gofpdf.Fpdf, rankings *leaderboard.UserRankings) {
	type item struct {
		label string
		value *int64
	}

	items := []item{
		{"Energy", rankings.Energy},
		{"Danceability", rankings.Danceability},
		{"Valence", rankings.Valence},
	}

	for _, it := range items {
		pdf.SetFont("Helvetica", "", 14)
		pdf.SetTextColor(220, 220, 220)
		pdf.CellFormat(0, 8, it.label, "", 1, "L", false, 0, "")

		pdf.SetFont("Helvetica", "B", 34)
		if it.value != nil {
			pdf.SetTextColor(255, 255, 255)
			pdf.CellFormat(0, 14, fmt.Sprintf("#%d", *it.value), "", 1, "L", false, 0, "")
		} else {
			pdf.CellFormat(0, 14, "N/A", "", 1, "L", false, 0, "")
		}

		pdf.Ln(6)
	}
}

func drawSummaryWrapped(pdf *gofpdf.Fpdf, tracks []track.Track) {
	total := len(tracks)
	avg := 0

	for _, t := range tracks {
		avg += t.Popularity
	}
	if total > 0 {
		avg /= total
	}

	drawBigStat(pdf, "Tracks Analyzed", fmt.Sprintf("%d", total))
	drawBigStat(pdf, "Avg Popularity", fmt.Sprintf("%d%%", avg))
}

func drawBigStat(pdf *gofpdf.Fpdf, label, value string) {
	pdf.SetFont("Helvetica", "", 14)
	pdf.SetTextColor(220, 220, 220)
	pdf.CellFormat(0, 8, label, "", 1, "L", false, 0, "")

	pdf.SetFont("Helvetica", "B", 40)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(0, 16, value, "", 1, "L", false, 0, "")

	pdf.Ln(10)
}
