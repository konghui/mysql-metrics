package cron

import (
	"fmt"
	"time"

	"github.com/open-falcon/common/model"
)

const (
	GAUGE = iota
	COUNTER
	NOEXIST
)

var Type = [...]string{
	"GAUGE",
	"COUNTER",
	"NOEXIST",
}

func dataType(key_ int) (rv string) {
	return Type[key_]
}

func NewMetric(name string, my *MysqlIns) (value *model.MetricValue) {
	metric := my.GetMetrics()
	if v, ok := metric[name]; ok {
		value = &model.MetricValue{
			Metric:    fmt.Sprintf("%s%s", my.GetPrefix(), name),
			Endpoint:  my.GetHostName(),
			Type:      dataType(v),
			Tags:      my.GetTag(),
			Timestamp: time.Now().Unix(),
			Step:      my.GetInterval(),
		}
		return
	}
	return

}
