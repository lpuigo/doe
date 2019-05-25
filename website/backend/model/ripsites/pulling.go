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

func (p *Pulling) GetTotalAggrDist() int {
	dist := 0
	for _, chunk := range p.Chuncks {
		dist += chunk.LoveDist + chunk.UndergroundDist + chunk.AerialDist + chunk.BuildingDist
	}
	return dist
}

func (p *Pulling) GetTotalDists() (love, underground, aerial, building int) {
	for _, chunk := range p.Chuncks {
		love += chunk.LoveDist
		underground += chunk.UndergroundDist
		aerial += chunk.AerialDist
		building += chunk.BuildingDist
	}
	return
}
