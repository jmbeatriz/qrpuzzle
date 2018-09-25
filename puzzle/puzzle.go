package puzzle

import (
	"fmt"
	"image"
	"image/draw"
)

type Puzzle struct {
	board      [][]Piece
	pieces     []Piece
	findPiece  func(puzzle *Puzzle, row int, column int, initPos int) (bool, int)
	iterations int
}

func NewPuzzle(rows int, columns int) Puzzle {
	puzzle := Puzzle{}
	board := make([][]Piece, columns)
	for idx := range board {
		board[idx] = make([]Piece, rows)
	}
	puzzle.board = board
	puzzle.iterations = 0
	return puzzle
}

func (puzzle *Puzzle) Board() *[][]Piece {
	return &puzzle.board
}

func (puzzle *Puzzle) BoardCell(column int, row int) *Piece {
	return &puzzle.board[column][row]
}

func (puzzle *Puzzle) PlacePiece(column int, row int, piece Piece) {
	puzzle.board[column][row] = piece
}


func (puzzle *Puzzle) Pieces() []Piece {
	return puzzle.pieces
}

func (puzzle *Puzzle) Piece(pos int) *Piece {
	return &puzzle.pieces[pos]
}

func (puzzle *Puzzle) SetPieces(pieces []Piece) {
	puzzle.pieces = pieces
}

func (puzzle *Puzzle) PiecesCount() int {
	return len(puzzle.pieces)
}

func (puzzle *Puzzle) SetFindPiece(function func(puzzle *Puzzle, row int, column int, initPos int) (bool, int)) {
	puzzle.findPiece = function
}

func (puzzle *Puzzle) ColumnsSize() int {
	return len(puzzle.board)
}

func (puzzle *Puzzle) RowsSize() int {
	return len(puzzle.board[1])
}

func (puzzle *Puzzle) IncrementIterations() {
	puzzle.iterations++
}

func (puzzle *Puzzle) ResetIterations() {
	puzzle.iterations = 0
}

func (puzzle *Puzzle) PieceSize() (width int, height int) {
	return puzzle.board[0][0].Size()
}

func (puzzle *Puzzle) PopulatePuzzle() {
	board := puzzle.board
	pieces := puzzle.pieces

	puzzle.ResetIterations()
	lastFoundPos := -1
	initPos := 0
	for row := 0; row < puzzle.RowsSize() && row >= 0; row++ {
		fmt.Println("Starting new row")
		for column := 0; column < puzzle.ColumnsSize() && column >= 0; column++ {

			fmt.Println("---------")
			fmt.Printf("Starting row:%d column:%d \n", row, column)
			fmt.Printf("Count used %d \n", puzzle.countUsedPieces())
			fmt.Printf("Start in Position %d \n", initPos)

			found, foundPos := puzzle.findPiece(puzzle, row, column, initPos)

			if found == true {
				fmt.Printf("Found row:%d column:%d in position: %d. Used %v \n", row, column, foundPos, pieces[foundPos].used)
				lastFoundPos = foundPos
				initPos = 0
			} else {
				fmt.Printf("Not Found: row:%d column:%d \n", row, column)
				found = true
				var previousRow int
				var previousColumn int

				if column == 0 {
					previousColumn = puzzle.ColumnsSize() - 1
					previousRow = row - 1
					lastFoundPos = board[previousColumn][previousRow].position
					// Go to last column of previous row
					row = row - 1
					column = puzzle.ColumnsSize() - 2
					fmt.Println("Going back to previous row")
				} else {
					previousColumn = column - 1
					previousRow = row
					lastFoundPos = board[previousColumn][previousRow].position
					// Go to previous column
					column = column - 2
				}

				fmt.Printf("Previos found row:%d column:%d with state '%v'. The piece position is %d \n", previousRow, previousColumn, pieces[lastFoundPos].used, lastFoundPos)
				puzzle.pieces[lastFoundPos].used = false
				initPos = lastFoundPos + 1
				fmt.Printf("The next iteration should start on position: %d.\n", initPos)

			}
			fmt.Printf("Count used %d \n", puzzle.countUsedPieces())

		}
	}

	fmt.Println("---------")
	fmt.Printf("Count used %d \n", puzzle.countUsedPieces())

	fmt.Printf("Number of pieces I tried to match %d \n", puzzle.iterations)
}

func (puzzle *Puzzle) ImageOnBoard() *image.Image {
	pieceWidth, pieceHeight := puzzle.PieceSize()
	imageWith := pieceWidth * puzzle.ColumnsSize()
	imageHeight := pieceHeight * puzzle.RowsSize()

	var myImage image.Image = image.NewRGBA(image.Rect(0, 0, imageWith, imageHeight))
	for row := 0; row < puzzle.RowsSize(); row++ {
		for column := 0; column < puzzle.ColumnsSize(); column++ {
			if puzzle.board[column][row].image != nil {

				desiredPosY := pieceWidth * row
				newPosY := puzzle.board[column][row].image.Bounds().Min.Y - desiredPosY
				desiredPosX := pieceHeight * column
				newPosX := puzzle.board[column][row].image.Bounds().Min.X - desiredPosX
				draw.Draw(myImage.(draw.Image), myImage.Bounds(), puzzle.board[column][row].image, image.Point{X: newPosX, Y: newPosY}, draw.Src)
			}
		}
	}

	return &myImage
}

func (puzzle *Puzzle) countUsedPieces() int {
	count := 0
	for i := range puzzle.pieces {
		if puzzle.pieces[i].used {
			count++
		}
	}
	return count
}
