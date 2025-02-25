package board

import (
	"fmt"
	"sort"
)

type HashMap map[Board]int

const (
	inf      = 100000000
	cacheHit = 1000
)

var (
	transposeTableUpper       = HashMap{}
	transposeTableLower       = HashMap{}
	formerTransposeTableUpper = HashMap{}
	formerTransposeTableLower = HashMap{}
	visitedNodes              uint
)

func Init() {
	BoardInit()
	EvaluateInit()
	BookInit()
	EndgameEvaluateInit()
}

// InputBoard オセロボード読み込み
func InputBoard(arr []int) {
	for i := 0; i < hw2; i++ {
		num, err := fmt.Scan()
		if err != nil {
			panic(err)
		}

		if num == '0' {
			arr[i] = black
		} else if num == '1' {
			arr[i] = white
		} else {
			arr[i] = vacant
		}
	}
}

// move ordering用評価値計算
func calcMoveOrderingScore(b Board) int {
	res := 0
	if upper, upperExists := formerTransposeTableUpper[b]; upperExists {
		res = cacheHit - upper
	} else if lower, lowerExists := formerTransposeTableLower[b]; lowerExists {
		res = -cacheHit - lower
	} else {
		res = -Evaluate(b)
	}
	return res
}

func calcMoveOrderingScoreFinal(b Board) int {
	res := 0
	var legal bool
	for global := 0; global < hw2; global++ {
		legal = legalArr[b.player][b.BoardIdx[placeIncluded[global][0]]][localPlace[placeIncluded[global][0]][global]] ||
			legalArr[b.player][b.BoardIdx[placeIncluded[global][1]]][localPlace[placeIncluded[global][1]][global]] ||
			legalArr[b.player][b.BoardIdx[placeIncluded[global][2]]][localPlace[placeIncluded[global][2]][global]]
		if placeIncluded[global][3] != -1 {
			legal = legal || legalArr[b.player][b.BoardIdx[placeIncluded[global][3]]][localPlace[placeIncluded[global][3]][global]]
		}
		if legal {
			res++
		}
	}
	return -res
}

func negaAlphaTranspose(b Board, depth int, passed bool, alpha int, beta int) int {
	visitedNodes++

	if depth == 0 {
		return Evaluate(b)
	}

	u := inf
	l := -inf

	if upper, upperExists := transposeTableUpper[b]; upperExists {
		u = upper
	}
	if lower, lowerExists := transposeTableLower[b]; lowerExists {
		l = lower
	}

	if u == l {
		return u
	}

	alpha = max(alpha, l)
	beta = min(beta, u)

	maxScore := -inf
	canPut := 0
	childNodes := []Board{}
	for coord := 0; coord < hw2; coord++ {
		if b.Legal(coord) {
			childNodes = append(childNodes, b.Move(coord))
			childNodes[canPut].value = calcMoveOrderingScore(childNodes[canPut])
			canPut++
		}
	}

	if canPut == 0 {
		if passed {
			return EndgameEvaluate(b)
		}

		b.player = 1 - b.player
		return -negaAlphaTranspose(b, depth, true, -beta, -alpha)
	}

	if canPut >= 2 {
		sort.Slice(childNodes, func(i, j int) bool {
			return childNodes[i].value > childNodes[j].value
		})
	}

	for _, nb := range childNodes {
		g := -negaAlphaTranspose(nb, depth-1, false, -beta, -alpha)
		if g >= beta {
			return g
		}
		alpha = max(alpha, g)
		maxScore = max(maxScore, g)
	}

	if maxScore < alpha {
		transposeTableUpper[b] = maxScore
	} else {
		transposeTableUpper[b] = maxScore
		transposeTableLower[b] = maxScore
	}
	return maxScore
}

