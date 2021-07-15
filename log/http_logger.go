package log

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"xiaomi_cloud/utils"
)

// GetHttpTransport 在默认http client基础上增加日志功能
func GetHttpTransport() http.RoundTripper {
	return &loggedRoundTripper{http.DefaultTransport}
}

type loggedRoundTripper struct {
	rt http.RoundTripper
}

func (c *loggedRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	startTime := time.Now()
	request.Header.Add("Account-Agent", "Mozilla/5.0 (X11; Linux x86_64) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.79 Safari/537.36")
	response, err := c.rt.RoundTrip(request)
	duration := time.Since(startTime)
	duration /= time.Millisecond
	if err != nil {
		Error(fmt.Sprintf("req_error method = %s duration = %d url = %s error = %s time = %s",
			request.Method, duration, request.URL.String(), err.Error(), time.Now().Format(utils.YmdHis)))
	} else {
		Info(fmt.Sprintf("req_success method = %s status = %d duration = %d url = %s time = %s",
			request.Method, response.StatusCode, duration, request.URL.String(), time.Now().Format(utils.YmdHis)))
	}
	fmt.Printf("%s\n", strings.Repeat("-", 50))
	return response, err
}
