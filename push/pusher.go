package push

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"unsafe"

	"ggball.com/smzdm/file"
	"ggball.com/smzdm/smzdm"
)

// 定义推送者，声明推送方法
type Pusher interface {
	Push(content string, contentType string)
}

type DingPusher struct {
	Token string
}

// 钉钉推送者实现推送方法
func (pusher DingPusher) PushDingDing(params interface{}) {
	Url, err := url.Parse("https://oapi.dingtalk.com/robot/send?access_token=" + pusher.Token)
	if err != nil {
		return
	}

	paramsJson, _ := json.Marshal(params)

	urlPath := Url.String()
	resp, err := http.Post(urlPath, "application/json;charset=utf-8", bytes.NewBuffer([]byte(string(paramsJson))))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}

	//fmt.Println(string(content))
	str := (*string)(unsafe.Pointer(&content)) //转化为string,优化内存
	fmt.Println(*str)

}

// 推送商品到钉钉
func PushProWithDingDing(pro []smzdm.Product, conf file.Config) {
	dingPusher := DingPusher{
		Token: conf.DingdingToken,
	}

	// 需要提前申明数组的容量
	links := make([]Link, len(pro))

	for index, item := range pro {
		link := Link{
			Title:      item.ArticlePrice + "!【" + item.ArticleTitle + "】" + "【什么值得买】" + "\n\r",
			MessageURL: item.ArticleUrl,
			PicURL:     item.ArticlePic,
		}
		links[index] = link
	}

	feedCard := FeedCard{
		Links: links,
	}

	params := DingFeedCardParam{
		MsgType:  "feedCard",
		FeedCard: feedCard,
	}

	dingPusher.PushDingDing(params)
}

// 推送文字到钉钉
func PushTextWithDingDing(resText string, conf file.Config) {
	dingPusher := DingPusher{
		Token: conf.DingdingToken,
	}

	text := Text{
		Content: resText + "【什么值得买】",
	}

	params := DingTextParam{
		MsgType: "text",
		Texts:   text,
	}

	dingPusher.PushDingDing(params)
}
