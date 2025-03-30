package gemini

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var (
	ctx    context.Context
	model  *genai.GenerativeModel
	logger *slog.Logger
)

func InitModel(slogger *slog.Logger) {
	logger = slogger
	ctx = context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		logger.Error("Error creating client: %+v\n", "err", err)
		log.Fatal("Error creating client:", err)
	}

	model = client.GenerativeModel("gemini-2.0-flash")
}

func Prompt(prompt string) (*genai.GenerateContentResponse, error) {
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		logger.Error("Error generating content: %+v\n", "err", err)
		log.Printf("Error generating content: %+v\n", err)
		return nil, err
	}
	if resp == nil {
		logger.Error("Response is nil")
		log.Println("Response is nil")
		return nil, nil
	}

	return resp, nil
}

func ExtractSummary(resp *genai.GenerateContentResponse) string {
	if len(resp.Candidates) == 0 {
		return ""
	}
	if len(resp.Candidates[0].Content.Parts) == 0 {
		return ""
	}
	// check if the part is of type genai.Text and return it as a string
	if text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
		return string(text)
	}
	return ""
}
