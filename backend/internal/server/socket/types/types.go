package types

type ConnectReq struct {
	Key string `form:"key" binding:"required,max=100"`
}
