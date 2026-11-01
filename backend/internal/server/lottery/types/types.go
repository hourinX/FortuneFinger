package types

//region 抽奖活动出入参

type IDPath struct {
	ID uint64 `uri:"id" binding:"required"`
}

type GetActivityListReq struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
	Name     string `form:"name" binding:"max=100"`
	OrgID    uint64 `form:"org_id"`
	Date     string `form:"date" binding:"omitempty,datetime=2006-01-02"`
	Phase    string `form:"phase" binding:"omitempty,oneof=pending active finished"`
	Status   uint8  `form:"status" binding:"omitempty,oneof=1 2"`
}

type GetActivityListResp struct {
	List  []ActivityListItem `json:"list"`
	Total int64              `json:"total"`
}

type ActivityListItem struct {
	ID                uint64 `json:"id,string"`
	Code              string `json:"code"`
	Name              string `json:"name"`
	OrgID             uint64 `json:"org_id,string"`
	OrgName           string `json:"org_name"`
	Venue             string `json:"venue"`
	ActivityStartTime string `json:"activity_start_time"`
	ActivityEndTime   string `json:"activity_end_time"`
	CheckinStartTime  string `json:"checkin_start_time"`
	CheckinEndTime    string `json:"checkin_end_time"`
	CheckinTemplate   string `json:"checkin_template"`
	Phase             string `json:"phase"`
	Status            uint8  `json:"status"`
}

type GetActivityDetailResp struct {
	ID                uint64 `json:"id,string"`
	Code              string `json:"code"`
	Name              string `json:"name"`
	OrgID             uint64 `json:"org_id,string"`
	Venue             string `json:"venue"`
	ActivityStartTime string `json:"activity_start_time"`
	ActivityEndTime   string `json:"activity_end_time"`
	CheckinStartTime  string `json:"checkin_start_time"`
	CheckinEndTime    string `json:"checkin_end_time"`
	CheckinTemplate   string `json:"checkin_template"`
	Status            uint8  `json:"status"`
	Remark            string `json:"remark"`
}

type ActivityReq struct {
	Name              string `json:"name" binding:"required,max=100"`
	OrgID             uint64 `json:"org_id,string" binding:"required"`
	Venue             string `json:"venue" binding:"required,max=255"`
	ActivityStartTime string `json:"activity_start_time" binding:"required,datetime=2006-01-02 15:04:05"`
	ActivityEndTime   string `json:"activity_end_time" binding:"required,datetime=2006-01-02 15:04:05"`
	CheckinStartTime  string `json:"checkin_start_time" binding:"required,datetime=2006-01-02 15:04:05"`
	CheckinEndTime    string `json:"checkin_end_time" binding:"required,datetime=2006-01-02 15:04:05"`
	CheckinTemplate   string `json:"checkin_template" binding:"required,oneof=invitation celebration luxe glass slot orbit grand-gala crimson-gala"`
	Status            uint8  `json:"status" binding:"required,oneof=1 2"`
	Remark            string `json:"remark" binding:"max=255"`
}

type CreateActivityReq struct {
	ActivityReq
}

type UpdateActivityReq struct {
	ActivityReq
}

type UpdateActivityCheckinTemplateReq struct {
	CheckinTemplate string `json:"checkin_template" binding:"required,oneof=invitation celebration luxe glass slot orbit grand-gala crimson-gala"`
}

//endregion

//region 抽奖轮次出入参

type GetRoundListReq struct {
	Page       int    `form:"page" binding:"omitempty,min=1"`
	PageSize   int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
	ActivityID uint64 `form:"activity_id"`
	Phase      string `form:"phase" binding:"omitempty,oneof=pending active finished"`
}

type GetRoundListResp struct {
	List  []RoundActivityItem `json:"list"`
	Total int64               `json:"total"`
}

type RoundActivityItem struct {
	ActivityID        uint64          `json:"activity_id,string"`
	ActivityName      string          `json:"activity_name"`
	ActivityStartTime string          `json:"activity_start_time"`
	ActivityEndTime   string          `json:"activity_end_time"`
	RoundCount        int             `json:"round_count"`
	Dates             []RoundDateItem `json:"dates"`
}

type RoundDateItem struct {
	DrawDate string          `json:"draw_date"`
	Rounds   []RoundListItem `json:"rounds"`
}

type RoundListItem struct {
	ID          uint64 `json:"id,string"`
	Sequence    int    `json:"sequence"`
	PrizeID     uint64 `json:"prize_id,string"`
	PrizeName   string `json:"prize_name"`
	Level       string `json:"level"`
	Name        string `json:"name"`
	DrawTime    string `json:"draw_time"`
	WinnerCount uint   `json:"winner_count"`
	Phase       string `json:"phase"`
}

type GetRoundDetailResp struct {
	ID          uint64 `json:"id,string"`
	ActivityID  uint64 `json:"activity_id,string"`
	PrizeID     uint64 `json:"prize_id,string"`
	Name        string `json:"name"`
	DrawDate    string `json:"draw_date"`
	DrawTime    string `json:"draw_time"`
	WinnerCount uint   `json:"winner_count"`
}

type RoundReq struct {
	ActivityID  uint64 `json:"activity_id,string" binding:"required"`
	PrizeID     uint64 `json:"prize_id,string" binding:"required"`
	Name        string `json:"name" binding:"required,max=100"`
	DrawDate    string `json:"draw_date" binding:"required,datetime=2006-01-02"`
	DrawTime    string `json:"draw_time" binding:"required,datetime=15:04:05"`
	WinnerCount uint   `json:"winner_count" binding:"required,min=1"`
}

type CreateRoundReq struct {
	RoundReq
}

type UpdateRoundReq struct {
	RoundReq
}

//endregion
