package main

import (
	"bytes"
	"encoding/json"
	"fishpi/logger"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gdamore/tcell/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/rivo/tview"
	"net/url"
	"regexp"
	"strings"
	"testing"
)

const corporate = `["index"][darkcyan]这是一的内容[#b4b4b4][""]
这是[#ff0000]空白内容
["second"][#bbbbbb]这是二的内容[""]`

func TestTView(t *testing.T) {
	app := tview.NewApplication()
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	textView.SetText(corporate)
	var i int
	textView.SetHighlightedFunc(func(added, removed, remaining []string) {
		i++
		text := textView.GetText(false)
		textView.SetText(fmt.Sprintf("%s第%d次点击 added:%s removed:%s remaining:%s", text, i, strings.Join(added, ","), strings.Join(removed, ","), strings.Join(remaining, ",")))
	})
	textView.SetToggleHighlights(false)

	textView.SetBackgroundColor(tcell.ColorDefault)
	if err := app.SetRoot(textView, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func TestDeleteSpan(t *testing.T) {
	text := `12354234<span id="elves123,.asd">Some text hereads</span>asdad`

	// 正则表达式模式
	pattern := `<span\s+[^>]*id\s*=\s*['"]([^'"]+)['"][^>]*>(.*?)</span>`

	// 创建正则表达式对象
	reg := regexp.MustCompile(pattern)

	// 替换匹配的内容为空字符串
	result := reg.ReplaceAllString(text, "")

	fmt.Println(result)
}

func TestIframeWeather(t *testing.T) {
	//str := `<iframe src="https://www.lingmx.com/card/index2.html?date=8/19,8/20,8/21,8/22,8/23&weatherCode=LIGHT_RAIN,LIGHT_RAIN,CLOUDY,CLOUDY,LIGHT_RAIN&max=32,33,35,36,36&min=26,26,26,27,27&t=厦门&st=31分钟后开始下小雨，但56分钟后会停" width="380" height="370" frameborder="0"></iframe>`
	msg := `<iframe src="https://www.lingmx.com/card/index2.html?date=8/19,8/20,8/21,8/22,8/23&weatherCode=LIGHT_RAIN,LIGHT_RAIN,CLOUDY,CLOUDY,LIGHT_RAIN&max=32,33,35,36,36&min=26,26,26,27,27&t=厦门&st=31分钟后开始下小雨，但56分钟后会停" width="380" height="370" frameborder="0"></iframe>`
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(msg))
	if err != nil {
		t.Fatalf("parse error: %s", err)
	}
	dom.Find(`iframe`).Each(func(i int, s *goquery.Selection) {
		src, exist := s.Attr("src")
		if !exist {
			return
		}
		u, e := url.Parse(src)
		if e != nil {
			msg = fmt.Sprintf("parse %s error: %s", src, err)
			return
		}
		msg = u.Query().Get("t") + "天气" + "\n"
		data := [][]string{
			strings.Split(u.Query().Get("weatherCode"), ","),
			strings.Split(u.Query().Get("max"), ","),
			strings.Split(u.Query().Get("min"), ","),
		}

		buffer := bytes.NewBufferString(msg)
		table := tablewriter.NewWriter(buffer)
		table.SetHeader(strings.Split(u.Query().Get("date"), ","))
		table.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})

		for _, v := range data {
			table.Append(v)
		}
		table.Render()
		msg = string(buffer.Bytes())
		msg += u.Query().Get("st")
	})

	t.Log(msg)
}

