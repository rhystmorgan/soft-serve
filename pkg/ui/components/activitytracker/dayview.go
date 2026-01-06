package activitytracker

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/soft-serve/pkg/backend"
)

// DayViewMode represents different day view modes
type DayViewMode int

const (
	DayViewModeOverview DayViewMode = iota
	DayViewModeRepository
	DayViewModeCommitDetail
)

// DayViewState holds the state for detailed day view
type DayViewState struct {
	Mode               DayViewMode
	SelectedRepository string
	SelectedCommitIdx  int
	RepositoryList     []string
	CommitsByRepo      map[string][]backend.CommitInfo
	TotalCommits       int
}

// initDayView initializes the day view state for a given date
func (m *Model) initDayView(date time.Time) *DayViewState {
	dateKey := date.Format("2006-01-02")
	activity, exists := m.activity[dateKey]

	state := &DayViewState{
		Mode:          DayViewModeOverview,
		CommitsByRepo: make(map[string][]backend.CommitInfo),
	}

	if !exists || activity.Count == 0 {
		return state
	}

	// Group commits by repository
	for _, commit := range activity.Commits {
		state.CommitsByRepo[commit.Repository] = append(state.CommitsByRepo[commit.Repository], commit)
	}

	// Sort repositories by commit count (descending)
	for repo := range state.CommitsByRepo {
		state.RepositoryList = append(state.RepositoryList, repo)
		// Sort commits within each repository by timestamp
		sort.Slice(state.CommitsByRepo[repo], func(i, j int) bool {
			return state.CommitsByRepo[repo][i].Timestamp.After(state.CommitsByRepo[repo][j].Timestamp)
		})
	}

	sort.Slice(state.RepositoryList, func(i, j int) bool {
		return len(state.CommitsByRepo[state.RepositoryList[i]]) > len(state.CommitsByRepo[state.RepositoryList[j]])
	})

	state.TotalCommits = activity.Count

	return state
}

// renderEnhancedDetailView renders the enhanced detailed day view with repository grouping
func (m *Model) renderEnhancedDetailView() string {
	dayState := m.initDayView(m.selectedDate)

	var content strings.Builder

	// Title with navigation indicators
	title := fmt.Sprintf("Activity for %s", m.selectedDate.Format("Monday, January 2, 2006"))
	content.WriteString(m.common.Styles.ActivityTracker.Title.Render(title))
	content.WriteString("\n\n")

	if dayState.TotalCommits == 0 {
		content.WriteString(m.common.Styles.NoContent.Render("No commits on this day"))
		content.WriteString("\n\n")
		content.WriteString(m.renderDayViewHelp())
		return m.common.Styles.ActivityTracker.Container.Render(content.String())
	}

	switch dayState.Mode {
	case DayViewModeOverview:
		content.WriteString(m.renderDayOverview(dayState))
	case DayViewModeRepository:
		content.WriteString(m.renderRepositoryDetail(dayState))
	case DayViewModeCommitDetail:
		content.WriteString(m.renderCommitDetail(dayState))
	}

	content.WriteString("\n")
	content.WriteString(m.renderDayViewHelp())

	return m.common.Styles.ActivityTracker.Container.Render(content.String())
}

// renderDayOverview renders the overview of all repositories for the day
func (m *Model) renderDayOverview(state *DayViewState) string {
	var content strings.Builder

	// Summary
	summary := fmt.Sprintf("%d commits across %d repositories", state.TotalCommits, len(state.RepositoryList))
	content.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(summary))
	content.WriteString("\n\n")

	// Repository breakdown
	for i, repo := range state.RepositoryList {
		commits := state.CommitsByRepo[repo]

		// Repository header
		repoHeader := fmt.Sprintf("%s", repo)
		commitCount := fmt.Sprintf("%d commits", len(commits))

		content.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render(repoHeader))
		content.WriteString("  ")
		content.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(commitCount))
		content.WriteString("\n")

		// Show first few commits
		maxShow := 3
		for j, commit := range commits {
			if j >= maxShow {
				remaining := len(commits) - maxShow
				content.WriteString(m.common.Styles.ActivityTracker.Legend.Render(
					fmt.Sprintf("  └─ ... and %d more commits", remaining)))
				content.WriteString("\n")
				break
			}

			// Commit preview
			timeStr := commit.Timestamp.Format("15:04")
			message := commit.Message
			if len(message) > 50 {
				message = message[:47] + "..."
			}

			commitLine := fmt.Sprintf("  ├─ %s %s %s",
				timeStr,
				commit.Hash[:8],
				message)

			content.WriteString(m.common.Styles.ActivityTracker.Legend.Render(commitLine))
			content.WriteString("\n")
		}

		if i < len(state.RepositoryList)-1 {
			content.WriteString("\n")
		}
	}

	return content.String()
}

