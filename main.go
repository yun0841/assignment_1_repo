/*
  Aurthors:
    Kuromi Chiikawa

  Revision history:
    Manual revision histories are hard to maintain and is always out-of-date.
    Please refer to the source control system's history
*/

package main

import (
	"fmt"
	"math"
	"sort"
)

type Solver struct {
	// A representation of the 9x9 Sudoku board
	// The board can be divided into 3x3 sub-regions
	// A B C
	// D E F
	// G H I
	// where A, B, C, D, E, F, G, H, I are 3x3 sub-regions
	board                [][]int // [row][col]
	curCombination       []int
	maxGcd               int
	perCombinationMaxGcd int
	rowsOrder            [9]int
	bestSol              [][]int
}

func (c *Solver) initBoard() {
	// . . . . . . . 2 .
	// . . . . . . . . 5
	// . 2 . . . . . . .
	// . . 0 . . . . . .
	// . . . . . . . . .
	// . . . 2 . . . . .
	// . . . . 0 . . . .
	// . . . . . 2 . . .
	// . . . . . . 5 . .

	c.board = [][]int{
		{-1, -1, -1, -1, -1, -1, -1, 2, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, 5},
		{-1, 2, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, 0, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, -1, -1, 2, -1, -1, -1, -1, -1},
		{-1, -1, -1, -1, 0, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1, 2, -1, -1, -1},
		{-1, -1, -1, -1, -1, -1, 5, -1, -1}}
}

func (c *Solver) arrayNumbersToInt(a []int) int {
	mul := int(math.Pow(10, float64(len(a))-1))
	result := 0
	for _, n := range a {
		result += n * mul
		mul /= 10
	}

	return result
}

func (c *Solver) findGcd(x int, y int) int {
	for y > 0 {
		x, y = y, x%y
	}

	return x
}

func (c *Solver) findGcdForArray(a []int) int {
	gcd := c.findGcd(a[0], a[1])
	for i := 2; i < len(a); i++ {
		gcd = c.findGcd(a[i], gcd)
	}

	return gcd
}

func (c *Solver) printBoard() {
	fmt.Println("------------------------")
	for i := range c.board {
		fmt.Printf("% 3d\n", c.board[i])
	}
	fmt.Println("------------------------")
}

func (c *Solver) printBestSol() {
	fmt.Println("------------------------")
	fmt.Println("Best solution:")
	for i := range c.bestSol {
		fmt.Printf("% 3d\n", c.bestSol[i])
	}
	fmt.Println("GCD: ", c.maxGcd)
	fmt.Println("------------------------")
}

func newSolver() Solver {
	c := Solver{}
	c.initBoard()
	c.maxGcd = 0
	c.perCombinationMaxGcd = 1

	// c.rowsOrder = [9]int{4, 3, 5, 0, 1, 2, 6, 7, 8}
	// c.rowsOrder = [9]int{4, 0, 1, 2, 3, 5, 6, 7, 8}
	// c.rowsOrder = [9]int{4, 0, 8, 1, 7, 2, 6, 3, 5}
	// c.rowsOrder = [9]int{0, 8, 1, 7, 2, 6, 3, 5, 4}
	c.rowsOrder = [9]int{0, 1, 2, 3, 5, 6, 7, 8, 4}
	// c.rowsOrder = [9]int{0, 1, 2, 3, 4, 5, 6, 7, 8}
	return c
}

func (c *Solver) checkCombination(numbers []int) bool {
	// check if numbers contain 0, 2, 5
	num_0_2_5 := 0
	for _, n := range numbers {
		if n != 2 && n != 5 && n != 0 {
			continue
		}

		num_0_2_5 += 1
		if num_0_2_5 == 3 {
			break
		}
	}
	return num_0_2_5 == 3
}

