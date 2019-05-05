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
