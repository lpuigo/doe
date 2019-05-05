package ripsites

type PullingChunk struct {
	StartingTronconName string
	StartingNodeName    string
	EndingTronconName   string
	EndingNodeName      string
	LoveDist            int
	UndergroundDist     int
	AerialDist          int
	BuildingDist        int
	State               State
}

type Pulling struct {
	CableName string
	Chuncks   []PullingChunk
	State     State
}
