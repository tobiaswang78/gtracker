package handle

import (
	"encoding/json"
	"itflow/cache"
	"itflow/db"
	"itflow/internal/datalog"
	"itflow/internal/project"
	"itflow/internal/response"
	"net/http"
	"strconv"

	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
)

func ProjectList(w http.ResponseWriter, r *http.Request) {

	projects := &project.ProjectList{}

	for k, v := range cache.CachePidName {
		pr := &project.Project{
			Id:          k,
			ProjectName: v,
		}
		projects.Plist = append(projects.Plist, pr)
	}

	send, _ := json.Marshal(projects)
	w.Write(send)
	return

}

func AddProject(w http.ResponseWriter, r *http.Request) {

	errorcode := &response.Response{}

	name := r.FormValue("name")

	getaritclesql := "insert into projectname(name) values(?)"
	var err error
	errorcode.Id, err = db.Mconn.Insert(getaritclesql, name)
	if err != nil {
		golog.Error(err)
		w.Write(errorcode.ErrorE(err))
		return
	}
	// 增加日志
	nickname := xmux.GetData(r).Get("nickname").(string)
	xmux.GetData(r).End = &datalog.AddLog{
		Ip:       r.RemoteAddr,
		Username: nickname,
		Classify: "project",
		Action:   "add",
	}

	// 更新缓存
	cache.CacheProjectPid[name] = errorcode.Id
	cache.CachePidName[errorcode.Id] = name
	send, _ := json.Marshal(errorcode)
	w.Write(send)
	return

}

func UpdateProject(w http.ResponseWriter, r *http.Request) {

	errorcode := &response.Response{}

	pr := xmux.GetData(r).Data.(*project.Project)

	getaritclesql := "update projectname set name=? where id=?"

	_, err := db.Mconn.Update(getaritclesql, pr.ProjectName, pr.Id)
	if err != nil {
		golog.Error(err)
		w.Write(errorcode.ErrorE(err))
		return
	}
	nickname := xmux.GetData(r).Get("nickname").(string)
	xmux.GetData(r).End = &datalog.AddLog{
		Ip:       r.RemoteAddr,
		Username: nickname,
		Classify: "project",
		Action:   "update",
	}

	// 更新缓存
	delete(cache.CacheProjectPid, cache.CachePidName[int64(pr.Id)])
	cache.CachePidName[int64(pr.Id)] = pr.ProjectName
	cache.CacheProjectPid[pr.ProjectName] = int64(pr.Id)

	send, _ := json.Marshal(errorcode)
	w.Write(send)
	return

}

func DeleteProject(w http.ResponseWriter, r *http.Request) {

	errorcode := &response.Response{}

	id := r.FormValue("id")
	pid, err := strconv.Atoi(id)
	if err != nil {
		golog.Error(err)
		w.Write(errorcode.ErrorE(err))
		return
	}
	// 是否有bug使用
	var count int
	err = db.Mconn.GetOne("select count(id) from bugs where pid=?", id).Scan(&count)
	if err != nil {
		golog.Error(err)
		w.Write(errorcode.ErrorE(err))
		return
	}

	if count > 0 {
		w.Write(errorcode.IsUse())
		return
	}

	getaritclesql := "delete from projectname where id=?"

	_, err = db.Mconn.Insert(getaritclesql, id)
	if err != nil {
		golog.Error(err)
		w.Write(errorcode.ErrorE(err))
		return
	}
	nickname := xmux.GetData(r).Get("nickname").(string)
	xmux.GetData(r).End = &datalog.AddLog{
		Ip:       r.RemoteAddr,
		Username: nickname,
		Classify: "project",
		Action:   "delete",
	}

	// 更新缓存
	delete(cache.CacheProjectPid, cache.CachePidName[int64(pid)])
	delete(cache.CachePidName, int64(pid))

	send, _ := json.Marshal(errorcode)
	w.Write(send)
	return

}
