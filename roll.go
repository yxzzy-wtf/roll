package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"
)

func globalRand(n int) int {
	r, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		panic(err)
	}
	return int(r.Int64())
}

func main() {
	if len(os.Args) == 0 {
		panic("No arguments (d2, d4, d6, d8, d10, d12, d20, d100)")
	}
	die := os.Args[1]
	rolls := 1

	if len(os.Args) == 3 {
		r, err := strconv.Atoi(os.Args[2])
		if err != nil {
			panic(err)
		}
		if r < 1 {
			panic("Cannot roll fewer than 1 times")
		}

		rolls = r
	}

	var dieFunc func() int
	var critGood = []int{}
	var critBad = []int{}
	var critNeut = []int{}
	switch die {
	case "d2":
		dieFunc = d2
	case "d4":
		dieFunc = d4
	case "d6":
		dieFunc = d6
	case "d8":
		dieFunc = d8
	case "d10":
		dieFunc = d10
		critGood = []int{10}
		critBad = []int{1}
	case "d12":
		dieFunc = d12
		critGood = []int{12}
		critBad = []int{1}
	case "d20":
		dieFunc = d20
		critGood = []int{20}
		critBad = []int{1}
	case "d100":
		dieFunc = d100
		critNeut = []int{11, 22, 33, 44, 55, 66, 77, 88, 99}
		critBad = []int{100}
		critGood = []int{1}
	default:
		panic("Invalid die choice (d2, d4, d6, d8, d10, d12, d20, d100)")
	}

	sum := 0
	for r := 0; r < rolls; r++ {
		bouncesMin := 2
		bouncesMax := 5

		if rolls == 1 {
			bouncesMin = 5
			bouncesMax = 11
		}
		sum += simulate(dieFunc, bouncesMin, bouncesMax, critNeut, critGood, critBad)
	}

	if rolls > 1 {
		fmt.Printf("Roll total: %v\n", sum)
	}

}

func simulate(dieFunc func() int, bouncesMin int, bouncesMax int, critNeut []int, critGood []int, critBad []int) int {
	bounces := globalRand(bouncesMax-bouncesMin) + bouncesMin

	fmt.Print("")

	for b := 0; b < bounces; b++ {
		preroll := dieFunc()
		fmt.Printf("\r>\033[34m%3d\033[0m", preroll) // blue while rolling

		rwait := globalRand(50) + 25
		time.Sleep(time.Millisecond * time.Duration(rwait))
	}

	actualRoll := dieFunc()
	colour := "\033[36m" // Blue
	if isIn(critNeut, actualRoll) {
		colour += "\033[1m\033[35m" // Purple
	} else if isIn(critBad, actualRoll) {
		colour = "\033[1m\033[31m" // Red
	} else if isIn(critGood, actualRoll) {
		colour = "\033[1m\033[32m" // Green
	}

	fmt.Printf("\r%v%4d\033[0m\n", colour, actualRoll)

	return actualRoll
}

func isIn(slice []int, item int) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

func d2() int {
	return globalRand(2) + 1
}

func d4() int {
	return globalRand(4) + 1
}

func d6() int {
	return globalRand(6) + 1
}

func d8() int {
	return globalRand(8) + 1
}

func d10() int {
	return globalRand(10) + 1
}

func d12() int {
	return globalRand(12) + 1
}

func d20() int {
	return globalRand(20) + 1
}

func d100() int {
	roll10 := globalRand(10)
	roll100 := globalRand(10) * 10
	roll := roll100 + roll10

	if roll == 0 {
		roll = 100
	}

	return roll
}
