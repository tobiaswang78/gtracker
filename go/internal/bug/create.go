package bug

import (
	"encoding/json"
	"errors"
	"itflow/cache"
	"itflow/internal/assist"
	"itflow/model"
)

// 前端编辑时需要的数据结构
type RespEditBug struct {
	Status      cache.Status    `json:"status"`
	Title       string          `json:"title"`
	Content     string          `json:"content"`
	Id          int64           `json:"id"`
	Selectusers assist.Names    `json:"selectuser"`
	Important   cache.Important `json:"important"`
	Level       cache.Level     `json:"level"`
	Projectname cache.Project   `json:"projectname"`
	Envname     cache.Env       `json:"envname"`
	Version     cache.Version   `json:"version"`
	Code        int             `json:"code"`
	Msg         string          `json:"message,omitempty"`
}

func (reb *RespEditBug) Marshal() []byte {
	if reb == nil {
		return nil
	}
	send, _ := json.Marshal(reb)
	return send
}

func (reb *RespEditBug) ToBug() (*model.Bug, error) {
	// 将获取的数据转为可以存表的数据
	bug := &model.Bug{}
	bug.ID = reb.Id
	reb.Envname = reb.Envname.Trim()
	reb.Important = reb.Important.Trim()
	bug.Content = reb.Content
	bug.Title = reb.Title
	reb.Level = reb.Level.Trim()
	reb.Projectname = reb.Projectname.Trim()
	reb.Version = reb.Version.Trim()
	if reb.Envname == "" || reb.Important == "" || reb.Level == "" ||
		reb.Projectname == "" || reb.Version == "" {
		return nil, errors.New("all name not by empty")
	}
	var ok bool
	bug.LevelId = reb.Level.Id()
	if bug.LevelId == 0 {
		return nil, errors.New("没有找到level key")
	}
	bug.ProjectId = reb.Projectname.Id()
	if bug.ProjectId == 0 {
		return nil, errors.New("没有找到project key")
	}
	//
	if bug.EnvId, ok = cache.CacheEnvEid[reb.Envname]; !ok {
		return nil, errors.New("没有找到env key")
	}
	//
	bug.ImportanceId = reb.Important.Id()
	if bug.ImportanceId == 0 {
		return nil, errors.New("没有找到important key")
	}
	if bug.VersionId, ok = cache.CacheVersionVid[reb.Version]; !ok {
		return nil, errors.New("没有找到version key")
	}

	bug.OprateUsers = reb.Selectusers.RealNameToUsersId()
	return bug, nil
}
