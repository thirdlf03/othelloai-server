package board

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

const (
	scW1            = 6400
	p31             = 3
	p32             = 9
	p33             = 27
	p34             = 81
	p35             = 243
	p36             = 729
	p37             = 2187
	p38             = 6561
	p39             = 19683
	nPattern        = 3
	nDense0         = 16
	nDense1         = 16
	nAddInput       = 3
	nAddDense0      = 8
	nAllInput       = 4
	maxMobility     = 30
	maxSurround     = 50
	maxEvaluatedIdx = 59049
)

var (
	patternSize = [nPattern]int{8, 10, 10}
	mobilityArr [2][nLine]int
	surroundArr [2][nLine]int
	patternArr  [nPattern][maxEvaluatedIdx]float64
	addArr      [maxMobility*2 + 1][maxSurround + 1][maxSurround + 1]float64
	finalDense  [nAllInput]float64
	finalBias   float64
)

func evaluateInit1() {
	var place, b, w int
	for idx := 0; idx < nLine; idx++ {
		b = CreateOneColor(idx, 0)
		w = CreateOneColor(idx, 1)
		mobilityArr[black][idx] = 0
		mobilityArr[white][idx] = 0
		surroundArr[black][idx] = 0
		surroundArr[white][idx] = 0
		for place = 0; place < hw; place++ {
			if place > 0 {
				if ((1 & (b >> (place - 1))) == 0) && ((1 & (w >> (place - 1))) == 0) {
					if (1 & (b >> place)) == 1 {
						surroundArr[black][idx]++
					} else if (1 & (w >> place)) == 1 {
						surroundArr[white][idx]++
					}
				}
			}
			if place < hw-1 {
				if ((1 & (b >> (place + 1))) == 0) && ((1 & (w >> (place + 1))) == 0) {
					if (1 & (b >> place)) == 1 {
						surroundArr[black][idx]++
					} else if (1 & (w >> place)) == 1 {
						surroundArr[white][idx]++
					}
				}
			}
		}
		for place = 0; place < hw; place++ {
			if legalArr[black][idx][place] {
				mobilityArr[black][idx]++
			}
			if legalArr[white][idx][place] {
				mobilityArr[white][idx]++
			}
		}
	}
}

func leakyRelu(x float64) float64 {
	return max(0.01*x, x)
}

func predictPattern(patternSize int, inArr []float64, dense0 [nDense0][20]float64, bias0 [nDense0]float64, dense1 [nDense1][nDense0]float64, bias1 [nDense1]float64, dense2 [nDense1]float64, bias2 float64) float64 {
	hidden0 := make([]float64, 16)
	var hidden1 float64
	for i := 0; i < nDense0; i++ {
		hidden0[i] = bias0[i]
		for j := 0; j < patternSize*2; j++ {
			hidden0[i] += dense0[i][j] * inArr[j]
		}
		hidden0[i] = leakyRelu(hidden0[i])
	}
	res := bias2
	for i := 0; i < nDense1; i++ {
		hidden1 = bias1[i]
		for j := 0; j < nDense0; j++ {
			hidden1 += dense1[i][j] * hidden0[j]
		}
		hidden1 = leakyRelu(hidden1)
		res += dense2[i] * hidden1
	}
	res = leakyRelu(res)
	return res
}

func calcPop(a, b, s int) int {
	return (a / pow3[s-1-b]) % 3
}

func calcRevIdx(patternIdx int, patternSize int, idx int) int {
	res := 0
	if patternIdx <= 1 {
		for i := 0; i < patternSize; i++ {
			res += pow3[i] * calcPop(idx, i, patternSize)
		}
	} else if patternIdx == 2 {
		res += p39 * calcPop(idx, 0, patternSize)
		res += p38 * calcPop(idx, 4, patternSize)
		res += p37 * calcPop(idx, 7, patternSize)
		res += p36 * calcPop(idx, 9, patternSize)
		res += p35 * calcPop(idx, 1, patternSize)
		res += p34 * calcPop(idx, 5, patternSize)
		res += p33 * calcPop(idx, 8, patternSize)
		res += p32 * calcPop(idx, 2, patternSize)
		res += p31 * calcPop(idx, 6, patternSize)
		res += calcPop(idx, 3, patternSize)
	}
	return res
}

