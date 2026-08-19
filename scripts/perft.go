package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"ranga/board"
	"strconv"
	"strings"
	"time"
)

const threshold = 100000

func PerftTestSuite(filepath string, outputPath string) error {
	text, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to perft test: %v", err)
	}

	rawTests := strings.Split(string(text), "\n")
	var tests []string
	for _, line := range rawTests {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			tests = append(tests, trimmed)
		}
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	sampleSize := min(len(tests), 1000)

	for i := range sampleSize {
		j := i + r.Intn(len(tests)-i)
		tests[i], tests[j] = tests[j], tests[i]
	}

	selectedTests := tests[:sampleSize]

	currentCount := 0
	var failedLines []string

	var runPerftTest = func(test string) error {
		currentCount++
		fields := strings.Split(test, "; ")
		if len(fields) < 2 {
			return nil
		}

		fen := fields[0]
		fmt.Printf("\r\033[K[%d/%d] Testing FEN: %s", currentCount, sampleSize, fen)
		os.Stdout.Sync()

		isLineFailed := false

		for _, field := range fields[1:] {
			fs := strings.Split(field, " ")
			if len(fs) < 2 {
				continue
			}

			depthStr := strings.TrimPrefix(fs[0], "D")
			depth, err := strconv.Atoi(depthStr)
			if err != nil {
				return fmt.Errorf("failed to parse depth component: %v", err)
			}

			nodes, err := strconv.Atoi(strings.TrimSpace(fs[1]))
			if err != nil {
				return fmt.Errorf("failed to parse node component: %v", err)
			}

			if nodes > threshold {
				return nil
			}

			b := board.NewBoard()
			if err := b.ParseFEN(fen); err != nil {
				return fmt.Errorf("failed to perft test: %v", err)
			}

			visited := board.Perft(context.Background(), &b, depth)
			if visited != int64(nodes) {
				isLineFailed = true
			}
		}

		if isLineFailed {
			failedLines = append(failedLines, test)
		}
		return nil
	}

	start := time.Now()

	for _, test := range selectedTests {
		if err := runPerftTest(test); err != nil {
			return fmt.Errorf("failed to perft test: %v", err)
		}
	}

	fmt.Printf("\n\nPerft Test Suite Completed: Randomly verified %d out of %d total positions in %v\n", sampleSize, len(rawTests), time.Since(start))

	outputData := strings.Join(failedLines, "\n") + "\n"
	if err := os.WriteFile(outputPath, []byte(outputData), 0644); err != nil {
		return fmt.Errorf("failed to write failed tests to file: %v", err)
	}
	fmt.Printf("Wrote %d failed test cases\n", len(failedLines))
	return nil
}

func main() {
	if err := PerftTestSuite("/Users/nchandur/workspace/ranga/data/perft-marcel.epd", "/Users/nchandur/workspace/ranga/data/perft-test-results.log"); err != nil {
		log.Fatal(err)
	}
}
