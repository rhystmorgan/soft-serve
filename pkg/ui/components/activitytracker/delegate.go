package activitytracker

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/soft-serve/pkg/backend"
)

// GridCell represents a single cell in the activity grid
type GridCell struct {
	Date     time.Time
	Activity *backend.ActivityData
	Week     int
	Weekday  int
	Selected bool
}

// GridDelegate handles interactions with grid cells
type GridDelegate struct {
	model *Model
}

// NewGridDelegate creates a new grid delegate
func NewGridDelegate(model *Model) *GridDelegate {
	return &GridDelegate{
		model: model,
	}
}

// GetCell returns the grid cell at the specified position
func (d *GridDelegate) GetCell(week, weekday int) *GridCell {
	now := time.Now()
	startDate := now.AddDate(0, 0, -364)
	dayIndex := week*7 + weekday

	if dayIndex >= 365 {
		return nil
	}

	date := startDate.AddDate(0, 0, dayIndex)
	dateKey := date.Format("2006-01-02")

	var activity *backend.ActivityData
	if d.model.activity != nil {
		activity = d.model.activity[dateKey]
	}

	return &GridCell{
		Date:     date,
		Activity: activity,
		Week:     week,
		Weekday:  weekday,
		Selected: week == d.model.selectedWeek && weekday == d.model.selectedWeekday,
	}
}

// GetCellStyle returns the appropriate style for a grid cell
func (d *GridDelegate) GetCellStyle(cell *GridCell) lipgloss.Style {
	if cell == nil {
		return d.model.common.Styles.ActivityTracker.Day.None
	}

	var baseStyle lipgloss.Style

	if cell.Activity != nil {
		level := backend.GetActivityLevel(cell.Activity.Count)
		switch level {
		case backend.ActivityNone:
			baseStyle = d.model.common.Styles.ActivityTracker.Day.None
		case backend.ActivityLow:
			baseStyle = d.model.common.Styles.ActivityTracker.Day.Low
		case backend.ActivityMedium:
			baseStyle = d.model.common.Styles.ActivityTracker.Day.Medium
		case backend.ActivityHigh:
			baseStyle = d.model.common.Styles.ActivityTracker.Day.High
		default:
			baseStyle = d.model.common.Styles.ActivityTracker.Day.None
		}
	} else {
		baseStyle = d.model.common.Styles.ActivityTracker.Day.None
	}

	// Apply selection styling
	if cell.Selected {
		baseStyle = baseStyle.Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("12"))
	}

	return baseStyle
}

// GetCellTooltip returns tooltip information for a grid cell
func (d *GridDelegate) GetCellTooltip(cell *GridCell) string {
	if cell == nil {
		return ""
	}

	dateStr := cell.Date.Format("Monday, January 2, 2006")

	if cell.Activity == nil || cell.Activity.Count == 0 {
		return fmt.Sprintf("%s\nNo commits", dateStr)
	}

	commitText := "commit"
	if cell.Activity.Count > 1 {
		commitText = "commits"
	}

	return fmt.Sprintf("%s\n%d %s", dateStr, cell.Activity.Count, commitText)
}

// GetCellSummary returns a brief summary for a grid cell
func (d *GridDelegate) GetCellSummary(cell *GridCell) string {
	if cell == nil || cell.Activity == nil {
		return "No activity"
	}

	return fmt.Sprintf("%d commits", cell.Activity.Count)
}

// IsValidPosition checks if the given week/weekday position is valid
func (d *GridDelegate) IsValidPosition(week, weekday int) bool {
	if week < 0 || week >= 53 || weekday < 0 || weekday >= 7 {
		return false
	}

	dayIndex := week*7 + weekday
	return dayIndex < 365
}

// GetDateFromPosition returns the date for a given week/weekday position
func (d *GridDelegate) GetDateFromPosition(week, weekday int) time.Time {
	now := time.Now()
	startDate := now.AddDate(0, 0, -364)
	dayIndex := week*7 + weekday
	return startDate.AddDate(0, 0, dayIndex)
}

// GetPositionFromDate returns the week/weekday position for a given date
func (d *GridDelegate) GetPositionFromDate(date time.Time) (week, weekday int) {
	now := time.Now()
	startDate := now.AddDate(0, 0, -364)
	daysDiff := int(date.Sub(startDate).Hours() / 24)

	if daysDiff < 0 || daysDiff >= 365 {
		return -1, -1
	}

	week = daysDiff / 7
	weekday = daysDiff % 7
	return week, weekday
}

