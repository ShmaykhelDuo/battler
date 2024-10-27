package match

type Player interface {
	SendState(state GameState) error
	SendError() error
	SendEnd() error
	RequestSkill() (int, error)
}
