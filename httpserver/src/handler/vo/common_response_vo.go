package vo

type CommonResponseVO struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}
