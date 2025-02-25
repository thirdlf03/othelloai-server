package game

import (
	"fmt"
)

const (
	hw     = 8
	hw2    = 64
	black  = 0
	white  = 1
	legal  = 2
	vacant = 3
)

var dy = [8]int{0, 1, 0, -1, 1, 1, -1, -1}
var dx = [8]int{1, 0, -1, 0, 1, -1, 1, -1}

type Othello struct {
	Grid    [hw][hw]int
	Player  int
	NStones [2]int
}

func inside(y, x int) bool {
	return 0 <= x && x < hw && 0 <= y && y < hw
}

func (o *Othello) Init() {
	for i := 0; i < hw; i++ {
		for j := 0; j < hw; j++ {
			o.Grid[i][j] = vacant
		}
	}
	o.Grid[3][3], o.Grid[4][4] = white, white
	o.Grid[3][4], o.Grid[4][3] = black, black
	o.Player = black
	o.NStones = [2]int{2, 2}
}

func (o *Othello) CheckLegal() bool {
	found := false
	for y := 0; y < hw; y++ {
		for x := 0; x < hw; x++ {
			if o.Grid[y][x] == legal {
				o.Grid[y][x] = vacant
			}
		}
	}

	for y := 0; y < hw; y++ {
		for x := 0; x < hw; x++ {
			if o.Grid[y][x] != vacant {
				continue
			}
			for d := 0; d < 8; d++ {
				ny, nx := y+dy[d], x+dx[d]
				if !inside(ny, nx) || o.Grid[ny][nx] != 1-o.Player {
					continue
				}
				for steps := 2; ; steps++ {
					ny, nx = y+dy[d]*steps, x+dx[d]*steps
					if !inside(ny, nx) || o.Grid[ny][nx] == vacant {
						break
					}
					if o.Grid[ny][nx] == o.Player {
						o.Grid[y][x] = legal
						found = true
						break
					}
				}
			}
		}
	}
	return found
}

func (o *Othello) Move(y, x int) bool {
	if !inside(y, x) || o.Grid[y][x] != legal {
		fmt.Println("非合法手です")
		return false
	}

	o.Grid[y][x] = o.Player
	nFlipped := 0
	for d := 0; d < 8; d++ {
		flipped := []struct{ y, x int }{}
		ny, nx := y+dy[d], x+dx[d]
		for inside(ny, nx) && o.Grid[ny][nx] == 1-o.Player {
			flipped = append(flipped, struct{ y, x int }{ny, nx})
			ny += dy[d]
			nx += dx[d]
		}
		if inside(ny, nx) && o.Grid[ny][nx] == o.Player {
			for _, f := range flipped {
				o.Grid[f.y][f.x] = o.Player
				nFlipped++
			}
		}
	}
	o.NStones[o.Player] += nFlipped + 1
	o.NStones[1-o.Player] -= nFlipped
	o.Player = 1 - o.Player
	return true
}

func (o *Othello) MoveStdin() {
	var coord string
	fmt.Print((map[bool]string{true: "黒", false: "白"})[o.Player == black] + " 着手: ")
	_, err := fmt.Scan(&coord)
	if err != nil || len(coord) < 2 {
		fmt.Println("座標を A1 や c5 のように入力してください")
		o.MoveStdin()
		return
	}
	y := int(coord[1] - '1')
	x := int(coord[0] - 'A')
	if !inside(y, x) {
		x = int(coord[0] - 'a')
		if !inside(y, x) {
			fmt.Println("座標を A1 や c5 のように入力してください")
			o.MoveStdin()
			return
		}
	}
	if !o.Move(y, x) {
		o.MoveStdin()
	}
}

func (o *Othello) Print() {
	fmt.Println("  A B C D E F G H")
	for y := 0; y < hw; y++ {
		fmt.Print(y+1, " ")
		for x := 0; x < hw; x++ {
			fmt.Print([]string{"X ", "O ", "* ", ". "}[o.Grid[y][x]])
		}
		fmt.Println()
	}
	fmt.Printf("黒 X %d - %d O 白\n", o.NStones[0], o.NStones[1])
}