func (c *Solver) fullCheckBoardForDuplicates() bool {
	// check whole col to look for any duplicate of numbers[start]
	{
		for r := 0; r < 9; r++ {
			for col := 0; col < 9; col++ {
				num_check := [10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
				k := c.board[r][col]
				num_check[k] += 1
				if num_check[k] > 1 {
					fmt.Println("Col contains duplicates")
					return false
				}
			}
		}
	}

	// check duplicates in region
	for r := 0; r < 9; r += 3 {
		for col := 0; col < 9; col += 3 {
			num_check := [10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

			indexes := [9]int{
				c.board[r][col],
				c.board[r][col+1],
				c.board[r][col+2],
				c.board[r+1][col],
				c.board[r+1][col+1],
				c.board[r+1][col+2],
				c.board[r+2][col],
				c.board[r+2][col+1],
				c.board[r+2][col+2],
			}

			for _, k := range indexes {
				num_check[k] += 1

				if num_check[k] > 1 {
					fmt.Println("Region contains duplicates")
					return false
				}
			}
		}
	}

	return true
}

func (c *Solver) permutateNumbersForOneRow(numbers *[]int, start int, rowOrderIndex int, curGcd int) {
	row := c.rowsOrder[rowOrderIndex]

	if start == len(*numbers) {
		// backup
		rowsCopy := make([]int, len(c.board[row]))
		copy(rowsCopy, c.board[row])

		rowNum := c.arrayNumbersToInt(*numbers)
		if rowOrderIndex != 1 {
			curGcd = c.findGcd(rowNum, curGcd)
		} else {
			firstRowNum := c.arrayNumbersToInt(c.board[c.rowsOrder[rowOrderIndex-1]])
			curGcd = c.findGcd(rowNum, firstRowNum)
		}

		if (rowOrderIndex > 0 && curGcd <= c.maxGcd) || !c.assignBoardRowByIndex(row, numbers) {
			// recover
			copy(c.board[row], rowsCopy)
			return
		}

		if c.perCombinationMaxGcd < curGcd {
			c.perCombinationMaxGcd = curGcd
			c.printBoard()
			fmt.Println("New max GCD for this combination: ", c.perCombinationMaxGcd)
			fmt.Println("num: ", numbers)
		}

		if rowOrderIndex <= 7 {
			originalNumbers := make([]int, len(c.curCombination))
			copy(originalNumbers, c.curCombination)
			c.permutateNumbersForOneRow(&originalNumbers, 0, rowOrderIndex+1, curGcd)
		} else {
			c.printBoard()
			fmt.Println("GCD: ", curGcd)

			if c.maxGcd < curGcd {
				c.maxGcd = curGcd

				c.bestSol = make([][]int, len(c.board))
				for i := range c.board {
					c.bestSol[i] = make([]int, len(c.board[i]))
					copy(c.bestSol[i], c.board[i])
				}

				c.printBoard()
				fmt.Println("GCD: ", curGcd)
				fmt.Println("Max GCD: ", c.maxGcd)

				if !c.fullCheckBoardForDuplicates() {
					fmt.Println("Check for duplicates FAILED")
				} else {
					fmt.Println("Check for duplicates OK")
				}
			}
		}

		// recover
		copy(c.board[row], rowsCopy)

		return
	}

	for i := start; i < len(*numbers); i++ {
		// swap the ith element with the start element
		(*numbers)[i], (*numbers)[start] = (*numbers)[start], (*numbers)[i]
		passCheck := true

		// check not overriding predefined numbers
		if c.board[row][start] != -1 && c.board[row][start] != (*numbers)[start] {
			passCheck = false
		}

		// check whole col to look for any duplicate of numbers[start]
		if passCheck {
			for r := 0; r < 9; r++ {
				if r == row {
					continue
				}
				if (*numbers)[start] == c.board[r][start] {
					passCheck = false
					break
				}
			}
		}

		// check region to look for any duplicate of numbers
		if passCheck {
			r := row - (row % 3)       // row
			col := start - (start % 3) // column
			num_check := [10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

			if c.board[row][start] == -1 {
				num_check[(*numbers)[start]] += 1
			}

			indexes := [9]int{
				c.board[r][col],
				c.board[r][col+1],
				c.board[r][col+2],
				c.board[r+1][col],
				c.board[r+1][col+1],
				c.board[r+1][col+2],
				c.board[r+2][col],
				c.board[r+2][col+1],
				c.board[r+2][col+2],
			}

			for _, k := range indexes {
				if k == -1 {
					continue
				}

				num_check[k] += 1

				if num_check[k] > 1 {
					passCheck = false
					break
				}
			}
		}

		if passCheck {
			c.permutateNumbersForOneRow(numbers, start+1, rowOrderIndex, curGcd)
		}

		// recover
		(*numbers)[i], (*numbers)[start] = (*numbers)[start], (*numbers)[i]
	}
}

// Try to assign numbers to a row of the board. Return true if the assignment is successful.
func (c *Solver) assignBoardRowByIndex(row int, numbers *[]int) bool {
	for i, _ := range *numbers {
		// coordinates of the current cell: row, i
		c.board[row][i] = (*numbers)[i]
	}

	return true
}

func (c *Solver) pick9From10Numbers(numbers []int) {
	for i := range len(numbers) {
		// swap the ith element with the last element
		numbers[i], numbers[len(numbers)-1] = numbers[len(numbers)-1], numbers[i]

		if !c.checkCombination(numbers[0:9]) {
			// recover
			numbers[i], numbers[len(numbers)-1] = numbers[len(numbers)-1], numbers[i]
			continue
		}

		// sort.Ints(numbers[0:9])
		sortedNumbers := make([]int, len(numbers)-1)
		copy(sortedNumbers, numbers[0:9])
		sort.Ints(sortedNumbers)
		c.curCombination = make([]int, len(numbers)-1)
		copy(c.curCombination, sortedNumbers)
		fmt.Println(i, ": checking numbers :", c.curCombination)

		// for each row
		//   permute the numbers picked
		//	 assign the numbers to the row and check if the assignment is valid
		c.perCombinationMaxGcd = 1
		c.permutateNumbersForOneRow(&sortedNumbers, 0, 0, 1)

		// recover
		numbers[i], numbers[len(numbers)-1] = numbers[len(numbers)-1], numbers[i]
	}
}

func main() {
	fmt.Println("hello, World!")

	solver := newSolver()
	solver.printBoard()
	solver.pick9From10Numbers([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	solver.printBestSol()
	solver.printBoard()
}
