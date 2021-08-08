package handle

import (
	"itflow/model"
	"itflow/response"
	"net/http"

	"github.com/hyahm/golog"
	"github.com/hyahm/xmux"
)

func UserKeyName(w http.ResponseWriter, r *http.Request) {
	// 获取用户keyvalue
	uid := xmux.GetInstance(r).Get("uid").(int64)
	kns, err := model.GetUserKeyName(uid)
	if err != nil {
		golog.Error(err)
		w.Write(response.ErrorE(err))
		return
	}
	res := response.Response{
		Data: kns,
	}
	w.Write(res.Marshal())
}
