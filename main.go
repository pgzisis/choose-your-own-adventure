package main

import (
	"fmt"
	"net/http"

	"github.com/pgzisis/choose-your-own-adventure/story"
)

func main() {
	http.Handle("/", story.NewStoryHandler())

	fmt.Println("Server listening...")
	http.ListenAndServe(":3000", nil)
}
