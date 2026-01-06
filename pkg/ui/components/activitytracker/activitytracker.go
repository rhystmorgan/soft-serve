package activitytracker

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/spinner"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/soft-serve/pkg/backend"
	"github.com/charmbracelet/soft-serve/pkg/ui/common"
	"github.com/charmbracelet/soft-serve/pkg/ui/pages/repo"
)

var waitBeforeLoading = time.Millisecond * 100

// ViewMode represents the current view mode
type ViewMode int

const (
	ViewModeGrid ViewMode = iota
	ViewModeDetail
	ViewModeStats
	ViewModeFilter
	ViewModeTimePeriod
)

// activityState represents the loading state
type activityState int

const (
	activityStateLoading activityState = iota
	activityStateReady
	activityStateError
)

// Model represents the activity tracker component
type Model struct {
	common           common.Common
	activity         map[string]*backend.ActivityData
	stats            *backend.ActivityStats
	spinner          spinner.Model
	state            activityState
	viewMode         ViewMode
	selectedDate     time.Time
	selectedWeek     int
	selectedWeekday  int
	showRepositories bool
	showContributors bool
	filter           *backend.ActivityFilter
	repositories     []string
	authors          []string
	dayViewState     *DayViewState
	currentYear      int
	loadingTime      time.Time
}

// New creates a new activity tracker model
func New(common common.Common) *Model {
	now := time.Now()
	sp := spinner.New(spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(common.Styles.Spinner))

	return &Model{
		common:          common,
		spinner:         sp,
		state:           activityStateLoading,
		viewMode:        ViewModeGrid,
		selectedDate:    now,
		selectedWeek:    0,
		selectedWeekday: int(now.Weekday()),
		currentYear:     now.Year(),
	}
}

func (m *Model) startLoading() tea.Cmd {
	m.loadingTime = time.Now()
	m.state = activityStateLoading
	return m.spinner.Tick
}

// LoadActivityMsg is sent when activity data should be loaded
type LoadActivityMsg struct{}

// ActivityLoadedMsg is sent when activity data has been loaded
type ActivityLoadedMsg struct {
	Activity     map[string]*backend.ActivityData
	Stats        *backend.ActivityStats
	Error        error
	Repositories []string
	Authors      []string
}

// Init implements tea.Model
func (m *Model) Init() tea.Cmd {
	// Reset all state variables to initial values
	now := time.Now()
	m.activity = nil
	m.stats = nil
	m.state = activityStateLoading
	m.viewMode = ViewModeGrid
	m.selectedDate = now
	m.selectedWeek = 0
	m.selectedWeekday = int(now.Weekday())
	m.showRepositories = false
	m.showContributors = false
	m.filter = nil
	m.repositories = nil
	m.authors = nil
	m.dayViewState = nil
	m.currentYear = now.Year()
	m.loadingTime = time.Time{} // Reset loading time

	return tea.Batch(
		m.startLoading(),
		func() tea.Msg { return LoadActivityMsg{} },
	)
}

// Update implements tea.Model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case spinner.TickMsg:
		if m.state == activityStateLoading && m.spinner.ID() == msg.ID {
			s, cmd := m.spinner.Update(msg)
			m.spinner = s
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	case LoadActivityMsg:
		m.state = activityStateLoading
		return m, tea.Batch(append(cmds, m.loadActivityCmd())...)
	case ActivityLoadedMsg:
		if msg.Error != nil {
			m.state = activityStateError
			return m, common.ErrorCmd(msg.Error)
		}
		m.state = activityStateReady
		m.activity = msg.Activity
		m.stats = msg.Stats
		if msg.Repositories != nil {
			m.repositories = msg.Repositories
		}
		if msg.Authors != nil {
			m.authors = msg.Authors
		}
		return m, tea.Batch(cmds...)
	case repo.EmptyRepoMsg:
		// Reset state when repository is empty
		m.activity = nil
		m.stats = nil
		m.state = activityStateReady
		m.viewMode = ViewModeGrid
		m.filter = nil
		m.repositories = nil
		m.authors = nil
		m.dayViewState = nil
		// Reset selection to current date
		now := time.Now()
		m.selectedDate = now
		m.selectedWeek = 0
		m.selectedWeekday = int(now.Weekday())
		m.currentYear = now.Year()
	case common.ErrorMsg:
		// Handle error messages from the application
		m.state = activityStateReady
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}
	return m, tea.Batch(cmds...)
}

