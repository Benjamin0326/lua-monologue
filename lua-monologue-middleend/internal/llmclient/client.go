package llmclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func CallLLM(id string, prompt string) (string, error) {
	data := map[string]string{
		"id":     id,
		"prompt": prompt,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	fmt.Println("📤 보내는 JSON:", string(jsonData))

	resp, err := http.Post("http://localhost:4321/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println("📦 응답 원문:", string(body)) // ← 이 줄 추가

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	responseText, ok := result["response"].(string)
	if !ok {
		return "", fmt.Errorf("응답 형식이 올바르지 않음")
	}

	return responseText, nil
}
