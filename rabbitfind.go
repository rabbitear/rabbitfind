package main

// TODO:
// - Put the suggestions in labels under the rabbitEntry
// - Use a regexp instead of all those strings.ReplaceAll's
// BUG:
// - could have less than sugslice[:10], in that for range loop.
import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var baseUrl = "http://suggestqueries.google.com/complete/search?output=firefox&hl=en&q="
var httpClient = &http.Client{Timeout: 10 * time.Second}
var IsPrintable = regexp.MustCompile(`^[[:print:]]+$`).MatchString

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

	box := container.NewVBox()
	box.Add(widget.NewLabel("Press ESC to search"))

	rabbitFindEntry := newRabbitEntry()
	rabbitFindEntry.SetPlaceHolder("Type stuff here...")

	rabbitFindEntry.OnChanged = func(s string) {
		if !IsPrintable(s) {
			return
		}
		s = strings.ReplaceAll(s, " ", "+")
		res, err := httpClient.Get(baseUrl + s)
		if err != nil {
			log.Println(err)
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err.Error())
		}
		tmpstr := string(body)
		tmpstr = strings.ReplaceAll(tmpstr, "\"", "")
		tmpstr = strings.ReplaceAll(tmpstr, "[", "")
		tmpstr = strings.ReplaceAll(tmpstr, "]", "")
		sugslice := strings.Split(tmpstr, ",")

		fmt.Printf("\n\n -=input> %s  -=len(sugslice)> %d\n", s, len(sugslice))
		for i, value := range sugslice[:10] {
			fmt.Printf("i=%d:%s  ", i, value)
		}
	}

	box.Add(rabbitFindEntry)

	w.SetContent(box)
	w.ShowAndRun()
}
