package cron

import (
	"fmt"
	"strconv"

	"github.com/golang/glog"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

type MysqlIns struct {
	host     string
	port     int
	proto    string
	userName string
	passwd   string
	tag      string
	conn     mysql.Conn
	hostname string
	prefix   string
	interval int64
	metrics  map[string]int
}

func (m *MysqlIns) String() string {
	return fmt.Sprintf("IP=>%s\tPort=>%d\tProto=>%s\tUserName=>%s\tPasswd=>%s\ttag=%s\thostname=%s\tinterval=%d\tprefix=%s", m.host, m.port, m.proto, m.userName, m.passwd, m.tag, m.hostname, m.interval, m.prefix)
}

func (my *MysqlIns) Connect() (err error) {
	if err = my.GetConnect().Connect(); err != nil {
		return
	}
	return
}

// convent the Dsn string to mysqlIns
func ParseDsn(dsn string) (my *MysqlIns, err error) {

	if MysqlDsnParttern.MatchString(dsn) {
		var port int
		paramList := MysqlDsnParttern.FindStringSubmatch(dsn)
		port, err = strconv.Atoi(paramList[5])

		if err != nil {
			return
		}

		my = &MysqlIns{userName: paramList[1], passwd: paramList[2], proto: paramList[3], host: paramList[4], port: port, tag: fmt.Sprintf("port=%d", port)}

	}
	return
}

func (my *MysqlIns) HostName() (rv string, err error) {

	rows, _, err := my.GetConnect().Query("show /*!50001 GLOBAL */ variables like 'hostname';")
	if err != nil {
		return
	}
	for _, row := range rows {
		rv = row.Str(1)
		return
	}
	glog.Warningf("get host name faild use ip instead.")
	rv = my.host
	return
}

func (my *MysqlIns) GetConnect() (conn mysql.Conn) {
	if my.conn == nil {
		my.conn = mysql.New(my.proto, "", fmt.Sprintf("%s:%d", my.host, my.port), my.userName, my.passwd)
	}
	conn = my.conn
	return
}

func (my *MysqlIns) SetHostName(host string) {
	my.hostname = host
}

func (my *MysqlIns) GetHostName() (rv string) {
	rv = my.hostname
	return
}

func (my *MysqlIns) GetTag() (tag string) {
	tag = my.tag
	return
}

func (my *MysqlIns) SetTag(tag string) {
	my.tag = tag
}

func (my *MysqlIns) GetPrefix() (prefix string) {
	prefix = my.prefix
	return
}

func (my *MysqlIns) SetInterval(interval int64) {
	my.interval = interval
}
func (my *MysqlIns) GetInterval() (interval int64) {
	interval = my.interval
	return
}

func (my *MysqlIns) SetMetrics(metric map[string]int) {
	my.metrics = metric
}
func (my *MysqlIns) GetMetrics() (metric map[string]int) {
	metric = my.metrics
	return
}

func (my *MysqlIns) SetPrefix(prefix string) {
	my.prefix = prefix
}
