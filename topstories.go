package hackernews

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/cmarkh/errs"
	"github.com/cmarkh/go-mail"
)

var (
	urlTopStories = "https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty"

	urlItem   = "https://hacker-news.firebaseio.com/v0/item/"
	urlSuffix = ".json?print=pretty"

	urlComments = "https://news.ycombinator.com/item?id="
)

type Story struct {
	Title       string
	Time        int64 //unix time
	TimeF       time.Time
	ID          int
	URL         string
	URLComments string
	Rank        int
}

// TopStories returns the top n stories from hacker news
func TopStories(num int) (stories []Story, err error) {
	ids, err := getStoryIDs(urlTopStories)
	if err != nil {
		err = errs.WrapErr(err)
		return
	}
	if len(ids) == 0 {
		return nil, errors.New("no stories found on Hacker News :(")
	}

	for i, id := range ids {
		if i >= num { //only get top n stories
			break
		}

		story, err := getStory(id)
		if err != nil {
			return nil, errs.WrapErr(err)
		}

		story.Rank = i + 1

		stories = append(stories, story)
	}

	return
}

func getStoryIDs(url string) (storyIDs []int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		err = errs.WrapErr(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = errs.WrapErr(err)
		return
	}

	err = json.Unmarshal(body, &storyIDs)
	if err != nil {
		err = errs.WrapErr(err)
		return
	}

	return
}

func getStory(id int) (story Story, err error) {
	resp, err := http.Get(urlItem + fmt.Sprint(id) + urlSuffix)
	if err != nil {
		err = errs.WrapErr(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = errs.WrapErr(err)
		return
	}

	err = json.Unmarshal(body, &story)
	if err != nil {
		err = errs.WrapErr(err)
		return
	}

	story.URLComments = urlComments + fmt.Sprint(id)
	story.TimeF = time.Unix(story.Time, 0)

	return
}

// Email emails the top n stories from hacker news
func Email(from mail.Account, to string, num int) (err error) {
	stories, err := TopStories(num)
	if err != nil {
		return errs.WrapErr(err)
	}

	tmpl, err := template.New("Hacker News").Parse(tmplEmail)
	if err != nil {
		err = errs.WrapErr(err)
		return
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, stories)
	if err != nil {
		err = errs.WrapErr(err)
		return
	}

	err = from.Send("Hacker News Top Stories", body.String(), to)
	if err != nil {
		err = errs.WrapErr(err)
		return
	}
	return
}

var tmplEmail = `
<html>
<style>
	a {
		color: black;
		text-decoration: none;
	}
</style>
	<body>
		{{range .}}
			<h3 style="margin-bottom: 0px;"><a href={{.URL}}>{{.Rank}}. {{.Title}}</a></h3>
			<a>{{.TimeF}}</a><a> </a><a href={{.URLComments}}>Comments</a>
		{{end}}
	</body>
</html>
`
