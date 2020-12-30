package main

import (
	"github.com/jung-kurt/gofpdf"
	"fmt"
	"strings"
	"time"
	"os/exec"
)

func main () {
	nx := 3
	ny := 2
	n := nx * ny

	fmt.Printf("Generating %d Sudokus in a %d x %d grid\n", n, nx, ny)
	sudokus := generateSudokus(n)
	timestamp := time.Now().Format("20060102-150405")

	filename := fmt.Sprintf("sudokus/sudokus-%v-%dx%d-%s.pdf", timestamp, nx, ny, "intermediate")
	createPDF(sudokus, timestamp, nx, ny, filename)
}

func generateSudokus(amount int) []string {
	out, err := exec.Command("sh", "-c", fmt.Sprintf("qqwing --generate %d --one-line --difficulty intermediate", amount)).Output()
	if err != nil {
		fmt.Printf("Error: %v", err)
		panic(err)
		return nil
	}
	fmt.Println(string(out))
	return strings.Split(string(out), "\n")
}

func smaller(a, b float64) float64 {
	if (a < b) {
		return a
	}
	return b
}

func createPDF(sudokus []string, timestamp string, nx, ny int, filename string) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetMargins(0, 0, 0)
	pdf.SetAutoPageBreak(false, 0)
	pdf.SetDrawColor(0, 0, 0)

	width, height := pdf.GetPageSize()

	margin := 4. //4 mm

	drawingWidth := width - 4*margin // 2 for the heading + 1 left + 1 right
	drawingHeight := height - 2*margin // 1 top, 1 bottom

	offsetY := (height-drawingHeight)/2

	L := smaller(drawingWidth/float64(nx), drawingHeight/float64(ny))*0.85 //small sudoku length
	fieldL := L/9

	thinLineWidth := L/300
	thickLineWidth := L/120

	pdf.SetFont("Helvetica", "", 12)

	//draw title
	pdf.MoveTo(0, height+margin)
	pdf.TransformBegin()
	pdf.TransformRotate(90, 0, height)
	title := filename
	pdf.CellFormat(height, 3*margin, title, "", 0, "MC", false, 0, "")
	pdf.TransformEnd()

	pdf.SetFont("Helvetica", "", fieldL*0.8*2.83) //2.83 points is a mm

	for X := 0; X < nx; X++ {
		for Y := 0; Y < ny; Y++ {

			x0 := 3*margin + float64(X)/float64(nx)*drawingWidth + (drawingWidth/float64(nx)-L)/2
			y0 :=  offsetY + float64(Y)/float64(ny)*drawingHeight + (drawingHeight/float64(ny)-L)/2

			// draw horizontal lines
			for ly := 0; ly < 10; ly++ {
				var w float64
				if ly % 3 == 0 {
					w = thickLineWidth
				} else {
					w = thinLineWidth
				}
				pdf.SetLineWidth(w)
				pdf.Line(x0-w/2, y0 + fieldL*float64(ly), x0+w/2 + L, y0 + fieldL*float64(ly))
			}
			// draw vertical lines
			for lx := 0; lx < 10; lx++ {
				var w float64
				if lx % 3 == 0 {
					w = thickLineWidth
				} else {
					w = thinLineWidth
				}
				pdf.SetLineWidth(w)
				pdf.Line(x0 + fieldL*float64(lx), y0-w/2, x0 + fieldL*float64(lx), y0+w/2 + L)
			}
			// draw numbers
			for i := 0; i < 9; i++ {
				for j := 0; j < 9; j++ {
					n := sudokus[Y*3+X][i*9+j]
					if string(n) != "." {
						dy := fieldL/20
						pdf.MoveTo(x0 + fieldL*float64(i), y0 + fieldL*float64(j) + dy)

						//parameters for drawing the number: cell w, h, number, no borders,
						//don't move, center verically & horizontally, no fill, no link x2
						pdf.CellFormat(fieldL, fieldL, string(n), "", 0, "CM", false, 0, "")
					}
				}
			}

		}
	}

	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Wrote sudokus to file %s successfully.\n", filename)
	}
}
