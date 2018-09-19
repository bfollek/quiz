package quiz

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/bfollek/gophercises/quiz/pkg/util"
)

type question struct {
	prompt string
	answer string
}

// Quiz holds the attributes of a quiz.
type Quiz struct {
	quizChan     chan question
	NumQuestions int
	timeLimit    int // seconds
}

// NewQuiz initializes a new Quiz. It also returns an error code.
func NewQuiz(fileName string, timeLimit int, shuffle bool) (*Quiz, error) {
	questions, err := loadQuestions(fileName)
	if err != nil {
		return nil, err
	}
	quiz := Quiz{NumQuestions: len(questions), timeLimit: timeLimit}
	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(quiz.NumQuestions, func(i, j int) { questions[i], questions[j] = questions[j], questions[i] })
	}
	quiz.quizChan = make(chan question, quiz.NumQuestions)
	// Send all questions to channel ahead of time
	for _, q := range questions {
		quiz.quizChan <- q
	}
	return &quiz, nil
}

func loadQuestions(fileName string) ([]question, error) {
	csvFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	questions := []question{}
	for _, flds := range records {
		questions = append(questions, question{prompt: flds[0], answer: normalize(flds[1])})
	}
	return questions, nil
}

func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// Run runs a quiz till the user answers all questions, or times out, or there's an error.
// It returns the number of correct answers and an error code.
func (quiz *Quiz) Run() (int, error) {
	var numAnswers, numCorrect int
	var timesUp bool
	score := make(chan int, quiz.NumQuestions)
	errors := make(chan error)
	err := promptToStart(quiz.timeLimit)
	if err != nil {
		return 0, err
	}
	go askQuestions(quiz.quizChan, score, errors)
	for !timesUp && (numAnswers < quiz.NumQuestions) {
		select {
		case <-time.After(time.Duration(quiz.timeLimit) * time.Second):
			timesUp = true
		case i := <-score:
			numCorrect += i
			numAnswers++
		case err = <-errors:
			return 0, err
		}
	}
	return numCorrect, nil
}

func promptToStart(limit int) error {
	reader := bufio.NewReader(os.Stdin)
	_, err := util.GetInput(reader, fmt.Sprintf("Press ENTER to start the quiz. You'll have %d seconds to complete it.", limit))
	return err
}

func askQuestions(quizChanR <-chan question, score chan<- int, errors chan<- error) {
	reader := bufio.NewReader(os.Stdin)
	for q := range quizChanR {
		answer, err := util.GetInput(reader, fmt.Sprintf("Question: %s. Your answer: ", q.prompt))
		if err != nil {
			errors <- err
			return
		}
		var i int
		if normalize(answer) == q.answer {
			i = 1
		}
		score <- i
	}
}
