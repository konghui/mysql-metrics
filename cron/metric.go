package cron

import (
	"fmt"
	"regexp"
	"time"

	"github.com/open-falcon/common/model"
)

const (
	MYSQL_DSN_PARTTERN = "(?P<user>[a-zA-Z]+):(?P<passwd>.*)@(?P<proto>[a-zA-Z]+)\\((?P<ip>[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+):(?P<port>[0-9]+)\\)"
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

var MysqlDsnParttern = regexp.MustCompile(MYSQL_DSN_PARTTERN)

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