// handleKeyPress handles keyboard input
func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.viewMode {
	case ViewModeGrid:
		return m.handleGridKeyPress(msg)
	case ViewModeDetail:
		return m.handleDetailKeyPress(msg)
	case ViewModeStats:
		return m.handleStatsKeyPress(msg)
	case ViewModeFilter:
		return m.handleFilterKeyPress(msg)
	case ViewModeTimePeriod:
		return m.handleTimePeriodKeyPress(msg)
	}
	return m, nil
}

// handleGridKeyPress handles keyboard input in grid mode
func (m *Model) handleGridKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "left", "h":
		if m.selectedWeek > 0 {
			m.selectedWeek--
			m.updateSelectedDate()
		}
	case "right", "l":
		if m.selectedWeek < 52 {
			m.selectedWeek++
			m.updateSelectedDate()
		}
	case "up", "k":
		if m.selectedWeekday > 0 {
			m.selectedWeekday--
			m.updateSelectedDate()
		}
	case "down", "j":
		if m.selectedWeekday < 6 {
			m.selectedWeekday++
			m.updateSelectedDate()
		}
	case "enter", " ":
		m.viewMode = ViewModeDetail
		m.dayViewState = nil // Reset day view state
	case "s":
		m.viewMode = ViewModeStats
	case "f":
		m.viewMode = ViewModeFilter
	case "t":
		m.viewMode = ViewModeTimePeriod
	case "r":
		m.showRepositories = !m.showRepositories
	case "c":
		m.showContributors = !m.showContributors
	case "y":
		// Jump to current year
		now := time.Now()
		m.currentYear = now.Year()
		m.selectedDate = now
		m.updateSelectedPosition()
		return m, m.loadActivityCmd()
	case "Y":
		// Jump to previous year
		m.currentYear--
		m.selectedDate = time.Date(m.currentYear, 12, 31, 0, 0, 0, 0, time.Local)
		m.updateSelectedPosition()
		return m, m.loadActivityCmd()
	case "g":
		// Go to beginning of year (vim-style)
		m.selectedDate = time.Date(m.currentYear, 1, 1, 0, 0, 0, 0, time.Local)
		m.updateSelectedPosition()
	case "G":
		// Go to end of year (vim-style)
		m.selectedDate = time.Date(m.currentYear, 12, 31, 0, 0, 0, 0, time.Local)
		m.updateSelectedPosition()
	case "0", "home":
		// Go to beginning of week
		m.selectedWeekday = 0
		m.updateSelectedDate()
	case "$", "end":
		// Go to end of week
		m.selectedWeekday = 6
		m.updateSelectedDate()
	case "w":
		// Jump forward by week
		if m.selectedWeek < 52 {
			m.selectedWeek++
			m.updateSelectedDate()
		}
	case "b":
		// Jump backward by week
		if m.selectedWeek > 0 {
			m.selectedWeek--
			m.updateSelectedDate()
		}
	case "e":
		// Export functionality (placeholder for Phase 4)
		// TODO: Implement export in Phase 4
	case "esc":
		m.viewMode = ViewModeGrid
		m.showRepositories = false
		m.showContributors = false
	}
	return m, nil
}

