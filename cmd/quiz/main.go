package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/bfollek/gophercises/quiz/internal/pkg/quiz"
)

const defaultFile = "data/problems.csv"

func main() {
	quiz, err := quiz.NewQuiz(parseArgs())
	if err != nil {
		log.Fatal(err)
	}
	numCorrect, err := quiz.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nThe quiz had %d questions. You got %d right.\n", quiz.NumQuestions, numCorrect)
}

func parseArgs() (string, int, bool) {
	var fileName = flag.String("file", defaultFile, "CSV file that has the questions and answers")
	var limit = flag.Int("limit", 30, "Time limit for the quiz, in seconds")
	var shuffle = flag.Bool("shuffle", false, "Shuffle the order of questions on each run (default no shuffle)")
	flag.Parse()
	return *fileName, *limit, *shuffle
}
