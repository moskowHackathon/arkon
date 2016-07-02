package routers

import (
	"github.com/astaxie/beego"
	"github.com/moskowHackathon/arkon/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/start/:id(\\w+)/", &controllers.StartController{})
	beego.Router("/question/:id(\\w+)/", &controllers.QuestionController{})
	beego.Router("/answer/:id(\\w+)/", &controllers.AnswerController{})
	beego.Router("/start/:id(\\w+)/", &controllers.CompleteController{})
}
