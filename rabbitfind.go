package main

// TODO:
// - Put the suggestions in labels under the rabbitSelectEntry
// - Use a regexp instead of all those strings.ReplaceAll's
// BUG:
// - [fixed] could have less than sugslice[:10], in that for range loop.
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

type rabbitSelectEntry struct {
	widget.SelectEntry
}

func (e *rabbitSelectEntry) onEsc() {
	fmt.Println(e.SelectEntry.Text)
	u, err := url.Parse("https://google.com/search?q=")
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	q.Set("q", e.SelectEntry.Text)
	u.RawQuery = q.Encode()
	fmt.Println(u)
	a := app.New()
	a.OpenURL(u)
	e.SelectEntry.SetText("")
}

func newRabbitSelectEntry() *rabbitSelectEntry {
	entry := &rabbitSelectEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *rabbitSelectEntry) TypedKey(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyEscape:
		e.onEsc()
	default:
		e.SelectEntry.TypedKey(key)
	}
}

func main() {
	a := app.New()
	w := a.NewWindow("RabbitFind")

	box := container.NewVBox()
	box.Add(widget.NewLabel("Press ESC to search"))

	rabbitFindSelectEntry := newRabbitSelectEntry()
	rabbitFindSelectEntry.SetPlaceHolder("Type stuff here...")

	rabbitFindSelectEntry.OnChanged = func(s string) {
		if !IsPrintable(s) {
			return
		}
		// container for suggestions
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
		sugslice := strings.FieldsFunc(tmpstr, func(r rune) bool {
			switch r {
			case '"':
				return true
			case '[':
				return true
			case ']':
				return true
			case ',':
				return true
			default:
				return false
			}
		})

		fmt.Printf("\n\n -=input> %s  -=len(sugslice)> %d\n", s, len(sugslice))
		if len(sugslice) > 10 {
			sugslice = sugslice[:10]
		}
		rabbitFindSelectEntry.SetOptions(sugslice)
		for i, value := range sugslice {
			fmt.Printf("i=%d:%s  ", i, value)
			if i == 10 {
				break
			}
		}
	}
	box.Add(rabbitFindSelectEntry)

	w.Resize(fyne.NewSize(500, 420))
	w.SetContent(box)
	w.ShowAndRun()
}
