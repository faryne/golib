package main

import (
	"encoding/json"
	"fmt"
	"github.com/faryne/golib/services/htmlparser"
	"net/http"
)

func main() {
	// setup parse rules
	var rules = []htmlparser.Rule{
		htmlparser.Rule{
			Identifier: "Book",
			Selector:   "div.bwbookitem",
			IsRepeated: true,
			Output: htmlparser.RuleOutput{
				Property: "",
				Target:   "html",
				//Callback: func(obj string) interface{} {
				//	return []interface{}{
				//		obj,
				//		strings.Index(obj, "魔王"),
				//	}
				//},
			},
			Children: []htmlparser.Rule{
				htmlparser.Rule{
					Identifier: "BookName",
					Selector:   "h4.bookname",
					IsRepeated: false,
					Output: htmlparser.RuleOutput{
						Property: "",
						Target:   "html",
						Callback: func(obj string) interface{} {
							return obj
						},
					},
				},
				htmlparser.Rule{
					Identifier: "CoverImage",
					Selector:   "div.bwbookcover > img",
					IsRepeated: false,
					Output: htmlparser.RuleOutput{
						Property: "data-src",
						Target:   "attr",
						Callback: func(obj string) interface{} {
							return obj
						},
					},
				},
				htmlparser.Rule{
					Identifier: "Authors",
					Selector:   "h5.booknamesub",
					IsRepeated: false,
					Output: htmlparser.RuleOutput{
						Property: "",
						Target:   "text",
						Callback: func(obj string) interface{} {
							return obj
						},
					},
				},
			},
		},
	}
	// initialize html parser
	parser := htmlparser.New(rules)
	// initialize http request
	req, err1 := http.NewRequest(http.MethodGet, "https://www.bookwalker.com.tw/more/fiction/1/3", nil)
	if err1 != nil {
		fmt.Println("err1: " + err1.Error())
		return
	}
	result, err := parser.Crawl(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//fmt.Printf("%+v\n", result)
	r, _ := json.Marshal(result)
	fmt.Println(string(r))
}
