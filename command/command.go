package command

import (
	"fmt"
	"github.com/alexeykirinyuk/take-smaller-tasks-tool/config"
	"github.com/alexeykirinyuk/take-smaller-tasks-tool/history"
	"github.com/alexeykirinyuk/take-smaller-tasks-tool/jira"
	"time"
)

func Execute() (*history.History, error) {
	c, err := config.Get()
	if err != nil {
		return nil, err
	}

	jiraService, err := jira.CreateService(&c)
	if err != nil {
		return nil, err
	}

	largeCount, err := jiraService.GetCountByEstimations(c.Large)
	if err != nil {
		return nil, err
	}

	smallCount, err := jiraService.GetCountByEstimations(c.Small)
	if err != nil {
		return nil, err
	}

	allCount, err := jiraService.GetCountOfAllEstimatedIssues()
	if err != nil {
		return nil, err
	}

	withWarning := smallCount+largeCount != allCount

	now := time.Now()

	h, err := history.Get()
	if err != nil {
		return nil, fmt.Errorf("error when get info from history: %w", err)
	}

	h.Items = append(h.Items, history.HistoryItem{
		Date:        now,
		LargeCount:  largeCount,
		SmallCount:  smallCount,
		AllCount:    allCount,
		HasWarnings: withWarning,
	})

	h = history.Justify(h)
	err = history.Save(h)
	if err != nil {
		return nil, err
	}

	return h, nil
}
