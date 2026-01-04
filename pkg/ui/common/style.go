package common

import (
	"github.com/charmbracelet/colorprofile"
	gansi "github.com/charmbracelet/glamour/v2/ansi"
	"github.com/charmbracelet/soft-serve/pkg/ui/styles"
)

// DefaultColorProfile is the default color profile used by the SSH server.
var DefaultColorProfile = colorprofile.ANSI256

func strptr(s string) *string {
	return &s
}

func boolptr(b bool) *bool {
	return &b
}

func uintptr(u uint) *uint {
	return &u
}

// CatppuccinStyleConfig returns a Glamour style configuration using Catppuccin colors.
func CatppuccinStyleConfig() gansi.StyleConfig {
	colours := styles.Colours
	defaultMargin := uint(2)

	return gansi.StyleConfig{
		Document: gansi.StyleBlock{
			StylePrimitive: gansi.StylePrimitive{
				BlockPrefix: "\n",
				BlockSuffix: "\n",
				Color:       strptr(colours.Text),
			},
			Margin: uintptr(defaultMargin),
		},
		BlockQuote: gansi.StyleBlock{
			StylePrimitive: gansi.StylePrimitive{},
			Indent:         uintptr(1),
			IndentToken:    strptr("│ "),
		},
		List: gansi.StyleList{
			LevelIndent: defaultMargin,
		},
		Heading: gansi.StyleBlock{
			StylePrimitive: gansi.StylePrimitive{
				BlockSuffix: "\n",
				Color:       strptr(colours.Peach),
				Bold:        boolptr(true),
			},
		},
		H1: gansi.StyleBlock{
			StylePrimitive: gansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				Color:           strptr(colours.Base),
				BackgroundColor: strptr(colours.Peach),
				Bold:            boolptr(true),
			},
		},
		H2: gansi.StyleBlock{
			StylePrimitive: gansi.StylePrimitive{
				Prefix: "## ",
			},
		},
		H3: gansi.StyleBlock{
			StylePrimitive: gansi.StylePrimitive{
				Prefix: "### ",
			},
		},
		H4: gansi.StyleBlock{
			StylePrimitive: gansi.StylePrimitive{
				Prefix: "#### ",
			},
		},
		H5: gansi.StyleBlock{
			StylePrimitive: gansi.StylePrimitive{
				Prefix: "##### ",
			},
		},
		H6: gansi.StyleBlock{
			StylePrimitive: gansi.StylePrimitive{
				Prefix: "###### ",
				Color:  strptr(colours.Subtext1),
				Bold:   boolptr(false),
			},
		},
		Strikethrough: gansi.StylePrimitive{
			CrossedOut: boolptr(true),
		},
		Emph: gansi.StylePrimitive{
			Italic: boolptr(true),
		},
		Strong: gansi.StylePrimitive{
			Bold: boolptr(true),
		},
		HorizontalRule: gansi.StylePrimitive{
			Color:  strptr(colours.Surface2),
			Format: "\n--------\n",
		},
		Item: gansi.StylePrimitive{
			BlockPrefix: "• ",
		},
		Enumeration: gansi.StylePrimitive{
			BlockPrefix: ". ",
		},
		Task: gansi.StyleTask{
			StylePrimitive: gansi.StylePrimitive{},
			Ticked:         "[✓] ",
			Unticked:       "[ ] ",
		},
		Link: gansi.StylePrimitive{
			Color:     strptr(colours.Blue),
			Underline: boolptr(true),
		},
		LinkText: gansi.StylePrimitive{
			Color: strptr(colours.Sapphire),
			Bold:  boolptr(true),
		},
		Image: gansi.StylePrimitive{
			Color:     strptr(colours.Pink),
			Underline: boolptr(true),
		},
		ImageText: gansi.StylePrimitive{
			Color:  strptr(colours.Subtext0),
			Format: "Image: {{.text}} →",
		},
		Code: gansi.StyleBlock{
			StylePrimitive: gansi.StylePrimitive{
				Prefix:          "\u00a0",
				Suffix:          "\u00a0",
				Color:           strptr(colours.Green),
				BackgroundColor: strptr(colours.Surface0),
			},
		},
		CodeBlock: gansi.StyleCodeBlock{
			StyleBlock: gansi.StyleBlock{
				StylePrimitive: gansi.StylePrimitive{
					Color: strptr(colours.Text),
				},
				Margin: uintptr(defaultMargin),
			},
			Chroma: &gansi.Chroma{
				Text: gansi.StylePrimitive{
					Color: strptr(colours.Text),
				},
				Error: gansi.StylePrimitive{
					Color:           strptr(colours.Base),
					BackgroundColor: strptr(colours.Red),
				},
				Comment: gansi.StylePrimitive{
					Color: strptr(colours.Overlay0),
				},
				CommentPreproc: gansi.StylePrimitive{
					Color: strptr(colours.Flamingo),
				},
				Keyword: gansi.StylePrimitive{
					Color: strptr(colours.Blue),
				},
				KeywordReserved: gansi.StylePrimitive{
					Color: strptr(colours.Mauve),
				},
				KeywordNamespace: gansi.StylePrimitive{
					Color: strptr(colours.Pink),
				},
				KeywordType: gansi.StylePrimitive{
					Color: strptr(colours.Mauve),
				},
				Operator: gansi.StylePrimitive{
					Color: strptr(colours.Sky),
				},
				Punctuation: gansi.StylePrimitive{
					Color: strptr(colours.Overlay2),
				},
				Name: gansi.StylePrimitive{
					Color: strptr(colours.Text),
				},
				NameBuiltin: gansi.StylePrimitive{
					Color: strptr(colours.Red),
				},
				NameTag: gansi.StylePrimitive{
					Color: strptr(colours.Mauve),
				},
				NameAttribute: gansi.StylePrimitive{
					Color: strptr(colours.Yellow),
				},
				NameClass: gansi.StylePrimitive{
					Color:     strptr(colours.Yellow),
					Underline: boolptr(true),
					Bold:      boolptr(true),
				},
				NameDecorator: gansi.StylePrimitive{
					Color: strptr(colours.Pink),
				},
				NameFunction: gansi.StylePrimitive{
					Color: strptr(colours.Green),
				},
				LiteralNumber: gansi.StylePrimitive{
					Color: strptr(colours.Peach),
				},
				LiteralString: gansi.StylePrimitive{
					Color: strptr(colours.Yellow),
				},
				LiteralStringEscape: gansi.StylePrimitive{
					Color: strptr(colours.Pink),
				},
				GenericDeleted: gansi.StylePrimitive{
					Color: strptr(colours.Red),
				},
				GenericEmph: gansi.StylePrimitive{
					Italic: boolptr(true),
				},
				GenericInserted: gansi.StylePrimitive{
					Color: strptr(colours.Green),
				},
				GenericStrong: gansi.StylePrimitive{
					Bold: boolptr(true),
				},
				GenericSubheading: gansi.StylePrimitive{
					Color: strptr(colours.Overlay1),
				},
				Background: gansi.StylePrimitive{
					BackgroundColor: strptr(colours.Base),
				},
			},
		},
		Table: gansi.StyleTable{
			StyleBlock: gansi.StyleBlock{
				StylePrimitive: gansi.StylePrimitive{},
			},
		},
		DefinitionDescription: gansi.StylePrimitive{
			BlockPrefix: "\n🠶 ",
		},
	}
}

// StyleConfig returns the default Glamour style configuration.
func StyleConfig() gansi.StyleConfig {
	return CatppuccinStyleConfig()
}

// StyleRenderer returns a new Glamour renderer.
func StyleRenderer() gansi.RenderContext {
	return StyleRendererWithStyles(StyleConfig())
}

// StyleRendererWithStyles returns a new Glamour renderer.
func StyleRendererWithStyles(styles gansi.StyleConfig) gansi.RenderContext {
	return gansi.NewRenderContext(gansi.Options{
		Styles: styles,
	})
}
