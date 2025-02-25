package board

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var book = make(map[string]int)

func BookInit() {
	file, err := os.Open("book/book.txt")
	if err != nil {
		fmt.Println("book file not exist")
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	b := Board{
		BoardIdx: [nBoardIdx]int{},
	}
	var y, x int
	n := 0
	firstBoard := [nBoardIdx]int{6560, 6560, 6560, 6425, 6371, 6560, 6560, 6560, 6560, 6560, 6560, 6425, 6371, 6560, 6560, 6560, 6560, 6560, 6560, 6560, 6398, 6452, 6398, 6560, 6560, 6560, 6560, 6560, 6560, 6560, 6560, 6479, 6344, 6479, 6560, 6560, 6560, 6560}

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 2 {
			continue
		}
		n++

		for i := 0; i < 4; i++ {
			b.player = 0
			copy(b.BoardIdx[:], firstBoard[:])

			for j := 0; j < len(line)-2; j += 2 {
				x = int(line[j]) - int('a')
				y = int(line[j+1]) - int('1')
				b = b.Move(transformMove(x, y, i))
			}
			x = int(line[len(line)-2]) - int('a')
			y = int(line[len(line)-1]) - int('1')
			book[strconv.FormatUint(b.Hash(), 10)] = transformMove(x, y, i)
		}
	}

	if err = scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	fmt.Printf("book board init %d\n", n)
}

func GetBook(b Board) int {
	if v, ok := book[strconv.FormatUint(b.Hash(), 10)]; ok {
		return v
	}
	return -1
}

func transformMove(x, y, rot int) int {
	switch rot {
	case 0:
		return y*hw + x
	case 1:
		return x*hw + y
	case 2:
		return (hw-1-y)*hw + (hw - 1 - x)
	case 3:
		return (hw-1-x)*hw + (hw - 1 - y)
	}
	return -1
}
