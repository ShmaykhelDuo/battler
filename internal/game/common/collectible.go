package common

import "github.com/ShmaykhelDuo/battler/internal/game"

// Collectible is a mixin that allows storing amounts.
type Collectible struct {
	amount int
}

// NewCollectible returns a new collectible with specified amount.
func NewCollectible(amount int) *Collectible {
	return &Collectible{amount: amount}
}

// Amount returns the collectible's amount.
func (c *Collectible) Amount() int {
	return c.amount
}

// Increase increases the collectible's amount.
func (c *Collectible) Increase(amount int) {
	c.amount += amount
}

// Decrease decreases the collectible's amount.
func (c *Collectible) Decrease(amount int) {
	c.amount -= amount
}

// HasExpired reports whether the effect has expired.
func (c *Collectible) HasExpired(turnState game.TurnState) bool {
	return c.amount == 0
}
