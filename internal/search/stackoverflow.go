package search

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

var agent *StackOverflowAgent

// Structs to serialize Stackoverflow data
type Answer struct {
	Score      int    `json:"score,omitempty"`
	Body       string `json:"body,omitempty"`
	IsAccepted bool   `json:"is_accepted,omitempty"`
}

type AnswersResponse struct {
	Items []Answer `json:"items,omitempty"`
}

type Question struct {
	Id          int    `json:"question_id,omitempty"`
	Title       string `json:"title,omitempty"`
	Body        string `json:"body,omitempty"`
	Score       int    `json:"score,omitempty"`
	ViewCount   int    `json:"view_count,omitempty"`
	AnswerCount int    `json:"answer_count,omitempty"`
}

type QuestionsResponse struct {
	Items []Question `json:"items,omitempty"`
}

// interface for http client
type AbstractHttpClient interface {
	Get(string, interface{})
}

// implementation of http client
type HttpClient struct{}

func (client *HttpClient) Get(url string, target interface{}) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(target)
	if err != nil {
		panic(err)
	}
	if err != nil {
		log.Fatal(err)
	}
}

// stackoverflow agent
type StackOverflowAgent struct {
	client AbstractHttpClient
	scheme string // "https"
	host   string // "api.stackexchange.com"
	filter string
	cache  *cache.Cache
}

func (agent *StackOverflowAgent) FillDefaults() {
	if agent.client == nil {
		agent.client = new(HttpClient)
	}
	if agent.scheme == "" {
		agent.scheme = "https"
	}
	if agent.host == "" {
		agent.host = "api.stackexchange.com"
	}
	if agent.filter == "" {
		agent.filter = "!6VvPDzQ)xXOrL"
	}
	if agent.cache == nil {
		agent.cache = cache.New(5*time.Minute, 10*time.Minute)
	}
}

func (agent *StackOverflowAgent) BuildUrl(path string, includeBody bool) string {
	// fill defaults before building url
	agent.FillDefaults()

	u := url.URL{
		Scheme: agent.scheme,
		Host:   agent.host,
		Path:   path,
	}
	q := u.Query()
	q.Set("order", "desc")
	q.Set("sort", "votes")
	q.Set("site", "stackoverflow")
	if includeBody {
		q.Set("filter", agent.filter)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func (agent *StackOverflowAgent) getQuestions(questionIDs []string, includeBody bool) []Question {
	// fill defaults before building url
	agent.FillDefaults()

	questionIDstr := strings.Join(questionIDs, ";")
	cacheKey := "getQuestions" + questionIDstr + strconv.FormatBool(includeBody)
	// find questions from cache
	cachedQuestions, found := agent.cache.Get(cacheKey)
	if found {
		return cachedQuestions.([]Question)
	}
	path := fmt.Sprintf("/2.2/questions/%s", questionIDstr)
	url := agent.BuildUrl(path, includeBody)
	questionsResponse := new(QuestionsResponse)
	agent.client.Get(url, questionsResponse)
	questions := questionsResponse.Items
	// cache questions
	agent.cache.Set(cacheKey, questions, cache.DefaultExpiration)
	return questions
}

func (agent *StackOverflowAgent) QuestionDetail(questionID string) Question {
	// fill defaults before building url
	agent.FillDefaults()
	questionIDs := []string{questionID}
	questions := agent.getQuestions(questionIDs, true)
	if len(questions) == 0 {
		panic("Something went wrong. Cannot get question.")
	}
	return questions[0]
}

func (agent *StackOverflowAgent) Questions(questionIDs []string) []Question {
	return agent.getQuestions(questionIDs, false)
}

func (agent *StackOverflowAgent) Answers(questionID string) []Answer {
	// fill defaults before building url
	agent.FillDefaults()

	cacheKey := "Answers" + questionID
	cachedAnswers, found := agent.cache.Get(cacheKey)
	if found {
		return cachedAnswers.([]Answer)
	}

	path := fmt.Sprintf("/2.2/questions/%s/answers", questionID)
	url := agent.BuildUrl(path, true)
	answersResponse := new(AnswersResponse)
	agent.client.Get(url, answersResponse)
	answers := answersResponse.Items
	agent.cache.Set(cacheKey, answers, cache.DefaultExpiration)
	return answers
}

func getStackOverflowAgent() *StackOverflowAgent {
	if agent == nil {
		agent = new(StackOverflowAgent)
	}
	return agent
}

func QuestionDetail(questionID string) Question {
	agent := getStackOverflowAgent()
	return agent.QuestionDetail(questionID)
}

func Questions(questionIDs []string) []Question {
	agent := getStackOverflowAgent()
	return agent.Questions(questionIDs)
}

func Answers(questionID string) []Answer {
	agent := getStackOverflowAgent()
	return agent.Answers(questionID)
}
