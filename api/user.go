package api

type (
	SendCodeReq struct {
		Phone string `form:"phone" binding:"required"`
	}

	SendCodeResp response
)

type (
	LoginReq struct {
		Phone    string `form:"phone" binding:"required"`
		Code     string `form:"code" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	LoginResp struct {
		Token string `json:"token"`
	}
)

type (
	SimpleUser struct {
		Id       uint   `json:"id"`
		NickName string `json:"nickname"`
		Icon     string `json:"icon"`
	}
	GetMeResp SimpleUser
)
