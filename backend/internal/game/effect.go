package game

// EffectType represents a type of an effect.
type EffectType int

const (
	EffectTypeBasic EffectType = iota
	EffectTypeProhibiting
	EffectTypeControl
	EffectTypeBuff
	EffectTypeDebuff
	EffectTypeState
	EffectTypeNumeric
)

// EffectDescription contains main features of an effect.
type EffectDescription struct {
	Name string     // effect's name
	Type EffectType // effect's type
}

// Effect represents an effect to be applied to a character.
type Effect interface {
	// Desc returns the effect's description.
	Desc() EffectDescription
}