// renderRepositoryDetail renders detailed view of a specific repository
func (m *Model) renderRepositoryDetail(state *DayViewState) string {
	var content strings.Builder

	if state.SelectedRepository == "" || len(state.RepositoryList) == 0 {
		return m.common.Styles.NoContent.Render("No repository selected")
	}

	commits, exists := state.CommitsByRepo[state.SelectedRepository]
	if !exists {
		return m.common.Styles.NoContent.Render("Repository not found")
	}

	// Repository header
	repoTitle := fmt.Sprintf("%s - %d commits", state.SelectedRepository, len(commits))
	content.WriteString(m.common.Styles.ActivityTracker.Title.Render(repoTitle))
	content.WriteString("\n\n")

	// Commit timeline
	for i, commit := range commits {
		isSelected := i == state.SelectedCommitIdx
		content.WriteString(m.renderCommitTimeline(commit, i+1, isSelected))
		content.WriteString("\n")
	}

	return content.String()
}

// renderCommitDetail renders detailed view of a specific commit
func (m *Model) renderCommitDetail(state *DayViewState) string {
	var content strings.Builder

	if state.SelectedRepository == "" {
		return m.common.Styles.NoContent.Render("No repository selected")
	}

	commits, exists := state.CommitsByRepo[state.SelectedRepository]
	if !exists || state.SelectedCommitIdx >= len(commits) {
		return m.common.Styles.NoContent.Render("Commit not found")
	}

	commit := commits[state.SelectedCommitIdx]

	// Commit header
	commitTitle := fmt.Sprintf("Commit Details")
	content.WriteString(m.common.Styles.ActivityTracker.Title.Render(commitTitle))
	content.WriteString("\n\n")

	// Commit information
	content.WriteString(m.renderCommitInfo(commit))

	return content.String()
}

// renderCommitTimeline renders a commit in timeline format
func (m *Model) renderCommitTimeline(commit backend.CommitInfo, index int, selected bool) string {
	var line strings.Builder

	// Selection indicator
	if selected {
		line.WriteString(m.common.Styles.ActivityTracker.StatValue.Render("► "))
	} else {
		line.WriteString("  ")
	}

	// Index
	line.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render(fmt.Sprintf("%2d. ", index)))

	// Time
	timeStr := commit.Timestamp.Format("15:04:05")
	line.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(timeStr))
	line.WriteString(" ")

	// Hash
	line.WriteString(m.common.Styles.ActivityTracker.Legend.Render(commit.Hash[:8]))
	line.WriteString(" ")

	// Author
	line.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render(fmt.Sprintf("by %s", commit.Author)))
	line.WriteString("\n")

	// Message (indented)
	message := commit.Message
	if len(message) > 80 {
		message = message[:77] + "..."
	}
	line.WriteString("     ")
	line.WriteString(m.common.Styles.ActivityTracker.Legend.Render(message))

	if selected {
		return m.common.Styles.ActivityTracker.SelectedItem.Render(line.String())
	}
	return line.String()
}

// renderCommitInfo renders detailed commit information
func (m *Model) renderCommitInfo(commit backend.CommitInfo) string {
	var info strings.Builder

	// Hash
	info.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render("Hash: "))
	info.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(commit.Hash))
	info.WriteString("\n")

	// Author
	info.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render("Author: "))
	info.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(commit.Author))
	info.WriteString("\n")

	// Repository
	info.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render("Repository: "))
	info.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(commit.Repository))
	info.WriteString("\n")

	// Timestamp
	info.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render("Time: "))
	info.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(commit.Timestamp.Format("15:04:05 MST")))
	info.WriteString("\n\n")

	// Message
	info.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render("Message:"))
	info.WriteString("\n")

	// Word wrap the message
	message := commit.Message
	maxWidth := 70
	words := strings.Fields(message)
	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		if currentLine.Len()+len(word)+1 > maxWidth {
			if currentLine.Len() > 0 {
				lines = append(lines, currentLine.String())
				currentLine.Reset()
			}
		}
		if currentLine.Len() > 0 {
			currentLine.WriteString(" ")
		}
		currentLine.WriteString(word)
	}
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	for _, line := range lines {
		info.WriteString("  ")
		info.WriteString(m.common.Styles.ActivityTracker.Legend.Render(line))
		info.WriteString("\n")
	}

	return info.String()
}

// renderDayViewHelp renders context-sensitive help for day view
func (m *Model) renderDayViewHelp() string {
	var help strings.Builder

	// Base navigation
	help.WriteString("← → navigate days • esc back to grid")

	// Context-specific help based on current view mode
	if m.dayViewState != nil {
		switch m.dayViewState.Mode {
		case DayViewModeOverview:
			if len(m.dayViewState.RepositoryList) > 0 {
				help.WriteString(" • enter view repository")
			}
		case DayViewModeRepository:
			help.WriteString(" • ↑↓ select commit • enter view details • backspace back to overview")
		case DayViewModeCommitDetail:
			help.WriteString(" • backspace back to repository")
		}
	}

	return m.common.Styles.ActivityTracker.Legend.Render(help.String())
}

// getDayViewTitle returns the appropriate title for the current day view mode
func (m *Model) getDayViewTitle() string {
	if m.dayViewState == nil {
		return "Day View"
	}

	switch m.dayViewState.Mode {
	case DayViewModeOverview:
		return "Day Overview"
	case DayViewModeRepository:
		if m.dayViewState.SelectedRepository != "" {
			return fmt.Sprintf("Repository: %s", m.dayViewState.SelectedRepository)
		}
		return "Repository View"
	case DayViewModeCommitDetail:
		return "Commit Details"
	default:
		return "Day View"
	}
}
