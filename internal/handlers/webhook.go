package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vinamra28/operator-reviewer/internal/models"
	"github.com/vinamra28/operator-reviewer/internal/services"
)

type WebhookHandler struct {
	gitlabService *services.GitLabService
	reviewService *services.ReviewService
	webhookSecret string
}

func NewWebhookHandler(gitlabService *services.GitLabService, reviewService *services.ReviewService, webhookSecret string) *WebhookHandler {
	return &WebhookHandler{
		gitlabService: gitlabService,
		reviewService: reviewService,
		webhookSecret: webhookSecret,
	}
}

func (h *WebhookHandler) HandleWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	if h.webhookSecret != "" {
		if !h.verifySignature(body, c.GetHeader("X-Gitlab-Token")) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
			return
		}
	}

	eventType := c.GetHeader("X-Gitlab-Event")
	if eventType != "Merge Request Hook" {
		c.JSON(http.StatusOK, gin.H{"message": "Event ignored"})
		return
	}

	var webhook models.GitLabWebhook
	if err := json.Unmarshal(body, &webhook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse webhook"})
		return
	}

	if webhook.ObjectAttributes.Action != "open" && webhook.ObjectAttributes.Action != "reopen" && webhook.ObjectAttributes.Action != "update" {
		c.JSON(http.StatusOK, gin.H{"message": "Action ignored"})
		return
	}

	go h.processMergeRequest(&webhook)

	c.JSON(http.StatusOK, gin.H{"message": "Webhook received"})
}

func (h *WebhookHandler) verifySignature(body []byte, signature string) bool {
	if h.webhookSecret == "" {
		return true
	}

	mac := hmac.New(sha256.New, []byte(h.webhookSecret))
	mac.Write(body)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func (h *WebhookHandler) processMergeRequest(webhook *models.GitLabWebhook) {
	projectID := webhook.Project.ID
	mrIID := webhook.ObjectAttributes.IID

	changes, err := h.gitlabService.GetMRChanges(projectID, mrIID)
	if err != nil {
		fmt.Printf("Error fetching MR changes: %v\n", err)
		return
	}

	if len(changes) == 0 {
		fmt.Println("No changes found in MR")
		return
	}

	review, err := h.reviewService.ReviewCode(changes, webhook.ObjectAttributes.Title, webhook.ObjectAttributes.Description)
	if err != nil {
		fmt.Printf("Error reviewing code: %v\n", err)
		return
	}

	for _, comment := range review.Comments {
		if err := h.gitlabService.PostMRComment(projectID, mrIID, comment); err != nil {
			fmt.Printf("Error posting comment: %v\n", err)
		}
	}

	if review.Summary != "" {
		summaryComment := fmt.Sprintf("## 🤖 AI Code Review Summary\n\n%s", review.Summary)
		if err := h.gitlabService.PostMRComment(projectID, mrIID, summaryComment); err != nil {
			fmt.Printf("Error posting summary comment: %v\n", err)
		}
	}
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}