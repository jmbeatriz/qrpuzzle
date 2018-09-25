package puzzle

import (
	"image"
	"image/color"
)

type Piece struct {
	image    image.Image
	initX    int
	initY    int
	endX     int
	endY     int
	used     bool
	position int
}

func NewPiece(image image.Image, initX int, initY int, endX int, endY int, pos int) Piece {
	piece := Piece{}
	piece.image = image
	piece.initX = initX
	piece.initY = initY
	piece.endX = endX
	piece.endY = endY
	piece.used = false
	piece.position = pos
	return piece
}

func (piece *Piece) Used() bool{
	return piece.used
}

func (piece *Piece) SetUsed(used bool) {
	piece.used = used
}

func (piece *Piece) GetColorUpRight() color.Color {
	return piece.image.At(piece.initX + 45, piece.initY + 5)
}

func (piece *Piece) GetColorUpLeft() color.Color {
	return piece.image.At(piece.initX + 5, piece.initY + 5)
}

func (piece *Piece) GetColorDownRight() color.Color {
	return piece.image.At(piece.initX + 45, piece.initY + 45)
}

func (piece *Piece) GetColorDownLeft() color.Color {
	return piece.image.At(piece.initX + 5, piece.initY + 45)
}

func (piece *Piece) GetColorCenter() color.Color {
	return piece.image.At(25, 25)
}

func (piece *Piece) Size() (width int, height int) {
	return piece.endX - piece.initX, piece.endY - piece.initY
}