package board

import (
	"fmt"
	"os"
)

const hw int = 8
const hw2 int = 64
const nBoardIdx int = 38
const nLine int = 6561
const black int = 0
const white int = 1
const vacant int = 2

var moveOffset = [nBoardIdx]int{
	1, 1, 1, 1, 1, 1, 1, 1,
	8, 8, 8, 8, 8, 8, 8, 8,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
}

var globalPlace = [nBoardIdx][hw]int{
	{0, 1, 2, 3, 4, 5, 6, 7},
	{8, 9, 10, 11, 12, 13, 14, 15},
	{16, 17, 18, 19, 20, 21, 22, 23},
	{24, 25, 26, 27, 28, 29, 30, 31},
	{32, 33, 34, 35, 36, 37, 38, 39},
	{40, 41, 42, 43, 44, 45, 46, 47},
	{48, 49, 50, 51, 52, 53, 54, 55},
	{56, 57, 58, 59, 60, 61, 62, 63},

	{0, 8, 16, 24, 32, 40, 48, 56},
	{1, 9, 17, 25, 33, 41, 49, 57},
	{2, 10, 18, 26, 34, 42, 50, 58},
	{3, 11, 19, 27, 35, 43, 51, 59},
	{4, 12, 20, 28, 36, 44, 52, 60},
	{5, 13, 21, 29, 37, 45, 53, 61},
	{6, 14, 22, 30, 38, 46, 54, 62},
	{7, 15, 23, 31, 39, 47, 55, 63},

	{5, 14, 23, -1, -1, -1, -1, -1},
	{4, 13, 22, 31, -1, -1, -1, -1},
	{3, 12, 21, 30, 39, -1, -1, -1},
	{2, 11, 20, 29, 38, 47, -1, -1},
	{1, 10, 19, 28, 37, 46, 55, -1},
	{0, 9, 18, 27, 36, 45, 54, 63},
	{8, 17, 26, 35, 44, 53, 62, -1},
	{16, 25, 34, 43, 52, 61, -1, -1},
	{24, 33, 42, 51, 60, -1, -1, -1},
	{32, 41, 50, 59, -1, -1, -1, -1},
	{40, 49, 58, -1, -1, -1, -1, -1},

	{2, 9, 16, -1, -1, -1, -1, -1},
	{3, 10, 17, 24, -1, -1, -1, -1},
	{4, 11, 18, 25, 32, -1, -1, -1},
	{5, 12, 19, 26, 33, 40, -1, -1},
	{6, 13, 20, 27, 34, 41, 48, -1},
	{7, 14, 21, 28, 35, 42, 49, 56},
	{15, 22, 29, 36, 43, 50, 57, -1},
	{23, 30, 37, 44, 51, 58, -1, -1},
	{31, 38, 45, 52, 59, -1, -1, -1},
	{39, 46, 53, 60, -1, -1, -1, -1},
	{47, 54, 61, -1, -1, -1, -1, -1},
}

var moveArr [2][nLine][hw][2]int
var legalArr [2][nLine][hw]bool
var flipArr [2][nLine][hw]int
var putArr [2][nLine][hw]int
var localPlace [nBoardIdx][hw2]int
var placeIncluded [hw2][4]int
var reverseBoard [nLine]int
var pow3 [11]int
var popDigit [nLine][hw]int
var popMid [nLine][hw][hw]int

type Board struct {
	BoardIdx [nBoardIdx]int
	player   int
	policy   int
	value    int
	NStones  int
}

func CreateOneColor(idx int, k int) int {
	res := 0
	for i := 0; i < hw; i++ {
		if idx%3 == k {
			res |= 1 << i
		}
		idx /= 3
	}
	return res
}

func trans(pt int, k int) int {
	if k == 0 {
		return pt << 1
	} else {
		return pt >> 1
	}
}

func moveLineHalf(p int, o int, place int, k int) int {
	var mask int
	res := 0
	pt := 1 << (hw - 1 - place)
	if (pt&p) != 0 || (pt&o) != 0 {
		return res
	}
	mask = trans(pt, k)
	for mask != 0 && (mask&o) != 0 {
		res++
		mask = trans(mask, k)
		if (mask & p) != 0 {
			return res
		}
	}
	return 0
}

