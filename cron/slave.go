package cron

import (
	"github.com/open-falcon/common/model"
)

var SlaveStatusToSend = []string{
	"Exec_Master_Log_Pos",
	"Read_Master_Log_Pos",
	"Relay_Log_Pos",
	"Seconds_Behind_Master",
	"Slave_IO_Running",
	"Slave_SQL_Running",
}

func (my *MysqlIns) slaveStatus() ([]*model.MetricValue, error) {

	isSlave := NewMetric("Is_slave", my)
	if isSlave == nil {
		return nil, nil
	}

	row, res, err := my.GetConnect().QueryFirst("SHOW SLAVE STATUS")
	if err != nil {
		return nil, err
	}

	// be master
	if row == nil {
		isSlave.Value = 0
		return []*model.MetricValue{isSlave}, nil
	}

	// be slave
	isSlave.Value = 1

	data := make([]*model.MetricValue, len(SlaveStatusToSend))
	for i, s := range SlaveStatusToSend {
		metric := NewMetric(s, my)
		if metric == nil {
			continue
		}
		data[i] = metric
		switch s {
		case "Slave_SQL_Running", "Slave_IO_Running":
			data[i].Value = 0
			v := row.Str(res.Map(s))
			if v == "Yes" {
				data[i].Value = 1
			}
		default:
			v, err := row.Int64Err(res.Map(s))
			if err != nil {
				data[i].Value = -1
			} else {
				data[i].Value = v
			}
		}
	}
	return append(data, isSlave), nil
}