// handleDetailKeyPress handles keyboard input in detail mode
func (m *Model) handleDetailKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.viewMode = ViewModeGrid
		m.dayViewState = nil
	case "left", "h":
		m.selectedDate = m.selectedDate.AddDate(0, 0, -1)
		m.updateSelectedPosition()
		m.dayViewState = nil // Reset day view state for new date
	case "right", "l":
		m.selectedDate = m.selectedDate.AddDate(0, 0, 1)
		m.updateSelectedPosition()
		m.dayViewState = nil // Reset day view state for new date
	case "enter":
		// Handle day view mode transitions
		if m.dayViewState != nil {
			switch m.dayViewState.Mode {
			case DayViewModeOverview:
				if len(m.dayViewState.RepositoryList) > 0 {
					m.dayViewState.Mode = DayViewModeRepository
					m.dayViewState.SelectedRepository = m.dayViewState.RepositoryList[0]
					m.dayViewState.SelectedCommitIdx = 0
				}
			case DayViewModeRepository:
				m.dayViewState.Mode = DayViewModeCommitDetail
			}
		}
	case "backspace":
		// Handle day view mode back navigation
		if m.dayViewState != nil {
			switch m.dayViewState.Mode {
			case DayViewModeCommitDetail:
				m.dayViewState.Mode = DayViewModeRepository
			case DayViewModeRepository:
				m.dayViewState.Mode = DayViewModeOverview
			}
		}
	case "up", "k":
		// Navigate within repository or commit list
		if m.dayViewState != nil && m.dayViewState.Mode == DayViewModeRepository {
			if m.dayViewState.SelectedCommitIdx > 0 {
				m.dayViewState.SelectedCommitIdx--
			}
		}
	case "down", "j":
		// Navigate within repository or commit list
		if m.dayViewState != nil && m.dayViewState.Mode == DayViewModeRepository {
			if m.dayViewState.SelectedRepository != "" {
				commits := m.dayViewState.CommitsByRepo[m.dayViewState.SelectedRepository]
				if m.dayViewState.SelectedCommitIdx < len(commits)-1 {
					m.dayViewState.SelectedCommitIdx++
				}
			}
		}
	}
	return m, nil
}

// handleStatsKeyPress handles keyboard input in stats mode
func (m *Model) handleStatsKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.viewMode = ViewModeGrid
	case "r":
		m.showRepositories = !m.showRepositories
	case "c":
		m.showContributors = !m.showContributors
	}
	return m, nil
}

// handleFilterKeyPress handles keyboard input in filter mode
func (m *Model) handleFilterKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.viewMode = ViewModeGrid
	case "enter":
		// Apply filter and return to grid
		m.viewMode = ViewModeGrid
		return m, m.loadFilteredActivityCmd()
	case "c":
		// Clear filter
		m.filter = nil
		m.viewMode = ViewModeGrid
		return m, m.loadActivityCmd()
	}
	return m, nil
}

// handleTimePeriodKeyPress handles keyboard input in time period mode
func (m *Model) handleTimePeriodKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.viewMode = ViewModeGrid
	case "1":
		// Current year
		now := time.Now()
		m.currentYear = now.Year()
		m.selectedDate = now
		m.updateSelectedPosition()
		m.viewMode = ViewModeGrid
		return m, m.loadActivityCmd()
	case "2":
		// Previous year
		m.currentYear = time.Now().Year() - 1
		m.selectedDate = time.Date(m.currentYear, 12, 31, 0, 0, 0, 0, time.Local)
		m.updateSelectedPosition()
		m.viewMode = ViewModeGrid
		return m, m.loadActivityCmd()
	case "3":
		// All time (last 365 days from today)
		now := time.Now()
		m.currentYear = now.Year()
		m.selectedDate = now
		m.updateSelectedPosition()
		m.viewMode = ViewModeGrid
		return m, m.loadActivityCmd()
	case "4":
		// Custom year selection (placeholder)
		// TODO: Implement custom year input
		m.viewMode = ViewModeGrid
	}
	return m, nil
}

// updateSelectedDate updates the selected date based on week/weekday position
func (m *Model) updateSelectedDate() {
	now := time.Now()
	startDate := now.AddDate(0, 0, -364)
	dayIndex := m.selectedWeek*7 + m.selectedWeekday
	m.selectedDate = startDate.AddDate(0, 0, dayIndex)
}

// updateSelectedPosition updates the week/weekday position based on selected date
func (m *Model) updateSelectedPosition() {
	now := time.Now()
	startDate := now.AddDate(0, 0, -364)
	daysDiff := int(m.selectedDate.Sub(startDate).Hours() / 24)
	if daysDiff >= 0 && daysDiff < 365 {
		m.selectedWeek = daysDiff / 7
		m.selectedWeekday = daysDiff % 7
	}
}

