package ui

import (
	"fmt"
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shurcooL/githubv4"
	"github.com/skanehira/ght/config"
	"github.com/skanehira/ght/github"
)

type Issue struct {
	Number    int
	State     string
	Author    string
	Title     string
	Labels    []string
	Assigness []string
}

type issueUI struct {
	issues  []Issue
	updater chan<- func()
	*tview.Table
}

func newIssueUI(updater chan<- func()) *issueUI {
	ui := &issueUI{
		Table:   tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(0, 1),
		updater: updater,
	}

	ui.SetTitle("issue list").SetTitleAlign(tview.AlignLeft)
	ui.updateIssueList()
	ui.SetBorder(true)
	return ui
}

func (ui *issueUI) updateIssueList() {
	table := ui.Clear()

	headers := []string{
		"Number",
		"State",
		"Title",
		"Author",
		"Labels",
		"Assignees",
	}

	for i, header := range headers {
		table.SetCell(0, i, &tview.TableCell{
			Text:            header,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tcell.ColorWhite,
			BackgroundColor: tcell.ColorDefault,
			Attributes:      tcell.AttrBold,
		})
	}

	go func() {
		ui.updater <- func() {
			v := map[string]interface{}{
				"owner":  githubv4.String(config.GitHub.Owner),
				"name":   githubv4.String(config.GitHub.Repo),
				"first":  githubv4.Int(25),
				"cursor": (*githubv4.String)(nil),
			}
			resp, err := github.GetIssue(v)
			if err != nil {
				log.Println(err)
				return
			}

			ui.issues = make([]Issue, len(resp.Nodes))
			for i, node := range resp.Nodes {
				issue := Issue{
					Number: int(node.Number),
					State:  string(node.State),
					Author: string(node.Author.Login),
					Title:  string(node.Title),
				}

				labels := make([]string, len(node.Labels.Nodes))
				for i, l := range node.Labels.Nodes {
					labels[i] = string(l.Name)
				}
				issue.Labels = labels

				assigness := make([]string, len(node.Assignees.Nodes))
				for i, a := range node.Assignees.Nodes {
					assigness[i] = string(a.Login)
				}
				issue.Assigness = assigness

				ui.issues[i] = issue
			}

			for i, issue := range ui.issues {
				table.SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf("#%d", issue.Number)).
					SetTextColor(tcell.ColorBlue))

				cell := tview.NewTableCell(issue.State)
				if issue.State == "OPEN" {
					cell.SetTextColor(tcell.ColorGreen)
				} else {
					cell.SetTextColor(tcell.ColorRed)
				}
				table.SetCell(i+1, 1, cell)

				table.SetCell(i+1, 2, tview.NewTableCell(issue.Title).
					SetTextColor(tcell.ColorWhite).SetExpansion(1))

				table.SetCell(i+1, 3, tview.NewTableCell(issue.Author).
					SetTextColor(tcell.ColorYellow))

				table.SetCell(i+1, 4, tview.NewTableCell(strings.Join(issue.Labels, ",")).
					SetTextColor(tcell.ColorAqua))

				table.SetCell(i+1, 5, tview.NewTableCell(strings.Join(issue.Assigness, ",")).
					SetTextColor(tcell.ColorOlive))
			}
		}
	}()
}
