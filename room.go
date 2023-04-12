package main

import (
	"log"
)

type Room struct {
	id    int
	board [10][10]int
}

// Move updates inner board state and move and returns the result of this move
func (r *Room) Move(role int, row, col int) int {
	r.board[row][col] = role
	log.Printf("board[%d][%d] = %d", row, col, role)

	// 0000 1 0000
	// vertically
	for i := valid(row - 4); i <= row; i++ {
		conn := 0
		for x := i; x <= valid(i+4); x++ {
			if r.board[x][col] == role {
				conn++
			}
		}
		if conn == 5 {
			return role
		}
	}

	// 1111 0 1111
	// horizontally
	for j := valid(col - 4); j <= col; j++ {
		conn := 0
		for y := j; y <= valid(j+4); y++ {
			if r.board[row][y] == role {
				conn++
			}
		}
		if conn == 5 {
			return role
		}
	}

	// 1
	//  1
	//   1
	//    1
	//     0
	//      1
	//       1
	//        1
	//         1
	{
		conn := 0
		for x, y := row, col; validate(x) && validate(y); x, y = x-1, y-1 {
			if r.board[x][y] == role {
				conn++
			} else {
				break
			}
		}
		for x, y := row+1, col+1; validate(x) && validate(y); x, y = x+1, y+1 {
			if r.board[x][y] == role {
				conn++
			} else {
				break
			}
		}
		if conn == 5 {
			return role
		}
	}

	//         1
	//        1
	//       1
	//      1
	//     0
	//    1
	//   1
	//  1
	// 1
	{
		conn := 0
		for x, y := row, col; validate(x) && validate(y); x, y = x+1, y-1 {
			if r.board[x][y] == role {
				conn++
			} else {
				break
			}
		}
		for x, y := row-1, col+1; validate(x) && validate(y); x, y = x-1, y+1 {
			if r.board[x][y] == role {
				conn++
			} else {
				break
			}
		}
		if conn == 5 {
			return role
		}
	}

	return 0
}

func valid(x int) int {
	if x < 0 {
		return 0
	}
	if x >= 10 {
		return 9
	}
	return x
}

func validate(x int) bool {
	return x >= 0 && x <= 9
}