// View implements tea.Model
func (m *Model) View() string {
	switch m.state {
	case activityStateLoading:
		if !m.loadingTime.IsZero() && m.loadingTime.Add(waitBeforeLoading).Before(time.Now()) {
			return m.renderLoading()
		}
		fallthrough
	case activityStateReady:
		return m.renderContent()
	case activityStateError:
		return m.common.Styles.Error.Render("Failed to load activity data")
	default:
		return m.renderContent()
	}
}

func (m *Model) renderLoading() string {
	msg := fmt.Sprintf("%s loading…", m.spinner.View())
	return m.common.Styles.SpinnerContainer.
		Height(m.common.Height).
		Render(msg)
}

func (m *Model) renderContent() string {
	if m.activity == nil {
		return m.common.Styles.NoContent.Render("No activity data available")
	}

	switch m.viewMode {
	case ViewModeDetail:
		return m.renderDetailView()
	case ViewModeStats:
		return m.renderStatsView()
	case ViewModeFilter:
		return m.renderFilterView()
	case ViewModeTimePeriod:
		return m.renderTimePeriodView()
	default:
		return m.renderGridView()
	}
}

// renderGridView renders the main grid view
func (m *Model) renderGridView() string {
	var content strings.Builder

	// Title
	content.WriteString(m.common.Styles.ActivityTracker.Title.Render("Git Activity"))
	content.WriteString("\n\n")

	// Activity grid
	content.WriteString(m.renderActivityGrid())
	content.WriteString("\n")

	// Legend
	content.WriteString(m.renderLegend())
	content.WriteString("\n")

	// Stats
	if m.stats != nil {
		content.WriteString(m.renderStats())
	}

	// Help text
	content.WriteString("\n")
	content.WriteString(m.renderHelp())

	// Expandable panels
	if m.showRepositories {
		content.WriteString("\n\n")
		content.WriteString(m.renderRepositoriesPanel())
	}

	if m.showContributors {
		content.WriteString("\n\n")
		content.WriteString(m.renderContributorsPanel())
	}

	return m.common.Styles.ActivityTracker.Container.Render(content.String())
}

// renderDetailView renders the detailed day view
func (m *Model) renderDetailView() string {
	// Initialize day view state if not already done
	if m.dayViewState == nil {
		m.dayViewState = m.initDayView(m.selectedDate)
	}

	return m.renderEnhancedDetailView()
}

// loadFilteredActivityCmd loads filtered activity data from the backend
func (m *Model) loadFilteredActivityCmd() tea.Cmd {
	return func() tea.Msg {
		b := m.common.Backend()
		if b == nil {
			return ActivityLoadedMsg{Error: fmt.Errorf("backend not available")}
		}

		var err error
		var activity map[string]*backend.ActivityData

		if m.filter != nil {
			activity, err = b.GetFilteredActivity(m.common.Context(), m.filter)
		} else {
			activity, err = b.GetGlobalActivity(m.common.Context())
		}

		if err != nil {
			return ActivityLoadedMsg{Error: err}
		}

		stats, err := b.GetActivityStats(m.common.Context())
		if err != nil {
			return ActivityLoadedMsg{Activity: activity, Error: err}
		}

		return ActivityLoadedMsg{
			Activity: activity,
			Stats:    stats,
		}
	}
}

// loadActivityCmd loads activity data from the backend
func (m *Model) loadActivityCmd() tea.Cmd {
	return func() tea.Msg {
		backend := m.common.Backend()
		if backend == nil {
			return ActivityLoadedMsg{Error: fmt.Errorf("backend not available")}
		}

		activity, err := backend.GetGlobalActivity(m.common.Context())
		if err != nil {
			return ActivityLoadedMsg{Error: err}
		}

		stats, err := backend.GetActivityStats(m.common.Context())
		if err != nil {
			return ActivityLoadedMsg{Activity: activity, Error: err}
		}

		// Also load repository and author lists for filtering
		repositories, _ := backend.GetRepositoryList(m.common.Context())
		authors, _ := backend.GetAuthorList(m.common.Context())

		return ActivityLoadedMsg{
			Activity:     activity,
			Stats:        stats,
			Repositories: repositories,
			Authors:      authors,
		}
	}
}

