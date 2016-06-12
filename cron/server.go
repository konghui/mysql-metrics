package cron

import (
	"github.com/open-falcon/common/model"
)

func (my *MysqlIns) GlobalStatus() ([]*model.MetricValue, error) {
	return my.mysqlState("SHOW /*!50001 GLOBAL */ STATUS")
}

func (my *MysqlIns) GlobalVariables() ([]*model.MetricValue, error) {
	return my.mysqlState("SHOW /*!50001 GLOBAL */ VARIABLES")
}

func (my *MysqlIns) mysqlState(sql string) ([]*model.MetricValue, error) {
	rows, _, err := my.GetConnect().Query(sql)
	if err != nil {
		return nil, err
	}

	data := make([]*model.MetricValue, len(rows))
	i := 0
	for _, row := range rows {
		key_ := row.Str(0)
		v, err := row.Int64Err(1)
		// Ignore non digital value
		if err != nil {
			continue
		}

		metric := NewMetric(key_, my)
		if metric == nil {
			continue
		}
		data[i] = metric
		data[i].Value = v
		i++
	}
	return data[:i], nil
}
