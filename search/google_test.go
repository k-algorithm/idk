package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockCollector struct{}

func (m MockCollector) Visit(url string) error {
	return nil
}
func (m MockCollector) Wait() {
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
	var mockCollector Collector = MockCollector{}
	idBuffer := []string{"1", "2", "3"}
	result := googleSearch(
		"test", mockCollector, &idBuffer, 0, 2, 0,
	)
	assert.Equal(t, result.QuestionIDs, idBuffer[:2], "Invalid QuestionIDs")
	assert.Equal(t, result.NextOffset, 0, "Invalid NextOffset")
	assert.Equal(t, result.NextQuestionIdx, 2, "Invalid NextQuestionIdx")
	assert.Equal(t, result.IsFinished, false, "Invalid IsFinished")
}
