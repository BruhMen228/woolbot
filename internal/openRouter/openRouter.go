package openrouter

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

func RequestToOpenRouterAi(apiKey string, model string, providerIgnore string, ctx string, msgQuery string) (*http.Response, error){
	data := map[string]interface{}{
		"model":    model,
		"provider": map[string]interface{}{
			"ignore":	[]string{providerIgnore},
		},
		"messages": []map[string]string{{"role": "system", "content": ctx}, {"role": "user", "content": msgQuery}},
		"stream":   false,
	}

	dataJson, err := json.Marshal(data)
	if err != nil {
		return nil, err

	}

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(dataJson))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + apiKey)

	client := &http.Client{
		Timeout: 90 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}