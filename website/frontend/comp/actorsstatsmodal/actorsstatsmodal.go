package actorsstatsmodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/actorsstatschart"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/model/actor"
	rs "github.com/lpuig/ewin/doe/website/frontend/model/ripsite"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
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
					<i class="fas fa-chart-area icon--left"></i>Tailles des équipes</span>
				</h2>
			</el-col>
		</el-row>
    </span>

    <!-- 
        Modal Body
        style="height: 100%;"
        
    -->
    <div style="height: 45vh; padding: 5px 25px; overflow-x: hidden;overflow-y: auto;">
        <!-- Last & First Name -->
        <el-row :gutter="10" align="middle" class="spaced" type="flex">
            <el-col :span="3" class="align-right">Nombres d'équipier :</el-col>
            <el-col :span="9">{{Actors.length}}</el-col>
        </el-row>
		<actors-stats-chart :stats="Stats"></actors-stats-chart>
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
		//hvue.Computed("currentActorRef", func(vm *hvue.VM) interface{} {
		//	aumm := ActorUpdateModalModelFromJS(vm.Object)
		//	aumm.CurrentActor.Ref = aumm.CurrentActor.LastName + " " + aumm.CurrentActor.FirstName
		//	return aumm.CurrentActor.Ref
		//}),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Model Methods

type ActorsStatsModalModel struct {
	*js.Object

	Visible bool     `js:"visible"`
	VM      *hvue.VM `js:"VM"`

	User   *fm.User       `js:"user"`
	Actors []*actor.Actor `js:"Actors"`
	Stats  *rs.TeamStats  `js:"Stats"`
}

func NewActorsStatsModalModel(vm *hvue.VM) *ActorsStatsModalModel {
	aumm := &ActorsStatsModalModel{Object: tools.O()}
	aumm.Visible = false
	aumm.VM = vm

	aumm.User = fm.NewUser()
	aumm.Actors = []*actor.Actor{}
	aumm.Stats = rs.NewTeamStats()

	return aumm
}

func ActorsStatsModalModelFromJS(o *js.Object) *ActorsStatsModalModel {
	return &ActorsStatsModalModel{Object: o}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

func (asmm *ActorsStatsModalModel) Show(actors []*actor.Actor, user *fm.User) {
	asmm.User = user
	asmm.Actors = actors
	asmm.CalcStats()
	asmm.Visible = true
}

func (asmm *ActorsStatsModalModel) Hide() {
	asmm.Visible = false
	asmm.Actors = []*actor.Actor{}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Business Methods

func (asmm *ActorsStatsModalModel) CalcStats() {
	asmm.Stats.Dates = []string{"2019-09-23", "2019-09-30", "2019-10-07", "2019-10-14", "2019-10-21", "2019-10-28"}
	asmm.Stats.Team = "Ewin Services"

	values := map[string]map[string][]float64{}
	values["actors"] = map[string][]float64{
		"test": []float64{15, 16, 12, 14, 21, 23},
	}

	asmm.Stats.Values = values
}
