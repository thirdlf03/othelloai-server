package main

import (
	"connectrpc.com/connect"
	"context"
	"log"
	"net/http"

	othellov1 "project/gen/othello/v1"
	"project/gen/othello/v1/othelloconnect"
)

func main() {
	client := othelloconnect.NewOthelloServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
		connect.WithGRPC(),
	)

	res, err := client.GetAIMove(
		context.Background(),
		connect.NewRequest(&othellov1.GetAIMoveRequest{
			Board: []int32{
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 1, 2, 0, 0, 0,
				0, 0, 0, 2, 1, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
			},
			Player: 2,
		}),
	)

	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res.Msg.GetY(), res.Msg.GetX())
}
