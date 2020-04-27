package jira

import (
	"fmt"
	"github.com/alexeykirinyuk/take-smaller-tasks-tool/config"
	"strings"

	jira2 "github.com/andygrunwald/go-jira"
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

func (s *Service) GetCountByEstimations(estimations []string) (int, error) {
	largeTasksJql := s.generateJql(estimations)

	return s.getCountByJql(largeTasksJql)
}

func (s *Service) GetCountOfAllEstimatedIssues() (int, error) {
	return s.getCountByJql(s.config.AllIssuesWithEstimation)
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
