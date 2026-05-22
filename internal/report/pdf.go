package report

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/jung-kurt/gofpdf"

	"gitlab.com/Uranury/tunescape/internal/leaderboard"
	"gitlab.com/Uranury/tunescape/internal/track"
)

func downloadImageToTemp(url string) (string, error) {
	if url == "" {
		return "", nil
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	
	tmpFile, err := os.CreateTemp("", "track_*.jpg")
	if err != nil {
		return "", err
	}
	defer func() {
		_ = tmpFile.Close()
	}()
	
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", err
	}
	
	return tmpFile.Name(), nil
}

func buildPDF(name string, tracks []track.Track, rankings *leaderboard.UserRankings) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)

	pdf.AddPage()
	pdf.SetFillColor(10, 10, 10)
	pdf.Rect(0, 0, 210, 297, "F")
	
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

	pdf.AddPage()
	pdf.SetFillColor(30, 58, 138)
	pdf.Rect(0, 0, 210, 297, "F")
	
	pdf.SetFont("Helvetica", "B", 22)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(0, 12, "Your Top Tracks", "", 1, "L", false, 0, "")
	pdf.Ln(10)
	
	for i := 0; i < 5 && i < len(tracks); i++ {
		t := tracks[i]
		
		yPos := pdf.GetY()
		
		if t.ImageURL != nil && *t.ImageURL != "" {
			imgPath, err := downloadImageToTemp(*t.ImageURL)
			if err == nil && imgPath != "" {
				_ = pdf.RegisterImage(imgPath, "")
				pdf.Image(imgPath, 15, yPos, 30, 30, false, "", 0, "")
				_ = os.Remove(imgPath)
			}
		}
		
		pdf.SetXY(55, yPos+5)
		pdf.SetFont("Helvetica", "B", 24)
		pdf.SetTextColor(255, 255, 255)
		pdf.CellFormat(0, 12, fmt.Sprintf("#%d", i+1), "", 1, "L", false, 0, "")
		
		pdf.SetXY(55, yPos+20)
		pdf.SetFont("Helvetica", "B", 14)
		pdf.CellFormat(0, 10, t.Name, "", 1, "L", false, 0, "")
		
		pdf.SetXY(55, yPos+32)
		pdf.SetFont("Helvetica", "", 10)
		pdf.SetTextColor(220, 220, 220)
		pdf.CellFormat(0, 8, fmt.Sprintf("Popularity: %d%%", t.Popularity), "", 1, "L", false, 0, "")
		
		pdf.SetY(yPos + 45)
	}

	pdf.AddPage()
	pdf.SetFillColor(80, 0, 120)
	pdf.Rect(0, 0, 210, 297, "F")
	
	pdf.SetFont("Helvetica", "B", 22)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(0, 12, "Your Music Profile", "", 1, "L", false, 0, "")
	pdf.Ln(10)
	
	rankItems := []struct {
		label string
		value *int64
	}{
		{"Energy", rankings.Energy},
		{"Danceability", rankings.Danceability},
		{"Valence", rankings.Valence},
	}
	
	for _, item := range rankItems {
		pdf.SetFont("Helvetica", "", 14)
		pdf.SetTextColor(220, 220, 220)
		pdf.CellFormat(0, 8, item.label, "", 1, "L", false, 0, "")
		
		pdf.SetFont("Helvetica", "B", 34)
		if item.value != nil {
			pdf.SetTextColor(255, 255, 255)
			pdf.CellFormat(0, 14, fmt.Sprintf("#%d", *item.value), "", 1, "L", false, 0, "")
		} else {
			pdf.CellFormat(0, 14, "N/A", "", 1, "L", false, 0, "")
		}
		pdf.Ln(10)
	}

	pdf.AddPage()
	pdf.SetFillColor(20, 20, 20)
	pdf.Rect(0, 0, 210, 297, "F")
	
	pdf.SetFont("Helvetica", "B", 22)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(0, 12, "Your Stats", "", 1, "L", false, 0, "")
	pdf.Ln(20)
	
	total := len(tracks)
	avg := 0
	for _, t := range tracks {
		avg += t.Popularity
	}
	if total > 0 {
		avg /= total
	}
	
	pdf.SetFont("Helvetica", "", 14)
	pdf.SetTextColor(220, 220, 220)
	pdf.CellFormat(0, 8, "Tracks Analyzed", "", 1, "L", false, 0, "")
	pdf.SetFont("Helvetica", "B", 40)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(0, 16, fmt.Sprintf("%d", total), "", 1, "L", false, 0, "")
	pdf.Ln(20)
	
	pdf.SetFont("Helvetica", "", 14)
	pdf.SetTextColor(220, 220, 220)
	pdf.CellFormat(0, 8, "Avg Popularity", "", 1, "L", false, 0, "")
	pdf.SetFont("Helvetica", "B", 40)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(0, 16, fmt.Sprintf("%d%%", avg), "", 1, "L", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}