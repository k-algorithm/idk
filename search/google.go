package search

import (
	"net/url"
	"strconv"
	"strings"
	"unicode"

	"github.com/gocolly/colly/v2"
)

const defaultUserAgent string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:34.0) Gecko/20100101 Firefox/34.0"

type Collector interface {
	Visit(url string) error
	Wait()
}

type CollyCollector struct {
	Collector *colly.Collector
}

func (c CollyCollector) Visit(url string) error {
	return c.Collector.Visit(url)
}
func (c CollyCollector) Wait() {
	c.Collector.Wait()
}

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
	UserAgent       string
}

func (p *GoogleParam) FillDefaults() {
	if p.PageSize == 0 {
		p.PageSize = 10
	}
	if p.UserAgent == "" {
		p.UserAgent = defaultUserAgent
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

func newCollector(userAgent string, idBuffer *[]string) *colly.Collector {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: coursera.org, www.coursera.org
		colly.AllowedDomains("google.com", "www.google.com"),
		// Set header
		colly.UserAgent((userAgent)),
	)
	// Extract details of the course
	c.OnHTML(`h3`, func(e *colly.HTMLElement) {
		p := e.DOM.Parent()
		link, ok := p.Attr("href")
		qid := ""
		if ok {
			qid = parseQuestionID(link)
		}
		(*idBuffer)[e.Index] = qid
	})
	return c
}

func googleSearch(
	query string,
	c Collector,
	idBuffer *[]string,
	pageOffset int,
	pageSize int,
	qidOffset int,
) GoogleResult {
	// init variables for pagination
	isFinished := false
	questionIDs := make([]string, 0, pageSize)
	nextQuestionIdx := 0
	currPageOffset := pageOffset

	for (len(questionIDs) < pageSize) && (!isFinished) {
		// clear buffer
		for i := 0; i < pageSize; i++ {
			(*idBuffer)[i] = ""
		}
		url := buildGoogleUrl(query, pageOffset)
		isFinished = true
		isStopped := false
		// Start scraping on google search
		c.Visit(url)
		c.Wait()
		// collect valid question ids
		for i, qid := range *idBuffer {
			if (currPageOffset == pageOffset) && (i < qidOffset) {
				continue
			}
			if qid != "" {
				isFinished = false
				if len(questionIDs) >= pageSize {
					isStopped = true
					nextQuestionIdx = i
					break
				}
				questionIDs = append(questionIDs, qid)
			}
		}
		if !isStopped {
			currPageOffset += pageSize
		}
	}

	return GoogleResult{
		QuestionIDs:     questionIDs,
		NextQuestionIdx: nextQuestionIdx,
		NextOffset:      currPageOffset,
		IsFinished:      isFinished,
	}
}

func Google(param GoogleParam) GoogleResult {
	// fill default params
	param.FillDefaults()

	// init id buffer and get collector
	// init with length 10 (# of google search result)
	idBuffer := make([]string, 10)
	collyCollector := newCollector(param.UserAgent, &idBuffer)
	c := &CollyCollector{Collector: collyCollector}
	return googleSearch(
		param.Query, c, &idBuffer, param.Offset, param.PageSize, param.NextQuestionIdx,
	)
}
