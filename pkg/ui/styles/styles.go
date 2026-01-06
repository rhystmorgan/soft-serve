package styles

import (
	"image/color"

	"github.com/charmbracelet/lipgloss/v2"
)

// XXX: For now, this is in its own package so that it can be shared between
// different packages without incurring an illegal import cycle.

// Styles defines styles for the UI.
type Styles struct {
	ActiveBorderColor   color.Color
	InactiveBorderColor color.Color

	App                  lipgloss.Style
	ServerName           lipgloss.Style
	TopLevelNormalTab    lipgloss.Style
	TopLevelActiveTab    lipgloss.Style
	TopLevelActiveTabDot lipgloss.Style

	MenuItem       lipgloss.Style
	MenuLastUpdate lipgloss.Style

	RepoSelector struct {
		Normal struct {
			Base    lipgloss.Style
			Title   lipgloss.Style
			Desc    lipgloss.Style
			Command lipgloss.Style
			Updated lipgloss.Style
		}
		Active struct {
			Base    lipgloss.Style
			Title   lipgloss.Style
			Desc    lipgloss.Style
			Command lipgloss.Style
			Updated lipgloss.Style
		}
	}

	Repo struct {
		Base       lipgloss.Style
		Title      lipgloss.Style
		Command    lipgloss.Style
		Body       lipgloss.Style
		Header     lipgloss.Style
		HeaderName lipgloss.Style
		HeaderDesc lipgloss.Style
	}

	Footer      lipgloss.Style
	Branch      lipgloss.Style
	HelpKey     lipgloss.Style
	HelpValue   lipgloss.Style
	HelpDivider lipgloss.Style
	URLStyle    lipgloss.Style

	Error      lipgloss.Style
	ErrorTitle lipgloss.Style
	ErrorBody  lipgloss.Style

	LogItem struct {
		Normal struct {
			Base    lipgloss.Style
			Hash    lipgloss.Style
			Title   lipgloss.Style
			Desc    lipgloss.Style
			Keyword lipgloss.Style
		}
		Active struct {
			Base    lipgloss.Style
			Hash    lipgloss.Style
			Title   lipgloss.Style
			Desc    lipgloss.Style
			Keyword lipgloss.Style
		}
	}

	Log struct {
		Commit         lipgloss.Style
		CommitHash     lipgloss.Style
		CommitAuthor   lipgloss.Style
		CommitDate     lipgloss.Style
		CommitBody     lipgloss.Style
		CommitStatsAdd lipgloss.Style
		CommitStatsDel lipgloss.Style
		Paginator      lipgloss.Style
	}

	Ref struct {
		Normal struct {
			Base     lipgloss.Style
			Item     lipgloss.Style
			ItemTag  lipgloss.Style
			ItemDesc lipgloss.Style
			ItemHash lipgloss.Style
		}
		Active struct {
			Base     lipgloss.Style
			Item     lipgloss.Style
			ItemTag  lipgloss.Style
			ItemDesc lipgloss.Style
			ItemHash lipgloss.Style
		}
		ItemSelector lipgloss.Style
		Paginator    lipgloss.Style
		Selector     lipgloss.Style
	}

	Tree struct {
		Normal struct {
			FileName lipgloss.Style
			FileDir  lipgloss.Style
			FileMode lipgloss.Style
			FileSize lipgloss.Style
		}
		Active struct {
			FileName lipgloss.Style
			FileDir  lipgloss.Style
			FileMode lipgloss.Style
			FileSize lipgloss.Style
		}
		Selector    lipgloss.Style
		FileContent lipgloss.Style
		Paginator   lipgloss.Style
		Blame       struct {
			Hash    lipgloss.Style
			Message lipgloss.Style
			Who     lipgloss.Style
		}
	}

	Stash struct {
		Normal struct {
			Message lipgloss.Style
		}
		Active struct {
			Message lipgloss.Style
		}
		Title    lipgloss.Style
		Selector lipgloss.Style
	}

	Spinner          lipgloss.Style
	SpinnerContainer lipgloss.Style

	NoContent lipgloss.Style

	StatusBar       lipgloss.Style
	StatusBarKey    lipgloss.Style
	StatusBarValue  lipgloss.Style
	StatusBarInfo   lipgloss.Style
	StatusBarBranch lipgloss.Style
	StatusBarHelp   lipgloss.Style

	Tabs         lipgloss.Style
	TabInactive  lipgloss.Style
	TabActive    lipgloss.Style
	TabSeparator lipgloss.Style

	Code struct {
		LineDigit lipgloss.Style
		LineBar   lipgloss.Style
	}

	ActivityTracker struct {
		Container lipgloss.Style
		Title     lipgloss.Style
		Grid      lipgloss.Style
		GridRow   lipgloss.Style
		Day       struct {
			None   lipgloss.Style
			Low    lipgloss.Style
			Medium lipgloss.Style
			High   lipgloss.Style
		}
		MonthLabel      lipgloss.Style
		WeekdayLabel    lipgloss.Style
		Legend          lipgloss.Style
		LegendItem      lipgloss.Style
		Stats           lipgloss.Style
		StatValue       lipgloss.Style
		StatLabel       lipgloss.Style
		SelectedDay     lipgloss.Style
		SelectedItem    lipgloss.Style
		DetailContainer lipgloss.Style
		CommitItem      lipgloss.Style
		FilterContainer lipgloss.Style
		FilterActive    lipgloss.Style
		FilterInactive  lipgloss.Style
		PanelTitle      lipgloss.Style
		PanelItem       lipgloss.Style
		HelpText        lipgloss.Style
	}
}

