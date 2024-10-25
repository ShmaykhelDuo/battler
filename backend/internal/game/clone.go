package game

func Clone(c, opp *Character) (clonedC, clonedOpp *Character) {
	clonedC = CloneCharacter(c)
	clonedOpp = CloneCharacter(opp)
	return
}

func CloneCharacter(c *Character) *Character {
	cloned := &Character{}
	*cloned = *c

	for i, s := range c.skills {
		clonedSkill := &Skill{}
		*clonedSkill = *s
		clonedSkill.c = cloned
		cloned.skills[i] = clonedSkill

		if c.lastUsedSkill == s {
			cloned.lastUsedSkill = clonedSkill
		}
	}

	cloned.effects = make(map[EffectDescription]Effect, len(c.effects))
	for i, e := range c.effects {
		cloned.effects[i] = e.Clone()
	}
	return cloned
}
