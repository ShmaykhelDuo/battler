package gamestate

import "github.com/ShmaykhelDuo/battler/internal/game"

type Colour int

const (
	ColourRed    Colour = 1
	ColourOrange Colour = 2
	ColourYellow Colour = 3
	ColourGreen  Colour = 4
	ColourCyan   Colour = 5
	ColourBlue   Colour = 6
	ColourViolet Colour = 7
	ColourPink   Colour = 8
	ColourGray   Colour = 9
	ColourBrown  Colour = 10
	ColourBlack  Colour = 11
	ColourWhite  Colour = 12
)

func NewColour(c game.Colour) Colour {
	switch c {
	case game.ColourRed:
		return ColourRed
	case game.ColourOrange:
		return ColourOrange
	case game.ColourYellow:
		return ColourYellow
	case game.ColourGreen:
		return ColourGreen
	case game.ColourCyan:
		return ColourCyan
	case game.ColourBlue:
		return ColourBlue
	case game.ColourViolet:
		return ColourViolet
	case game.ColourPink:
		return ColourPink
	case game.ColourGray:
		return ColourGray
	case game.ColourBrown:
		return ColourBrown
	case game.ColourBlack:
		return ColourBlack
	case game.ColourWhite:
		return ColourWhite
	default:
		return 0
	}
}