func negaScout(b Board, depth int, passed bool, alpha int, beta int) int {
	visitedNodes++

	if depth == 0 {
		return Evaluate(b)
	}

	u := inf
	l := -inf
	if upper, upperExists := transposeTableUpper[b]; upperExists {
		u = upper
	}
	if lower, lowerExists := transposeTableLower[b]; lowerExists {
		l = lower
	}

	if u == l {
		return u
	}

	alpha = max(alpha, l)
	beta = min(beta, u)

	maxScore := -inf
	canPut := 0
	childNodes := []Board{}
	for coord := 0; coord < hw2; coord++ {
		if b.Legal(coord) {
			childNodes = append(childNodes, b.Move(coord))
			childNodes[canPut].value = calcMoveOrderingScore(childNodes[canPut])
			canPut++
		}
	}

	if canPut == 0 {
		if passed {
			return EndgameEvaluate(b)
		}
		b.player = 1 - b.player
		return -negaScout(b, depth, true, -beta, -alpha)
	}

	if canPut >= 2 {
		sort.Slice(childNodes, func(i, j int) bool {
			return childNodes[i].value > childNodes[j].value
		})
	}

	g := -negaScout(childNodes[0], depth-1, false, -beta, -alpha)
	if g > alpha {
		transposeTableLower[b] = g
	}
	alpha = max(alpha, g)
	maxScore = max(maxScore, g)

	for i := 1; i < canPut; i++ {
		g = -negaAlphaTranspose(childNodes[i], depth-1, false, -alpha-1, -alpha)
		if g >= beta {
			if g > l {
				transposeTableLower[b] = g
			}
			return g
		}
		if g > alpha {
			alpha = g
			g = -negaScout(childNodes[i], depth-1, false, -beta, -alpha)
			if g >= beta {
				if g > l {
					transposeTableLower[b] = g
				}
				return g
			}
		}
		alpha = max(alpha, g)
		maxScore = max(maxScore, g)
	}

	if maxScore < alpha {
		transposeTableUpper[b] = maxScore
	} else {
		transposeTableUpper[b] = maxScore
		transposeTableLower[b] = maxScore
	}
	return maxScore
}

func Search(b Board, depth int) int {
	visitedNodes = 0
	transposeTableUpper = HashMap{}
	transposeTableLower = HashMap{}
	formerTransposeTableUpper = HashMap{}
	formerTransposeTableLower = HashMap{}
	childNodes := []Board{}
	canPut := 0
	for coord := 0; coord < hw2; coord++ {
		if b.Legal(coord) {
			childNodes = append(childNodes, b.Move(coord))
			canPut++
		}
	}

	var searchDepth, res, score, alpha, beta int
	var startDepth = max(1, depth-3)
	for searchDepth = startDepth; searchDepth <= depth; searchDepth++ {
		alpha = -inf
		beta = inf
		if canPut >= 2 {
			for _, nb := range childNodes {
				nb.value = calcMoveOrderingScore(nb)
			}
			sort.Slice(childNodes, func(i, j int) bool {
				return childNodes[i].value > childNodes[j].value
			})
		}

		score = -negaScout(childNodes[0], searchDepth-1, false, -beta, -alpha)
		alpha = score
		res = childNodes[0].policy

		for i := 1; i < canPut; i++ {
			score = -negaAlphaTranspose(childNodes[i], searchDepth-1, false, -alpha-1, -alpha)
			if alpha < score {
				alpha = score
				score = -negaScout(childNodes[i], searchDepth-1, false, -beta, -alpha)
				res = childNodes[i].policy
			}
			alpha = max(alpha, score)
		}
		fmt.Print("depth: ", searchDepth, " score: ", alpha, " policy: ", res, " nodes: ", visitedNodes, "\n")
		transposeTableUpper, formerTransposeTableUpper = formerTransposeTableUpper, transposeTableUpper
		transposeTableUpper = HashMap{}
		transposeTableLower, formerTransposeTableLower = formerTransposeTableLower, transposeTableLower
		transposeTableLower = HashMap{}
	}
	return res
}

