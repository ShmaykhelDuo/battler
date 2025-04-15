package game

var CharacterLevelCaps = map[CharacterRarity][]int{
	CharacterRarityST: {1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 3, 4, 4, 5, 5, 5, 5},
	CharacterRarityAD: {1, 1, 2, 2, 2, 2, 3, 3, 4, 4, 4, 4, 4, 5, 5, 6, 6, 6, 6},
	CharacterRaritySP: {1, 2, 4, 5, 5, 5, 6, 8, 9, 11, 11, 11, 11, 11, 13, 14, 16, 16, 16},
	CharacterRarityRP: {2, 6, 10, 13, 13, 14, 14, 16, 20, 23, 24, 24, 25, 27, 29, 31, 35, 37, 37},
	CharacterRarityLF: {4, 7, 12, 16, 17, 18, 20, 23, 27, 30, 31, 31, 31, 33, 34, 38, 41, 43, 44},
}

var ResultModifiers = map[CharacterRarity]int64{
	CharacterRarityST: 1,
	CharacterRarityAD: 2,
	CharacterRaritySP: 3,
	CharacterRarityRP: 4,
	CharacterRarityLF: 5,
}