// renderStatsView renders the detailed statistics view
func (m *Model) renderStatsView() string {
	var content strings.Builder

	// Title
	content.WriteString(m.common.Styles.ActivityTracker.Title.Render("Activity Statistics"))
	content.WriteString("\n\n")

	if m.stats != nil {
		// Main stats
		content.WriteString(m.renderDetailedStats())
		content.WriteString("\n\n")

		// Repository stats
		content.WriteString(m.renderRepositoriesPanel())
		content.WriteString("\n\n")

		// Contributor stats
		content.WriteString(m.renderContributorsPanel())
	}

	// Help text
	content.WriteString("\n")
	content.WriteString(m.common.Styles.ActivityTracker.Legend.Render("r toggle repositories • c toggle contributors • esc back to grid"))

	return m.common.Styles.ActivityTracker.Container.Render(content.String())
}

// SetSize implements common.Component
// getMargins calculates the margins needed for proper spacing with UI elements
func (m *Model) getMargins() (int, int) {
	wm := 0
	hm := 0

	// Account for any UI elements that reduce available space
	// This follows the pattern used by other tab components
	if m.common.Styles != nil {
		// Add vertical margin for any frame elements
		hm = m.common.Styles.Tabs.GetVerticalFrameSize()
	}

	return wm, hm
}

func (m *Model) SetSize(width, height int) {
	m.common.SetSize(width, height)
	wm, hm := m.getMargins()

	// Apply margins to ensure proper spacing
	// The available space is reduced by the calculated margins
	m.common.Width = width - wm
	m.common.Height = height - hm
}

// renderActivityGrid renders the 365-day activity grid
func (m *Model) renderActivityGrid() string {
	var grid strings.Builder

	// Get sorted dates for the past 365 days
	now := time.Now()
	dates := make([]time.Time, 365)
	for i := 0; i < 365; i++ {
		dates[i] = now.AddDate(0, 0, -364+i) // Start from 364 days ago to today
	}

	// Render month labels
	grid.WriteString(m.renderMonthLabels(dates))
	grid.WriteString("\n")

	// Render weekday labels and grid
	weekdays := []string{"", "M", "", "W", "", "F", ""}

	for weekday := 0; weekday < 7; weekday++ {
		var row strings.Builder

		// Weekday label
		label := ""
		if weekday < len(weekdays) {
			label = weekdays[weekday]
		}
		row.WriteString(m.common.Styles.ActivityTracker.WeekdayLabel.Render(label))
		row.WriteString(" ")

		// Days for this weekday
		for week := 0; week < 53; week++ {
			dayIndex := week*7 + weekday
			if dayIndex < len(dates) {
				date := dates[dayIndex]
				dateKey := date.Format("2006-01-02")

				var dayStyle lipgloss.Style
				if activity, exists := m.activity[dateKey]; exists {
					level := backend.GetActivityLevel(activity.Count)
					dayStyle = m.getDayStyle(level)
				} else {
					dayStyle = m.common.Styles.ActivityTracker.Day.None
				}

				// Highlight selected day
				if week == m.selectedWeek && weekday == m.selectedWeekday {
					dayStyle = dayStyle.Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("12"))
				}

				row.WriteString(dayStyle.Render("  "))
			} else {
				// Empty space for days that don't exist
				row.WriteString("  ")
			}
		}

		grid.WriteString(m.common.Styles.ActivityTracker.GridRow.Render(row.String()))
		grid.WriteString("\n")
	}

	return m.common.Styles.ActivityTracker.Grid.Render(grid.String())
}

// renderMonthLabels renders the month labels above the grid
func (m *Model) renderMonthLabels(dates []time.Time) string {
	var labels strings.Builder

	// Space for weekday labels
	labels.WriteString("   ")

	currentMonth := -1
	for week := 0; week < 53; week++ {
		dayIndex := week * 7
		if dayIndex < len(dates) {
			month := int(dates[dayIndex].Month())
			if month != currentMonth {
				monthName := dates[dayIndex].Format("Jan")
				labels.WriteString(m.common.Styles.ActivityTracker.MonthLabel.Render(monthName))
				currentMonth = month
			} else {
				labels.WriteString("   ")
			}
		} else {
			labels.WriteString("   ")
		}
	}

	return labels.String()
}

