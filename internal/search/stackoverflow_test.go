package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockClient struct{}

func (client *MockClient) Get(url string, target interface{}) {}

func TestStackOverflowAgentBuildUrl(t *testing.T) {
	agent := StackOverflowAgent{}
	e1 := "https://api.stackexchange.com/a?order=desc&site=stackoverflow&sort=votes"
	assert.Equal(t, e1, agent.BuildUrl("a", false), "Invalid stackoverflow URL")
	e2 := "https://api.stackexchange.com/a?filter=%216VvPDzQ%29xXOrL&order=desc&site=stackoverflow&sort=votes"
	assert.Equal(t, e2, agent.BuildUrl("a", true), "Invalid stackoverflow URL")
}

func TestStackOverflowQuestionsCache(t *testing.T) {
	agent := &StackOverflowAgent{
		client: &MockClient{},
	}
	agent.FillDefaults()
	assert.Equal(t, agent.cache.ItemCount(), 0, "Cache should be empty.")
	questionIdSlice := []string{"123"}
	cacheKey := "getQuestions123false"
	questions := agent.Questions(questionIdSlice)
	assert.Equal(t, agent.cache.ItemCount(), 1, "Cache miss.")
	cachedValue, cached := agent.cache.Get(cacheKey)
	assert.Equal(t, cached, true, "Cache key missing: "+cacheKey)
	assert.Equal(t, cachedValue, questions, "Cache value differs.")
}

func TestStackOverflowAnswersCache(t *testing.T) {
	agent := &StackOverflowAgent{
		client: &MockClient{},
	}
	agent.FillDefaults()
	assert.Equal(t, agent.cache.ItemCount(), 0, "Cache should be empty.")
	questionId := "123"
	cacheKey := "Answers123"
	answers := agent.Answers(questionId)
	assert.Equal(t, agent.cache.ItemCount(), 1, "Cache miss.")
	cachedValue, cached := agent.cache.Get(cacheKey)
	assert.Equal(t, cached, true, "Cache key missing: "+cacheKey)
	assert.Equal(t, cachedValue, answers, "Cache value differs.")
}
