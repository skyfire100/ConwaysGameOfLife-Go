package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Constants defining the size of the game board and delay between frames
const (
	width          = 92
	height         = 47
	delayInSeconds = 0.1 // 100 milliseconds
)

// ANSI color codes for coloring the output
var (
	greenColorCode = "\033[32m"
	resetColorCode = "\033[0m"
)

// The main function that runs the game of life
func main() {
	rand.Seed(time.Now().UnixNano()) // seed the random number generator

	// Create the initial game board
	board := make([][]bool, height)
	for i := range board {
		board[i] = make([]bool, width)
		for j := range board[i] {
			board[i][j] = rand.Intn(2) == 0
		}
	}

	prevBoards := make([][][]bool, 0) // stores previous board states
	loopCount := 0                    // keeps track of the number of times the same board state is encountered

	// Game loop
	for {
		printBoard(board)             // Print the current state of the board
		board = nextGeneration(board) // Generate the next state of the board

		// Check if the board is stuck in a loop
		for _, prevBoard := range prevBoards {
			if equal(prevBoard, board) {
				loopCount++
				if loopCount == 3 { // If the same board state is encountered for the third time
					fmt.Println("Board stuck in loop. Resetting board...")
					time.Sleep(2 * time.Second)      // Pause for 3 seconds before resetting
					board = resetBoard()             // Reset the board to a new random state
					prevBoards = make([][][]bool, 0) // Clear the history of previous board states
					loopCount = 0                    // Reset the loop counter
				}
				break
			}
		}
		if loopCount == 0 {
			prevBoards = append(prevBoards, board) // Add the current board state to the history of previous board states
		}
		time.Sleep(time.Duration(delayInSeconds * float64(time.Second))) // Wait before the next frame
	}
}

// Helper function to compare two game boards
func equal(a, b [][]bool) bool {
	if len(a) != len(b) || len(a[0]) != len(b[0]) {
		return false
	}
	for i := range a {
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

// Helper function to generate a new random game board
func resetBoard() [][]bool {
	board := make([][]bool, height)
	for i := range board {
		board[i] = make([]bool, width)
		for j := range board[i] {
			board[i][j] = rand.Intn(2) == 0
		}
	}
	return board
}

// Function to print the board
func printBoard(board [][]bool) {
	// Clear the console before printing the board
	fmt.Print("\033c")

	// Print the top border
	topBorder := "+" + strings.Repeat("-", len(board[0])*2) + "+"
	fmt.Println(topBorder)

	// Print the board
	for i := range board {
		row := "|"
		for j := range board[i] {
			if board[i][j] {
				row += greenColorCode + "██" + resetColorCode // Use green color code for live cells
			} else {
				row += "  " // Use empty spaces for dead cells
			}
		}
		row += "|"
		fmt.Println(row)
	}

	// Print the bottom border
	fmt.Println(topBorder)
	fmt.Println()
}

// Function to compute the next generation of the board
func nextGeneration(board [][]bool) [][]bool {
	// Create a new board to store the next generation
	nextBoard := make([][]bool, len(board))
	for i := range nextBoard {
		nextBoard[i] = make([]bool, len(board[i]))
	}

	// Compute the state of each cell in the next generation
	for i := range board {
		for j := range board[i] {
			neighbors := countNeighbors(board, i, j)               // Count the number of live neighbors for this cell
			if board[i][j] && (neighbors == 2 || neighbors == 3) { // If the cell is live and has 2 or 3 live neighbors, it survives
				nextBoard[i][j] = true
			} else if !board[i][j] && neighbors == 3 { // If the cell is dead and has exactly 3 live neighbors, it becomes alive
				nextBoard[i][j] = true
			}
		}
	}

	// Return the next generation board
	return nextBoard
}

// Function to count the number of live neighbors for a given cell
func countNeighbors(board [][]bool, row, col int) int {
	count := 0
	for i := row - 1; i <= row+1; i++ {
		for j := col - 1; j <= col+1; j++ {
			if i < 0 || i >= len(board) || j < 0 || j >= len(board[i]) {
				continue // Skip cells that are outside the board
			}
			if i == row && j == col {
				continue // Skip the cell itself
			}
			if board[i][j] {
				count++ // Count live neighbors
			}
		}
	}
	return count
}
