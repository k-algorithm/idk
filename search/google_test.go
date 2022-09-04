package search

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockCollector struct {
	buffer      *[]string
	numResult   int
	numReturned int
	isFinished  bool
}

func (m *MockCollector) Visit(url string) error {
	i := 0
	for m.numReturned < m.numResult && i < len(*m.buffer) {
		(*m.buffer)[i] = strconv.Itoa(m.numReturned)
		m.numReturned += 1
		i += 1
	}
	return nil
}
func (m *MockCollector) Wait() {
}

func TestParseQuestionID(t *testing.T) {
	url := "/url?esrc=s&q=&rct=j&sa=U&url=https://stackoverflow.com/questions/24622388/lorem-ipsum"
	assert.Equal(t, "24622388", parseQuestionID(url), "Invalid ID received.")

	url = "/url?esrc=s&q=&rct=j&sa=U&url=https://stackoverflow.com/collectives/go"
	assert.Equal(t, "", parseQuestionID(url), "Invalid ID received.")
}

func TestBuildGoogleUrl(t *testing.T) {
	expected := "https://google.com/search?q=site%3A+stackoverflow.com+golang&start=11"
	start := 11
	query := "golang"
	assert.Equal(t, expected, buildGoogleUrl(query, start), "Invalid google URL")
}

func TestGoogleSearch(t *testing.T) {
	// test1: pagesize: 2, result: 10, not finished
	idBuffer := make([]string, 10)
	mockCollector := &MockCollector{&idBuffer, 10, 0, false}
	expected := []string{"0", "1"}
	result := googleSearch(
		"test", mockCollector, &idBuffer, 0, 2, 0,
	)
	assert.Equal(t, expected, result.QuestionIDs, "Invalid QuestionIDs")
	assert.Equal(t, 0, result.NextOffset, "Invalid NextOffset")
	assert.Equal(t, 2, result.NextQuestionIdx, "Invalid NextQuestionIdx")
	assert.Equal(t, false, result.IsFinished, "Invalid IsFinished")

	// test2: pagesize: 3, result: 3, not finished
	idBuffer = make([]string, 3)
	mockCollector = &MockCollector{&idBuffer, 3, 0, false}
	expected = []string{"0", "1", "2"}
	result = googleSearch(
		"test", mockCollector, &idBuffer, 0, 3, 0,
	)
	assert.Equal(t, expected, result.QuestionIDs, "Invalid QuestionIDs")
	assert.Equal(t, 3, result.NextOffset, "Invalid NextOffset")
	assert.Equal(t, 0, result.NextQuestionIdx, "Invalid NextQuestionIdx")
	assert.Equal(t, false, result.IsFinished, "Invalid IsFinished")

	// test3: pagesize: 4, result: 3, not finished
	idBuffer = make([]string, 4)
	mockCollector = &MockCollector{&idBuffer, 3, 0, false}
	expected = []string{"0", "1", "2"}
	result = googleSearch(
		"test", mockCollector, &idBuffer, 0, 4, 0,
	)
	assert.Equal(t, expected, result.QuestionIDs, "Invalid QuestionIDs")
	assert.Equal(t, 0, result.NextQuestionIdx, "Invalid NextQuestionIdx")
	assert.Equal(t, true, result.IsFinished, "Invalid IsFinished")
}
