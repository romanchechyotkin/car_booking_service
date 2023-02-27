package user

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
	"time"
)

var createUserRequestMetrics = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "user_create_one",
	Subsystem:  "http",
	Name:       "request",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
}, []string{"status"})

func CreateUserObserveRequest(d time.Duration, status int) {
	createUserRequestMetrics.WithLabelValues(strconv.Itoa(status)).Observe(d.Seconds())
}
