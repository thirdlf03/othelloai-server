package main

import (
	"connectrpc.com/connect"
	"context"
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log"
	"net/http"
	"project/board"
	"project/game"
	othellov1 "project/gen/othello/v1"
	"project/gen/othello/v1/othelloconnect"
)

const (
	hw             = 8
	hw2            = 64
	black          = 0
	white          = 1
	vacant         = 3
	completedDepth = 16
)

type OthelloServer struct{}

func ai(o *game.Othello, b *board.Board, aiPlayer int) (int, int) {
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
		return y, x
	}
	if b.NStones >= hw2-completedDepth {
		policy = board.SearchFinal(*b)
	} else {
		policy = board.Search(*b, 8)
	}
	if policy == -1 {
		fmt.Println("AIがパスします")
		o.Player = 1 - o.Player
		return -1, -1
	}
	y, x := policy/hw, policy%hw
	return y, x
}

var othello game.Othello
var b board.Board

func convertToOthelloGrid(arr []int32, o *game.Othello) {
	for i := 0; i < hw; i++ {
		for j := 0; j < hw; j++ {
			if arr[i*hw+j] == 0 {
				o.Grid[i][j] = vacant
			} else if arr[i*hw+j] == 2 {
				o.Grid[i][j] = white
			} else {
				o.Grid[i][j] = black
			}
			fmt.Printf("%d", arr[i*hw+j])
		}
		fmt.Println(" ")
	}
}

func (o *OthelloServer) GetAIMove(
	ctx context.Context,
	req *connect.Request[othellov1.GetAIMoveRequest],
) (*connect.Response[othellov1.GetAIMoveResponse], error) {
	var player int

	log.Println("Request headers: ", req.Header())

	body := req.Msg.Board
	convertToOthelloGrid(body, &othello)

	fmt.Println("Request body: ")
	fmt.Println(req.Msg.Player)

	if int(req.Msg.Player) == 1 {
		player = black
	} else {
		player = white
	}

	y, x := ai(&othello, &b, player)
	res := connect.NewResponse[othellov1.GetAIMoveResponse](
		&othellov1.GetAIMoveResponse{
			Y: int32(y),
			X: int32(x),
		})
	res.Header().Set("Othello-Version", "v1")
	fmt.Println("AIの手: ", y, x, "color: ", req.Msg.Player)
	return res, nil
}

func main() {
	othello.Init()
	board.Init()

	server := &OthelloServer{}
	mux := http.NewServeMux()
	path, handler := othelloconnect.NewOthelloServiceHandler(server)
	mux.Handle(path, handler)
	fmt.Println("Server started at localhost:8080")
	err := http.ListenAndServe(
		"localhost:8080",
		h2c.NewHandler(mux, &http2.Server{}),
	)
	if err != nil {
		fmt.Println(err)
	}
}
