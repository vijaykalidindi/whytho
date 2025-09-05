package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/vinamra28/operator-reviewer/internal/models"
	"google.golang.org/api/option"
)

type ReviewService struct {
	client *genai.Client
}

func NewReviewService(apiKey string) *ReviewService {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		panic(fmt.Sprintf("Failed to create Gemini client: %v", err))
	}

	return &ReviewService{
		client: client,
	}
}

func (r *ReviewService) ReviewCode(changes []models.MRChange, title, description string) (*models.CodeReview, error) {
	ctx := context.Background()
	model := r.client.GenerativeModel("gemini-1.5-flash")

	model.SetTemperature(0.1)

	var codeContent strings.Builder
	codeContent.WriteString(fmt.Sprintf("## Merge Request Details\n"))
	codeContent.WriteString(fmt.Sprintf("**Title:** %s\n", title))
	codeContent.WriteString(fmt.Sprintf("**Description:** %s\n\n", description))

	for _, change := range changes {
		if change.DeletedFile {
			continue
		}

		codeContent.WriteString(fmt.Sprintf("## File: %s\n", change.NewPath))
		if change.NewFile {
			codeContent.WriteString("(New file)\n")
		}
		if change.RenamedFile {
			codeContent.WriteString(fmt.Sprintf("(Renamed from: %s)\n", change.OldPath))
		}
		codeContent.WriteString("```diff\n")
		codeContent.WriteString(change.Diff)
		codeContent.WriteString("\n```\n\n")
	}

	prompt := fmt.Sprintf(`You are an expert code reviewer. Please review the following merge request changes and provide:

1. A brief summary of the changes
2. Specific actionable feedback for improvements
3. Identify potential bugs, security issues, or performance problems
4. Suggest best practices if applicable

Please format your response as follows:
- Start with a summary paragraph
- Then provide specific comments, each starting with "COMMENT:"

Here are the changes to review:

%s

Focus on:
- Code quality and maintainability
- Security vulnerabilities
- Performance issues
- Best practices
- Potential bugs
- Documentation needs

Be constructive and specific in your feedback.`, codeContent.String())

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate review: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response generated")
	}

	reviewText := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			reviewText += string(txt)
		}
	}

	return r.parseReview(reviewText), nil
}

func (r *ReviewService) parseReview(reviewText string) *models.CodeReview {
	lines := strings.Split(reviewText, "\n")
	
	var summary strings.Builder
	var comments []string
	
	inSummary := true
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		if strings.HasPrefix(line, "COMMENT:") {
			inSummary = false
			comment := strings.TrimPrefix(line, "COMMENT:")
			comment = strings.TrimSpace(comment)
			if comment != "" {
				comments = append(comments, comment)
			}
		} else if inSummary && line != "" {
			if summary.Len() > 0 {
				summary.WriteString(" ")
			}
			summary.WriteString(line)
		}
	}
	
	return &models.CodeReview{
		Summary:  summary.String(),
		Comments: comments,
	}
}