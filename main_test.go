package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"testing"
)

const corporate = `["1"][darkcyan]这是一的内容[#b4b4b4][""]这是空白内容["2"]这是二的内容[""]`

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
		textView.SetText(fmt.Sprintf("%s第%d次点击 added:%s removed:%s remaining:%s", text, i, added, removed, remaining))
	})
	textView.SetToggleHighlights(false)

	textView.SetBackgroundColor(tcell.ColorDefault)
	if err := app.SetRoot(textView, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
