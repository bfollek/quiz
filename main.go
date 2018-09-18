package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
)

type qA struct {
	q string
	a string
}

func main() {
	qAs, err := load()
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(os.Stdin)
	var rv int
	numRight := 0
	for _, nxt := range qAs {
		if rv, err = askQuestion(reader, nxt); err != nil {
			log.Fatal(err)
		}
		numRight += rv
	}
	fmt.Printf("You answered %d questions. You got %d right.\n", len(qAs), numRight)
}

// askQuestion returns 1 if the answer was correct, else 0.
func askQuestion(reader *bufio.Reader, nxt qA) (int, error) {
	fmt.Printf("Question: %s. Your answer: ", nxt.q)
	text, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	if text[:len(text)-1] == nxt.a { // Drop trailing newline
		return 1, nil
	}
	return 0, nil
}

func load() ([]qA, error) {
	qAs := []qA{}
	var fileName = flag.String("file", "problems.csv", "CSV file that has the questions and answers")
	flag.Parse()
	csvFile, err := os.Open(*fileName)
	if err != nil {
		return qAs, err
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	records, err := reader.ReadAll()
	if err != nil {
		return qAs, err
	}
	for _, flds := range records {
		qAs = append(qAs, qA{q: flds[0], a: flds[1]})
	}
	return qAs, nil
}
