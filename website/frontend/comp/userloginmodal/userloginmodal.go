package userloginmodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteedit"
	"github.com/lpuig/ewin/doe/website/frontend/comp/worksiteinfo"
	fm "github.com/lpuig/ewin/doe/website/frontend/model"
	"github.com/lpuig/ewin/doe/website/frontend/tools"
	"github.com/lpuig/ewin/doe/website/frontend/tools/elements/message"
	"honnef.co/go/js/xhr"
	"strconv"
)

const (
	compname        = "user-login-modal"
	template string = `

<el-dialog 
		:visible.sync="visible" 
		width="450px"
>
	<!-- 
		Modal Title
	-->
    <span slot="title">
		<el-row :gutter="10" type="flex" align="middle">
			<el-col :span="24">
				<h2 style="margin: 0 0">
					<i class="fas fa-sign-in-alt icon--left"></i>Connexion Utilisateur
				</h2>
			</el-col>
		</el-row>
    </span>

	<!-- 
		Modal Body
	-->
    <el-form id="userForm" :model="user" size="mini" style="margin: 30px">
        <el-form-item label="Login" :label-width="labelSize">
            <el-input v-model="user.Name" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="Mot de Passe" :label-width="labelSize">
            <el-input type="password" v-model="user.Pwd" autocomplete="off"></el-input>
        </el-form-item>
    </el-form>

	<!-- 
		Body Action Bar
	-->
    <span slot="footer">
        <el-button size="mini" @click="visible = false">Abandon</el-button>
        <el-button type="primary" size="mini" @click="Submit">Confirmer</el-button>
    </span>
</el-dialog>
`
)

type UserLoginModalModel struct {
	*js.Object

	Visible bool     `js:"visible"`
	VM      *hvue.VM `js:"VM"`

	User      *fm.User `js:"user"`
	LabelSize string   `js:"labelSize"`
}

func NewUserLoginModalModel(vm *hvue.VM) *UserLoginModalModel {
	ulmm := &UserLoginModalModel{Object: tools.O()}
	ulmm.Visible = false
	ulmm.VM = vm

	ulmm.User = fm.NewUser()
	ulmm.LabelSize = "120px"

	return ulmm
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

func Register() {
	hvue.NewComponent(compname,
		ComponentOptions()...,
	)
}

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component(compname, ComponentOptions()...)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		worksiteedit.RegisterComponent(),
		worksiteinfo.RegisterComponent(),
		hvue.Template(template),
		hvue.Props("user"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewUserLoginModalModel(vm)
		}),
		hvue.MethodsOf(&UserLoginModalModel{}),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

func (ulmm *UserLoginModalModel) Show(u *fm.User) {
	ulmm.User = u
	ulmm.Visible = true
}

func (ulmm *UserLoginModalModel) Hide() {
	ulmm.Visible = false
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Action Button Methods

func (ulmm *UserLoginModalModel) Submit() {
	ulmm.VM.Emit("update:user", ulmm.User)

	go ulmm.submitLogin()
	//ulmm.Hide()
}

func (ulmm *UserLoginModalModel) submitLogin() {
	f := js.Global.Get("FormData").New()
	f.Call("append", "user", ulmm.User.Name)
	f.Call("append", "pwd", ulmm.User.Pwd)
	print("submitLogin", f)
	req := xhr.NewRequest("POST", "/api/login")
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(f)
	if err != nil {
		message.ErrorStr(ulmm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != 200 {
		message.SetDuration(tools.WarningMsgDuration)
		msg := "Quelquechose c'est mal passé !\n"
		msg += "Le server retourne un code " + strconv.Itoa(req.Status) + "\n"
		message.ErrorMsgStr(ulmm.VM, msg, req.Response, true)
		return
	}
	message.SetDuration(tools.SuccessMsgDuration)
	message.SuccesStr(ulmm.VM, "User connecté")
}
