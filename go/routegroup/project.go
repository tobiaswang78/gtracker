package routegroup

import (
	"itflow/handle"
	"itflow/internal/project"
	"itflow/midware"

	"github.com/hyahm/xmux"
)

var Project *xmux.GroupRoute

func init() {
	Project = xmux.NewGroupRoute().ApiCreateGroup("project", "项目相关接口（建议给项目发起者添加操作）", "项目管理")
	Project.ApiCodeField("code").ApiCodeMsg("0", "成功")
	Project.ApiCodeField("code").ApiCodeMsg("20", "token过期")
	Project.ApiCodeField("code").ApiCodeMsg("2", "系统错误")
	Project.ApiCodeField("code").ApiCodeMsg("", "其他错误,请查看返回的msg")
	Project.ApiReqHeader("X-Token", "xxxxxxxxxxxxxxxxxxxxxxxxxx")
	Project.Pattern("/project/list").Post(handle.ProjectList).ApiDescribe("获取所有列表")

	Project.Pattern("/project/add").Post(handle.AddProject).Bind(&project.ReqProject{}).
		AddMidware(midware.JsonToStruct).End(midware.EndLog).
		ApiDescribe("增加项目")

	Project.Pattern("/project/update").Post(handle.UpdateProject).Bind(&project.ReqProject{}).AddMidware(midware.JsonToStruct).
		End(midware.EndLog).ApiDescribe("修改项目")

	Project.Pattern("/project/delete").Get(handle.DeleteProject).End(midware.EndLog).
		ApiDescribe("删除项目")

	Project.Pattern("/get/project").Post(handle.GetProject).ApiDescribe("获取所有项目名")
	Project.Pattern("/get/myproject").Post(handle.GetMyProject).ApiDescribe("获取自己权限的项目名")
}
