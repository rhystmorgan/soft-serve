package keymap

import "github.com/charmbracelet/bubbles/v2/key"

// KeyMap is a map of key bindings for the UI.
type KeyMap struct {
	Quit       key.Binding
	Up         key.Binding
	Down       key.Binding
	UpDown     key.Binding
	LeftRight  key.Binding
	Arrows     key.Binding
	GotoTop    key.Binding
	GotoBottom key.Binding
	Select     key.Binding
	Section    key.Binding
	Back       key.Binding
	PrevPage   key.Binding
	NextPage   key.Binding
	Help       key.Binding

	SelectItem key.Binding
	BackItem   key.Binding

	Copy key.Binding

	// Activity Tracker specific bindings
	ActivityStats      key.Binding
	ActivityFilter     key.Binding
	ActivityTimePeriod key.Binding
	ActivityRefresh    key.Binding
	ActivityExport     key.Binding
	ActivityYear       key.Binding
	ActivityPrevYear   key.Binding
	ActivityWeekJump   key.Binding
	ActivityWeekBack   key.Binding
}

// DefaultKeyMap returns the default key map.
func DefaultKeyMap() *KeyMap {
	km := new(KeyMap)

	km.Quit = key.NewBinding(
		key.WithKeys(
			"q",
			"ctrl+c",
		),
		key.WithHelp(
			"q",
			"quit",
		),
	)

	km.Up = key.NewBinding(
		key.WithKeys(
			"up",
			"k",
		),
		key.WithHelp(
			"↑",
			"up",
		),
	)

	km.Down = key.NewBinding(
		key.WithKeys(
			"down",
			"j",
		),
		key.WithHelp(
			"↓",
			"down",
		),
	)

	km.UpDown = key.NewBinding(
		key.WithKeys(
			"up",
			"down",
			"k",
			"j",
		),
		key.WithHelp(
			"↑↓",
			"navigate",
		),
	)

	km.LeftRight = key.NewBinding(
		key.WithKeys(
			"left",
			"h",
			"right",
			"l",
		),
		key.WithHelp(
			"←→",
			"navigate",
		),
	)

	km.Arrows = key.NewBinding(
		key.WithKeys(
			"up",
			"right",
			"down",
			"left",
			"k",
			"j",
			"h",
			"l",
		),
		key.WithHelp(
			"↑←↓→",
			"navigate",
		),
	)

	km.GotoTop = key.NewBinding(
		key.WithKeys(
			"home",
			"g",
		),
		key.WithHelp(
			"g/home",
			"goto top",
		),
	)

	km.GotoBottom = key.NewBinding(
		key.WithKeys(
			"end",
			"G",
		),
		key.WithHelp(
			"G/end",
			"goto bottom",
		),
	)

	km.Select = key.NewBinding(
		key.WithKeys(
			"enter",
		),
		key.WithHelp(
			"enter",
			"select",
		),
	)

	km.Section = key.NewBinding(
		key.WithKeys(
			"tab",
			"shift+tab",
		),
		key.WithHelp(
			"tab",
			"section",
		),
	)

	km.Back = key.NewBinding(
		key.WithKeys(
			"esc",
		),
		key.WithHelp(
			"esc",
			"back",
		),
	)

	km.PrevPage = key.NewBinding(
		key.WithKeys(
			"pgup",
			"b",
			"u",
		),
		key.WithHelp(
			"pgup",
			"prev page",
		),
	)

	km.NextPage = key.NewBinding(
		key.WithKeys(
			"pgdown",
			"f",
			"d",
		),
		key.WithHelp(
			"pgdn",
			"next page",
		),
	)

	km.Help = key.NewBinding(
		key.WithKeys(
			"?",
		),
		key.WithHelp(
			"?",
			"toggle help",
		),
	)

	km.SelectItem = key.NewBinding(
		key.WithKeys(
			"l",
			"right",
		),
		key.WithHelp(
			"→/l",
			"select",
		),
	)

	km.BackItem = key.NewBinding(
		key.WithKeys(
			"h",
			"left",
			"backspace",
		),
		key.WithHelp(
			"←/h",
			"back",
		),
	)

	km.Copy = key.NewBinding(
		key.WithKeys(
			"c",
			"ctrl+c",
		),
		key.WithHelp(
			"c",
			"copy text",
		),
	)

	// Activity Tracker key bindings
	km.ActivityStats = key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "statistics"),
	)

	km.ActivityFilter = key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "filter"),
	)

	km.ActivityTimePeriod = key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "time period"),
	)

	km.ActivityRefresh = key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	)

	km.ActivityExport = key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "export"),
	)

	km.ActivityYear = key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "current year"),
	)

	km.ActivityPrevYear = key.NewBinding(
		key.WithKeys("Y"),
		key.WithHelp("Y", "previous year"),
	)

	km.ActivityWeekJump = key.NewBinding(
		key.WithKeys("w"),
		key.WithHelp("w", "next week"),
	)

	km.ActivityWeekBack = key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "prev week"),
	)

	return km
}
