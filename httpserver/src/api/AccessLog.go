package api

import (
	"bytes"
	commonLog "git.garena.com/jinghua.yang/entry-task-common/log"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w CustomResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func AccessLogHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.RequestURI[0:7] == "/static" || c.Request.RequestURI == "/uc/uploadProfile" {
			c.Next()
			return
		}
		var requestBodyBytes []byte
		if c.Request.Body != nil {
			requestBodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBodyBytes))

		res := &CustomResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = res
		c.Next()

		commonLog.SugarLogger.Infof("url=%s, status=%d,request=%s, resp=%s", c.Request.URL, c.Writer.Status(), string(requestBodyBytes), res.body.String())
	}
}
