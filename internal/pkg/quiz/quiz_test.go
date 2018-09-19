package quiz

import "testing"

const testFile = "testdata/test.csv"
const noSuchFile = "testdata/nosuchfile.csv"

func TestNormalize(t *testing.T) {
	var cases = []struct {
		input    string
		expected string
	}{
		{"foo", "foo"},
		{"BAR", "bar"},
		{"   BaZ  ", "baz"},
	}
	for _, c := range cases {
		result := normalize(c.input)
		if c.expected != result {
			t.Errorf("Fail: expected %s, got %s", c.expected, result)
		}
	}
}

func TestLoadQuestions(t *testing.T) {
	var cases = []struct {
		file        string
		expectedLen int
	}{
		{testFile, 6},
	}
	for _, c := range cases {
		questions, err := loadQuestions(c.file)
		if err != nil {
			t.Errorf("Fail: unexpected error: %v", err)
		}
		resultLen := len(questions)
		if c.expectedLen != resultLen {
			t.Errorf("Fail: expected length %d, got %d", c.expectedLen, resultLen)
		}
	}
	_, err := loadQuestions(noSuchFile)
	if err == nil {
		t.Error("Fail: Expected error because file doesn't exist, but err was nil")
	}
}

func TestNewQuiz(t *testing.T) {
	var cases = []struct {
		file                 string
		expectedNumQuestions int
	}{
		{testFile, 6},
	}
	for _, c := range cases {
		quiz, err := NewQuiz(c.file, 5, false)
		if err != nil {
			t.Errorf("Fail: unexpected error: %v", err)
		}
		if c.expectedNumQuestions != quiz.NumQuestions {
			t.Errorf("Fail: expected %d questions, got %d", c.expectedNumQuestions, quiz.NumQuestions)
		}
		lenChan := len(quiz.quizChan)
		if c.expectedNumQuestions != lenChan {
			t.Errorf("Fail: expected %d questions in channel, got %d", c.expectedNumQuestions, lenChan)
		}
		_, err = NewQuiz(noSuchFile, 10, false)
		if err == nil {
			t.Error("Fail: Expected error because file doesn't exist, but err was nil")
		}
	}
}

func TestNewQuizShuffled(t *testing.T) {
	// Conceivably, the shuffling could preserve the original order. In that case, the test would
	// fail, but it would be a false negative. To make that risk miniscule, run a few tries to give
	// the shuffling more chances.
	var cases = []struct {
		file  string
		tries int
	}{
		{testFile, 5},
	}
	for _, c := range cases {
		unshuffledQuestions, _ := loadQuestions(c.file)
		for i := 0; i < c.tries; i++ {
			quiz, _ := NewQuiz(c.file, 30, true)
			for _, unshuf := range unshuffledQuestions {
				shuf := <-quiz.quizChan
				if unshuf.prompt != shuf.prompt {
					return // Order is different, so shuffling worked
				}
			}
		}
	}
	t.Error("Fail: Shuffled questions were in same order as unshuffled questions")
}