func TestJsonWeather(t *testing.T) {
	msg := `{"date":"4/16,4/17,4/18","st":"未来24小时多云","min":"15.49,17.49,20.49","msgType":"weather","t":"厦门","max":"25.41,26.49,26.49","weatherCode":"PARTLY_CLOUDY_DAY,CLOUDY,CLOUDY","type":"weather"}`

	weather := new(Weather)
	if err := json.Unmarshal([]byte(msg), weather); err != nil {
		t.Fatalf("parse json error: %s", err)
	}

	weatherCodes := strings.Split(weather.WeatherCode, ",")

	weatherMap := map[string]string{
		"CLEAR_DAY":           "晴",
		"CLEAR_NIGHT":         "晴",
		"PARTLY_CLOUDY_DAY":   "多云 ",
		"PARTLY_CLOUDY_NIGHT": "多云",
		"CLOUDY":              "阴",
		"LIGHT_HAZE":          "轻度雾霾",
		"MODERATE_HAZE":       "中度雾霾",
		"HEAVY_HAZE":          "重度雾霾",
		"LIGHT_RAIN":          "小雨",
		"MODERATE_RAIN":       "中雨",
		"HEAVY_RAIN":          "大雨",
		"STORM_RAIN":          "暴雨",
		"FOG":                 "雾",
		"LIGHT_SNOW":          "小雪",
		"MODERATE_SNOW":       "中雪",
		"HEAVY_SNOW":          "大雪",
		"STORM_SNOW":          "暴雪",
		"DUST":                "浮尘",
		"SAND":                "沙尘",
		"WIND":                "大风",
	}

	var weatherWords []string
	for _, v := range weatherCodes {
		var str string
		if w, ok := weatherMap[v]; ok {
			str = w
		} else {
			str = v
		}
		weatherWords = append(weatherWords, str)
	}

	data := [][]string{
		weatherWords,
		strings.Split(weather.Max, ","),
		strings.Split(weather.Min, ","),
	}

	msg = fmt.Sprintf("%s天气\n", weather.T)
	buffer := bytes.NewBufferString(msg)
	table := tablewriter.NewWriter(buffer)
	table.SetHeader(strings.Split(weather.Date, ","))
	table.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
	msg = string(buffer.Bytes())
	msg += weather.St

	t.Log(msg)
}

type Weather struct {
	Date        string `json:"date"`
	St          string `json:"st"`
	Min         string `json:"min"`
	MsgType     string `json:"msgType"`
	T           string `json:"t"`
	Max         string `json:"max"`
	WeatherCode string `json:"weatherCode"`
	Type        string `json:"type"`
}

//{
//    "date": "4/16,4/17,4/18",
//    "st": "未来24小时多云",
//    "min": "15.49,17.49,20.49",
//    "msgType": "weather",
//    "t": "厦门",
//    "max": "25.41,26.49,26.49",
//    "weatherCode": "PARTLY_CLOUDY_DAY,CLOUDY,CLOUDY",
//    "type": "weather"
//}

// {
//    "coverURL": "http://p1.music.126.net/uUbc7XBIK_GRHWGwxdGtSw==/109951168750937805.jpg",
//    "msgType": "music",
//    "from": "",
//    "source": "http://music.163.com/song/media/outer/url?id=110043",
//    "title": "单身情歌",
//    "type": "music"
//}

/*

   let CodeMap = {
       CLEAR_DAY: "晴",
       CLEAR_NIGHT: "晴",
       PARTLY_CLOUDY_DAY: "多云 ",
       PARTLY_CLOUDY_NIGHT: "多云",
       CLOUDY: "阴",
       LIGHT_HAZE: "轻度雾霾",
       MODERATE_HAZE: "中度雾霾",
       HEAVY_HAZE: "重度雾霾",
       LIGHT_RAIN: "小雨",
       MODERATE_RAIN: "中雨",
       HEAVY_RAIN: "大雨",
       STORM_RAIN: "暴雨",
       FOG: "雾",
       LIGHT_SNOW: "小雪",
       MODERATE_SNOW: "中雪",
       HEAVY_SNOW: "大雪",
       STORM_SNOW: "暴雪",
       DUST: "浮尘",
       SAND: "沙尘",
       WIND: "大风",
   }
*/

func TestLogger(t *testing.T) {
	l := logger.New()
	l.Debug("问题出现了")
}