// renderLegend renders the activity level legend
func (m *Model) renderLegend() string {
	var legend strings.Builder

	legend.WriteString(m.common.Styles.ActivityTracker.Legend.Render("Less "))

	// Activity level squares
	levels := []backend.ActivityLevel{
		backend.ActivityNone,
		backend.ActivityLow,
		backend.ActivityMedium,
		backend.ActivityHigh,
	}

	for _, level := range levels {
		style := m.getDayStyle(level)
		legend.WriteString(m.common.Styles.ActivityTracker.LegendItem.Render(style.Render("  ")))
	}

	legend.WriteString(m.common.Styles.ActivityTracker.Legend.Render(" More"))

	return legend.String()
}

// renderStats renders activity statistics
func (m *Model) renderStats() string {
	var stats strings.Builder

	stats.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render("Total commits: "))
	stats.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(fmt.Sprintf("%d", m.stats.TotalCommits)))
	stats.WriteString("  ")

	stats.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render("Active days: "))
	stats.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(fmt.Sprintf("%d", m.stats.ActiveDays)))
	stats.WriteString("  ")

	stats.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render("Current streak: "))
	stats.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(fmt.Sprintf("%d", m.stats.CurrentStreak)))
	stats.WriteString("  ")

	stats.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render("Longest streak: "))
	stats.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(fmt.Sprintf("%d", m.stats.LongestStreak)))

	return m.common.Styles.ActivityTracker.Stats.Render(stats.String())
}

// getDayStyle returns the appropriate style for an activity level
func (m *Model) getDayStyle(level backend.ActivityLevel) lipgloss.Style {
	switch level {
	case backend.ActivityNone:
		return m.common.Styles.ActivityTracker.Day.None
	case backend.ActivityLow:
		return m.common.Styles.ActivityTracker.Day.Low
	case backend.ActivityMedium:
		return m.common.Styles.ActivityTracker.Day.Medium
	case backend.ActivityHigh:
		return m.common.Styles.ActivityTracker.Day.High
	default:
		return m.common.Styles.ActivityTracker.Day.None
	}
}

// renderHelp renders help text
func (m *Model) renderHelp() string {
	helpText := "↑↓←→ navigate • enter view day • s stats • f filter • t time period"
	return m.common.Styles.ActivityTracker.Legend.Render(helpText)
}

// renderTimePeriodView renders the time period selection view
func (m *Model) renderTimePeriodView() string {
	var content strings.Builder

	// Title
	content.WriteString(m.common.Styles.ActivityTracker.Title.Render("Select Time Period"))
	content.WriteString("\n\n")

	// Current selection
	currentPeriod := fmt.Sprintf("Currently viewing: %d", m.currentYear)
	content.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(currentPeriod))
	content.WriteString("\n\n")

	// Options
	content.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render("Available periods:"))
	content.WriteString("\n\n")

	options := []string{
		"1. Current year (" + fmt.Sprintf("%d", time.Now().Year()) + ")",
		"2. Previous year (" + fmt.Sprintf("%d", time.Now().Year()-1) + ")",
		"3. Last 365 days (rolling)",
		"4. Custom year (coming soon)",
	}

	for _, option := range options {
		content.WriteString("  ")
		content.WriteString(m.common.Styles.ActivityTracker.Legend.Render(option))
		content.WriteString("\n")
	}

	content.WriteString("\n")

	// Help text
	content.WriteString(m.common.Styles.ActivityTracker.Legend.Render("1-4 select period • esc back to grid"))

	return m.common.Styles.ActivityTracker.Container.Render(content.String())
}

