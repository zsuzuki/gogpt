package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type requestData struct {
	Prompt string `json:"prompt"`
	// 他のオプションやパラメータを設定
}

func downloadImage(url string, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	imgData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, imgData, 0644)
	if err != nil {
		return err
	}

	return nil
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

	generatedImages := result["data"].([]interface{})
	for i, img := range generatedImages {
		generatedImageURL := img.(map[string]interface{})["url"].(string)
		fmt.Printf("生成された画像 %d のURL：%s\n", i+1, generatedImageURL)

		dateStr := time.Now().Format("2006-01-02-15-04")
		filename := fmt.Sprintf("image-%s-%d.png", dateStr, i+1)

		err = downloadImage(generatedImageURL, filename)
		if err != nil {
			fmt.Printf("画像 %d のダウンロードに失敗しました。%v\n", i+1, err)
			continue
		}

		fmt.Printf("画像 %d がダウンロードされ、ファイル名 %s で保存されました。\n", i+1, filename)
	}
}