func BoardInit() {
	var idx, b, w, place, i, j, k, lPlace, incIdx int
	pow3[0] = 1
	for idx = 1; idx < 11; idx++ {
		pow3[idx] = pow3[idx-1] * 3
	}
	for i = 0; i < nLine; i++ {
		for j = 0; j < hw; j++ {
			popDigit[i][j] = (i / pow3[hw-1-j]) % 3
		}
	}
	for i = 0; i < nLine; i++ {
		for j = 0; j < hw; j++ {
			for k = 0; k < hw; k++ {
				popMid[i][j][k] = (i - (i/pow3[j])*pow3[j]) / pow3[k]
			}
		}
	}
	for idx = 0; idx < nLine; idx++ {
		b = CreateOneColor(idx, 0)
		w = CreateOneColor(idx, 1)
		for place = 0; place < hw; place++ {
			reverseBoard[idx] *= 3
			if (b>>(place))&1 == 1 {
				reverseBoard[idx] += 0
			} else if (w>>(place))&1 == 1 {
				reverseBoard[idx] += 1
			} else {
				reverseBoard[idx] += 2
			}
		}
		for place = 0; place < hw; place++ {
			moveArr[black][idx][place][0] = moveLineHalf(b, w, place, 0)
			moveArr[black][idx][place][1] = moveLineHalf(b, w, place, 1)
			if moveArr[black][idx][place][0] != 0 || moveArr[black][idx][place][1] != 0 {
				legalArr[black][idx][place] = true
			} else {
				legalArr[black][idx][place] = false
			}
			moveArr[white][idx][place][0] = moveLineHalf(w, b, place, 0)
			moveArr[white][idx][place][1] = moveLineHalf(w, b, place, 1)
			if moveArr[white][idx][place][0] != 0 || moveArr[white][idx][place][1] != 0 {
				legalArr[white][idx][place] = true
			} else {
				legalArr[white][idx][place] = false
			}
		}
		for place = 0; place < hw; place++ {
			flipArr[black][idx][place] = idx
			flipArr[white][idx][place] = idx
			putArr[black][idx][place] = idx
			putArr[white][idx][place] = idx
			if (b>>(hw-1-place))&1 == 1 {
				flipArr[white][idx][place] += pow3[hw-1-place]
			} else if (w>>(hw-1-place))&1 == 1 {
				flipArr[black][idx][place] -= pow3[hw-1-place]
			} else {
				putArr[black][idx][place] -= pow3[hw-1-place] * 2
				putArr[white][idx][place] -= pow3[hw-1-place]
			}
		}
	}
	for place = 0; place < hw2; place++ {
		incIdx = 0
		for idx = 0; idx < nBoardIdx; idx++ {
			for lPlace = 0; lPlace < hw; lPlace++ {
				if globalPlace[idx][lPlace] == place {
					placeIncluded[place][incIdx] = idx
					incIdx++
				}
			}
		}
		if incIdx == 3 {
			placeIncluded[place][incIdx] = -1
		}
	}
	for idx = 0; idx < nBoardIdx; idx++ {
		for place = 0; place < hw2; place++ {
			localPlace[idx][place] = -1
			for lPlace = 0; lPlace < hw; lPlace++ {
				if globalPlace[idx][lPlace] == place {
					localPlace[idx][place] = lPlace
				}
			}
		}
	}
	fmt.Fprintln(os.Stderr, "Board initialized")
}

func (b Board) Less(another Board) bool {
	return b.value > another.value
}

func (b Board) Equal(another Board) bool {
	if b.player != another.player {
		return false
	}
	for i := 0; i < hw; i++ {
		if b.BoardIdx[i] != another.BoardIdx[i] {
			return false
		}
	}
	return true
}

func (b Board) NotEqual(another Board) bool {
	return !b.Equal(another)
}

func (b Board) Hash() uint64 {
	return uint64(b.BoardIdx[0] +
		b.BoardIdx[1]*17 +
		b.BoardIdx[2]*289 +
		b.BoardIdx[3]*4913 +
		b.BoardIdx[4]*83521 +
		b.BoardIdx[5]*1419857 +
		b.BoardIdx[6]*24137549 +
		b.BoardIdx[7]*410338673)
}