// renderFilterView renders the filter configuration view
func (m *Model) renderFilterView() string {
	var content strings.Builder

	// Title
	content.WriteString(m.common.Styles.ActivityTracker.Title.Render("Filter Activity"))
	content.WriteString("\n\n")

	// Current filter status
	if m.filter == nil {
		content.WriteString(m.common.Styles.ActivityTracker.Legend.Render("No filter applied - showing all activity"))
	} else {
		content.WriteString(m.common.Styles.ActivityTracker.Legend.Render("Filter applied:"))
		content.WriteString("\n")

		if m.filter.StartDate != nil {
			content.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render("Start date: "))
			content.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(m.filter.StartDate.Format("2006-01-02")))
			content.WriteString("\n")
		}

		if m.filter.EndDate != nil {
			content.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render("End date: "))
			content.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(m.filter.EndDate.Format("2006-01-02")))
			content.WriteString("\n")
		}

		if len(m.filter.Repositories) > 0 {
			content.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render("Repositories: "))
			content.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(strings.Join(m.filter.Repositories, ", ")))
			content.WriteString("\n")
		}

		if len(m.filter.Authors) > 0 {
			content.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render("Authors: "))
			content.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(strings.Join(m.filter.Authors, ", ")))
			content.WriteString("\n")
		}
	}

	content.WriteString("\n")

	// Available repositories
	if len(m.repositories) > 0 {
		content.WriteString(m.common.Styles.ActivityTracker.Title.Render("Available Repositories:"))
		content.WriteString("\n")
		for i, repo := range m.repositories {
			if i >= 10 { // Limit display
				content.WriteString(m.common.Styles.ActivityTracker.Legend.Render("... and more"))
				break
			}
			content.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(fmt.Sprintf("• %s", repo)))
			content.WriteString("\n")
		}
		content.WriteString("\n")
	}

	// Available authors
	if len(m.authors) > 0 {
		content.WriteString(m.common.Styles.ActivityTracker.Title.Render("Available Authors:"))
		content.WriteString("\n")
		for i, author := range m.authors {
			if i >= 10 { // Limit display
				content.WriteString(m.common.Styles.ActivityTracker.Legend.Render("... and more"))
				break
			}
			content.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(fmt.Sprintf("• %s", author)))
			content.WriteString("\n")
		}
	}

	// Help text
	content.WriteString("\n")
	content.WriteString(m.common.Styles.ActivityTracker.Legend.Render("enter apply filter • c clear filter • esc back to grid"))

	return m.common.Styles.ActivityTracker.Container.Render(content.String())
}

// renderCommitItem renders a single commit item
func (m *Model) renderCommitItem(commit backend.CommitInfo) string {
	var item strings.Builder

	// Time
	timeStr := commit.Timestamp.Format("15:04")
	item.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render(timeStr))
	item.WriteString(" ")

	// Hash
	item.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(commit.Hash))
	item.WriteString(" ")

	// Repository
	item.WriteString(m.common.Styles.ActivityTracker.Legend.Render(fmt.Sprintf("[%s]", commit.Repository)))
	item.WriteString(" ")

	// Author
	item.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render(commit.Author))
	item.WriteString(" ")

	// Message
	message := commit.Message
	if len(message) > 60 {
		message = message[:57] + "..."
	}
	item.WriteString(m.common.Styles.ActivityTracker.Legend.Render(message))

	return item.String()
}

// renderDetailedStats renders detailed statistics
func (m *Model) renderDetailedStats() string {
	var stats strings.Builder

	stats.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render("Total commits: "))
	stats.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(fmt.Sprintf("%d", m.stats.TotalCommits)))
	stats.WriteString("\n")

	stats.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render("Active days: "))
	stats.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(fmt.Sprintf("%d", m.stats.ActiveDays)))
	stats.WriteString("\n")

	stats.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render("Current streak: "))
	stats.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(fmt.Sprintf("%d days", m.stats.CurrentStreak)))
	stats.WriteString("\n")

	stats.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render("Longest streak: "))
	stats.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(fmt.Sprintf("%d days", m.stats.LongestStreak)))

	return stats.String()
}

// renderRepositoriesPanel renders the repositories statistics panel
func (m *Model) renderRepositoriesPanel() string {
	if m.stats == nil || len(m.stats.TopRepositories) == 0 {
		return m.common.Styles.NoContent.Render("No repository data available")
	}

	var panel strings.Builder

	panel.WriteString(m.common.Styles.ActivityTracker.Title.Render("Top Repositories"))
	panel.WriteString("\n")

	for i, repo := range m.stats.TopRepositories {
		if i >= 10 { // Limit to top 10
			break
		}

		panel.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render(fmt.Sprintf("%2d. ", i+1)))
		panel.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(repo.Name))
		panel.WriteString(m.common.Styles.ActivityTracker.Legend.Render(fmt.Sprintf(" (%d commits)", repo.Commits)))
		panel.WriteString("\n")
	}

	return panel.String()
}

