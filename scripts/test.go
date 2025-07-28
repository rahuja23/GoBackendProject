package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

func updatePost(postID int, p UpdatePostPayload, wg *sync.WaitGroup) {
	defer wg.Done()
	url := fmt.Sprintf("http://localhost:8000/v1/posts/%d", postID)

	// Create the JSON Payload
	b, _ := json.Marshal(p)
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("Error Creating Request: ", err)

	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error Updating Post: ", err)
	}
	defer resp.Body.Close()
	fmt.Println("Update Post response Status:", resp.Status)
}

func main() {

	var wg sync.WaitGroup
	postID := 2
	wg.Add(2)
	title := "NEW TITLE FROM USER A"
	content := "NEW CONTENT FROM USER B"
	go updatePost(postID, UpdatePostPayload{Title: &title}, &wg)
	time.Sleep(1 * time.Second)
	go updatePost(postID, UpdatePostPayload{Content: &content}, &wg)
	wg.Wait()
}
