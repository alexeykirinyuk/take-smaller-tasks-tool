package command

import (
	"fmt"
	"github.com/alexeykirinyuk/take-smaller-tasks-tool/config"
	"github.com/alexeykirinyuk/take-smaller-tasks-tool/history"
	"github.com/alexeykirinyuk/take-smaller-tasks-tool/jira"
	"github.com/alexeykirinyuk/take-smaller-tasks-tool/notification"
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

	largeIssues, err := jiraService.GetLargeIssues()
	if err != nil {
		return nil, err
	}

	smallCount, err := jiraService.GetSmallIssuesCount()
	if err != nil {
		return nil, err
	}

	allCount, err := jiraService.GetCountOfAllEstimatedIssues()
	if err != nil {
		return nil, err
	}

	withWarning := smallCount+len(largeIssues) != allCount

	now := time.Now()

	h, err := history.Get()
	if err != nil {
		return nil, fmt.Errorf("error when get info from history: %w", err)
	}

	h.Items = append(h.Items, history.HistoryItem{
		Date:        now,
		LargeCount:  len(largeIssues),
		SmallCount:  smallCount,
		AllCount:    allCount,
		HasWarnings: withWarning,
		LargeIssues: largeIssues,
	})

	h = history.Justify(h)

	if c.EmailNotificationsEnabled {
		notificator := notification.CreteNotificator(c.SMTP)

		err = notificator.Notify(h)
		if err != nil {
			return nil, err
		}
	}

	err = history.Save(h)
	if err != nil {
		return nil, err
	}

	return h, nil
}
