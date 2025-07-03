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
		ID       *uint64 `json:"id" redis:"id"`
		NickName *string `json:"nickName" redis:"nickName"`
		Icon     *string `json:"icon" redis:"icon"`
	}
	MeResp struct {
		Response
		Data SimpleUser `json:"data"`
	}
)