func preEvaluationPatter(patternIdx int, evaluateIdx int, patternSize int, dense0 [nDense0][20]float64, bias0 [nDense0]float64, dense1 [nDense1][nDense0]float64, bias1 [nDense1]float64, dense2 [nDense1]float64, bias2 float64) {
	var digit int
	arr := make([]float64, 20)
	var tmpPatternArr [maxEvaluatedIdx]float64
	for idx := 0; idx < pow3[patternSize]; idx++ {
		for i := 0; i < patternSize; i++ {
			digit = (idx / pow3[patternSize-1-i]) % 3
			if digit == 0 {
				arr[i] = 1.0
				arr[patternSize+1] = 0.0
			} else if digit == 1 {
				arr[i] = 0.0
				arr[patternSize+1] = 1.0
			} else {
				arr[i] = 0.0
				arr[patternSize+1] = 0.0
			}
		}
		patternArr[evaluateIdx][idx] = predictPattern(patternSize, arr, dense0, bias0, dense1, bias1, dense2, bias2)
		tmpPatternArr[calcRevIdx(patternIdx, patternSize, idx)] = patternArr[evaluateIdx][idx]
	}
	for idx := 0; idx < pow3[patternSize]; idx++ {
		patternArr[evaluateIdx][idx] += tmpPatternArr[idx]
	}
}

func predictAdd(mobility int, sur0 int, sur1 int, dense0 [nAddDense0][nAddInput]float64, bias0 [nAddDense0]float64, dense1 [nAddDense0]float64, bias1 float64) float64 {
	var hidden0 [nAddDense0]float64
	var inArr [nAddInput]float64
	inArr[0] = float64(mobility) / 30.0
	inArr[1] = (float64(sur0) - 15.0) / 15.0
	inArr[2] = (float64(sur1) - 15.0) / 15.0
	for i := 0; i < nAddDense0; i++ {
		hidden0[i] = bias0[i]
		for j := 0; j < nAddInput; j++ {
			hidden0[i] += dense0[i][j] * inArr[j]
		}
		hidden0[i] = leakyRelu(hidden0[i])
	}
	res := bias1
	for i := 0; i < nAddDense0; i++ {
		res += dense1[i] * hidden0[i]
	}
	res = leakyRelu(res)
	return res
}

func preEvaluationAdd(dense0 [nAddDense0][nAddInput]float64, bias0 [nAddDense0]float64, dense1 [nAddDense0]float64, bias1 float64) {
	for mobility := -maxMobility; mobility <= maxMobility; mobility++ {
		for sur0 := 0; sur0 <= maxSurround; sur0++ {
			for sur1 := 0; sur1 <= maxSurround; sur1++ {
				addArr[mobility+maxMobility][sur0][sur1] = predictAdd(mobility, sur0, sur1, dense0, bias0, dense1, bias1)
			}
		}
	}
}

