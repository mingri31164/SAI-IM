package group

import (
	"net/http"

	"SAI-IM/apps/social/api/internal/logic/group"
	"SAI-IM/apps/social/api/internal/svc"
	"SAI-IM/apps/social/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GroupListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupListRep
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := group.NewGroupListLogic(r.Context(), svcCtx)
		resp, err := l.GroupList(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
