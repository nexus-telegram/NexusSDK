package handler

import (
	"encoding/json"
	"fmt"
	"github.com/nexus-telegram/NexusSDK/httpclient"
	"github.com/nexus-telegram/NexusSDK/types"
	"github.com/nexus-telegram/NexusSDK/utils"
	"go.uber.org/zap"
	"io"
)

type GameDataRequest struct {
	Game     string             `json:"game"`
	Telegram types.TelegramData `json:"telegram"`
	APIKey   string             `json:"api-key"`
	Proxy    types.Proxy        `json:"proxy"`
}

func refreshGameData(client *httpclient.HTTPClient, game string, apiKey string, telegram types.TelegramData, proxyConfig types.Proxy) ([]byte, error) {
	var nexusApiBaseURL = "http://34.95.182.203:1337/api"
	url := fmt.Sprintf("%s/telegram/game-data", nexusApiBaseURL)
	requestBody := GameDataRequest{
		Game:     game,
		Telegram: telegram,
		APIKey:   apiKey,
		Proxy:    proxyConfig,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log := utils.GetLogger()
		log.Error("Failed to marshal request body", zap.Error(err))
		return nil, err
	}
	resp, err := client.Post(url, jsonData)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
