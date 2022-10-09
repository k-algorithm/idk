package search

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var defaultFilter = "!6VvPDzQ)xXOrL"

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

func buildStackOverflowUrl(path string, includeBody bool) string {
	u := url.URL{
		Scheme: "https",
		Host:   "api.stackexchange.com",
		Path:   path,
	}
	q := u.Query()
	q.Set("order", "desc")
	q.Set("sort", "votes")
	q.Set("site", "stackoverflow")
	if includeBody {
		q.Set("filter", defaultFilter)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func getJsonResponse(url string, target interface{}) {
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

func getQuestions(questionIDs []string, includeBody bool) []Question {
	questionIDstr := strings.Join(questionIDs, ";")
	path := fmt.Sprintf("/2.2/questions/%s", questionIDstr)
	url := buildStackOverflowUrl(path, includeBody)
	questionsResponse := new(QuestionsResponse)
	getJsonResponse(url, questionsResponse)
	return questionsResponse.Items
}

func QuestionDetail(questionID string) Question {
	questionIDs := []string{questionID}
	questions := getQuestions(questionIDs, true)
	if len(questions) == 0 {
		panic("Something went wrong. Cannot get question.")
	}
	return questions[0]
}

func Questions(questionIDs []string) []Question {
	return getQuestions(questionIDs, false)
}

func Answers(questionID string) []Answer {
	path := fmt.Sprintf("/2.2/questions/%s/answers", questionID)
	url := buildStackOverflowUrl(path, true)
	answersResponse := new(AnswersResponse)
	getJsonResponse(url, answersResponse)
	return answersResponse.Items
}
