package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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
