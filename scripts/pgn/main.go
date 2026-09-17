package main

import (
	"log"
)

func main() {

	dir := "../../data/selftest-12/"
	source := dir + "games.pgn"
	destination := dir + "games.txt"

	fens, errs := ParseGamesFromFile(source)

	for _, err := range errs {
		if err != nil {
			log.Println(err)
		}
	}

	if err := WriteGamesToFile(destination, fens); err != nil {
		log.Fatal(err)
	}

	if err := WriteErrsToFile(dir+"errors.txt", errs); err != nil {
		log.Fatal(err)
	}

}
