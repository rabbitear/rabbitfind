package main

import (
	"fmt"
	"log"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

type rabbitEntry struct {
	widget.Entry
}

func (e *rabbitEntry) onEsc() {
	fmt.Println(e.Entry.Text)
	u, err := url.Parse("https://google.com/search?q=")
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	q.Set("q", e.Entry.Text)
	u.RawQuery = q.Encode()
	fmt.Println(u)
	a := app.New()
	a.OpenURL(u)
	e.Entry.SetText("")
}

func newRabbitEntry() *rabbitEntry {
	entry := &rabbitEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *rabbitEntry) TypedKey(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyEscape:
		e.onEsc()
	default:
		e.Entry.TypedKey(key)
	}
}

func main() {
	a := app.New()
	w := a.NewWindow("RabbitFind")

	entry := newRabbitEntry()

	w.SetContent(entry)
	w.ShowAndRun()
}
