package ssa

import "log"

func branchingPhiCheck(f *Func) {
	for _, b := range f.Blocks {
		if b.Kind != BlockIf || len(b.Preds) != 2 {
			continue
		}

		var ind *Value // induction variable

		c := b.Controls[0]
		switch c.Op {
		case OpLeq64, OpLeq32, OpLeq16, OpLeq8:
			fallthrough
		case OpLess64, OpLess32, OpLess16, OpLess8:
			ind = c.Args[0]
		default:
			continue
		}

		if checkForIncrementInBranchingPhi(ind) {
			log.Println("Detected induction variable with phi argument that relies on equivalent adds.")
		}

	}

}

func isPhiWithEquivalentAdds(phi *Value, ind *Value) bool {
	if phi.Op != OpPhi {
		return false
	}

	val1, val2 := phi.Args[0], phi.Args[1]

	//Check if both arguments of the phi are adds
	if !(val1.Op == OpAdd64 || val1.Op == OpAdd32 || val1.Op == OpAdd16 || val1.Op == OpAdd8) {
		return false
	}
	if !(val2.Op == OpAdd64 || val2.Op == OpAdd32 || val2.Op == OpAdd16 || val2.Op == OpAdd8) {
		return false
	}

	//Check if the adds are equivalent (operation and arguments)
	if val1.Op != val2.Op {
		return false
	}
	if !(val1.Args[0] == val2.Args[0] && val1.Args[1] == val2.Args[1]) {
		return false
	}

	//Check if the adds refer back to the original phi variable
	if !((val1.Args[0] == ind || val1.Args[1] == ind) && (val2.Args[0] == ind || val2.Args[1] == ind)) {
		return false
	}

	return true

}

func checkForIncrementInBranchingPhi(ind *Value) bool {
	if ind.Op != OpPhi {
		return false
	}

	if n := ind.Args[1]; isPhiWithEquivalentAdds(n, ind) {
		return true
	} else if n := ind.Args[0]; isPhiWithEquivalentAdds(n, ind) {
		return true
	}

	return false
}
