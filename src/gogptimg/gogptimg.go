package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type requestData struct {
	Prompt string `json:"prompt"`
	Num int `json:"n"`
	Size string `json:"size"`
	// 他のオプションやパラメータを設定
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("プロンプトを引数として入力してください。")
		return
	}

	prompt := os.Args[1]
	apiKey, exists := os.LookupEnv("OPENAI_API_KEY")
	if !exists {
		fmt.Println("環境変数 'OPENAI_API_KEY' が設定されていません。")
		return
	}

	headers := map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", apiKey)},
		"Content-Type":  {"application/json"},
	}

	data := requestData{
		Prompt: prompt,
		Num: 2,
		Size: "1024x1024",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("JSON変換に失敗しました。", err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/images/generations", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("リクエストの作成に失敗しました。", err)
		return
	}

	req.Header = headers
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("リクエストの送信に失敗しました。", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("レスポンスの読み取りに失敗しました。", err)
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("JSON解析に失敗しました。", err)
		return
	}

	generatedImage := result["data"]
	fmt.Println("生成された画像：", generatedImage)
}

