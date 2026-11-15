package main

import (
	"fmt"
	"github.com/apache/fory/go/fory"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

// UpReportPropertiesDeviceMessage 上报属性设备消息
type UpReportPropertiesDeviceMessage struct {
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
	// 设备SN
	DeviceSn string `json:"deviceSn,omitempty"`
	// 设备名称
	DeviceName string `json:"deviceName,omitempty"`
	// 设备属性
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// UpReportPropertiesGatewayMessage 上报属性网关消息
type UpReportPropertiesGatewayMessage struct {
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

	// 设备属性
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// DownCommandGatewayMessage 下行指令网关消息
type DownCommandGatewayMessage struct {
	// TenantCode 租户编码
	TenantCode string `json:"tenantCode,omitempty"`

	// GatewaySn 网关SN
	GatewaySn string `json:"gatewaySn,omitempty"`

	// Timestamp 时间戳（Java Long，毫秒）
	Timestamp int64 `json:"timestamp,omitempty"`

	// RequestId 请求ID
	RequestId int64 `json:"requestId,omitempty"`

	// Metadata 扩展元数据
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// DeviceSn 设备SN
	DeviceSn string `json:"deviceSn,omitempty"`

	// DeviceName 设备名称
	DeviceName string `json:"deviceName,omitempty"`

	// Command 指令
	Command string `json:"command,omitempty"`

	// Param 参数值
	Param map[string]interface{} `json:"param,omitempty"`
}

// UpCommandReplyGatewayMessage 上行指令回复网关消息
type UpCommandReplyGatewayMessage struct {
	// TenantCode 租户编码
	TenantCode string `json:"tenantCode,omitempty"`

	// GatewaySn 网关SN
	GatewaySn string `json:"gatewaySn,omitempty"`

	// Timestamp 时间戳（Java Long，毫秒）
	Timestamp int64 `json:"timestamp,omitempty"`

	// RequestId 请求ID
	RequestId int64 `json:"requestId,omitempty"`

	// Metadata 扩展元数据
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// DeviceSn 设备SN
	DeviceSn string `json:"deviceSn,omitempty"`

	// DeviceName 设备名称
	DeviceName string `json:"deviceName,omitempty"`

	// Command 指令
	Command string `json:"command,omitempty"`

	// Payload 内容
	Payload map[string]interface{} `json:"payload,omitempty"`

	// Message 信息
	Message string `json:"message,omitempty"`
}

func main() {
	f := fory.New(fory.WithXlang(true))
	if err := f.RegisterStruct(UpReportPropertiesGatewayMessage{}, 244); err != nil {
		fmt.Println(err)
	}
	if err := f.RegisterStruct(UpReportPropertiesDeviceMessage{}, 245); err != nil {
		fmt.Println(err)
	}
	if err := f.RegisterStruct(DownCommandGatewayMessage{}, 246); err != nil {
		fmt.Println(err)
	}
	if err := f.RegisterStruct(UpCommandReplyGatewayMessage{}, 247); err != nil {
		fmt.Println(err)
	}
	options := mqtt.NewClientOptions().AddBroker("tcp://0.0.0.0:41883").SetClientID("test-001").SetUsername("root").SetPassword("laokou123").SetConnectTimeout(5 * time.Second)
	client := mqtt.NewClient(options)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
    subToken := client.Subscribe("down/laokouyun/KEG10000/command", 1, func(c mqtt.Client, m mqtt.Message) {
        var msg DownCommandGatewayMessage
        if err := f.Unmarshal(m.Payload(), &msg); err != nil {
            fmt.Println("反序列化失败:", err)
            return
        }
        fmt.Printf("收到下行指令: sn=%s, command=%s, requestId=%d, param=%v\n", msg.DeviceSn, msg.Command, msg.RequestId, msg.Param)
        	message := UpCommandReplyGatewayMessage{
        			TenantCode: "laokouyun",
        			GatewaySn:  "KEG10000",
        			Timestamp:  time.Now().UnixMilli(),
        			RequestId:  msg.RequestId,
        			Metadata:   map[string]interface{}{},
        			Command: msg.Command,
        			Payload: map[string]interface{}{},
        			Message: "OK",
        	}
        	b3, err := f.Serialize(&message)
        	if err != nil {
        		fmt.Println(err)
        	}
        	client.Publish("up/laokouyun/KEG10000/command/reply", 1, false, b3)
    })
    if subToken.Wait() && subToken.Error() != nil {
        fmt.Println("订阅失败:", subToken.Error())
        return
    }
	defer client.Disconnect(250)
	sendDeviceMessage(f, client)
	time.Sleep(1 * time.Second)
	sendGatewayMessage(f, client)
	time.Sleep(100 * time.Second)
}

func sendDeviceMessage(f *fory.Fory, client mqtt.Client) {
	// 设备属性
	topic1 := "up/laokouyun/KEG10000/CNC100/properties/report"
	msg := UpReportPropertiesDeviceMessage{
    TenantCode: "laokouyun",
    				GatewaySn:  "KEG10000",
    				Timestamp:  time.Now().UnixMilli(),
    				RequestId:  100,
    				Metadata:   map[string]interface{}{},
    							DeviceSn:   "CNC100",
                    			DeviceName: "数控机床",
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
}

func sendGatewayMessage(f *fory.Fory, client mqtt.Client) {
	// 网关属性
	topic2 := "up/laokouyun/KEG10000/properties/report"
	message := UpReportPropertiesGatewayMessage{
			TenantCode: "laokouyun",
			GatewaySn:  "KEG10000",
			Timestamp:  time.Now().UnixMilli(),
			RequestId:  101,
			Metadata:   map[string]interface{}{},
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
