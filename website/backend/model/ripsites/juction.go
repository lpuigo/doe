package ripsites

type Operation struct {
	Type        string
	TronconName string
	NbFiber     int
	NbSplice    int
	State       State
}

type Junction struct {
	TronconName string
	Operations  []Operation
	State       State
}
