package ripsites

type PullingChunk struct {
	TronconName      string
	StartingNodeName string
	EndingNodeName   string
	LoveDist         int
	UndergroundDist  int
	AerialDist       int
	BuildingDist     int
	State            State
}

type Pulling struct {
	CableName string
	Chuncks   []PullingChunk
	State     State
}
