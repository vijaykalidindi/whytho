package services

import (
	"fmt"

	"github.com/vinamra28/operator-reviewer/internal/models"
	"github.com/xanzy/go-gitlab"
)

type GitLabService struct {
	client *gitlab.Client
}

func NewGitLabService(token, baseURL string) *GitLabService {
	git, err := gitlab.NewClient(token, gitlab.WithBaseURL(baseURL))
	if err != nil {
		panic(fmt.Sprintf("Failed to create GitLab client: %v", err))
	}

	return &GitLabService{
		client: git,
	}
}

func (g *GitLabService) GetMRChanges(projectID, mrIID int) ([]models.MRChange, error) {
	changes, _, err := g.client.MergeRequests.GetMergeRequestChanges(projectID, mrIID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get MR changes: %w", err)
	}

	var mrChanges []models.MRChange
	for _, change := range changes.Changes {
		mrChange := models.MRChange{
			OldPath:     change.OldPath,
			NewPath:     change.NewPath,
			AMode:       change.AMode,
			BMode:       change.BMode,
			NewFile:     change.NewFile,
			RenamedFile: change.RenamedFile,
			DeletedFile: change.DeletedFile,
			Diff:        change.Diff,
		}
		mrChanges = append(mrChanges, mrChange)
	}

	return mrChanges, nil
}

func (g *GitLabService) PostMRComment(projectID, mrIID int, comment string) error {
	note := &gitlab.CreateMergeRequestNoteOptions{
		Body: &comment,
	}

	_, _, err := g.client.Notes.CreateMergeRequestNote(projectID, mrIID, note)
	if err != nil {
		return fmt.Errorf("failed to post MR comment: %w", err)
	}

	return nil
}

func (g *GitLabService) GetMRDetails(projectID, mrIID int) (*gitlab.MergeRequest, error) {
	mr, _, err := g.client.MergeRequests.GetMergeRequest(projectID, mrIID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get MR details: %w", err)
	}

	return mr, nil
}