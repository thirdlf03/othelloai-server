package board

const (
	step = 100
)

var countArr [nLine]int
var countAllArr [nLine]int

func EndgameEvaluateInit() {
	for idx := 0; idx < nLine; idx++ {
		b := CreateOneColor(idx, 0)
		w := CreateOneColor(idx, 1)
		countArr[idx] = 0
		countAllArr[idx] = 0
		for place := 0; place < hw; place++ {
			countArr[idx] += 1 & (b >> place)
			countArr[idx] -= 1 & (w >> place)
			countAllArr[idx] += 1 & (b >> place)
			countAllArr[idx] += 1 & (w >> place)
		}
	}
}

func EndgameEvaluate(b Board) int {
	count := 0
	nVacant := hw2
	for i := 0; i < hw; i++ {
		count += countArr[b.BoardIdx[i]]
		nVacant -= countAllArr[b.BoardIdx[i]]
	}
	if b.player == 1 {
		count = -count
	}
	if count > 0 {
		count += nVacant
	} else if count < 0 {
		count -= nVacant
	}
	return count * step
}
