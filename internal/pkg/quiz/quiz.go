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
	questions, err := loadQuestions(fileName, shuffle)
	if err != nil {
		return nil, err
	}
	quiz := Quiz{NumQuestions: len(questions), timeLimit: timeLimit}
	quiz.quizChan = make(chan question, quiz.NumQuestions)
	// Send all questions to channel ahead of time
	for _, q := range questions {
		quiz.quizChan <- q
	}
	close(quiz.quizChan) // So that range will loop to completion
	return &quiz, nil
}

func loadQuestions(fileName string, shuffle bool) ([]question, error) {
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
	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(questions), func(i, j int) { questions[i], questions[j] = questions[j], questions[i] })
	}
	return questions, nil
}

func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// Run runs a quiz till the user answers all questions, or times out, or there's an error.
// It returns the number of correct answers and an error code.
func (quiz *Quiz) Run() (int, error) {
	err := promptToStart(quiz.timeLimit)
	if err != nil {
		return 0, err
	}
	var numRight int
	var rightAnswer bool
	answers := make(chan bool, quiz.NumQuestions)
	errors := make(chan error)
	go askQuestions(quiz.quizChan, answers, errors)
	for stillQuestions, stillTime := true, true; stillQuestions && stillTime; {
		select {
		case <-time.After(time.Duration(quiz.timeLimit) * time.Second):
			stillTime = false
		case rightAnswer, stillQuestions = <-answers:
			if stillQuestions && rightAnswer {
				numRight++
			}
		case err = <-errors:
			return 0, err
		}
	}
	return numRight, nil
}

func promptToStart(limit int) error {
	reader := bufio.NewReader(os.Stdin)
	_, err := util.GetInput(reader, fmt.Sprintf("Press ENTER to start the quiz. You'll have %d seconds to complete it.", limit))
	return err
}

func askQuestions(quizChanR <-chan question, answers chan<- bool, errors chan<- error) {
	reader := bufio.NewReader(os.Stdin)
	for q := range quizChanR {
		answer, err := util.GetInput(reader, fmt.Sprintf("Question: %s. Your answer: ", q.prompt))
		if err != nil {
			errors <- err
			return
		}
		answers <- normalize(answer) == q.answer
	}
	close(answers) // So that Run() knows we're done
}
