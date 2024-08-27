package metrics

import (
	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"time"
)

type LogPoint struct {
	Location string    `lp:"location"`
	Species  string    `lp:"species"`
	Count    int       `lp:"count"`
	Time     time.Time `lp:"time"`
}
type Reporter struct {
	influxDataBase string
	influxDBToken  string
	influxDBUrl    string
	client         *influxdb3.Client
}

var reporter *Reporter

func (*Reporter) RpcMethodMetric(methodName string, status string, time time.Time, cost time.Duration) {
	line := struct {
		Measurement string    `lp:"measurement"`
		Method      string    `lp:"tag,method"`
		Status      string    `lp:"tag,status"`
		Time        time.Time `lp:"timestamp"`
		Cost        int       `lp:"field,cost"`
	}{
		"rpc",
		methodName,
		status,
		time,
		int(cost.Milliseconds()),
	}
	data := []any{line}
	err := reporter.client.WriteData(context.Background(), data)
	if err != nil {
		panic(err)
	}

}

// INFLUXDB_TOKEN=KP7XM7Z9UP4jpegDNz6q6naHvQzCVcTZEBKFU1SmUE9nPTHtOyPFI4o0M7JitZhPAoYA-aHoSXLMS8XIthVA-w==
func Init() {
	//url := "https://us-east-1-1.aws.cloud2.influxdata.com"
	//token := "uA5ZLOV4HIofjgRMzYB8j0oxdNSSBhFeayWkJuTnHqHEJ3Ldda8IjJa7eyHXg8NOTGkmhj6t3EL7qdzQ2pdkwQ=="
	//client, err := influxdb3.New(influxdb3.ClientConfig{
	//	Host:     url,
	//	Token:    token,
	//	Database: "paymaster",
	//})
	//if err != nil {
	//	panic(err)
	//}
	//defer func(client *influxdb3.Client) {
	//	err = client.Close()
	//	if err != nil {
	//		panic(err)
	//	}
	//}(client)
	//s1 := struct {
	//	Measurement string    `lp:"measurement"`
	//	Sensor      string    `lp:"tag,location"`
	//	Temp        float64   `lp:"field,temperature"`
	//	Hum         int       `lp:"field,humidity"`
	//	Time        time.Time `lp:"timestamp"`
	//	Description string    `lp:"-"`
	//}{
	//	"stat",
	//	"Paris",
	//	23.5,
	//	55,
	//	time.Now(),
	//	"Paris weather conditions",
	//}
	//data := []any{s1}
	//err = client.WriteData(context.Background(), data)
	//if err != nil {
	//	panic(err)
	//}
}
