package history

import (
	"encoding/json"
	"fmt"
	"github.com/alexeykirinyuk/take-smaller-tasks-tool/jira"
	"os"
	"strings"
	"time"
)

const historyFileName = "history.json"

type History struct {
	Items []HistoryItem
}

type HistoryItem struct {
	Date        time.Time
	LargeCount  int
	SmallCount  int
	AllCount    int
	HasWarnings bool
	LargeTasks  []*jira.Issue
}

func Get() (history *History, err error) {
	file, err := os.Open(historyFileName)

	if err != nil {
		if os.IsNotExist(err) {
			return &History{Items: []HistoryItem{}}, nil
		}

		err = fmt.Errorf("error when open history file: %s", err)
		return
	}

	defer func() {
		err = file.Close()
	}()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&history)

	if err != nil {
		err = fmt.Errorf("error when parsing history: %s", err)
		return
	}

	return
}

func Save(history *History) error {
	file, err := os.OpenFile(historyFileName, os.O_TRUNC|os.O_WRONLY, os.ModePerm)

	if err != nil && os.IsExist(err) {
		return fmt.Errorf("error when open history file: %s", err)
	}

	if err != nil && os.IsNotExist(err) {
		file, err = os.Create(historyFileName)
		if err != nil {
			return fmt.Errorf("error when create history file: %s", err)
		}
	}

	defer func() {
		err = file.Close()
	}()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(history)

	if err != nil {
		return fmt.Errorf("error when serialize history: %s", err)
	}

	return nil
}

func Justify(history *History) *History {
	keys := make(map[string]bool)
	var resultItems []HistoryItem
	for _, entry := range history.Items {
		if _, value := keys[entry.Date.Format("2016-01-02")]; !value {
			keys[entry.Date.Format("2016-01-02")] = true
			resultItems = append(resultItems, entry)
		}
	}

	return &History{Items: resultItems}
}

func (h *History) String() string {
	b := strings.Builder{}

	b.WriteString("Date\t\tLarge\\Small\t\tAll Estimated\t\tWith Warnings\t\tLarge Tasks\r\n")
	for _, item := range h.Items {
		var largeIssues []string
		for _, i := range item.LargeTasks {
			largeIssues = append(largeIssues, i.String())
		}

		s := fmt.Sprintf("%s\t%d\\%d\t\t\t%d\t\t\t%t\t\t\t%s\r\n",
			item.Date.Format("02.01.2006"),
			item.LargeCount,
			item.SmallCount,
			item.AllCount,
			item.HasWarnings,
			strings.Join(largeIssues, ", "))
		b.WriteString(s)
	}

	return b.String()
}
