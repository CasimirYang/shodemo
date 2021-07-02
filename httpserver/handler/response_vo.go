package handler

type ResponseVO struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}
