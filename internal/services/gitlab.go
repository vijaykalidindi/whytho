package services

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/vinamra28/operator-reviewer/internal/models"
	"github.com/xanzy/go-gitlab"
)

type GitLabService struct {
	client *gitlab.Client
}

func NewGitLabService(token, baseURL string) *GitLabService {
	logrus.WithField("base_url", baseURL).Info("Creating GitLab client")
	git, err := gitlab.NewClient(token, gitlab.WithBaseURL(baseURL))
	if err != nil {
		logrus.WithError(err).WithField("base_url", baseURL).Fatal("Failed to create GitLab client")
		panic(fmt.Sprintf("Failed to create GitLab client: %v", err))
	}

	logrus.Info("GitLab client created successfully")
	return &GitLabService{
		client: git,
	}
}

func (g *GitLabService) GetMRChanges(projectID, mrIID int) ([]models.MRChange, error) {
	logrus.WithFields(logrus.Fields{
		"project_id": projectID,
		"mr_iid": mrIID,
	}).Debug("Fetching merge request changes")

	changes, _, err := g.client.MergeRequests.GetMergeRequestChanges(projectID, mrIID, nil)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"project_id": projectID,
			"mr_iid": mrIID,
		}).Error("Failed to fetch merge request changes from GitLab API")
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

	logrus.WithFields(logrus.Fields{
		"project_id": projectID,
		"mr_iid": mrIID,
		"changes_count": len(mrChanges),
	}).Debug("Successfully fetched merge request changes")

	return mrChanges, nil
}

func (g *GitLabService) PostMRComment(projectID, mrIID int, comment string) error {
	logrus.WithFields(logrus.Fields{
		"project_id": projectID,
		"mr_iid": mrIID,
	}).Debug("Posting comment to merge request")

	note := &gitlab.CreateMergeRequestNoteOptions{
		Body: &comment,
	}

	_, _, err := g.client.Notes.CreateMergeRequestNote(projectID, mrIID, note)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"project_id": projectID,
			"mr_iid": mrIID,
		}).Error("Failed to post comment to GitLab")
		return fmt.Errorf("failed to post MR comment: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"project_id": projectID,
		"mr_iid": mrIID,
	}).Debug("Comment posted successfully to merge request")

	return nil
}

func (g *GitLabService) GetMRDetails(projectID, mrIID int) (*gitlab.MergeRequest, error) {
	logrus.WithFields(logrus.Fields{
		"project_id": projectID,
		"mr_iid": mrIID,
	}).Debug("Fetching merge request details")

	mr, _, err := g.client.MergeRequests.GetMergeRequest(projectID, mrIID, nil)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"project_id": projectID,
			"mr_iid": mrIID,
		}).Error("Failed to fetch merge request details from GitLab API")
		return nil, fmt.Errorf("failed to get MR details: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"project_id": projectID,
		"mr_iid": mrIID,
		"mr_title": mr.Title,
	}).Debug("Successfully fetched merge request details")

	return mr, nil
}