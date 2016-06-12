package cron

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/open-falcon/common/model"
)

func (my *MysqlIns) innodbStatus() ([]*model.MetricValue, error) {
	status, _, err := my.GetConnect().QueryFirst("SHOW /*!50000 ENGINE */ INNODB STATUS")
	if err != nil {
		return nil, err
	}
	ctn := status.Str(2)
	rows := strings.Split(ctn, "\n")
	return parseInnodbStatus(my, rows)
}

func parseInnodbStatus(my *MysqlIns, rows []string) ([]*model.MetricValue, error) {
	var section string
	data := make([]*model.MetricValue, 0)
	for _, row := range rows {
		switch {
		case match("^BACKGROUND THREAD$", row):
			section = "BACKGROUND THREAD"
			continue
		case match("^DEAD LOCK ERRORS$", row), match("^LATEST DETECTED DEADLOCK$", row):
			section = "DEAD LOCK ERRORS"
			continue
		case match("^FOREIGN KEY CONSTRAINT ERRORS$", row), match("^LATEST FOREIGN KEY ERROR$", row):
			section = "FOREIGN KEY CONSTRAINT ERRORS"
			continue
		case match("^SEMAPHORES$", row):
			section = "SEMAPHORES"
			continue
		case match("^TRANSACTIONS$", row):
			section = "TRANSACTIONS"
			continue
		case match("^FILE I/O$", row):
			section = "FILE I/O"
			continue
		case match("^INSERT BUFFER AND ADAPTIVE HASH INDEX$", row):
			section = "INSERT BUFFER AND ADAPTIVE HASH INDEX"
			continue
		case match("^LOG$", row):
			section = "LOG"
			continue
		case match("^BUFFER POOL AND MEMORY$", row):
			section = "BUFFER POOL AND MEMORY"
			continue
		case match("^ROW OPERATIONS$", row):
			section = "ROW OPERATIONS"
			continue
		}

		if section == "SEMAPHORES" {
			matches := regexp.MustCompile(`^Mutex spin waits\s+(\d+),\s+rounds\s+(\d+),\s+OS waits\s+(\d+)`).FindStringSubmatch(row)
			if len(matches) == 4 {
				spin_waits, _ := strconv.Atoi(matches[1])
				Innodb_mutex_spin_waits := NewMetric("Innodb_mutex_spin_waits", my)
				if Innodb_mutex_spin_waits == nil {
					continue
				}
				Innodb_mutex_spin_waits.Value = spin_waits
				data = append(data, Innodb_mutex_spin_waits)

				spin_rounds, _ := strconv.Atoi(matches[2])
				Innodb_mutex_spin_rounds := NewMetric("Innodb_mutex_spin_rounds", my)
				if Innodb_mutex_spin_rounds == nil {
					continue
				}
				Innodb_mutex_spin_rounds.Value = spin_rounds
				data = append(data, Innodb_mutex_spin_rounds)

				os_waits, _ := strconv.Atoi(matches[3])
				Innodb_mutex_os_waits := NewMetric("Innodb_mutex_os_waits", my)
				if Innodb_mutex_os_waits == nil {
					continue
				}
				Innodb_mutex_os_waits.Value = os_waits
				data = append(data, Innodb_mutex_os_waits)
			}
		}
	}
	return data, nil
}

func match(pattern, s string) bool {
	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		return false
	}
	return matched
}
