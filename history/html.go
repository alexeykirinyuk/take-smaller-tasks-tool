package history

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

type tableItem struct {
	Date        string
	LargeCount  int
	SmallCount  int
	AllCount    int
	LargeIssues string
}

const templateString = `
	<h1>Key Result: Take Smaller Tasks Statistics</h1>
	<table>
		<th>date</th>
		<th>large</th>
		<th>small</th>
		<th>all</th>
		<th>large</th>
		{{range .}}
			<tr>
				<td>{{.Date}}</td>
				<td>{{.LargeCount}}</td>
				<td>{{.SmallCount}}</td>
				<td>{{.AllCount}}</td>
				<td>{{.LargeIssues}}</td>
			</tr>
		{{end}}
	</table>`

func (h *History) Html() (string, error) {
	tableItems := h.getTableItems()

	t := template.New("Take Smaller Tasks Template")
	t, err := t.Parse(templateString)
	if err != nil {
		return "", fmt.Errorf("error when parse template: %s", err)
	}

	var buffer bytes.Buffer
	if err := t.Execute(&buffer, tableItems); err != nil {
		return "", fmt.Errorf("error when generate html: %s", err)
	}

	return buffer.String(), nil
}

func (h *History) getTableItems() []tableItem {
	var items []tableItem
	for _, item := range h.Items {
		var largeIssues []string
		for _, i := range item.LargeIssues {
			largeIssues = append(largeIssues, i.String())
		}

		items = append(items, tableItem{
			Date:        item.Date.Format("02.01.2006"),
			LargeCount:  item.LargeCount,
			SmallCount:  item.SmallCount,
			AllCount:    item.AllCount,
			LargeIssues: strings.Join(largeIssues, ", "),
		})
	}
	return items
}
