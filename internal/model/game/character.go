package game

import "github.com/ShmaykhelDuo/battler/internal/game"

type AvailableCharacter struct {
	Number          int
	Level           int
	LevelExperience int
	MatchCount      int
	WinCount        int
}

type Character struct {
	Character *game.CharacterData
	Rarity    CharacterRarity
}
