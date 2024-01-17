package consumer

type Message interface {
	isMessage()
}

type TelemetryMessage struct {
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Name      string                 `json:"name,omitempty"`
	Tags      MessageTags            `json:"tags,omitempty"`
	Timestamp int64                  `json:"timestamp,omitempty"`
}

type MessageTags struct {
	Host          string `json:"host,omitempty"`
	InterfaceName string `json:"interface_name,omitempty"`
	Node          string `json:"node"`
	Path          string `json:"path,omitempty"`
	Source        string `json:"source,omitempty"`
	Subscription  string `json:"subscription,omitempty"`
}

type DelayMessage struct {
	TelemetryMessage
	Average  float64 `json:"delay_measurement_session/last_advertisement_information/advertised_values/average,omitempty"`
	Maximum  float64 `json:"delay_measurement_session/last_advertisement_information/advertised_values/maximum,omitempty"`
	Minimum  float64 `json:"delay_measurement_session/last_advertisement_information/advertised_values/minimum,omitempty"`
	Variance float64 `json:"delay_measurement_session/last_advertisement_information/advertised_values/variance,omitempty"`
}

type LossMessage struct {
	TelemetryMessage
	LossPercentage float64 `json:"interface_status_and_data/enabled/packet_loss_percentage,omitempty"`
}

type BandwidthMessage struct {
	TelemetryMessage
	Bandwidth float64 `json:"interface_status_and_data/enabled/bandwidth,omitempty"`
}

func (TelemetryMessage) isMessage() {}
