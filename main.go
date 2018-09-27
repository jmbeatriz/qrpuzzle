package main

import (
	"bytes"
	"fmt"
	"github.com/jmbeatriz/qrpuzzle/puzzle"
	"github.com/jmbeatriz/qrpuzzle/colorutils"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var imageSource string = "/Users/jmbeatriz/Downloads/image_desafio.png"
var rows int = 20
var columns int = 20
var pieceSize int = 50

func main() {
	http.HandleFunc("/", puzzleHandler)
	log.Println("Listening on 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func puzzleHandler(w http.ResponseWriter, _ *http.Request) {
	start := time.Now().Nanosecond() / 1000000

	//Init Puzzle
	rawPuzzleImage := getRawImage(imageSource)
	puzzle := puzzle.NewPuzzle(rows, columns)
	puzzle.SetPieces(getPiecesFromImage(rawPuzzleImage, rows, columns, pieceSize))
	puzzle.SetFindPiece(findPiece)

	//Resolve the puzzle
	puzzle.PopulatePuzzle()
	puzzleImage := puzzle.ImageOnBoard()

	writeImage(w, puzzleImage)

	end := time.Now().Nanosecond() / 1000000
	fmt.Printf("END. Time elapsed %d millis", end-start)
}

func findPiece(puzzle *puzzle.Puzzle, row int, column int, initPos int) (bool, int) {
	if column >= 0 && column < puzzle.ColumnsSize() && row >= 0 && row <= puzzle.RowsSize() {
		for p := initPos; p < puzzle.PiecesCount(); p++ {

			if puzzle.Piece(p).Used() != true {

				puzzle.IncrementIterations()

				//First Column
				if column == 0 {
					if !colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorUpLeft()) || !colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorDownLeft()) {
						continue
					}

					// Rows from second to last in first Column
					if row > 0 {
						if !colorutils.IsSameColor(puzzle.BoardCell(column, row-1).GetColorDownRight(), puzzle.Piece(p).GetColorUpRight()) {
							continue
						}

						//Middle of first Column
						if row < 19 {
							if colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorUpRight()) || colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorDownRight()) || !colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorUpLeft()) || !colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorDownLeft()) {
								continue
							}
						}
					}
				}

				//Last Column
				if column == 19 {
					if !colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorUpRight()) || !colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorDownRight()) {
						continue
					}

					// Rows from second to last in last Column
					if row > 0 {
						if !colorutils.IsSameColor(puzzle.BoardCell(column, row-1).GetColorDownLeft(), puzzle.Piece(p).GetColorUpLeft()) {
							continue
						}

						//Middle of last Column
						if row > 0 && row < 19 {
							if colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorUpLeft()) || colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorDownLeft()) || !colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorUpRight()) || !colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorDownRight()) {
								continue
							}

						}
					}
				}

				//Fisrt Row
				if row == 0 {
					if !colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorUpRight()) || !colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorUpLeft()) {
						continue
					}

					// Columns from second to last in first Row
					if column > 0 {
						if !colorutils.IsSameColor(puzzle.BoardCell(column-1, row).GetColorDownRight(), puzzle.Piece(p).GetColorDownLeft()) {
							continue
						}

						// Middle of first Row
						if column < 19 {
							if colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorDownRight()) || colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorDownLeft()) || !colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorUpRight()) || !colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorUpLeft()) {
								continue
							}
						}
					}
				}

				//Last Row
				if row == 19 {
					if !colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorDownRight()) || !colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorDownLeft()) {
						continue
					}

					// Columns from second to last in last Row
					if column > 0 {
						if !colorutils.IsSameColor(puzzle.BoardCell(column-1, row).GetColorUpRight(), puzzle.Piece(p).GetColorUpLeft()) || !colorutils.IsSameColor(puzzle.BoardCell(column,row-1).GetColorDownLeft(), puzzle.Piece(p).GetColorUpLeft()) || !colorutils.IsSameColor(puzzle.BoardCell(column, row-1).GetColorDownRight(), puzzle.Piece(p).GetColorUpRight()) {
							continue
						}

						// Middle of last Row
						if column < 19 {
							if colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorUpRight()) || colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorUpLeft()) || !colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorDownRight()) || !colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorDownLeft()) {
								continue
							}
						}
					}
				}

				//Center of the puzzle
				if row > 0 && column > 0 && row < 19 && column < 19 {
					if colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorUpRight()) || colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorDownRight()) || colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorUpLeft()) || colorutils.IsBlackOrWhite(puzzle.Piece(p).GetColorDownLeft()) {
						continue
					}

					if !colorutils.IsSameColor(puzzle.BoardCell(column, row-1).GetColorDownRight(), puzzle.Piece(p).GetColorUpRight()) || !colorutils.IsSameColor(puzzle.BoardCell(column, row-1).GetColorDownLeft(), puzzle.Piece(p).GetColorUpLeft()) || !colorutils.IsSameColor(puzzle.BoardCell(column-1, row).GetColorUpRight(), puzzle.Piece(p).GetColorUpLeft()) || !colorutils.IsSameColor(puzzle.BoardCell(column-1, row).GetColorDownRight(), puzzle.Piece(p).GetColorDownLeft()) {
						continue
					}

				}

				puzzle.PlacePiece(column, row, *puzzle.Piece(p))

				puzzle.Piece(p).SetUsed(true)
				return true, p

			}

		}

	}

	return false, 0
}

func getRawImage(path string) image.Image {
	reader, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	return m
}

func getPiecesFromImage(picture image.Image, rows int, columns int, pieceSize int) []puzzle.Piece {
	pieces := make([]puzzle.Piece, columns*rows)
	pos := 0

	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			initX := i * pieceSize
			initY := j * pieceSize
			endX := (i + 1) * pieceSize
			endY := (j + 1) * pieceSize

			mySubImage := picture.(interface {
				SubImage(r image.Rectangle) image.Image
			}).SubImage(image.Rect(initX, initY, endX, endY))

			piece := puzzle.NewPiece(mySubImage, initX, initY, endX, endY, pos)
			pieces[pos] = piece
			pos++

		}
	}
	return pieces
}

// writeImage encodes an image 'img' in jpeg format and writes it into ResponseWriter.
func writeImage(w http.ResponseWriter, img *image.Image) {
	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, *img); err != nil {
		log.Fatalln("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}
