package search

import (
	"net/url"
	"strconv"
	"strings"
	"unicode"

	"github.com/gocolly/colly/v2"
)

var defaultUserAgent string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:34.0) Gecko/20100101 Firefox/34.0"

type GoogleResult struct {
	QuestionIDs     []string
	NextQuestionIdx int
	NextOffset      int
	IsFinished      bool
}

type GoogleParam struct {
	Query           string
	PageSize        int
	Offset          int
	NextQuestionIdx int
}

func (p *GoogleParam) FillDefaults() {
	if p.PageSize == 0 {
		p.PageSize = 10
	}
}

func parseQuestionID(url string) string {
	// expected url in format - /url?esrc=s&q=&rct=j&sa=U&url=https://....
	rawSlice := strings.Split(url, "url=")
	// unexpected url format
	if len(rawSlice) != 2 {
		return ""
	}
	urlSlice := strings.Split(rawSlice[1], "/")
	// filter out non-question pages
	if urlSlice[3] != "questions" {
		return ""
	}
	// filter out non-question pages (filter out "tagged" or other cases)
	if !unicode.IsDigit(rune(urlSlice[4][0])) {
		return ""
	}
	return urlSlice[4]
}

func buildGoogleUrl(query string, start int) string {
	u := url.URL{
		Scheme: "https",
		Host:   "google.com",
		Path:   "search",
	}
	q := u.Query()
	q.Set("q", "site: stackoverflow.com "+query)
	q.Set("start", strconv.Itoa(start))
	u.RawQuery = q.Encode()
	return u.String()
}

func Google(param GoogleParam) GoogleResult {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: coursera.org, www.coursera.org
		colly.AllowedDomains("google.com", "www.google.com"),
		// Set header
		colly.UserAgent((defaultUserAgent)),
	)

	qidBuffer := make([]string, param.PageSize)
	isFinished := false

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		isFinished = true
	})
	// Extract details of the course
	c.OnHTML(`h3`, func(e *colly.HTMLElement) {
		p := e.DOM.Parent()
		link, ok := p.Attr("href")
		qid := ""
		if ok {
			qid = parseQuestionID(link)
		}
		qidBuffer[e.Index] = qid
	})

	offset := param.Offset
	questionIDs := make([]string, 0, param.PageSize)
	nextQuestionIdx := 0
outerLoop:
	for (len(questionIDs) < param.PageSize) && (!isFinished) {
		url := buildGoogleUrl(param.Query, offset)
		// Start scraping on google search
		c.Visit(url)
		c.Wait()
		for i, qid := range qidBuffer {
			if (offset == param.Offset) && (i < param.NextQuestionIdx) {
				continue
			}
			if qid != "" {
				isFinished = false
				if len(questionIDs) >= param.PageSize {
					nextQuestionIdx = i
					break outerLoop
				}
				questionIDs = append(questionIDs, qid)
			}
		}
		// clear buffer
		for i := 0; i < param.PageSize; i++ {
			qidBuffer[i] = ""
		}
		offset += param.PageSize
	}

	return GoogleResult{
		QuestionIDs:     questionIDs,
		NextQuestionIdx: nextQuestionIdx,
		NextOffset:      offset,
		IsFinished:      isFinished,
	}
}
