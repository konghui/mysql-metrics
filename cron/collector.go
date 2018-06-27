package cron

import (
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/konghui/mysql-metrics/g"
	"github.com/open-falcon/common/model"
)

func MysqlAlive(my *MysqlIns, ok bool) {

	data := NewMetric("mysql_alive_local", my)
	if data == nil {
		return
	}
	if ok {
		data.Value = 1
	}
	g.SendMetrics([]*model.MetricValue{data})
}

func Append(data []*model.MetricValue, value []*model.MetricValue) (rv []*model.MetricValue) {
	if value == nil {
		rv = data
		return
	}
	rv = append(data, value...)
	return
}

func (my *MysqlIns) FetchData() (err error) {
	var hostname string
	defer func() {
		MysqlAlive(my, err == nil)
	}()

	if err = my.Connect(); err != nil {
		return
	}

	if hostname, err = my.HostName(); err != nil {
		return
	}

	// cfg tag endpoint overwrite hostname
	if strings.Contains(my.tag, "endpoint") {
		tags := strings.Split(my.tag, ",")
		tagm := map[string]string{}
		for _, t := range tags {
			ts := strings.Split(t, "=")
			tagm[ts[0]] = ts[1]
		}
		if len(tagm["endpoint"]) > 0 {
			hostname = tagm["endpoint"]
		}
	}

	my.SetHostName(hostname)

	defer my.GetConnect().Close()

	data := make([]*model.MetricValue, 0)
	globalStatus, err := my.GlobalStatus()

	if err != nil {
		return
	}
	data = Append(data, globalStatus)

	globalVars, err := my.GlobalVariables()
	if err != nil {
		return
	}
	data = Append(data, globalVars)

	innodbState, err := my.innodbStatus()
	if err != nil {
		return
	}
	data = Append(data, innodbState)

	slaveState, err := my.slaveStatus()
	if err != nil {
		return
	}
	data = Append(data, slaveState)

	g.SendMetrics(data)
	return
}

func Collect() {
	if !g.Config().Transfer.Enable {
		glog.Warningf("Open falcon transfer is not enabled!!!")
		return
	}

	if g.Config().Transfer.Addr == "" {
		glog.Warningf("Open falcon transfer addr is null!!!")
		return
	}
	db := g.Config().Daemon.Db
	if !g.Config().Daemon.Enable {
		glog.Warningf("Daemon collect not enabled in cfg.json!!!")

		if len(db) < 1 {
			glog.Warningf("Not set addrs of daemon in cfg.json!!!")
		}
		return
	}

	go collect(db)
}

func collect(db []string) {
	var interval int64 = g.Config().Transfer.Interval
	var tout = g.Config().Daemon.Timeout
	timeout := time.Duration(tout) * time.Second
	timer := time.NewTicker(time.Duration(interval) * time.Second)
	metrics := g.Config().Metrics
	prefix := g.Config().Daemon.Prefix
	glog.Infof("MySQL Monitor for falcon")

	for {
		<-timer.C
		for _, v := range db {
			myIns, err := ParseDsn(v)
			if err != nil {
				glog.Warningf(err.Error())
			}

			myIns.SetInterval(interval)
			myIns.SetMetrics(metrics)
			myIns.SetPrefix(prefix)
			myIns.GetConnect().SetTimeout(timeout)

			err = myIns.FetchData()
			if err != nil {
				glog.Warningf(err.Error())
			}
		}
	}
}
