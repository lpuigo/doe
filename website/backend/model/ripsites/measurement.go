package ripsites

type Measurement struct {
	DestNodeName string
	NbFiber      int // DestNode.Operation.Attente.NbFiber
	NbOK         int
	NbWarn1      int
	NbWarn2      int
	NbKO         int
	//NbEvent      int // == len(NodeNames)
	Dist      int // DestNode.DistFromPM
	NodeNames []string
	State     State
}

func (m *Measurement) GetNbMeas() int {
	return m.NbFiber
}

func (m *Measurement) NbSplice() int {
	return len(m.NodeNames)
}
