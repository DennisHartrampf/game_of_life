package main

import "fmt"

func main() {
	game := CreateGame(20, 20)
	for !game.finished {
		game.MakeMove()
		fmt.Println(game)
	}
}

type Board struct {
	w, h  int
	board [][]bool
}

type Game struct {
	board      Board
	generation int
	finished   bool
}

func (g *Game) String() string {
	return fmt.Sprintf("Generation: %v, Finished: %v", g.generation, g.finished)
}

func CreateGame(w, h int) *Game {
	game := Game{board: Board{w, h, make([][]bool, h)}}
	for i := 0; i < h; i++ {
		game.board.board[i] = make([]bool, w)
	}
	return &game
}

func (g *Game) MakeMove() {
	newBoard := Board{
		g.board.w,
		g.board.h,
		make([][]bool, g.board.w)}
	g.board = newBoard
	g.generation++
	if (g.generation > 10) {
		g.finished = true
	}
}