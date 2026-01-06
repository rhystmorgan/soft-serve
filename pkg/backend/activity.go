package backend

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/charmbracelet/soft-serve/pkg/proto"
)

// ActivityData represents commit activity for a specific day
type ActivityData struct {
	Date    time.Time
	Count   int
	Commits []CommitInfo
}

// CommitInfo represents basic commit information
type CommitInfo struct {
	Hash       string
	Author     string
	Message    string
	Repository string
	Timestamp  time.Time
}

// ActivityLevel represents the intensity level of activity
type ActivityLevel int

const (
	ActivityNone ActivityLevel = iota
	ActivityLow
	ActivityMedium
	ActivityHigh
)

// GetActivityLevel returns the activity level based on commit count
func GetActivityLevel(count int) ActivityLevel {
	switch {
	case count == 0:
		return ActivityNone
	case count <= 3:
		return ActivityLow
	case count <= 8:
		return ActivityMedium
	default:
		return ActivityHigh
	}
}

// GetGlobalActivity returns commit activity data for all repositories over the past year
func (b *Backend) GetGlobalActivity(ctx context.Context) (map[string]*ActivityData, error) {
	// Check cache first
	if cachedData, found := b.GetCachedActivity(ctx); found {
		return cachedData, nil
	}

	// Get all repositories
	repos, err := b.Repositories(ctx)
	if err != nil {
		return nil, err
	}

	// Initialize activity map for the past 365 days
	activityMap := make(map[string]*ActivityData)
	now := time.Now()

	// Initialize all days in the past year with zero activity
	for i := 0; i < 365; i++ {
		date := now.AddDate(0, 0, -i)
		dateKey := date.Format("2006-01-02")
		activityMap[dateKey] = &ActivityData{
			Date:    date,
			Count:   0,
			Commits: make([]CommitInfo, 0),
		}
	}

	// Process each repository
	for _, repo := range repos {
		if err := b.processRepositoryActivity(ctx, repo, activityMap); err != nil {
			b.logger.Errorf("failed to process activity for repo %s: %v", repo.Name(), err)
			continue
		}
	}

	// Cache the results
	stats, _ := b.calculateStatsFromActivity(activityMap)
	b.SetCachedActivity(activityMap, stats)

	return activityMap, nil
}

// processRepositoryActivity processes commit activity for a single repository
func (b *Backend) processRepositoryActivity(ctx context.Context, repo proto.Repository, activityMap map[string]*ActivityData) error {
	// Open the git repository
	gitRepo, err := repo.Open()
	if err != nil {
		return err
	}

	// Get HEAD reference
	head, err := gitRepo.HEAD()
	if err != nil {
		// Repository might be empty, skip it
		return nil
	}

	// Get commits from the past year
	// We'll use a reasonable page size and iterate through commits
	const pageSize = 100
	page := 1
	oneYearAgo := time.Now().AddDate(-1, 0, 0)

	for {
		commits, err := gitRepo.CommitsByPage(head, page, pageSize)
		if err != nil || len(commits) == 0 {
			break
		}

		foundOldCommit := false
		for _, commit := range commits {
			// Stop if we've gone back more than a year
			if commit.Author.When.Before(oneYearAgo) {
				foundOldCommit = true
				break
			}

			// Add commit to activity map
			dateKey := commit.Author.When.Format("2006-01-02")
			if activity, exists := activityMap[dateKey]; exists {
				activity.Count++
				activity.Commits = append(activity.Commits, CommitInfo{
					Hash:       commit.ID.String()[:8], // Short hash
					Author:     commit.Author.Name,
					Message:    commit.Summary(),
					Repository: repo.Name(),
					Timestamp:  commit.Author.When,
				})
			}
		}

		if foundOldCommit {
			break
		}

		page++

		// Safety check to prevent infinite loops
		if page > 1000 {
			break
		}
	}

	return nil
}

// GetActivityStats returns summary statistics for the activity data
func (b *Backend) GetActivityStats(ctx context.Context) (*ActivityStats, error) {
	// Check if we have cached stats
	activityCache.mu.RLock()
	if time.Since(activityCache.Timestamp) < 5*time.Minute && activityCache.Stats != nil {
		stats := activityCache.Stats
		activityCache.mu.RUnlock()
		return stats, nil
	}
	activityCache.mu.RUnlock()

	// Get activity data (this will also cache it)
	activityMap, err := b.GetGlobalActivity(ctx)
	if err != nil {
		return nil, err
	}

	return b.calculateStatsFromActivity(activityMap)
}

// ActivityStats represents summary statistics for git activity
type ActivityStats struct {
	TotalCommits    int
	ActiveDays      int
	CurrentStreak   int
	LongestStreak   int
	TopRepositories []RepoActivity
	TopContributors []ContributorActivity
}

// RepoActivity represents commit activity for a repository
type RepoActivity struct {
	Name    string
	Commits int
}

// ContributorActivity represents commit activity for a contributor
type ContributorActivity struct {
	Name    string
	Commits int
}

// ActivityFilter represents filtering options for activity data
type ActivityFilter struct {
	StartDate    *time.Time
	EndDate      *time.Time
	Repositories []string
	Authors      []string
}

// ActivityCache represents cached activity data
type ActivityCache struct {
	Data      map[string]*ActivityData
	Stats     *ActivityStats
	Timestamp time.Time
	mu        sync.RWMutex
}

// Global cache instance
var activityCache = &ActivityCache{}

