package ripsites

type Operation struct {
	Type        string
	TronconName string
	NbFiber     int
	NbSplice    int
	State       State
}

type Junction struct {
	NodeName   string
	Operations []Operation
	State      State
}

func (j *Junction) GetNbFiber() int {
	nbFiber := 0
	for _, ope := range j.Operations {
		nbFiber += ope.NbFiber
	}
	return nbFiber
}
