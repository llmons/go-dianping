package v1

type (
	IsFollowResp struct {
		Response
		Data bool `json:"data"`
	}
)

type (
	FollowCommonsResp struct {
		Response
		Data []*SimpleUser `json:"data"`
	}
)
