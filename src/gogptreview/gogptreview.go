package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type RequestPayload struct {
	Messages    []Message `json:"messages"`
	Model       string    `json:"model"`
	Temperature float32   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponsePayload struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func main() {
	apiKey, exists := os.LookupEnv("OPENAI_API_KEY")
	if !exists {
		fmt.Println("環境変数 'OPENAI_API_KEY' が設定されていません。")
		return
	}

	apiUrl := "https://api.openai.com/v1/chat/completions"

	text := "以下のコードのバグ、想定できる問題点を見つけてください\n"
	for _, arg := range os.Args[1:] {
		fmt.Println("input file: ", arg)
		inputFile, err := os.Open(arg)
		if err != nil {
			fmt.Println("Error opening input file:", err)
			continue
		}
		defer inputFile.Close()

		scanner := bufio.NewScanner(inputFile)
		for scanner.Scan() {
			text = strings.Join([]string{text, scanner.Text()}, "\n")
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error scanning input file:", err)
		}
	}

	messages := []Message{
		{
			Role:    "system",
			Content: "私はコードのレビュワーです。コードを提示してください。",
		},
		{
			Role:    "user",
			Content: text,
		},
	}

	requestPayload := RequestPayload{
		Model:       "gpt-3.5-turbo",
		Temperature: 0.5,
		Messages:    messages}
	requestBody, err := json.Marshal(requestPayload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var responsePayload ResponsePayload
	err = json.Unmarshal(responseBody, &responsePayload)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	if len(responsePayload.Choices) > 0 {
		responseText := responsePayload.Choices[0].Message.Content
		fmt.Println("ChatGPT response:", responseText)
	} else {
		fmt.Println("No response received.")
	}
}
