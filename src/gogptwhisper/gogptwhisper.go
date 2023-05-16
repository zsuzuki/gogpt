package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type Response struct {
	Text string `json:"text"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing argument: audio file path")
		os.Exit(1)
	}

	filePath := os.Args[1]
	url := "https://api.openai.com/v1/audio/transcriptions"

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add the file
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fw, err := w.CreateFormFile("file", f.Name())
	if err != nil {
		panic(err)
	}
	if _, err = io.Copy(fw, f); err != nil {
		panic(err)
	}

	// Add the model field
	if err = w.WriteField("model", "whisper-1"); err != nil {
		panic(err)
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can create a post request
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		panic(err)
	}

	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	// Add authorization header
	req.Header.Add("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

	// Do the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	// Check the response
	if res.StatusCode != http.StatusOK {
		panic(res.Status)
	}

	// Parse the response body
	var response Response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		panic(err)
	}

	// Print the "text" field of the response
	fmt.Println(response.Text)
}

