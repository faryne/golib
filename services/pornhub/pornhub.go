package pornhub

import (
	"github.com/faryne/golib/services/htmlparser"
	"net/http"
	"net/url"
)

type PornhubStruct struct {
	BaseUrl string
	Rules   []htmlparser.Rule
}

func New() *PornhubStruct {
	rules := []htmlparser.Rule{
		htmlparser.Rule{
			Identifier: "Videos",
			Selector:   ".pcVideoListItem.js-pop.videoblock.videoBox",
			IsRepeated: true,
			Output:     htmlparser.RuleOutput{},
			Children: []htmlparser.Rule{
				htmlparser.Rule{
					Identifier: "ThumbImage",
					Selector:   ".fade.videoPreviewBg.linkVideoThumb.js-linkVideoThumb.img.fadeUp > img",
					IsRepeated: false,
					Output: htmlparser.RuleOutput{
						Property: "data-src",
						Target:   "attr",
					},
				},
				htmlparser.Rule{
					Identifier: "Title",
					Selector:   "span.title > a",
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
					Identifier: "Duration",
					Selector:   ".marker-overlays.js-noFade > var",
					IsRepeated: false,
					Output:     htmlparser.RuleOutput{},
					Children:   nil,
				},
			},
		},
	}
	return &PornhubStruct{
		BaseUrl: "https://cn.pornhub.com/video/search",
		Rules:   rules,
	}
}

func (p *PornhubStruct) Search(keyword string) {
	u := url.Values{}
	u.Set("search", keyword)
	query := u.Encode()

	req, _ := http.NewRequest(http.MethodGet, p.BaseUrl+"?"+query, nil)
	parser := htmlparser.New(p.Rules)
	result, _ := parser.Crawl(req)
}
