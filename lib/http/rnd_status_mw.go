package libhttp

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func RandomStatusMiddleware(min, max int, logger *logrus.Entry) func(handler http.Handler) http.Handler {
	rand.Seed(time.Now().UnixNano())

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			value := min + rand.Intn(max-min+1)
			logger.Infof("the dice has the number %d", value)

			if value%7 == 0 {
				logger.Warnf("%d %% 7 = 0 ;imitation of the internal server error, target endpoint does not triggered", value)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if value%9 == 0 {
				logger.Warn("%d %% 9 = 0 ; imitation of the bad request error, target endpoint does not triggered")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
