package v1

type (
	SendCodeReq struct {
		Phone string `form:"phone" binding:"required"`
	}
)

type (
	LoginReq struct {
		Phone    string  `json:"phone" binding:"required"`
		Code     string  `json:"code" binding:"required"`
		Password *string `json:"password"`
	}
	LoginRespData struct {
		Token string `json:"token"`
	}
	LoginResp struct {
		Response
		Data LoginRespData `json:"data"`
	}
)

type (
	SimpleUser struct {
		ID       int64   `json:"id"`
		NickName *string `json:"nickname"`
		Icon     *string `json:"icon"`
	}
	GetMeRespData SimpleUser
	GetMeResp     struct {
		Response
		Data GetMeRespData `json:"data"`
	}
)
