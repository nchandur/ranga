package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"ranga/internal/board"
	"strconv"
	"strings"
	"time"
)

func PerftTestSuite(inputPath string, sample, threshold int) error {
	text, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to perft test: %v", err)
	}

	rawLines := strings.Split(string(text), "\n")
	var tests []string
	for _, line := range rawLines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			tests = append(tests, trimmed)
		}
	}

	totalPositions := len(tests)

	if sample > 0 && sample < len(tests) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := range sample {
			j := i + r.Intn(len(tests)-i)
			tests[i], tests[j] = tests[j], tests[i]
		}
		tests = tests[:sample]
	}

	positionsSearched := 0
	successCount := 0
	failCount := 0
	var failedLines []string

	runPerftTest := func(test string) error {
		fields := strings.Split(test, "; ")
		if len(fields) < 2 {
			return nil
		}

		fen := fields[0]

		testedThisLine := false
		lineFailed := false

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
				continue
			}
			testedThisLine = true

			b := board.NewBoard()
			if err := b.ParseFEN(fen); err != nil {
				return fmt.Errorf("failed to perft test: %v", err)
			}

			visited := board.Perft(context.Background(), &b, depth)
			if visited != int64(nodes) {
				lineFailed = true
			}
		}

		if testedThisLine {
			positionsSearched++
			fmt.Printf("\r\033[K[%d] Testing FEN: %s", positionsSearched, fen)
			os.Stdout.Sync()

			if lineFailed {
				failCount++
				failedLines = append(failedLines, test)
			} else {
				successCount++
			}
		}

		return nil
	}

	start := time.Now()

	for _, test := range tests {
		if err := runPerftTest(test); err != nil {
			return fmt.Errorf("failed to perft test: %v", err)
		}
	}

	fmt.Printf("\n\nPerft Test Suite Completed in %v\n", time.Since(start))
	fmt.Printf("Total positions in file: %d\n", totalPositions)
	fmt.Printf("Positions searched: %d\n", positionsSearched)
	fmt.Printf("Successful:         %d\n", successCount)
	fmt.Printf("Failed:             %d\n", failCount)

	if len(failedLines) > 0 {
		outputPath := filepath.Join(filepath.Dir(inputPath), "perft-test-results.log")
		outputData := strings.Join(failedLines, "\n") + "\n"
		if err := os.WriteFile(outputPath, []byte(outputData), 0644); err != nil {
			return fmt.Errorf("failed to write failed tests to file: %v", err)
		}
		fmt.Printf("Wrote %d failed test case(s) to %s\n", len(failedLines), outputPath)
	}

	return nil
}

func main() {
	inputPath := flag.String("input", "", "path to the perft EPD test file (required)")
	sample := flag.Int("sample", 0, "number of positions to randomly sample; 0 or omitted runs all positions")
	threshold := flag.Int("threshold", 1000, "number of maximum nodes expected")

	flag.Parse()

	if *inputPath == "" {
		log.Fatal("missing required -input flag, e.g. -input=/path/to/perft.epd")
	}

	if *sample < 0 {
		log.Fatal("--sample must be a non-negative integer")
	}

	if *threshold < 0 {
		log.Fatal("--threshold must be a non-negative integer")
	}

	if err := PerftTestSuite(*inputPath, *sample, *threshold); err != nil {
		log.Fatal(err)
	}
}
