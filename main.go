package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Reply string `json:"response"`
}

func main() {
	botToken := "BOT_TOKEN_HERE"

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	updates, _ := bot.UpdatesViaLongPolling(nil)
	bh, _ := th.NewBotHandler(bot, updates)

	bh.HandleMessage(func(bot *telego.Bot, message telego.Message) {
		go func() {
			if message.Text != "" {
				thinkingMsg, err := bot.SendMessage(tu.Message(
					tu.ID(message.Chat.ID),
					"⌛ Loading...",
				))
				if err != nil {
					fmt.Println("Error sending loading message:", err)
					return
				}

				thinkingMsgID := thinkingMsg.MessageID

				modifiedPrompt := message.Text + "\nAlways respond in English."
				response, err := queryOllamaModel(modifiedPrompt)
				if err != nil {
					_, editErr := bot.EditMessageText(&telego.EditMessageTextParams{
						ChatID:    tu.ID(message.Chat.ID),
						MessageID: thinkingMsgID,
						Text:      "Error connecting to the model. Please try again later.",
					})
					if editErr != nil {
						fmt.Println("Error editing message:", editErr)
					}
					fmt.Println("Error querying Ollama model:", err)
					return
				}

				if response == "" {
					_, editErr := bot.EditMessageText(&telego.EditMessageTextParams{
						ChatID:    tu.ID(message.Chat.ID),
						MessageID: thinkingMsgID,
						Text:      "Model didn't provide a response. Try another prompt.",
					})
					if editErr != nil {
						fmt.Println("Error editing message:", editErr)
					}
					return
				}

				escapedResponse := escapeMarkdownV2(response)
				_, err = bot.EditMessageText(&telego.EditMessageTextParams{
					ChatID:    tu.ID(message.Chat.ID),
					MessageID: thinkingMsgID,
					Text:      escapedResponse,
					ParseMode: telego.ModeMarkdownV2,
				})
				if err != nil {
					fmt.Println("Error editing message:", err)
				}
			}
		}()
	})

	defer bh.Stop()
	defer bot.StopLongPolling()

	bh.Start()
}

func queryOllamaModel(prompt string) (string, error) {
	url := "http://localhost:11434/api/generate"
	requestBody, _ := json.Marshal(OllamaRequest{
		// model: llama3.2
		Model:  "llama3.2",
		Prompt: prompt,
	})
	start := time.Now()
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	fmt.Println("HTTP request duration:", time.Since(start))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var ollamaResponse OllamaResponse
	err = json.Unmarshal(body, &ollamaResponse)
	if err != nil {
		return "", err
	}

	return ollamaResponse.Reply, nil
}

func escapeMarkdownV2(s string) string {
	replacements := map[string]string{
		"_": "\\_",
		"[": "\\[",
		"]": "\\]",
		"(": "\\(",
		")": "\\)",
		"~": "\\~",
		">": "\\>",
		"#": "\\#",
		"+": "\\+",
		"-": "\\-",
		"=": "\\=",
		"{": "\\{",
		"}": "\\}",
		".": "\\.",
		"!": "\\!",
		"*": "\\*",
	}

	for old, new := range replacements {
		s = strings.ReplaceAll(s, old, new)
	}

	return s
}
