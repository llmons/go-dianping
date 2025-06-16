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
	LoginRespData struct {
		Token string `json:"token"`
	}
 LoginResp struct {
	 response
  Data LoginRespData `json:"data"`
	}
)

type (
	SimpleUser struct {
		Id       uint   `json:"id"`
		NickName string `json:"nickname"`
		Icon     string `json:"icon"`
	}
 GetMeRespData SimpleUser
	GetMeResp struct{
  response
  Data GetMeRespData `json:"data"`
 }
)
