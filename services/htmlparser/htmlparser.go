package htmlparser

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	h "net/http"
	"strings"
)

type ParseResult struct {
	Error   error
	Content interface{}
}

type RuleOutput struct {
	Property string
	Target   string
	Callback func(obj string) interface{}
}

type Rule struct {
	Identifier string     `validate:"required"`
	Selector   string     `validate:"required"`
	IsRepeated bool       `validate:"required"`
	Output     RuleOutput `validate:"required"`
	Children   []Rule
}

type parser struct {
	Rules      []Rule
	HttpClient h.Client
}

var output = make(map[string]interface{}, 0)

// Init Crawler
func New(rules []Rule) *parser {
	return &parser{
		Rules:      rules,
		HttpClient: h.Client{},
	}
}

// Starting to crawl
func (p *parser) Crawl(req *h.Request) (map[string]interface{}, error) {
	// send http request & get response
	resp, err := p.HttpClient.Do(req)
	defer resp.Body.Close()
	// 將指標轉為實體並指定給 body
	body := *resp
	if err != nil {
		return output, err
	}

	// initialize goquery
	doc, err := goquery.NewDocumentFromReader(body.Body)
	if err != nil {
		return output, err
	}
	for _, rule := range p.Rules {
		output[rule.Identifier] = p.parse(rule, doc.Find(rule.Selector))
	}
	return output, err
}
func (p *parser) parse(rule Rule, selection *goquery.Selection) interface{} {
	if rule.IsRepeated == false {
		if len(rule.Children) <= 0 {
			return p.clear(rule, selection)
		}
		var output = make(map[string]interface{})
		if rule.Children != nil {
			for _, r := range rule.Children {
				output[rule.Identifier] = p.parse(r, selection.Find(r.Selector))
			}
		}
		return output
	}

	var tmp = make([]ParseResult, 0)
	var tmp1 = make(map[string][]interface{}, 0)
	selection.Each(func(i int, s *goquery.Selection) {
		if len(rule.Children) <= 0 {
			tmp = append(tmp, p.clear(rule, s))
		} else {
			for _, r := range rule.Children {
				tmp1[r.Identifier] = append(tmp1[r.Identifier], p.parse(r, s.Find(r.Selector)))
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
