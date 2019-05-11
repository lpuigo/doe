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

func (p *Pulling) GetTotalDist() int {
	dist := 0
	for _, chunk := range p.Chuncks {
		dist += chunk.LoveDist + chunk.UndergroundDist + chunk.AerialDist + chunk.BuildingDist
	}
	return dist
}