// GetAdjacentCell returns the cell adjacent to the current position
func (d *GridDelegate) GetAdjacentCell(week, weekday int, direction Direction) *GridCell {
	newWeek, newWeekday := d.getAdjacentPosition(week, weekday, direction)
	if !d.IsValidPosition(newWeek, newWeekday) {
		return nil
	}
	return d.GetCell(newWeek, newWeekday)
}

// Direction represents movement directions in the grid
type Direction int

const (
	DirectionUp Direction = iota
	DirectionDown
	DirectionLeft
	DirectionRight
)

// getAdjacentPosition calculates the adjacent position based on direction
func (d *GridDelegate) getAdjacentPosition(week, weekday int, direction Direction) (int, int) {
	switch direction {
	case DirectionUp:
		return week, weekday - 1
	case DirectionDown:
		return week, weekday + 1
	case DirectionLeft:
		return week - 1, weekday
	case DirectionRight:
		return week + 1, weekday
	default:
		return week, weekday
	}
}

// GetWeekRange returns the start and end dates for a given week
func (d *GridDelegate) GetWeekRange(week int) (start, end time.Time) {
	now := time.Now()
	startDate := now.AddDate(0, 0, -364)

	weekStart := startDate.AddDate(0, 0, week*7)
	weekEnd := weekStart.AddDate(0, 0, 6)

	return weekStart, weekEnd
}

// GetMonthBoundaries returns the weeks where month boundaries occur
func (d *GridDelegate) GetMonthBoundaries() map[int]string {
	boundaries := make(map[int]string)
	now := time.Now()
	startDate := now.AddDate(0, 0, -364)

	currentMonth := -1
	for week := 0; week < 53; week++ {
		weekStart := startDate.AddDate(0, 0, week*7)
		month := int(weekStart.Month())

		if month != currentMonth {
			boundaries[week] = weekStart.Format("Jan")
			currentMonth = month
		}
	}

	return boundaries
}

// GetActivitySummary returns a summary of activity for the entire grid
func (d *GridDelegate) GetActivitySummary() ActivitySummary {
	summary := ActivitySummary{}

	if d.model.activity == nil {
		return summary
	}

	for week := 0; week < 53; week++ {
		for weekday := 0; weekday < 7; weekday++ {
			cell := d.GetCell(week, weekday)
			if cell != nil && cell.Activity != nil {
				summary.TotalCommits += cell.Activity.Count
				summary.ActiveDays++

				if cell.Activity.Count > summary.MaxDayCommits {
					summary.MaxDayCommits = cell.Activity.Count
					summary.MaxDay = cell.Date
				}
			}
		}
	}

	if summary.ActiveDays > 0 {
		summary.AveragePerDay = float64(summary.TotalCommits) / float64(summary.ActiveDays)
	}

	return summary
}

// ActivitySummary contains summary statistics for the activity grid
type ActivitySummary struct {
	TotalCommits  int
	ActiveDays    int
	AveragePerDay float64
	MaxDayCommits int
	MaxDay        time.Time
}

// GetCurrentStreak calculates the current commit streak
func (d *GridDelegate) GetCurrentStreak() int {
	if d.model.activity == nil {
		return 0
	}

	streak := 0
	now := time.Now()

	// Start from today and work backwards
	for i := 0; i < 365; i++ {
		date := now.AddDate(0, 0, -i)
		dateKey := date.Format("2006-01-02")

		if activity, exists := d.model.activity[dateKey]; exists && activity.Count > 0 {
			streak++
		} else {
			break
		}
	}

	return streak
}

// GetLongestStreak calculates the longest commit streak in the current period
func (d *GridDelegate) GetLongestStreak() int {
	if d.model.activity == nil {
		return 0
	}

	maxStreak := 0
	currentStreak := 0
	now := time.Now()

	// Check all days in the period
	for i := 364; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateKey := date.Format("2006-01-02")

		if activity, exists := d.model.activity[dateKey]; exists && activity.Count > 0 {
			currentStreak++
			if currentStreak > maxStreak {
				maxStreak = currentStreak
			}
		} else {
			currentStreak = 0
		}
	}

	return maxStreak
}
