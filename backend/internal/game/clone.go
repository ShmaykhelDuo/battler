package game

import clone "github.com/huandu/go-clone/generic"

func Clone(c, opp *Character) (clonedC, clonedOpp *Character) {
	clonedC = cloneCharacter(c)
	clonedOpp = cloneCharacter(opp)
	return
}

func cloneCharacter(c *Character) *Character {
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

	cloned.effects = make([]Effect, len(c.effects))
	for i, e := range c.effects {
		cloned.effects[i] = clone.Clone(e)
	}
	return cloned
}
