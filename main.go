package main

import (
	"fmt"
	"os"
	"project/board"
	"project/game"
)

const hw = 8
const hw2 = 64
const black = 0
const white = 1
const vacant = 3
const completedDepth = 16

func ai(o *game.Othello, b *board.Board, aiPlayer int) {
	arr := make([]int, hw2)
	fmt.Println("AIが思考中...")

	for y := 0; y < hw; y++ {
		for x := 0; x < hw; x++ {
			idx := y*hw + x
			if o.Grid[y][x] == black {
				arr[idx] = black
			} else if o.Grid[y][x] == white {
				arr[idx] = white
			} else {
				arr[idx] = vacant
			}
		}
	}
	b.TranslateFromArr(arr, aiPlayer)
	board.Evaluate(*b)
	policy := board.GetBook(*b)
	if policy != -1 {
		y, x := policy/hw, policy%hw
		o.Move(y, x)
		return
	}
	if b.NStones >= hw2-completedDepth {
		policy = board.SearchFinal(*b)
	} else {
		policy = board.Search(*b, 8)
	}
	if policy == -1 {
		fmt.Println("AIがパスします")
		o.Player = 1 - o.Player
		return
	}
	y, x := policy/hw, policy%hw
	fmt.Printf("AIの手: %c%d\n", 'A'+x, y+1)
	fmt.Println("AIの手を反映中...")
	fmt.Println(" ")
	o.Move(y, x)
}

func main() {
	var aiPlayer int
	fmt.Print("AIの手番 0: 黒(先手) 1: 白(後手) : ")
	_, err := fmt.Scan(&aiPlayer)
	if err != nil || (aiPlayer != 0 && aiPlayer != 1) {
		fmt.Println("0か1を入力してください")
		os.Exit(1)
	}

	o := game.Othello{}
	o.Init()

	board.Init()
	b := board.Board{}

	for {
		if !o.CheckLegal() {
			o.Player = 1 - o.Player

			if !o.CheckLegal() {
				break
			}
		}
		o.Print()
		if o.Player == aiPlayer {
			ai(&o, &b, aiPlayer)
		} else {
			o.MoveStdin()
		}
	}

	// 最終結果の表示
	o.Print()
	if o.NStones[black] > o.NStones[white] {
		fmt.Println("黒(X)の勝ち!")
	} else if o.NStones[black] < o.NStones[white] {
		fmt.Println("白(O)の勝ち!")
	} else {
		fmt.Println("引き分け!")
	}
}
