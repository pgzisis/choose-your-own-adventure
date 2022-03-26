package story

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
)

type story map[string]chapter

type chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"paragraphs"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type storyHandler struct {
	story story
}

// NewStoryHandler creates the story handler
func NewStoryHandler() *storyHandler {
	jsonFileName := parseJsonFlag()
	file := openFile(*jsonFileName)
	story := decodeJson(file)

	return &storyHandler{story: story}
}

func parseJsonFlag() *string {
	jsonFileName := flag.String("json", "gopher.json", "json file name")
	flag.Parse()

	return jsonFileName
}

func openFile(name string) *os.File {
	file, err := os.Open("gopher.json")
	if err != nil {
		log.Fatalln(err)
	}

	return file
}

func decodeJson(file *os.File) story {
	decoder := json.NewDecoder(file)
	var story story
	err := decoder.Decode(&story)
	if err != nil {
		log.Fatalln(err)
	}

	return story
}

func (storyHandler *storyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	chapterName := getChapterName(r)

	if chapter, ok := storyHandler.story[chapterName]; ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		template, err := template.ParseFiles("template.html")
		if err != nil {
			log.Fatalln(err)
		}

		template.Execute(w, chapter)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func getChapterName(r *http.Request) string {
	path := r.URL.Path
	chapterName := path[1:]
	if chapterName == "" {
		chapterName = "intro"
	}

	return chapterName
}
