package main

import "math/rand"

func GenerateSymbolArray(symbols map[string]uint) []string {
	symbolArray := []string{}

	for symbol, count := range symbols {
		for i := uint(0); i < count; i++ {
			symbolArray = append(symbolArray, symbol)
		}
	}

	return symbolArray
}

func getRandomNumber(min int, max int) int {
	return min + rand.Intn(max-min+1)
}

func GetSpin(reel []string, rows int, cols int) [][]string {
	result := [][]string{}

	// Initialize the result slice with empty slices for each row
	for i := 0; i < rows; i++ {
		result = append(result, []string{})
	}

	for col := 0; col < cols; col++ {
		selected := map[int]bool{}

		for row := 0; row < rows; row++ {
			for true {
				randomIndex := getRandomNumber(0, len(reel)-1)

				_, exists := selected[randomIndex]

				if !exists {
					selected[randomIndex] = true
					result[row] = append(result[row], reel[randomIndex])
					break
				}
			}
		}
	}

	return result
}

func CheckWin(spin [][]string, multipliers map[string]uint) []uint {
	winningLines := []uint{}

	// Check for winning combinations (simplified example - check rows)
	for _, row := range spin {
		win := true
		checkSymbol := row[0]

		// Check if all symbols in the row are the same starting from the second symbol
		for _, symbol := range row[1:] {
			if symbol != checkSymbol {
				win = false
				break
			}
		}

		if win {
			winningLines = append(winningLines, multipliers[checkSymbol])
		} else {
			winningLines = append(winningLines, 0)
		}
	}

	return winningLines
}