// renderContributorsPanel renders the contributors statistics panel
func (m *Model) renderContributorsPanel() string {
	if m.stats == nil || len(m.stats.TopContributors) == 0 {
		return m.common.Styles.NoContent.Render("No contributor data available")
	}

	var panel strings.Builder

	panel.WriteString(m.common.Styles.ActivityTracker.Title.Render("Top Contributors"))
	panel.WriteString("\n")

	for i, contributor := range m.stats.TopContributors {
		if i >= 10 { // Limit to top 10
			break
		}

		panel.WriteString(m.common.Styles.ActivityTracker.StatLabel.Render(fmt.Sprintf("%2d. ", i+1)))
		panel.WriteString(m.common.Styles.ActivityTracker.StatValue.Render(contributor.Name))
		panel.WriteString(m.common.Styles.ActivityTracker.Legend.Render(fmt.Sprintf(" (%d commits)", contributor.Commits)))
		panel.WriteString("\n")
	}

	return panel.String()
}

// TabComponent interface implementation

// TabName implements common.TabComponent
func (m *Model) TabName() string {
	return "Activity"
}

// StatusBarValue implements common.TabComponent
func (m *Model) StatusBarValue() string {
	if m.state == activityStateLoading {
		return "Loading activity data..."
	}
	if m.state == activityStateError {
		return "Error loading activity"
	}
	if m.stats != nil {
		return fmt.Sprintf("%d commits", m.stats.TotalCommits)
	}
	return ""
}

// StatusBarInfo implements common.TabComponent
func (m *Model) StatusBarInfo() string {
	if m.viewMode == ViewModeDetail && m.selectedDate != (time.Time{}) {
		return m.selectedDate.Format("Jan 2, 2006")
	}
	return fmt.Sprintf("Year %d", m.currentYear)
}

// SpinnerID implements common.TabComponent
func (m *Model) SpinnerID() int {
	return m.spinner.ID()
}

// Path implements common.TabComponent
func (m *Model) Path() string {
	switch m.viewMode {
	case ViewModeDetail:
		return fmt.Sprintf("activity/%s", m.selectedDate.Format("2006-01-02"))
	case ViewModeStats:
		return "activity/stats"
	case ViewModeFilter:
		return "activity/filter"
	default:
		return "activity"
	}
}

// help.KeyMap interface implementation

// ShortHelp implements help.KeyMap
func (m *Model) ShortHelp() []key.Binding {
	switch m.viewMode {
	case ViewModeGrid:
		return []key.Binding{
			m.common.KeyMap.UpDown,
			m.common.KeyMap.LeftRight,
			m.common.KeyMap.Select,
		}
	case ViewModeDetail:
		return []key.Binding{
			m.common.KeyMap.LeftRight,
			m.common.KeyMap.Select,
			m.common.KeyMap.Back,
		}
	case ViewModeStats:
		return []key.Binding{
			m.common.KeyMap.Back,
		}
	default:
		return []key.Binding{}
	}
}

// FullHelp implements help.KeyMap
func (m *Model) FullHelp() [][]key.Binding {
	switch m.viewMode {
	case ViewModeGrid:
		return [][]key.Binding{
			{
				m.common.KeyMap.Up,
				m.common.KeyMap.Down,
				m.common.KeyMap.LeftRight,
			},
			{
				m.common.KeyMap.Select,
				m.common.KeyMap.Back,
			},
		}
	case ViewModeDetail:
		return [][]key.Binding{
			{
				m.common.KeyMap.LeftRight,
				m.common.KeyMap.UpDown,
			},
			{
				m.common.KeyMap.Select,
				m.common.KeyMap.Back,
			},
		}
	default:
		return [][]key.Binding{{m.common.KeyMap.Back}}
	}
}
