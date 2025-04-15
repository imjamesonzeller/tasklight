package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.design/x/hotkey"
)

type App struct {
	ctx       context.Context
	isVisible bool
}

func NewApp() *App {
	return &App{}
}

//go:embed .env
var envContent string

type Config struct {
	NotionDBID   string
	NotionSecret string
	OpenAIAPIKey string
}

var AppConfig Config

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.isVisible = true

	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//}

	envMap, err := godotenv.Unmarshal(envContent)
	if err != nil {
		log.Println("Error loading embedded .env:", err)
		return
	}

	for k, v := range envMap {
		if err := os.Setenv(k, v); err != nil {
			log.Printf("Could not set %s: %v", k, err)
		}
	}

	AppConfig = Config{
		NotionDBID:   os.Getenv("NOTION_DB_ID"),
		NotionSecret: os.Getenv("NOTION_SECRET"),
		OpenAIAPIKey: os.Getenv("OPENAI_API_KEY"),
	}

	go func() {
		runtime.LockOSThread() // Lock this goroutine to an OS thread (required by macOS for hotkeys)
		a.RegisterHotKey()
	}()
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

type TaskInformation struct {
	Title string  `json:"title"`
	Date  *string `json:"date"`
}

func SendNotionAPIRequest(taskInformation TaskInformation) string {
	// Request Body with Boilerplate and filled in information
	// Build the "Name" property
	nameProp := map[string]interface{}{
		"type": "title",
		"title": []map[string]interface{}{
			{
				"type": "text",
				"text": map[string]interface{}{
					"content": taskInformation.Title,
				},
			},
		},
	}

	// Create the full "properties" object dynamically
	properties := map[string]interface{}{
		"Name": nameProp,
	}

	// If a date is present, add the "Date" property
	if taskInformation.Date != nil {
		properties["Due Date"] = map[string]interface{}{
			"date": map[string]interface{}{
				"start": *taskInformation.Date,
				// "end": nil,         // optional
				// "time_zone": nil,  // optional
			},
		}
	}

	// Create full post body
	postBody := map[string]interface{}{
		"parent": map[string]string{
			"type":        "database_id",
			"database_id": AppConfig.NotionDBID,
		},
		"properties": properties,
	}

	// Convert postBody to JSONObject
	jsonData, err := json.Marshal(postBody)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return err.Error()
	}

	// Request Details
	url := "https://api.notion.com/v1/pages"
	authorization := "Bearer " + AppConfig.NotionSecret
	fmt.Println(authorization)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))

	// Add Headers to Request
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Notion-Version", "2022-06-28")

	// Response from API (Actually send POST)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer response.Body.Close()
	responseStatus := response.Status

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)
	log.Printf(sb)
	return responseStatus
}

func ProcessedMessageFromAI(input string) TaskInformation {
	// TODO: Add handling of No API key and failure of OpenAI request so it just sends input back.

	client := openai.NewClient(option.WithAPIKey(AppConfig.OpenAIAPIKey))

	today := time.Now().Format("2006-01-02") // ISO 8601 format
	prompt := fmt.Sprintf(`You are a helpful task parsing assistant. Your job is to parse natural language
                                  task descriptions into structured data.
                                  Today's date is is %s.
                                  Extract the task title and date from this sentence: "%s".
                                  Return only a JSON object in this exact format: { "title": ..., "date": ... }.
                                  If no date is mentioned, set the "date" value to null.
                                  The date should be in ISO 8601 format when present.`, today, input)

	resp, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4oMini,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
	})
	if err != nil {
		panic(err)
	}

	// Extract and parse JSON
	var task TaskInformation
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &task)
	if err != nil {
		panic(err)
	}

	return task
}

func (a *App) ProcessMessage(message string) {
	// Process message through OpenAi
	taskInformation := ProcessedMessageFromAI(message)

	// Send Request to Notion using received information
	responseStatus := SendNotionAPIRequest(taskInformation)
	if responseStatus != "200 OK" {
		wruntime.EventsEmit(a.ctx, "Backend:ErrorEvent", responseStatus)
	} else { // <-- Successful Request to API
		a.ToggleVisibility()
	}

	return
}

func (a *App) RegisterHotKey() {
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl}, hotkey.KeySpace)

	if err := hk.Register(); err != nil {
		fmt.Printf("Failed to register hotkey: %v\n", err)
		return
	}

	fmt.Println("Hotkey registered!")

	go func() {
		for {
			select {
			case <-hk.Keydown():
				fmt.Println("Hotkey down")
				wruntime.EventsEmit(a.ctx, "Backend:GlobalHotkeyEvent", time.Now().String())
				a.ToggleVisibility()
			case <-hk.Keyup():
				fmt.Println("Hotkey up")
			}
		}
	}()
}

func (a *App) ToggleVisibility() {
	if a.isVisible {
		wruntime.Hide(a.ctx)
	} else {
		wruntime.WindowShow(a.ctx)
		wruntime.WindowSetAlwaysOnTop(a.ctx, true)
	}

	a.isVisible = !a.isVisible
}
