package main

import (
	"flag"
	"log"
	"path/filepath"
)

func main() {

	dir := flag.String("dir", "", "path to games directory")
	flag.Parse()

	if *dir == "" {
		log.Fatal("path to games directory required")
	}

	source := filepath.Join(*dir, "games.pgn")
	destination := filepath.Join(*dir, "games.txt")
	errors := filepath.Join(*dir, "errors.log")

	fens, errs := ParseGamesFromFile(source)

	for _, err := range errs {
		if err != nil {
			log.Println(err)
		}
	}

	if err := WriteGamesToFile(destination, fens); err != nil {
		log.Fatal(err)
	}

	if err := WriteErrsToFile(errors, errs); err != nil {
		log.Fatal(err)
	}

}
