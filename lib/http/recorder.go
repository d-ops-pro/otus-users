package libhttp

import "net/http"

type responseRecorder struct {
	status int
	writer http.ResponseWriter
}

func NewResponseRecorder(w http.ResponseWriter) *responseRecorder {
	return &responseRecorder{writer: w}
}

func (r responseRecorder) Status() int {
	// status hasn't been set, so it should be OK
	if r.status == 0 {
		return http.StatusOK
	}

	return r.status
}

func (r responseRecorder) Header() http.Header {
	return r.writer.Header()
}

func (r responseRecorder) Write(b []byte) (int, error) {
	return r.writer.Write(b)
}

func (r responseRecorder) WriteHeader(s int) {
	r.status = s
	r.writer.WriteHeader(s)
}
