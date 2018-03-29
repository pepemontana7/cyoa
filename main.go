package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"html/template"
)

var GopherStory map[string]Arc

type Arc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type StoryHandler struct {
	Story map[string]Arc
}

func (sh StoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("story.html")
	check(err)
	var k string
	if len(r.URL.Path) == 1 {
		k = "intro"
	} else {
		k = r.URL.Path[1:]
	}
	t.Execute(w, sh.Story[k])
	//fmt.Fprintf(w, "hello, you've hit %s\n", r.URL.Path)
	//fmt.Fprintf(w, "Story intro: %s\n", sh.Story["intro"])
}
func main() {

	storyFile := flag.String("story", "gopher.json", "Gopher json stories")
	flag.Parse()

	st, err := os.Open(*storyFile)
	check(err)
	defer st.Close()

	stData, err := ioutil.ReadAll(st)
	check(err)

	err = json.Unmarshal(stData, &GopherStory)
	check(err)
	//Gets fields in a Struct
	/*for k, v := range GopherStory {
		fmt.Println(k)
		val := reflect.ValueOf(&v).Elem()
		for i := 0; i < val.NumField(); i++ {
			typeField := val.Type().Field(i)
			fmt.Printf(" Field Name: %s ", typeField.Name)
		}
	}*/

	h := http.NewServeMux()
	h.Handle("/", StoryHandler{GopherStory})
	for k := range GopherStory {
		fmt.Printf("adding handle for %s \n", k)
		h.Handle("/"+k, StoryHandler{GopherStory})
	}

	err = http.ListenAndServe(":9999", h)
	check(err)

}