func negaAlphaTransposeFinal(b Board, passed bool, alpha int, beta int) int {
	visitedNodes++
	u := inf
	l := -inf
	if _, upperExists := transposeTableUpper[b]; upperExists {
		u = transposeTableUpper[b]
	}
	if _, lowerExists := transposeTableLower[b]; lowerExists {
		l = transposeTableLower[b]
	}

	if u == l {
		return u
	}

	alpha = max(alpha, l)
	beta = min(beta, u)
	maxScore := -inf
	childNodes := []Board{}
	canPut := 0
	for coord := 0; coord < hw2; coord++ {
		if b.Legal(coord) {
			childNodes = append(childNodes, b.Move(coord))
			childNodes[canPut].value = calcMoveOrderingScoreFinal(childNodes[canPut])
			canPut++
		}
	}

	if canPut == 0 {
		if passed {
			return EndgameEvaluate(b)
		}
		b.player = 1 - b.player
		return -negaAlphaTransposeFinal(b, true, -beta, -alpha)
	}

	if canPut >= 2 {
		sort.Slice(childNodes, func(i, j int) bool {
			return childNodes[i].value > childNodes[j].value
		})
	}

	for _, nb := range childNodes {
		g := -negaAlphaTransposeFinal(nb, false, -beta, -alpha)
		if g >= beta {
			if g > l {
				transposeTableLower[b] = g
			}
			return g
		}
		alpha = max(alpha, g)
		maxScore = max(maxScore, g)
	}

	if maxScore < alpha {
		transposeTableUpper[b] = maxScore
	} else {
		transposeTableUpper[b] = maxScore
		transposeTableLower[b] = maxScore
	}
	return maxScore
}

func negaScoutFinal(b Board, passed bool, alpha int, beta int) int {
	visitedNodes++
	u := inf
	l := -inf
	if _, upperExists := transposeTableUpper[b]; upperExists {
		u = transposeTableUpper[b]
	}
	if _, lowerExists := transposeTableLower[b]; lowerExists {
		l = transposeTableLower[b]
	}
	if u == l {
		return u
	}

	alpha = max(alpha, l)
	beta = min(beta, u)

	childNodes := []Board{}
	canPut := 0
	maxScore := -inf
	for coord := 0; coord < hw2; coord++ {
		if b.Legal(coord) {
			childNodes = append(childNodes, b.Move(coord))
			childNodes[canPut].value = calcMoveOrderingScoreFinal(childNodes[canPut])
			canPut++
		}
	}

	if canPut == 0 {
		if passed {
			return EndgameEvaluate(b)
		}
		b.player = 1 - b.player
		return -negaScoutFinal(b, true, -beta, -alpha)
	}

	if canPut >= 2 {
		sort.Slice(childNodes, func(i, j int) bool {
			return childNodes[i].value > childNodes[j].value
		})
	}

	g := -negaScoutFinal(childNodes[0], false, -beta, -alpha)
	if g >= beta {
		if g > l {
			transposeTableLower[b] = g
		}
		return g
	}
	alpha = max(alpha, g)
	maxScore = max(maxScore, g)

	if maxScore < alpha {
		transposeTableUpper[b] = maxScore
	} else {
		transposeTableUpper[b] = maxScore
		transposeTableLower[b] = maxScore
	}
	return maxScore
}

func SearchFinal(b Board) int {
	visitedNodes = 0
	transposeTableUpper = HashMap{}
	transposeTableLower = HashMap{}
	formerTransposeTableUpper = HashMap{}
	formerTransposeTableLower = HashMap{}
	childNodes := []Board{}
	canPut := 0
	for coord := 0; coord < hw2; coord++ {
		if b.Legal(coord) {
			childNodes = append(childNodes, b.Move(coord))
			childNodes[canPut].value = calcMoveOrderingScore(childNodes[canPut])
			canPut++
		}
	}

	if canPut >= 2 {
		sort.Slice(childNodes, func(i, j int) bool {
			return childNodes[i].value > childNodes[j].value
		})
	}
	alpha := inf
	beta := -inf
	res := -1
	score := -negaScoutFinal(childNodes[0], false, -beta, -alpha)
	alpha = score
	res = childNodes[0].policy
	for i := 1; i < canPut; i++ {
		score = -negaAlphaTransposeFinal(childNodes[i], false, -alpha-step, -alpha)
		if alpha < score {
			alpha = score
			score = -negaScoutFinal(childNodes[i], false, -beta, -alpha)
			res = childNodes[i].policy
		}
		alpha = max(alpha, score)
	}
	fmt.Print("depth: ", 1, " score: ", alpha, " policy: ", res, " nodes: ", visitedNodes, "\n")
	return res
}