// DefaultStyles returns default styles for the UI.
func DefaultStyles() *Styles {
	highlightColor := lipgloss.Color("210")
	highlightColorDim := lipgloss.Color("174")
	selectorColor := lipgloss.Color("167")
	hashColor := lipgloss.Color("185")

	s := new(Styles)

	s.ActiveBorderColor = lipgloss.Color("62")
	s.InactiveBorderColor = lipgloss.Color("241")

	s.App = lipgloss.NewStyle().
		Margin(1, 2)

	s.ServerName = lipgloss.NewStyle().
		Height(1).
		MarginLeft(1).
		MarginBottom(1).
		Padding(0, 1).
		Background(lipgloss.Color("57")).
		Foreground(lipgloss.Color("229")).
		Bold(true)

	s.TopLevelNormalTab = lipgloss.NewStyle().
		MarginRight(2)

	s.TopLevelActiveTab = s.TopLevelNormalTab.
		Foreground(lipgloss.Color("36"))

	s.TopLevelActiveTabDot = lipgloss.NewStyle().
		Foreground(lipgloss.Color("36"))

	s.RepoSelector.Normal.Base = lipgloss.NewStyle().
		PaddingLeft(1).
		Border(lipgloss.Border{Left: " "}, false, false, false, true).
		Height(3)

	s.RepoSelector.Normal.Title = lipgloss.NewStyle().Bold(true)

	s.RepoSelector.Normal.Desc = lipgloss.NewStyle().
		Foreground(lipgloss.Color("243"))

	s.RepoSelector.Normal.Command = lipgloss.NewStyle().
		Foreground(lipgloss.Color("132"))

	s.RepoSelector.Normal.Updated = lipgloss.NewStyle().
		Foreground(lipgloss.Color("243"))

	s.RepoSelector.Active.Base = s.RepoSelector.Normal.Base.
		BorderStyle(lipgloss.Border{Left: "┃"}).
		BorderForeground(lipgloss.Color("176"))

	s.RepoSelector.Active.Title = s.RepoSelector.Normal.Title.
		Foreground(lipgloss.Color("212"))

	s.RepoSelector.Active.Desc = s.RepoSelector.Normal.Desc.
		Foreground(lipgloss.Color("246"))

	s.RepoSelector.Active.Updated = s.RepoSelector.Normal.Updated.
		Foreground(lipgloss.Color("212"))

	s.RepoSelector.Active.Command = s.RepoSelector.Normal.Command.
		Foreground(lipgloss.Color("204"))

	s.MenuItem = lipgloss.NewStyle().
		PaddingLeft(1).
		Border(lipgloss.Border{
			Left: " ",
		}, false, false, false, true).
		Height(3)

	s.MenuLastUpdate = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Align(lipgloss.Right)

	s.Repo.Base = lipgloss.NewStyle()

	s.Repo.Title = lipgloss.NewStyle().
		Padding(0, 2)

	s.Repo.Command = lipgloss.NewStyle().
		Foreground(lipgloss.Color("168"))

	s.Repo.Body = lipgloss.NewStyle().
		Margin(1, 0)

	s.Repo.Header = lipgloss.NewStyle().
		MaxHeight(2).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color("236"))

	s.Repo.HeaderName = lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		Bold(true)

	s.Repo.HeaderDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color("243"))

	s.Footer = lipgloss.NewStyle().
		MarginTop(1).
		Padding(0, 1).
		Height(1)

	s.Branch = lipgloss.NewStyle().
		Foreground(lipgloss.Color("203")).
		Background(lipgloss.Color("236")).
		Padding(0, 1)

	s.HelpKey = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	s.HelpValue = lipgloss.NewStyle().
		Foreground(lipgloss.Color("239"))

	s.HelpDivider = lipgloss.NewStyle().
		Foreground(lipgloss.Color("237")).
		SetString(" • ")

	s.URLStyle = lipgloss.NewStyle().
		MarginLeft(1).
		Foreground(lipgloss.Color("168"))

	s.Error = lipgloss.NewStyle().
		MarginTop(2)

	s.ErrorTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("204")).
		Bold(true).
		Padding(0, 1)

	s.ErrorBody = lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		MarginLeft(2)

	s.LogItem.Normal.Base = lipgloss.NewStyle().
		Border(lipgloss.Border{
			Left: " ",
		}, false, false, false, true).
		PaddingLeft(1)

	s.LogItem.Active.Base = s.LogItem.Normal.Base.
		Border(lipgloss.Border{
			Left: "┃",
		}, false, false, false, true).
		BorderForeground(selectorColor)

	s.LogItem.Active.Hash = s.LogItem.Normal.Hash.
		Foreground(hashColor)

	s.LogItem.Active.Hash = lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor)

	s.LogItem.Normal.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("105"))

	s.LogItem.Active.Title = lipgloss.NewStyle().
		Foreground(highlightColor).
		Bold(true)

	s.LogItem.Normal.Desc = lipgloss.NewStyle().
		Foreground(lipgloss.Color("246"))

	s.LogItem.Active.Desc = lipgloss.NewStyle().
		Foreground(lipgloss.Color("95"))

	s.LogItem.Active.Keyword = s.LogItem.Active.Desc.
		Foreground(highlightColorDim)

	s.LogItem.Normal.Hash = lipgloss.NewStyle().
		Foreground(hashColor)

	s.LogItem.Active.Hash = lipgloss.NewStyle().
		Foreground(highlightColor)

	s.Log.Commit = lipgloss.NewStyle().
		Margin(0, 2)

	s.Log.CommitHash = lipgloss.NewStyle().
		Foreground(hashColor).
		Bold(true)

	s.Log.CommitBody = lipgloss.NewStyle().
		MarginTop(1).
		MarginLeft(2)

	s.Log.CommitStatsAdd = lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Bold(true)

	s.Log.CommitStatsDel = lipgloss.NewStyle().
		Foreground(lipgloss.Color("203")).
		Bold(true)

	s.Log.Paginator = lipgloss.NewStyle().
		Margin(0).
		Align(lipgloss.Center)

	s.Ref.Normal.Item = lipgloss.NewStyle()

	s.Ref.ItemSelector = lipgloss.NewStyle().
		Foreground(selectorColor).
		SetString("> ")

	s.Ref.Active.Item = lipgloss.NewStyle().
		Foreground(highlightColorDim)

	s.Ref.Normal.Base = lipgloss.NewStyle()

	s.Ref.Active.Base = lipgloss.NewStyle()

	s.Ref.Normal.ItemTag = lipgloss.NewStyle().
		Foreground(lipgloss.Color("39"))

	s.Ref.Active.ItemTag = lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor)

	s.Ref.Active.Item = lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor)

	s.Ref.Normal.ItemDesc = lipgloss.NewStyle().
		Faint(true)

	s.Ref.Active.ItemDesc = lipgloss.NewStyle().
		Foreground(highlightColor).
		Faint(true)

	s.Ref.Normal.ItemHash = lipgloss.NewStyle().
		Foreground(hashColor).
		Bold(true)

	s.Ref.Active.ItemHash = lipgloss.NewStyle().
		Foreground(highlightColor).
		Bold(true)

	s.Ref.Paginator = s.Log.Paginator

	s.Ref.Selector = lipgloss.NewStyle()

	s.Tree.Selector = s.Tree.Normal.FileName.
		Width(1).
		Foreground(selectorColor)

	s.Tree.Normal.FileName = lipgloss.NewStyle().
		MarginLeft(1)

	s.Tree.Active.FileName = s.Tree.Normal.FileName.
		Bold(true).
		Foreground(highlightColor)

	s.Tree.Normal.FileDir = lipgloss.NewStyle().
		Foreground(lipgloss.Color("39"))

	s.Tree.Active.FileDir = lipgloss.NewStyle().
		Foreground(highlightColor)

	s.Tree.Normal.FileMode = s.Tree.Active.FileName.
		Width(10).
		Foreground(lipgloss.Color("243"))

	s.Tree.Active.FileMode = s.Tree.Normal.FileMode.
		Foreground(highlightColorDim)

	s.Tree.Normal.FileSize = s.Tree.Normal.FileName.
		Foreground(lipgloss.Color("243"))

	s.Tree.Active.FileSize = s.Tree.Normal.FileName.
		Foreground(highlightColorDim)

	s.Tree.FileContent = lipgloss.NewStyle()

	s.Tree.Paginator = s.Log.Paginator

	s.Tree.Blame.Hash = lipgloss.NewStyle().
		Foreground(hashColor).
		Bold(true)

	s.Tree.Blame.Message = lipgloss.NewStyle()

	s.Tree.Blame.Who = lipgloss.NewStyle().
		Faint(true)

	s.Spinner = lipgloss.NewStyle().
		MarginTop(1).
		MarginLeft(2).
		Foreground(lipgloss.Color("205"))

	s.SpinnerContainer = lipgloss.NewStyle()

	s.NoContent = lipgloss.NewStyle().
		MarginTop(1).
		MarginLeft(2).
		Foreground(lipgloss.Color("242"))

	s.StatusBar = lipgloss.NewStyle().
		Height(1)

	s.StatusBarKey = lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Background(lipgloss.Color("206")).
		Foreground(lipgloss.Color("228"))

	s.StatusBarValue = lipgloss.NewStyle().
		Padding(0, 1).
		Background(lipgloss.Color("235")).
		Foreground(lipgloss.Color("243"))

	s.StatusBarInfo = lipgloss.NewStyle().
		Padding(0, 1).
		Background(lipgloss.Color("212")).
		Foreground(lipgloss.Color("230"))

	s.StatusBarBranch = lipgloss.NewStyle().
		Padding(0, 1).
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230"))

	s.StatusBarHelp = lipgloss.NewStyle().
		Padding(0, 1).
		Background(lipgloss.Color("237")).
		Foreground(lipgloss.Color("243"))

	s.Tabs = lipgloss.NewStyle().
		Height(1)

	s.TabInactive = lipgloss.NewStyle()

	s.TabActive = lipgloss.NewStyle().
		Underline(true).
		Foreground(lipgloss.Color("36"))

	s.TabSeparator = lipgloss.NewStyle().
		SetString("│").
		Padding(0, 1).
		Foreground(lipgloss.Color("238"))

	s.Code.LineDigit = lipgloss.NewStyle().Foreground(lipgloss.Color("239"))

	s.Code.LineBar = lipgloss.NewStyle().Foreground(lipgloss.Color("236"))

	s.Stash.Normal.Message = lipgloss.NewStyle().MarginLeft(1)

	s.Stash.Active.Message = s.Stash.Normal.Message.Foreground(selectorColor)

	s.Stash.Title = lipgloss.NewStyle().
		Foreground(hashColor).
		Bold(true)

	s.Stash.Selector = lipgloss.NewStyle().
		Width(1).
		Foreground(selectorColor)

	return s
}

