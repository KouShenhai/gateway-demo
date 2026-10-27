package main

import (
	"fmt"
	"github.com/apache/fory/go/fory"
	mqtt "github.com/eclipse/paho.mqtt.golang"
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

func main() {
	f := fory.New(fory.WithXlang(true))
	if err := f.RegisterStruct(UpReportPropertiesDeviceMessage{}, 1); err != nil {
		fmt.Println(err)
	}
	options := mqtt.NewClientOptions().AddBroker("tcp://127.0.0.1:41883").SetClientID("test-001").SetUsername("root").SetPassword("laokou123").SetConnectTimeout(5 * time.Second)
	client := mqtt.NewClient(options)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
	defer client.Disconnect(250)
	topic := "up/laokouyun/KEG100/DEV100/properties/report"
	msg := UpReportPropertiesDeviceMessage{
		DeviceMessage: DeviceMessage{
			GatewayMessage: GatewayMessage{
				TenantCode: "laokouyun",
				GatewaySn:  "KEG100",
				Timestamp:  time.Now().UnixMilli(),
				RequestId:  1001,
				Metadata:   map[string]interface{}{},
			},
			DeviceSn:   "DEV100",
			DeviceName: "温度传感器",
		},
		Properties: map[string]interface{}{
			"temperature": 25.6,
			"humidity":    60,
		},
	}
	bytes, err := f.Serialize(msg)
	if err != nil {
		fmt.Println(err)
	}
	client.Publish(topic, 0, false, bytes)
}
