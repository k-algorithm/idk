package main

import (
	"log"

	"github.com/k-algorithm/idk/search"
)

func main() {
	result := search.Google(search.GoogleParam{
		Query:    "golang",
		PageSize: 10,
	})
	if len(result.QuestionIDs) == 0 {
		log.Println("No results..")
		return
	}
	questions := search.Questions(result.QuestionIDs)
	log.Println(questions)
	qid := result.QuestionIDs[0]
	log.Println("QuestionID:", qid)
	question := search.QuestionDetail(qid)
	log.Println("Title:", question.Title)
	log.Println("Body:", question.Body)
	answers := search.Answers(qid)

	for i, answer := range answers {
		log.Println("[Answer", i, "]", "score:", answer.Score)
		log.Println(answer.Body)
	}
}
