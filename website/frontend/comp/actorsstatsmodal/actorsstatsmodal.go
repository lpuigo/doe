package actorsstatsmodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/actorsstatschart"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor"
	rs "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/date"
)

const template string = `<el-dialog
        :before-close="Hide"
        :visible.sync="visible"
        width="70%"
>
    <!-- 
        Modal Title
    -->
    <span slot="title">
		<el-row :gutter="10" align="middle" type="flex">
			<el-col :span="12">
				<h2 style="margin: 0 0">
					<i class="fas fa-chart-area icon--left"></i>Nombre d'Acteurs
				</h2>
			</el-col>
		</el-row>
    </span>

    <!-- 
        Modal Body
        style="height: 100%;"
        
    -->
    <div style="height: 65vh; padding: 5px 25px; overflow-x: hidden;overflow-y: auto;">
        <el-row type="flex" align="middle">
            <el-col :offset="13" :span="1">
                <el-button icon="fas fa-chevron-down" size="mini" @click="ResetCurrentDate()"></el-button>
            </el-col>
            <el-col :span="10">
                <el-button icon="fas fa-chevron-left" size="mini" @click="CurrentDateBefore()"></el-button>
                <span style="margin: 0px 10px">{{CurrentDateRange()}}</span>
                <el-button icon="fas fa-chevron-right" size="mini" @click="CurrentDateAfter()"></el-button>
            </el-col>
        </el-row>
		<actors-stats-chart ref="Chart" :stats="Stats" style="margin-top: 5px"></actors-stats-chart>
    </div>

    <!-- 
        Modal Footer Action Bar
    -->
    <span slot="footer">
		<el-row :gutter="15">
			<el-col :span="24" style="text-align: right">
				<el-button @click="Hide" size="mini">Fermer</el-button>
			</el-col>
		</el-row>
	</span>
</el-dialog>`

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("actors-stats-modal", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		actorsstatschart.RegisterComponent(),
		hvue.Template(template),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewActorsStatsModalModel(vm)
		}),
		hvue.MethodsOf(&ActorsStatsModalModel{}),
		//hvue.Computed("UpdatedStats", func(vm *hvue.VM) interface{} {
		//	asmm := ActorsStatsModalModelFromJS(vm.Object)
		//	asmm.UpdateChart()
		//	return asmm.Stats
		//}),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Model Methods

type ActorsStatsModalModel struct {
	*js.Object

	Visible bool     `js:"visible"`
	VM      *hvue.VM `js:"VM"`

	User        *fm.User       `js:"user"`
	Actors      []*actor.Actor `js:"Actors"`
	Stats       *rs.TeamStats  `js:"Stats"`
	CurrentDate string         `js:"CurrentDate"`
	DateRange   int            `js:"DateRange"`
}

func NewActorsStatsModalModel(vm *hvue.VM) *ActorsStatsModalModel {
	asmm := &ActorsStatsModalModel{Object: tools.O()}
	asmm.Visible = false
	asmm.VM = vm

	asmm.User = fm.NewUser()
	asmm.Actors = []*actor.Actor{}
	asmm.Stats = rs.NewTeamStats()
	asmm.DateRange = 12
	asmm.ResetCurrentDate()

	return asmm
}

func ActorsStatsModalModelFromJS(o *js.Object) *ActorsStatsModalModel {
	return &ActorsStatsModalModel{Object: o}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

func (asmm *ActorsStatsModalModel) Show(actors []*actor.Actor, user *fm.User) {
	asmm.User = user
	asmm.Actors = actors
	asmm.UpdateChart()
	asmm.Visible = true
}

func (asmm *ActorsStatsModalModel) Hide() {
	asmm.Visible = false
	asmm.Stats = rs.NewTeamStats()
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Action Methods

func (asmm *ActorsStatsModalModel) UpdateChart() {
	asmm.CalcStats()
	if asmm.Visible {
		actorsstatschart.ActorsStatsChartFromJS(asmm.VM.Refs("Chart")).SetChart(asmm.Stats)
	}
}

func (asmm *ActorsStatsModalModel) ResetCurrentDate() {
	asmm.CurrentDate = date.GetMonday(date.TodayAfter(-asmm.DateRange / 2 * 7))
	asmm.UpdateChart()
}

func (asmm *ActorsStatsModalModel) CurrentDateBefore() {
	asmm.CurrentDate = date.After(asmm.CurrentDate, -7)
	asmm.UpdateChart()
}

func (asmm *ActorsStatsModalModel) CurrentDateAfter() {
	asmm.CurrentDate = date.After(asmm.CurrentDate, 7)
	asmm.UpdateChart()
}

func (asmm *ActorsStatsModalModel) CurrentDateRange() string {
	return "du " + date.DateString(asmm.CurrentDate) + " au " + date.DateString(asmm.CurrentRangeEnd())
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Business Methods

func (asmm *ActorsStatsModalModel) CurrentRangeEnd() string {
	return date.After(asmm.CurrentDate, asmm.DateRange*7-7)
}

//func (asmm *ActorsStatsModalModel) DateOf(i int) string {
//	return date.Day(date.After(asmm.CurrentDate, i))
//}
//
func (asmm *ActorsStatsModalModel) CalcDates() []string {
	res := make([]string, asmm.DateRange)
	day := asmm.CurrentDate
	res[0] = day
	for i := 1; i < asmm.DateRange; i++ {
		day = date.After(day, 7)
		res[i] = day
	}
	return res
}

func (asmm *ActorsStatsModalModel) CalcStatsValues(stats *rs.TeamStats) {
	cs := actor.NewCalendarSeeker(stats.Dates)
	for _, actr := range asmm.Actors {
		cs.Append(actr)
	}
	nbEmployees, nbActing := cs.CalcStats()

	values := map[string]map[string][]float64{}
	values["employees"] = map[string][]float64{
		"employés": nbEmployees,
	}
	values["acting"] = map[string][]float64{
		"présents": nbActing,
	}
	stats.Values = values
}

func (asmm *ActorsStatsModalModel) CalcStats() {
	asmm.Stats = rs.NewTeamStats()
	asmm.Stats.Dates = asmm.CalcDates()
	asmm.Stats.Team = "Ewin Services"
	asmm.CalcStatsValues(asmm.Stats)
}
