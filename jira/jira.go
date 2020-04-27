package jira

import (
	"fmt"
	"github.com/alexeykirinyuk/take-smaller-tasks-tool/config"
	jira2 "github.com/andygrunwald/go-jira"
	"strings"
)

type Service struct {
	config *config.Configuration
	client *jira2.Client
}

func CreateService(config *config.Configuration) (service *Service, err error) {
	tp := jira2.BasicAuthTransport{
		Username: config.Login,
		Password: config.Token,
	}

	client, err := jira2.NewClient(tp.Client(), config.Domain)
	if err != nil {
		err = fmt.Errorf("error when create jira client: %s", err)
		return
	}

	service = &Service{client: client, config: config}
	return
}

type Issue struct {
	Key      string
	Assignee string
}

func (i *Issue) String() string {
	res := i.Key
	if i.Assignee != "" {
		res += fmt.Sprintf(" (%s)", i.Assignee)
	}

	return res
}

func (s *Service) GetLargeIssues() ([]*Issue, error) {
	largeTasksJql := s.generateJql(s.config.Large)

	var res []*Issue
	issues, err := s.getIssuesByJql(largeTasksJql)
	if err != nil {
		return nil, err
	}

	for _, issue := range issues {
		resultIssue := &Issue{
			Key: issue.Key,
		}

		if issue.Fields.Assignee != nil {
			resultIssue.Assignee = issue.Fields.Assignee.DisplayName
		}

		res = append(res, resultIssue)
	}

	return res, nil
}

func (s *Service) GetSmallIssuesCount() (int, error) {
	smallTasksJql := s.generateJql(s.config.Small)

	return s.getCountByJql(smallTasksJql)
}

func (s *Service) GetCountOfAllEstimatedIssues() (int, error) {
	return s.getCountByJql(s.config.AllIssuesWithEstimation)
}

func (s *Service) getIssuesByJql(jql string) ([]jira2.Issue, error) {
	issues, _, err := s.client.Issue.Search(jql, &jira2.SearchOptions{})
	if err != nil {
		return []jira2.Issue{}, fmt.Errorf("erro when search issues by filter (%s): %s", jql, err)
	}

	return issues, nil
}

func (s *Service) getCountByJql(jql string) (int, error) {
	_, resp, err := s.client.Issue.Search(jql, &jira2.SearchOptions{MaxResults: 1})
	if err != nil {
		return 0, fmt.Errorf("erro when search issues by filter (%s): %s", jql, err)
	}

	return resp.Total, nil
}

func (s *Service) generateJql(estimations []string) string {
	var withBraces []string
	for _, est := range estimations {
		withBraces = append(withBraces, fmt.Sprintf(`"%s"`, est))
	}

	return fmt.Sprintf(s.config.IssuesByEstimation, strings.Join(withBraces, ","))
}