func (b *Board) Print() {
	var i, j, tmp int
	var res string
	for i = 0; i < hw; i++ {
		tmp = b.BoardIdx[i]
		res = ""
		for j = 0; j < hw; j++ {
			if tmp%3 == 0 {
				res = "X " + res
			} else if tmp%3 == 1 {
				res = "O " + res
			} else {
				res = ". " + res
			}
			tmp /= 3
		}
		fmt.Fprintln(os.Stderr, res)
	}
	fmt.Fprintln(os.Stderr)
}

func (b *Board) Legal(gPlace int) bool {
	res := false
	for i := 0; i < 3; i++ {
		res = res || legalArr[b.player][b.BoardIdx[placeIncluded[gPlace][i]]][localPlace[placeIncluded[gPlace][i]][gPlace]]
	}
	if placeIncluded[gPlace][3] != -1 {
		res = res || legalArr[b.player][b.BoardIdx[placeIncluded[gPlace][3]]][localPlace[placeIncluded[gPlace][3]][gPlace]]
	}
	return res
}

func (b *Board) Move(gPlace int) Board {
	var res Board
	for i := 0; i < nBoardIdx; i++ {
		res.BoardIdx[i] = b.BoardIdx[i]
	}
	b.moveP(&res, gPlace, 0)
	b.moveP(&res, gPlace, 1)
	b.moveP(&res, gPlace, 2)
	if placeIncluded[gPlace][3] != -1 {
		b.moveP(&res, gPlace, 3)
	}
	for i := 0; i < 3; i++ {
		res.BoardIdx[placeIncluded[gPlace][i]] = putArr[b.player][res.BoardIdx[placeIncluded[gPlace][i]]][localPlace[placeIncluded[gPlace][i]][gPlace]]
	}
	if placeIncluded[gPlace][3] != -1 {
		res.BoardIdx[placeIncluded[gPlace][3]] = putArr[b.player][res.BoardIdx[placeIncluded[gPlace][3]]][localPlace[placeIncluded[gPlace][3]][gPlace]]
	}
	res.player = 1 - b.player
	res.NStones = b.NStones + 1
	res.policy = gPlace
	return res
}

func (b *Board) TranslateToArr(res []int) {
	var i, j int
	for i = 0; i < hw; i++ {
		for j = 0; j < hw; j++ {
			res[i*hw+j] = popDigit[b.BoardIdx[i]][j]
		}
	}
}

func (b *Board) TranslateFromArr(arr []int, player int) {
	var i, j int
	for i = 0; i < nBoardIdx; i++ {
		b.BoardIdx[i] = nLine - 1
	}
	b.NStones = hw2
	for i = 0; i < hw2; i++ {
		for j = 0; j < 4; j++ {
			if placeIncluded[i][j] == -1 {
				continue
			}
			if arr[i] == black {
				b.BoardIdx[placeIncluded[i][j]] -= 2 * pow3[hw-1-localPlace[placeIncluded[i][j]][i]]
			} else if arr[i] == white {
				b.BoardIdx[placeIncluded[i][j]] -= pow3[hw-1-localPlace[placeIncluded[i][j]][i]]
			} else if j == 0 {
				b.NStones--
			}
		}
	}
	b.player = player
}

func (b *Board) flip(res *Board, gPlace int) {
	var i int
	for i = 0; i < 3; i++ {
		res.BoardIdx[placeIncluded[gPlace][i]] =
			flipArr[b.player][res.BoardIdx[placeIncluded[gPlace][i]]][localPlace[placeIncluded[gPlace][i]][gPlace]]
	}
	if placeIncluded[gPlace][3] != -1 {
		res.BoardIdx[placeIncluded[gPlace][3]] =
			flipArr[b.player][res.BoardIdx[placeIncluded[gPlace][3]]][localPlace[placeIncluded[gPlace][3]][gPlace]]
	}
}

func (b *Board) moveP(res *Board, gPlace int, i int) {
	var j, place int
	place = localPlace[placeIncluded[gPlace][i]][gPlace]
	for j = 1; j <= moveArr[b.player][b.BoardIdx[placeIncluded[gPlace][i]]][place][0]; j++ {
		b.flip(res, gPlace-moveOffset[placeIncluded[gPlace][i]]*j)
	}
	for j = 1; j <= moveArr[b.player][b.BoardIdx[placeIncluded[gPlace][i]]][place][1]; j++ {
		b.flip(res, gPlace+moveOffset[placeIncluded[gPlace][i]]*j)
	}
}