// GetFilteredActivity returns filtered activity data
func (b *Backend) GetFilteredActivity(ctx context.Context, filter *ActivityFilter) (map[string]*ActivityData, error) {
	// Get base activity data
	activityMap, err := b.GetGlobalActivity(ctx)
	if err != nil {
		return nil, err
	}

	if filter == nil {
		return activityMap, nil
	}

	filteredMap := make(map[string]*ActivityData)

	for dateKey, activity := range activityMap {
		// Filter by date range
		if filter.StartDate != nil && activity.Date.Before(*filter.StartDate) {
			continue
		}
		if filter.EndDate != nil && activity.Date.After(*filter.EndDate) {
			continue
		}

		// Create filtered activity for this day
		filteredActivity := &ActivityData{
			Date:    activity.Date,
			Count:   0,
			Commits: make([]CommitInfo, 0),
		}

		// Filter commits
		for _, commit := range activity.Commits {
			includeCommit := true

			// Filter by repository
			if len(filter.Repositories) > 0 {
				found := false
				for _, repo := range filter.Repositories {
					if commit.Repository == repo {
						found = true
						break
					}
				}
				if !found {
					includeCommit = false
				}
			}

			// Filter by author
			if includeCommit && len(filter.Authors) > 0 {
				found := false
				for _, author := range filter.Authors {
					if commit.Author == author {
						found = true
						break
					}
				}
				if !found {
					includeCommit = false
				}
			}

			if includeCommit {
				filteredActivity.Commits = append(filteredActivity.Commits, commit)
				filteredActivity.Count++
			}
		}

		filteredMap[dateKey] = filteredActivity
	}

	return filteredMap, nil
}

// GetCachedActivity returns cached activity data if available and fresh
func (b *Backend) GetCachedActivity(ctx context.Context) (map[string]*ActivityData, bool) {
	activityCache.mu.RLock()
	defer activityCache.mu.RUnlock()

	// Cache is valid for 5 minutes
	if time.Since(activityCache.Timestamp) < 5*time.Minute && activityCache.Data != nil {
		return activityCache.Data, true
	}

	return nil, false
}

// SetCachedActivity updates the activity cache
func (b *Backend) SetCachedActivity(data map[string]*ActivityData, stats *ActivityStats) {
	activityCache.mu.Lock()
	defer activityCache.mu.Unlock()

	activityCache.Data = data
	activityCache.Stats = stats
	activityCache.Timestamp = time.Now()
}

// GetRepositoryList returns a list of all repository names
func (b *Backend) GetRepositoryList(ctx context.Context) ([]string, error) {
	repos, err := b.Repositories(ctx)
	if err != nil {
		return nil, err
	}

	var repoNames []string
	for _, repo := range repos {
		repoNames = append(repoNames, repo.Name())
	}

	sort.Strings(repoNames)
	return repoNames, nil
}

// GetAuthorList returns a list of all commit authors
func (b *Backend) GetAuthorList(ctx context.Context) ([]string, error) {
	activityMap, err := b.GetGlobalActivity(ctx)
	if err != nil {
		return nil, err
	}

	authorSet := make(map[string]bool)
	for _, activity := range activityMap {
		for _, commit := range activity.Commits {
			authorSet[commit.Author] = true
		}
	}

	var authors []string
	for author := range authorSet {
		authors = append(authors, author)
	}

	sort.Strings(authors)
	return authors, nil
}

// calculateStatsFromActivity calculates statistics from activity data
func (b *Backend) calculateStatsFromActivity(activityMap map[string]*ActivityData) (*ActivityStats, error) {
	stats := &ActivityStats{
		TotalCommits:    0,
		ActiveDays:      0,
		CurrentStreak:   0,
		LongestStreak:   0,
		TopRepositories: make([]RepoActivity, 0),
		TopContributors: make([]ContributorActivity, 0),
	}

	repoCommits := make(map[string]int)
	authorCommits := make(map[string]int)

	// Sort dates for streak calculation
	var sortedDates []string
	for date := range activityMap {
		sortedDates = append(sortedDates, date)
	}
	sort.Strings(sortedDates)

	// Calculate statistics
	currentStreak := 0
	longestStreak := 0

	for i := len(sortedDates) - 1; i >= 0; i-- {
		date := sortedDates[i]
		activity := activityMap[date]

		stats.TotalCommits += activity.Count

		if activity.Count > 0 {
			stats.ActiveDays++
			currentStreak++
			if currentStreak > longestStreak {
				longestStreak = currentStreak
			}

			// Count by repository and author
			for _, commit := range activity.Commits {
				repoCommits[commit.Repository]++
				authorCommits[commit.Author]++
			}
		} else {
			if i == len(sortedDates)-1 {
				// Only count current streak if it includes today
				stats.CurrentStreak = currentStreak
			}
			currentStreak = 0
		}
	}

	stats.LongestStreak = longestStreak

	// Convert maps to sorted slices
	for repo, count := range repoCommits {
		stats.TopRepositories = append(stats.TopRepositories, RepoActivity{
			Name:    repo,
			Commits: count,
		})
	}

	for author, count := range authorCommits {
		stats.TopContributors = append(stats.TopContributors, ContributorActivity{
			Name:    author,
			Commits: count,
		})
	}

	// Sort by commit count
	sort.Slice(stats.TopRepositories, func(i, j int) bool {
		return stats.TopRepositories[i].Commits > stats.TopRepositories[j].Commits
	})

	sort.Slice(stats.TopContributors, func(i, j int) bool {
		return stats.TopContributors[i].Commits > stats.TopContributors[j].Commits
	})

	// Limit to top 10
	if len(stats.TopRepositories) > 10 {
		stats.TopRepositories = stats.TopRepositories[:10]
	}
	if len(stats.TopContributors) > 10 {
		stats.TopContributors = stats.TopContributors[:10]
	}

	return stats, nil
}
