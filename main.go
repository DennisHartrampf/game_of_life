package main

import (
	"fmt"
	tm "github.com/buger/goterm"
	"time"
	"math/rand"
)

func main() {
	rand.Seed(time.Now().Unix())
	game := CreateGame(20, 20)
	//game.board.CreateBlinker(3, 3)
	//game.board.CreateBlock(0, 0)
	//game.board.CreateBlock(18, 0)
	//game.board.CreateBlock(0, 18)
	//game.board.CreateBlock(18, 18)
	game.board.FillRandomly()
	for !game.finished {
		game.PrintToConsole()
		game.MakeMove()
	}
	for {
		var i = 1
		i++
	}
	// temporarily unreachable
	if (len(game.board.aliveCells) == 0) {
		fmt.Printf("The game is over after %v generations because all cells died.", game.generation)
	} else {
		fmt.Printf("The game is over after %v generations because there will be no more changes.", game.generation)
	}
}

func (g Game) PrintToConsole() {
	tm.Clear()
	tm.MoveCursor(1, 1)
	tm.Print(g.String())

	// Create Box with 30% width of current screen, and height of 20 lines
	box := tm.NewBox(g.board.w + 2, g.board.h + 2, 0)
	box.Border = "# # # # # #"

	// Add some content to the box
	// Note that you can add ANY content, even tables
	fmt.Fprint(box)

	// Move Box to approx center of the screen
	tm.Print(tm.MoveTo(box.String(), 1, 3))
	for c := range(g.board.aliveCells) {
		tm.MoveCursor(c.y + 1 + 3, c.x + 1 + 1)
		tm.Print(tm.Background(tm.Color("X", tm.GREEN), tm.WHITE))
	}
	tm.Flush()
	time.Sleep(time.Millisecond * 50)
}

func (board *Board) CreateBlinker(x, y int) {
	board.aliveCells[Coordinate{x, y}] = true
	board.aliveCells[Coordinate{x + 1, y}] = true
	board.aliveCells[Coordinate{x + 2, y}] = true
}

func (board *Board) CreateBlock(x, y int) {
	board.aliveCells[Coordinate{x, y}] = true
	board.aliveCells[Coordinate{x + 1, y}] = true
	board.aliveCells[Coordinate{x, y + 1}] = true
	board.aliveCells[Coordinate{x + 1, y + 1}] = true
}

func (board *Board) FillRandomly() {
	for x := 0; x < board.w; x++ {
		for y := 0; y < board.h; y++ {
			if rand.Intn(2) == 0 {
				board.aliveCells[Coordinate{x, y}] = true
			}
		}
	}
}

func (g *Game) Revive(board *Board, c Coordinate) {
	if (g.InRange(c)) {
		board.aliveCells[c] = true
	}
}

func (g *Game) Kill(board *Board, c Coordinate) {
	if (g.InRange(c)) {
		delete(board.aliveCells, c)
	}
}

func (g Game) InRange(c Coordinate) bool {
	if (c.x < 0 || c.x >= g.board.w) {
		return false
	}
	if (c.y < 0 || c.y >= g.board.h) {
		return false
	}
	return true
}

func (g Game) IsAlive(c Coordinate) bool {
	var _, ok = g.board.aliveCells[c]
	return ok
}

type CoordinateSet map[Coordinate]bool

func (c CoordinateSet) Equals(other CoordinateSet) bool {
	if (c == nil) {
		return other == nil
	}
	if (len(c) != len(other)) {
		return false
	}
	for k, v := range (c) {
		if (other[k] != v) {
			return false
		}
	}

	return true
}

func (c Coordinate) String() string {
	return fmt.Sprintf("(%v, %v)", c.x, c.y)
}

func (c CoordinateSet) String() string {
	var s string = ""
	for k := range (c) {
		s = s + k.String() + ", "
	}
	return s
}

type Coordinate struct {
	x, y int
}

func (c Coordinate) Neighbours() CoordinateSet {
	neighbours := make(CoordinateSet)
	neighbours[Coordinate{c.x - 1, c.y - 1}] = true
	neighbours[Coordinate{c.x, c.y - 1}] = true
	neighbours[Coordinate{c.x + 1, c.y - 1}] = true
	neighbours[Coordinate{c.x - 1, c.y}] = true
	neighbours[Coordinate{c.x + 1, c.y}] = true
	neighbours[Coordinate{c.x - 1, c.y + 1}] = true
	neighbours[Coordinate{c.x, c.y + 1}] = true
	neighbours[Coordinate{c.x + 1, c.y + 1}] = true
	return neighbours
}

type Board struct {
	w, h       int
	aliveCells CoordinateSet
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
	game := Game{board: Board{w, h, make(CoordinateSet)}}
	return &game
}

func (g *Game) MakeMove() {
	newBoard := Board{
		g.board.w,
		g.board.h,
		make(CoordinateSet)}

	var aliveCount map[Coordinate]int = make(map[Coordinate]int)

	for k := range (g.board.aliveCells) {
		for n := range (k.Neighbours()) {
			aliveCount[n]++
		}
	}

	for k, v := range (aliveCount) {
		if (g.InRange(k)) {
			g.Rule1(&newBoard, k, v)
			g.Rule2(&newBoard, k, v)
			g.Rule3(&newBoard, k, v)
			g.Rule4(&newBoard, k, v)
		}
	}

	if (g.board.aliveCells.Equals(newBoard.aliveCells)) {
		g.finished = true
	} else {
		g.board = newBoard
		g.generation++
	}
}

func (g Game) Rule1(newBoard *Board, c Coordinate, aliveNeighbours int) {
	// Kill cell if less than 2 live neighbours -> do nothing
}

func (g Game) Rule2(newBoard *Board, c Coordinate, aliveNeighbours int) {
	// Cell stays alive if it has 2 or 3 live neighbours
	if (g.IsAlive(c) && (aliveNeighbours == 2 || aliveNeighbours == 3)) {
		g.Revive(newBoard, c)
	}
}

func (g Game) Rule3(newBoard *Board, c Coordinate, aliveNeighbours int) {
	// Kill cell if it has more than three live neighbours -> do nothing
}

func (g Game) Rule4(newBoard *Board, c Coordinate, aliveNeighbours int) {
	// Revive dead cell with exactly 3 live neighbours
	if (!g.IsAlive(c) && aliveNeighbours == 3) {
		g.Revive(newBoard, c)
	}
}