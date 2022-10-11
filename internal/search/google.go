package search

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/gocolly/colly/v2"
)

const defaultUserAgent string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:34.0) Gecko/20100101 Firefox/34.0"
const defaultGooglePageSize int = 10

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
	NextGooglePage  int
	IsFinished      bool
}

type GoogleParam struct {
	Query                 string
	PageSize              int
	GoogleStartPage       int
	GoogleNextQuestionIdx int
	GoogleMaxNumTrial     int
	GooglePageSize        int
	UserAgent             string
	DebugMode             bool
}

func (p *GoogleParam) FillDefaults() {
	if p.PageSize == 0 {
		p.PageSize = 10
	}
	if p.UserAgent == "" {
		p.UserAgent = defaultUserAgent
	}
	if p.GooglePageSize == 0 {
		p.GooglePageSize = defaultGooglePageSize
	}
	if p.GoogleMaxNumTrial == 0 {
		p.GoogleMaxNumTrial = int(p.PageSize/p.GooglePageSize) + 1
	}
}

func parseQuestionID(url string) string {
	// expected url in format - /url?esrc=s&q=&rct=j&sa=U&url=https://....
	rawSlice := strings.Split(url, "url=")
	// unexpected url format
	if len(rawSlice) != 2 {
		return ""
	}
	stackoverflowUrl := rawSlice[1]
	if !strings.HasPrefix(stackoverflowUrl, "https://stackoverflow.com") {
		log.Println("prefix no", stackoverflowUrl)
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

func buildGoogleUrl(query string, start int, num int) string {
	u := url.URL{
		Scheme: "https",
		Host:   "google.com",
		Path:   "search",
	}
	q := u.Query()
	q.Set("q", "site:https://stackoverflow.com/questions "+query)
	q.Set("start", strconv.Itoa(start))
	q.Set("num", strconv.Itoa(num))
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
	c.OnError(func(r *colly.Response, err error) {
		log.Println(err)
	})
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
	userPageSize int,
	googleStartPage int,
	questionIdOffset int,
	googlePageSize int,
	googleMaxNumTrial int,
) GoogleResult {
	// init variables for pagination
	questionIDs := make([]string, 0, userPageSize)
	nextQuestionIdx := 0
	currGooglePage := googleStartPage
	for numTrial := 0; (len(questionIDs) < userPageSize) && (numTrial < googleMaxNumTrial); numTrial += 1 {
		// clear buffer
		for i := 0; i < len(*idBuffer); i++ {
			(*idBuffer)[i] = ""
		}
		// build google url
		url := buildGoogleUrl(query, currGooglePage, googlePageSize)
		isStopped := false
		// Start scraping on google search
		c.Visit(url)
		c.Wait()
		// collect valid question ids
		for i, qid := range *idBuffer {
			if (currGooglePage == googleStartPage) && (i < questionIdOffset) {
				continue
			}
			if qid != "" {
				if len(questionIDs) >= userPageSize {
					isStopped = true
					nextQuestionIdx = i
					break
				}
				questionIDs = append(questionIDs, qid)
			}
		}
		// update current google page if not stopped
		if !isStopped {
			currGooglePage += googlePageSize
		}
	}

	return GoogleResult{
		QuestionIDs:     questionIDs,
		NextQuestionIdx: nextQuestionIdx,
		NextGooglePage:  currGooglePage,
		IsFinished:      len(questionIDs) < userPageSize,
	}
}

func loadVector() GoogleResult {
	jsonFile, err := os.Open("./vects/google_result.json")
	decoder := json.NewDecoder(jsonFile)
	googleResult := new(GoogleResult)
	err = decoder.Decode(googleResult)
	if err != nil {
		panic(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	return *googleResult
}

func Google(param GoogleParam) GoogleResult {
	// fill default params
	param.FillDefaults()

	if param.DebugMode {
		return loadVector()
	}

	// init id buffer and get collector
	// init with length 10 (# of google search result)
	idBuffer := make([]string, param.GooglePageSize)
	collyCollector := newCollector(param.UserAgent, &idBuffer)
	c := &CollyCollector{Collector: collyCollector}
	return googleSearch(
		param.Query,
		c,
		&idBuffer,
		param.PageSize,
		param.GoogleStartPage,
		param.GoogleNextQuestionIdx,
		param.GooglePageSize,
		param.GoogleMaxNumTrial,
	)
}
