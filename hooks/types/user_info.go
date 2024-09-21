package types

type AbilityScope struct {
	Web struct {
		Access         bool `json:"access"`
		RoleSetting    bool `json:"roleSetting"`
		ManageGeometry bool `json:"manageGeometry"`
	} `json:"web"`
	Mobile struct {
		Access        bool `json:"access"`
		VideosProcess bool `json:"videosProcess"`
	} `json:"mobile"`
}

type UserInfo struct {
	OrgId        string       `json:"orgId"`
	OrgName      string       `json:"orgName"`
	UserName     string       `json:"userName"`
	UserId       string       `json:"userId"`
	AbilityScope AbilityScope `json:"abilityScope"`
}