func evaluateInit2() {
	file, err := os.Open("models/model.txt")
	if err != nil {
		fmt.Println("evaluation file not exist")
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var (
		dense0    [nDense0][20]float64
		bias0     [nDense0]float64
		dense1    [nDense1][nDense0]float64
		bias1     [nDense1]float64
		dense2    [nDense1]float64
		bias2     float64
		addDense0 [nAddDense0][nAddInput]float64
		addBias0  [nAddDense0]float64
		addDense1 [nAddDense0]float64
		addBias1  float64
	)

	for patternIdx := 0; patternIdx < nPattern; patternIdx++ {
		for i := 0; i < nDense0; i++ {
			for j := 0; j < patternSize[patternIdx]*2; j++ {
				scanner.Scan()
				line := scanner.Text()
				dense0[i][j], _ = strconv.ParseFloat(line, 64)
			}
		}
		for i := 0; i < nDense0; i++ {
			scanner.Scan()
			line := scanner.Text()
			bias0[i], _ = strconv.ParseFloat(line, 64)
		}
		for i := 0; i < nDense1; i++ {
			for j := 0; j < nDense0; j++ {
				scanner.Scan()
				line := scanner.Text()
				dense1[i][j], _ = strconv.ParseFloat(line, 64)
			}
		}
		for i := 0; i < nDense1; i++ {
			scanner.Scan()
			line := scanner.Text()
			bias1[i], _ = strconv.ParseFloat(line, 64)
		}
		for i := 0; i < nDense1; i++ {
			scanner.Scan()
			line := scanner.Text()
			dense2[i], _ = strconv.ParseFloat(line, 64)
		}
		scanner.Scan()
		line := scanner.Text()
		bias2, _ = strconv.ParseFloat(line, 64)
		preEvaluationPatter(patternIdx, patternIdx, patternSize[patternIdx], dense0, bias0, dense1, bias1, dense2, bias2)
	}

	for i := 0; i < nAddDense0; i++ {
		for j := 0; j < nAddInput; j++ {
			scanner.Scan()
			line := scanner.Text()
			addDense0[i][j], _ = strconv.ParseFloat(line, 64)
		}
	}
	for i := 0; i < nAddDense0; i++ {
		scanner.Scan()
		line := scanner.Text()
		addBias0[i], _ = strconv.ParseFloat(line, 64)
	}
	for i := 0; i < nAddDense0; i++ {
		scanner.Scan()
		line := scanner.Text()
		addDense1[i], _ = strconv.ParseFloat(line, 64)
	}
	scanner.Scan()
	line := scanner.Text()
	addBias1, _ = strconv.ParseFloat(line, 64)
	preEvaluationAdd(addDense0, addBias0, addDense1, addBias1)

	for i := 0; i < nAllInput; i++ {
		scanner.Scan()
		line = scanner.Text()
		finalDense[i], _ = strconv.ParseFloat(line, 64)
	}
	scanner.Scan()
	line = scanner.Text()
	finalBias, _ = strconv.ParseFloat(line, 64)
	fmt.Println("evaluation file read")
}

func EvaluateInit() {
	evaluateInit1()
	evaluateInit2()
}

func calcMobility(b Board) int {
	playerFactor := 1
	if b.player == 1 {
		playerFactor = -1
	}

	sum := mobilityArr[b.player][b.BoardIdx[0]] + mobilityArr[b.player][b.BoardIdx[1]] +
		mobilityArr[b.player][b.BoardIdx[2]] + mobilityArr[b.player][b.BoardIdx[3]] +
		mobilityArr[b.player][b.BoardIdx[4]] + mobilityArr[b.player][b.BoardIdx[5]] +
		mobilityArr[b.player][b.BoardIdx[6]] + mobilityArr[b.player][b.BoardIdx[7]] +
		mobilityArr[b.player][b.BoardIdx[8]] + mobilityArr[b.player][b.BoardIdx[9]] +
		mobilityArr[b.player][b.BoardIdx[10]] + mobilityArr[b.player][b.BoardIdx[11]] +
		mobilityArr[b.player][b.BoardIdx[12]] + mobilityArr[b.player][b.BoardIdx[13]] +
		mobilityArr[b.player][b.BoardIdx[14]] + mobilityArr[b.player][b.BoardIdx[15]] +
		mobilityArr[b.player][b.BoardIdx[16]-p35+1] + mobilityArr[b.player][b.BoardIdx[26]-p35+1] +
		mobilityArr[b.player][b.BoardIdx[27]-p35+1] + mobilityArr[b.player][b.BoardIdx[37]-p35+1] +
		mobilityArr[b.player][b.BoardIdx[17]-p34+1] + mobilityArr[b.player][b.BoardIdx[25]-p34+1] +
		mobilityArr[b.player][b.BoardIdx[28]-p34+1] + mobilityArr[b.player][b.BoardIdx[36]-p34+1] +
		mobilityArr[b.player][b.BoardIdx[18]-p33+1] + mobilityArr[b.player][b.BoardIdx[24]-p33+1] +
		mobilityArr[b.player][b.BoardIdx[29]-p33+1] + mobilityArr[b.player][b.BoardIdx[35]-p33+1] +
		mobilityArr[b.player][b.BoardIdx[19]-p32+1] + mobilityArr[b.player][b.BoardIdx[23]-p32+1] +
		mobilityArr[b.player][b.BoardIdx[30]-p32+1] + mobilityArr[b.player][b.BoardIdx[34]-p32+1] +
		mobilityArr[b.player][b.BoardIdx[20]-p31+1] + mobilityArr[b.player][b.BoardIdx[22]-p31+1] +
		mobilityArr[b.player][b.BoardIdx[31]-p31+1] + mobilityArr[b.player][b.BoardIdx[33]-p31+1] +
		mobilityArr[b.player][b.BoardIdx[21]] + mobilityArr[b.player][b.BoardIdx[32]]

	return playerFactor * sum
}

func sFill5(b int) int {
	if popDigit[b][2] != 2 {
		return b - p35 + 1
	} else {
		return b
	}
}

func sFill4(b int) int {
	if popDigit[b][3] != 2 {
		return b - p34 + 1
	} else {
		return b
	}
}

func sFill3(b int) int {
	if popDigit[b][4] != 2 {
		return b - p33 + 1
	} else {
		return b
	}
}

func sFill2(b int) int {
	if popDigit[b][5] != 2 {
		return b - p32 + 1
	} else {
		return b
	}
}

func sFill1(b int) int {
	if popDigit[b][6] != 2 {
		return b - p31 + 1
	} else {
		return b
	}
}

func calcSurround(b Board, p int) int {
	return surroundArr[p][b.BoardIdx[0]] + surroundArr[p][b.BoardIdx[1]] + surroundArr[p][b.BoardIdx[2]] + surroundArr[p][b.BoardIdx[3]] +
		surroundArr[p][b.BoardIdx[4]] + surroundArr[p][b.BoardIdx[5]] + surroundArr[p][b.BoardIdx[6]] + surroundArr[p][b.BoardIdx[7]] +
		surroundArr[p][b.BoardIdx[8]] + surroundArr[p][b.BoardIdx[9]] + surroundArr[p][b.BoardIdx[10]] + surroundArr[p][b.BoardIdx[11]] +
		surroundArr[p][b.BoardIdx[12]] + surroundArr[p][b.BoardIdx[13]] + surroundArr[p][b.BoardIdx[14]] + surroundArr[p][b.BoardIdx[15]] +
		surroundArr[p][sFill5(b.BoardIdx[16])] + surroundArr[p][sFill5(b.BoardIdx[26])] + surroundArr[p][sFill5(b.BoardIdx[27])] + surroundArr[p][sFill5(b.BoardIdx[37])] +
		surroundArr[p][sFill4(b.BoardIdx[17])] + surroundArr[p][sFill4(b.BoardIdx[25])] + surroundArr[p][sFill4(b.BoardIdx[28])] + surroundArr[p][sFill4(b.BoardIdx[36])] +
		surroundArr[p][sFill3(b.BoardIdx[18])] + surroundArr[p][sFill3(b.BoardIdx[24])] + surroundArr[p][sFill3(b.BoardIdx[29])] + surroundArr[p][sFill3(b.BoardIdx[35])] +
		surroundArr[p][sFill2(b.BoardIdx[19])] + surroundArr[p][sFill2(b.BoardIdx[23])] + surroundArr[p][sFill2(b.BoardIdx[30])] + surroundArr[p][sFill2(b.BoardIdx[34])] +
		surroundArr[p][sFill1(b.BoardIdx[20])] + surroundArr[p][sFill1(b.BoardIdx[22])] + surroundArr[p][sFill1(b.BoardIdx[31])] + surroundArr[p][sFill1(b.BoardIdx[33])] +
		surroundArr[p][b.BoardIdx[21]] + surroundArr[p][b.BoardIdx[32]]
}

func edge2x(b []int, x int, y int) float64 {
	return patternArr[1][popDigit[b[x]][1]*p39+b[y]*p31+popDigit[b[x]][6]]
}

func triangle0(b []int, w, x, y, z int) float64 {
	return patternArr[2][b[w]/p34*p36+b[x]/p35*p33+b[y]/p36*p31+b[z]/p37]
}

func triangle1(b []int, w, x, y, z int) float64 {
	return patternArr[2][reverseBoard[b[w]]/p34*p36+reverseBoard[b[x]]/p35*p33+reverseBoard[b[y]]/p36*p31+reverseBoard[b[z]]/p37]
}

func calcPattern(b Board) float64 {
	return finalDense[0]*(patternArr[0][b.BoardIdx[21]]+patternArr[0][b.BoardIdx[32]]) +
		finalDense[1]*(edge2x(b.BoardIdx[:], 1, 0)+edge2x(b.BoardIdx[:], 6, 7)+edge2x(b.BoardIdx[:], 9, 8)+edge2x(b.BoardIdx[:], 14, 15)) +
		finalDense[2]*(triangle0(b.BoardIdx[:], 0, 1, 2, 3)+triangle0(b.BoardIdx[:], 7, 6, 5, 4)+triangle0(b.BoardIdx[:], 15, 14, 13, 12)+triangle1(b.BoardIdx[:], 15, 14, 13, 12))
}

func Evaluate(b Board) int {
	mobility := min(maxMobility*2, max(0, maxMobility+calcMobility(b)))
	sur0 := min(maxSurround, calcSurround(b, 0))
	sur1 := min(maxSurround, calcSurround(b, 1))
	var res float64
	if b.player == 1 {
		res = -1.0 * (finalBias + calcPattern(b) + finalDense[3]*addArr[mobility][sur0][sur1])
	} else {
		res = finalBias + calcPattern(b) + finalDense[3]*addArr[mobility][sur0][sur1]
	}

	ans := int(math.Max(-1.0, math.Min(1.0, res)) * scW1)
	return ans
}
