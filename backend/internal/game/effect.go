package game

// EffectDescription contains main features of an effect.
type EffectDescription struct {
	Name string // effect's name
}

// Effect represents an effect to be applied to a character.
type Effect interface {
	// Desc returns the effect's description.
	Desc() EffectDescription
}