// CatppuccinStyles returns Catppuccin-themed styles for the UI.
func CatppuccinStyles() *Styles {
	// Catppuccin Mocha color palette
	peach := lipgloss.Color("#fab387")    // Primary UI accent
	yellow := lipgloss.Color("#f9e2af")   // Git hashes
	blue := lipgloss.Color("#89b4fa")     // Directories, tabs
	sapphire := lipgloss.Color("#74c7ec") // Activity tracker low
	lavender := lipgloss.Color("#b4befe") // Branches, activity tracker high
	mauve := lipgloss.Color("#cba6f7")    // Tags
	sky := lipgloss.Color("#89dceb")      // Commands, secondary actions
	green := lipgloss.Color("#a6e3a1")    // Success, additions
	red := lipgloss.Color("#f38ba8")      // Errors, deletions

	// Text hierarchy
	text := lipgloss.Color("#cdd6f4")     // Primary text
	subtext1 := lipgloss.Color("#bac2de") // Secondary text
	subtext0 := lipgloss.Color("#a6adc8") // Tertiary text
	overlay1 := lipgloss.Color("#7f849c") // Muted text
	overlay0 := lipgloss.Color("#6c7086") // Very muted

	// Backgrounds
	crust := lipgloss.Color("#11111b")    // Main background
	surface0 := lipgloss.Color("#313244") // Borders
	surface1 := lipgloss.Color("#45475a") // Subtle backgrounds

	s := new(Styles)

	s.ActiveBorderColor = peach
	s.InactiveBorderColor = overlay0

	s.App = lipgloss.NewStyle().
		Margin(1, 2)

	s.ServerName = lipgloss.NewStyle().
		Height(1).
		MarginLeft(1).
		MarginBottom(1).
		Padding(0, 1).
		Background(peach).
		Foreground(crust).
		Bold(true)

	s.TopLevelNormalTab = lipgloss.NewStyle().
		MarginRight(2)

	s.TopLevelActiveTab = s.TopLevelNormalTab.
		Foreground(blue)

	s.TopLevelActiveTabDot = lipgloss.NewStyle().
		Foreground(blue)

	s.RepoSelector.Normal.Base = lipgloss.NewStyle().
		PaddingLeft(1).
		Border(lipgloss.Border{Left: " "}, false, false, false, true).
		Height(3)

	s.RepoSelector.Normal.Title = lipgloss.NewStyle().
		Foreground(text).
		Bold(true)

	s.RepoSelector.Normal.Desc = lipgloss.NewStyle().
		Foreground(subtext1)

	s.RepoSelector.Normal.Command = lipgloss.NewStyle().
		Foreground(sky)

	s.RepoSelector.Normal.Updated = lipgloss.NewStyle().
		Foreground(overlay1)

	s.RepoSelector.Active.Base = s.RepoSelector.Normal.Base.
		BorderStyle(lipgloss.Border{Left: "┃"}).
		BorderForeground(peach)

	s.RepoSelector.Active.Title = s.RepoSelector.Normal.Title.
		Foreground(peach)

	s.RepoSelector.Active.Desc = s.RepoSelector.Normal.Desc.
		Foreground(subtext0)

	s.RepoSelector.Active.Updated = s.RepoSelector.Normal.Updated.
		Foreground(peach)

	s.RepoSelector.Active.Command = s.RepoSelector.Normal.Command.
		Foreground(sky)

	s.MenuItem = lipgloss.NewStyle().
		PaddingLeft(1).
		Border(lipgloss.Border{
			Left: " ",
		}, false, false, false, true).
		Height(3)

	s.MenuLastUpdate = lipgloss.NewStyle().
		Foreground(overlay1).
		Align(lipgloss.Right)

	s.Repo.Base = lipgloss.NewStyle()

	s.Repo.Title = lipgloss.NewStyle().
		Padding(0, 2)

	s.Repo.Command = lipgloss.NewStyle().
		Foreground(sky)

	s.Repo.Body = lipgloss.NewStyle().
		Margin(1, 0)

	s.Repo.Header = lipgloss.NewStyle().
		MaxHeight(2).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(surface0)

	s.Repo.HeaderName = lipgloss.NewStyle().
		Foreground(peach).
		Bold(true)

	s.Repo.HeaderDesc = lipgloss.NewStyle().
		Foreground(subtext1)

	s.Footer = lipgloss.NewStyle().
		MarginTop(1).
		Padding(0, 1).
		Height(1)

	s.Branch = lipgloss.NewStyle().
		Foreground(lavender).
		Background(surface1).
		Padding(0, 1)

	s.HelpKey = lipgloss.NewStyle().
		Foreground(overlay1)

	s.HelpValue = lipgloss.NewStyle().
		Foreground(overlay0)

	s.HelpDivider = lipgloss.NewStyle().
		Foreground(surface0).
		SetString(" • ")

	s.URLStyle = lipgloss.NewStyle().
		MarginLeft(1).
		Foreground(sky)

	s.Error = lipgloss.NewStyle().
		MarginTop(2)

	s.ErrorTitle = lipgloss.NewStyle().
		Foreground(crust).
		Background(red).
		Bold(true).
		Padding(0, 1)

	s.ErrorBody = lipgloss.NewStyle().
		Foreground(text).
		MarginLeft(2)

	s.LogItem.Normal.Base = lipgloss.NewStyle().
		Border(lipgloss.Border{
			Left: " ",
		}, false, false, false, true).
		PaddingLeft(1)

	s.LogItem.Active.Base = s.LogItem.Normal.Base.
		Border(lipgloss.Border{
			Left: "┃",
		}, false, false, false, true).
		BorderForeground(peach)

	s.LogItem.Normal.Hash = lipgloss.NewStyle().
		Foreground(yellow)

	s.LogItem.Active.Hash = lipgloss.NewStyle().
		Foreground(yellow).
		Bold(true)

	s.LogItem.Normal.Title = lipgloss.NewStyle().
		Foreground(text)

	s.LogItem.Active.Title = lipgloss.NewStyle().
		Foreground(peach).
		Bold(true)

	s.LogItem.Normal.Desc = lipgloss.NewStyle().
		Foreground(subtext1)

	s.LogItem.Active.Desc = lipgloss.NewStyle().
		Foreground(subtext0)

	s.LogItem.Active.Keyword = s.LogItem.Active.Desc.
		Foreground(peach)

	s.Log.Commit = lipgloss.NewStyle().
		Margin(0, 2)

	s.Log.CommitHash = lipgloss.NewStyle().
		Foreground(yellow).
		Bold(true)

	s.Log.CommitBody = lipgloss.NewStyle().
		MarginTop(1).
		MarginLeft(2)

	s.Log.CommitStatsAdd = lipgloss.NewStyle().
		Foreground(green).
		Bold(true)

	s.Log.CommitStatsDel = lipgloss.NewStyle().
		Foreground(red).
		Bold(true)

	s.Log.Paginator = lipgloss.NewStyle().
		Margin(0).
		Align(lipgloss.Center)

	s.Ref.Normal.Item = lipgloss.NewStyle().
		Foreground(text)

	s.Ref.ItemSelector = lipgloss.NewStyle().
		Foreground(peach).
		SetString("> ")

	s.Ref.Active.Item = lipgloss.NewStyle().
		Foreground(peach)

	s.Ref.Normal.Base = lipgloss.NewStyle()

	s.Ref.Active.Base = lipgloss.NewStyle()

	s.Ref.Normal.ItemTag = lipgloss.NewStyle().
		Foreground(mauve)

	s.Ref.Active.ItemTag = lipgloss.NewStyle().
		Bold(true).
		Foreground(mauve)

	s.Ref.Normal.ItemDesc = lipgloss.NewStyle().
		Foreground(subtext1)

	s.Ref.Active.ItemDesc = lipgloss.NewStyle().
		Foreground(peach)

	s.Ref.Normal.ItemHash = lipgloss.NewStyle().
		Foreground(yellow).
		Bold(true)

	s.Ref.Active.ItemHash = lipgloss.NewStyle().
		Foreground(yellow).
		Bold(true)

	s.Ref.Paginator = s.Log.Paginator

	s.Ref.Selector = lipgloss.NewStyle()

	s.Tree.Selector = lipgloss.NewStyle().
		Width(1).
		Foreground(peach)

	s.Tree.Normal.FileName = lipgloss.NewStyle().
		MarginLeft(1).
		Foreground(text)

	s.Tree.Active.FileName = s.Tree.Normal.FileName.
		Bold(true).
		Foreground(peach)

	s.Tree.Normal.FileDir = lipgloss.NewStyle().
		Foreground(blue)

	s.Tree.Active.FileDir = lipgloss.NewStyle().
		Foreground(blue).
		Bold(true)

	s.Tree.Normal.FileMode = lipgloss.NewStyle().
		Width(10).
		Foreground(overlay1)

	s.Tree.Active.FileMode = s.Tree.Normal.FileMode.
		Foreground(subtext0)

	s.Tree.Normal.FileSize = lipgloss.NewStyle().
		Foreground(overlay1)

	s.Tree.Active.FileSize = lipgloss.NewStyle().
		Foreground(subtext0)

	s.Tree.FileContent = lipgloss.NewStyle()

	s.Tree.Paginator = s.Log.Paginator

	s.Tree.Blame.Hash = lipgloss.NewStyle().
		Foreground(yellow).
		Bold(true)

	s.Tree.Blame.Message = lipgloss.NewStyle().
		Foreground(text)

	s.Tree.Blame.Who = lipgloss.NewStyle().
		Foreground(subtext1)

	s.Spinner = lipgloss.NewStyle().
		MarginTop(1).
		MarginLeft(2).
		Foreground(peach)

	s.SpinnerContainer = lipgloss.NewStyle()

	s.NoContent = lipgloss.NewStyle().
		MarginTop(1).
		MarginLeft(2).
		Foreground(overlay1)

	s.StatusBar = lipgloss.NewStyle().
		Height(1)

	s.StatusBarKey = lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Background(peach).
		Foreground(crust)

	s.StatusBarValue = lipgloss.NewStyle().
		Padding(0, 1).
		Background(surface1).
		Foreground(text)

	s.StatusBarInfo = lipgloss.NewStyle().
		Padding(0, 1).
		Background(blue).
		Foreground(crust)

	s.StatusBarBranch = lipgloss.NewStyle().
		Padding(0, 1).
		Background(lavender).
		Foreground(crust)

	s.StatusBarHelp = lipgloss.NewStyle().
		Padding(0, 1).
		Background(surface0).
		Foreground(subtext1)

	s.Tabs = lipgloss.NewStyle().
		Height(1)

	s.TabInactive = lipgloss.NewStyle().
		Foreground(subtext1)

	s.TabActive = lipgloss.NewStyle().
		Underline(true).
		Foreground(blue)

	s.TabSeparator = lipgloss.NewStyle().
		SetString("│").
		Padding(0, 1).
		Foreground(surface0)

	s.Code.LineDigit = lipgloss.NewStyle().
		Foreground(overlay0)

	s.Code.LineBar = lipgloss.NewStyle().
		Foreground(surface0)

	s.Stash.Normal.Message = lipgloss.NewStyle().
		MarginLeft(1).
		Foreground(text)

	s.Stash.Active.Message = s.Stash.Normal.Message.
		Foreground(peach)

	s.Stash.Title = lipgloss.NewStyle().
		Foreground(yellow).
		Bold(true)

	s.Stash.Selector = lipgloss.NewStyle().
		Width(1).
		Foreground(peach)

	// ActivityTracker styles with Catppuccin blue progression
	s.ActivityTracker.Container = lipgloss.NewStyle().
		Padding(1, 2).
		Margin(1, 0)

	s.ActivityTracker.Title = lipgloss.NewStyle().
		Foreground(peach).
		Bold(true).
		MarginBottom(1)

	s.ActivityTracker.Grid = lipgloss.NewStyle().
		MarginBottom(1)

	s.ActivityTracker.GridRow = lipgloss.NewStyle()

	// Activity levels using Catppuccin blue colors
	s.ActivityTracker.Day.None = lipgloss.NewStyle().
		Width(2).
		Height(1).
		Background(surface0). // #313244 - no activity
		Foreground(surface0)

	s.ActivityTracker.Day.Low = lipgloss.NewStyle().
		Width(2).
		Height(1).
		Background(sapphire). // #74c7ec - low activity
		Foreground(sapphire)

	s.ActivityTracker.Day.Medium = lipgloss.NewStyle().
		Width(2).
		Height(1).
		Background(blue). // #89b4fa - medium activity
		Foreground(blue)

	s.ActivityTracker.Day.High = lipgloss.NewStyle().
		Width(2).
		Height(1).
		Background(lavender). // #b4befe - high activity
		Foreground(lavender)

	s.ActivityTracker.MonthLabel = lipgloss.NewStyle().
		Foreground(subtext1).
		Width(3).
		Align(lipgloss.Center)

	s.ActivityTracker.WeekdayLabel = lipgloss.NewStyle().
		Foreground(subtext1).
		Width(2).
		Align(lipgloss.Center)

	s.ActivityTracker.Legend = lipgloss.NewStyle().
		MarginTop(1).
		Foreground(subtext1)

	s.ActivityTracker.LegendItem = lipgloss.NewStyle().
		MarginRight(1)

	s.ActivityTracker.Stats = lipgloss.NewStyle().
		MarginTop(2).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(surface1)

	s.ActivityTracker.StatValue = lipgloss.NewStyle().
		Foreground(peach).
		Bold(true)

	s.ActivityTracker.StatLabel = lipgloss.NewStyle().
		Foreground(subtext1)

	// Interactive styles for Phase 2 features
	s.ActivityTracker.SelectedDay = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(peach).
		Width(2).
		Height(1)

	s.ActivityTracker.SelectedItem = lipgloss.NewStyle().
		Background(surface1).
		Foreground(text).
		Bold(true)

	s.ActivityTracker.DetailContainer = lipgloss.NewStyle().
		Padding(1, 2).
		Margin(1, 0).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(surface1)

	s.ActivityTracker.CommitItem = lipgloss.NewStyle().
		Padding(0, 1).
		MarginBottom(1).
		Border(lipgloss.Border{Left: "│"}, false, false, false, true).
		BorderForeground(surface1)

	s.ActivityTracker.FilterContainer = lipgloss.NewStyle().
		Padding(1, 2).
		Margin(1, 0).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(blue)

	s.ActivityTracker.FilterActive = lipgloss.NewStyle().
		Foreground(green).
		Bold(true)

	s.ActivityTracker.FilterInactive = lipgloss.NewStyle().
		Foreground(surface1)

	s.ActivityTracker.PanelTitle = lipgloss.NewStyle().
		Foreground(blue).
		Bold(true).
		MarginBottom(1).
		Border(lipgloss.Border{Bottom: "─"}, false, false, true, false).
		BorderForeground(surface1)

	s.ActivityTracker.PanelItem = lipgloss.NewStyle().
		Padding(0, 1).
		MarginLeft(1)

	s.ActivityTracker.HelpText = lipgloss.NewStyle().
		Foreground(overlay0).
		Italic(true).
		MarginTop(1)

	return s
}

// GetStyles returns the appropriate styles based on environment or preference.
// Defaults to Catppuccin theme, but can fall back to DefaultStyles.
func GetStyles() *Styles {
	// For now, always return Catppuccin theme
	// In the future, this could check environment variables or config
	return CatppuccinStyles()
}
