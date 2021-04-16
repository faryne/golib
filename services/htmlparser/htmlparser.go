package htmlparser

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/faryne/golib/services/http"
	h "net/http"
	"strings"
)

type ParseResult struct {
	Error   error
	Content interface{}
}
type Response struct {
	Response map[string]interface{}
}

type Rule struct {
	Identifier string `validate:"required"`
	Selector   string `validate:"required"`
	IsRepeated bool   `validate:"required"`
	Output     struct {
		Property string
		Target   string
		Callback func(obj string) (interface{}, error)
	} `validate:"required"`
	Children []Rule
}

type parser struct {
	Rules      []Rule
	HttpClient *http.HttpClient
}

// Init Crawler
func New(rules []Rule) *parser {
	return &parser{
		Rules:      rules,
		HttpClient: http.New(),
	}
}

// Starting to crawl
func (p *parser) Crawl(req h.Request) (Response, error) {
	var output = Response{}
	// send http request & get response
	response, err := p.HttpClient.SendRequest(req)
	response.Header.Get("Content-Type")
	if err != nil {
		return output, err
	}
	// initialize goquery
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return output, err
	}
	for _, rule := range p.Rules {
		output.Response[rule.Identifier] = p.parse(rule, doc.Find(rule.Selector))
	}
	return output, err
}
func (p *parser) parse(rule Rule, selection *goquery.Selection) interface{} {
	if rule.IsRepeated == false {
		if len(rule.Children) <= 0 {
			return p.clear(rule, selection)
		}
		var output = make(map[string]interface{})
		for _, r := range rule.Children {
			output[rule.Identifier] = p.parse(r, selection.Find(r.Selector))
		}
		return output
	}

	var tmp = make([]ParseResult, 0)
	var tmp1 = make(map[string][]interface{}, 0)
	selection.Each(func(i int, s *goquery.Selection) {
		if len(rule.Children) <= 0 {
			tmp = append(tmp, p.clear(rule, selection))
		} else {
			for _, r := range rule.Children {
				tmp1[r.Identifier] = append(tmp1[r.Identifier], p.parse(r, selection.Find(r.Selector)))
			}
		}
	})
	if len(rule.Children) <= 0 {
		return tmp
	}
	return tmp1
}

func (p *parser) clear(rule Rule, selection *goquery.Selection) ParseResult {
	dataType := strings.ToUpper(rule.Output.Property)
	var content string

	switch dataType {
	case "html":
		var err error
		content, err = selection.Html()
		if err != nil {
			return ParseResult{
				Error:   err,
				Content: nil,
			}
		}
	case "attr":
		var isExisted bool
		content, isExisted = selection.Attr(rule.Output.Property)
		if isExisted == false {
			return ParseResult{
				Error:   errors.New(fmt.Sprintf("Property %s not existed")),
				Content: nil,
			}
		}
	case "text":
	default:
		content = selection.Text()
	}

	if rule.Output.Callback != nil {
		return ParseResult{
			Error:   nil,
			Content: rule.Output.Callback(content),
		}
	}
	return ParseResult{
		Error:   nil,
		Content: content,
	}
}
