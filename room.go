package main

import "log"

const (
	UNOCCUPIED  = 0
	WHITE_STONE = 1
	BLACK_STONE = 2
)

type stone int

type Room struct {
	id    int
	board [10][10]int
}

// TODO: check who wins
// res int // 0: not end; 1: player1/white wins; 2: player2/black wins
func (r *Room) Move(role int, row, col int) int {
	r.board[row][col] = role
	log.Printf("board[%d][%d] = %d", row, col, role)

	return 0
}
