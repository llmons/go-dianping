package v1

type (
	SendCodeReq struct {
		Phone string `form:"phone" binding:"required"`
	}
)

type (
	LoginReq struct {
		Phone    string  `json:"phone" example:"13456789001" binding:"required"`
		Code     string  `json:"code" example:"123456" binding:"required"`
		Password *string `json:"password" example:""`
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
	MeRespData SimpleUser
	MeResp     struct {
		Response
		Data MeRespData `json:"data"`
	}
)
