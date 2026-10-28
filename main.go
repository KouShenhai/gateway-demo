package main

import (
	"bufio"
	"fmt"
	"github.com/apache/fory/go/fory"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"time"
)

// GatewayMessage 网关消息基类
type GatewayMessage struct {
	// 租户编码
	TenantCode string `json:"tenantCode,omitempty"`
	// 网关SN
	GatewaySn string `json:"gatewaySn,omitempty"`
	// 时间戳
	Timestamp int64 `json:"timestamp,omitempty"`
	// 请求ID
	RequestId int64 `json:"requestId,omitempty"`
	// 扩展元数据
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// DeviceMessage 设备消息（嵌入 GatewayMessage）
type DeviceMessage struct {
	GatewayMessage

	// 设备SN
	DeviceSn string `json:"deviceSn,omitempty"`
	// 设备名称
	DeviceName string `json:"deviceName,omitempty"`
}

// UpReportPropertiesDeviceMessage 上报属性设备消息
type UpReportPropertiesDeviceMessage struct {
	DeviceMessage

	// 设备属性
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// UpReportPropertiesGatewayMessage 上报属性网关消息
type UpReportPropertiesGatewayMessage struct {
	GatewayMessage

	// 设备属性
	Properties map[string]interface{} `json:"properties,omitempty"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	f := fory.New(fory.WithXlang(true))
	if err := f.RegisterStruct(GatewayMessage{}, 1); err != nil {
		fmt.Println(err)
	}
	if err := f.RegisterStruct(DeviceMessage{}, 2); err != nil {
		fmt.Println(err)
	}
	if err := f.RegisterStruct(UpReportPropertiesDeviceMessage{}, 3); err != nil {
		fmt.Println(err)
	}
	if err := f.RegisterStruct(UpReportPropertiesGatewayMessage{}, 4); err != nil {
		fmt.Println(err)
	}
	options := mqtt.NewClientOptions().AddBroker("tcp://0.0.0.0:41883").SetClientID("test-001").SetUsername("root").SetPassword("laokou123").SetConnectTimeout(5 * time.Second)
	client := mqtt.NewClient(options)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
	defer client.Disconnect(250)
	if input == "1" {
		// 设备属性
		topic1 := "up/laokouyun/KEG100/CNC100/properties/report"
		msg := UpReportPropertiesDeviceMessage{
			DeviceMessage: DeviceMessage{
				GatewayMessage: GatewayMessage{
					TenantCode: "laokouyun",
					GatewaySn:  "KEG100",
					Timestamp:  time.Now().UnixMilli(),
					RequestId:  100,
					Metadata:   map[string]interface{}{},
				},
				DeviceSn:   "CNC100",
				DeviceName: "数控机床",
			},
			Properties: map[string]interface{}{
				"programNumber":            "CNC100",
				"programName":              "SHAFT_ROUGHING_V2.NC",
				"readTime":                 "2026-10-28 14:06:47",
				"spindleRatio":             85,
				"spindleLoad":              623,
				"spindleName":              "SPINDLE-100",
				"spindleSpeed":             4380,
				"spindleActualSpeed":       5010,
				"workTime":                 18640,
				"bootTime":                 "2026-10-28 07:30:12",
				"cuttingTime":              9420,
				"feedRate":                 78,
				"feedRatio":                100,
				"singleWorkPieces":         4,
				"workPieces":               1357,
				"workMode":                 "AUTO",
				"runningState":             "running",
				"emergencyStopState":       true,
				"alarmState":               true,
				"programRegisteredNumber":  128,
				"programRegistrableNumber": 372,
				"memoryUsed":               524288,
				"memoryAvailable":          786432,
				"spindleSetSpeed":          4500,
				"spindleTemperature":       46,
				"machinePositionY":         1520,
				"spindleCurrent":           1850,
				"feedSetRate":              800,
				"workpieceNumber":          1357,
				"toolNumber":               7,
				"blockNumber":              1204,
				"servoLoadX":               32,
				"servoLoadY":               28,
				"servoLoadZ":               41,
				"surplusPositionX":         15,
				"surplusPositionY":         8,
				"surplusPositionZ":         3,
				"workpiecePositionX":       2350,
				"workpiecePositionY":       1200,
				"workpiecePositionZ":       450,
				"machinePositionX":         2500,
				"machinePositionZ":         480,
				"currentPositionX":         2350,
				"currentPositionY":         1200,
				"currentPositionZ":         450,
			},
		}
		b1, err := f.Serialize(&msg)
		if err != nil {
			fmt.Println(err)
		}
		client.Publish(topic1, 0, false, b1)
	} else if input == "2" {
		// 网关属性
		topic2 := "up/laokouyun/KEG100/properties/report"
		message := UpReportPropertiesGatewayMessage{
			GatewayMessage: GatewayMessage{
				TenantCode: "laokouyun",
				GatewaySn:  "KEG100",
				Timestamp:  time.Now().UnixMilli(),
				RequestId:  101,
				Metadata:   map[string]interface{}{},
			},
			Properties: map[string]interface{}{
				"diskTotal": "128GB",
				"diskUsed":  "76.8GB",
				"diskUsage": "60%",
				"cpuUsage":  "35.6%",
				"memTotal":  "16GB",
				"memUsed":   "9.2GB",
				"memUsage":  "57.5%",
				"uptime":    "15天8小时32分钟",
				"version":   "v1.0.0",
				"interval":  30,
			},
		}
		b2, err := f.Serialize(&message)
		if err != nil {
			fmt.Println(err)
		}
		client.Publish(topic2, 0, false, b2)
	}
}
