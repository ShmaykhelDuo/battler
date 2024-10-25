package ml2

type State struct {
	PrevMoves []int
	First     bool
}

func (s State) ToSlice() []int {
	intRec := make([]int, 21)
	if s.First {
		intRec[0] = 1
	}
	ourAdd := 0
	oppAdd := 1
	if !s.First {
		oppAdd = 0
		ourAdd = 1
	}
	for i := range 10 {
		if i*2+ourAdd < len(s.PrevMoves) {
			intRec[1+i] = s.PrevMoves[i*2+ourAdd]
		} else {
			intRec[1+i] = -1
		}

		if i*2+oppAdd < len(s.PrevMoves) {
			intRec[11+i] = s.PrevMoves[i*2+oppAdd]
		} else {
			intRec[11+i] = -1
		}
	}
	return intRec
}
