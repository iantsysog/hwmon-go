package model

type Kind string

const (
	KindTemp    Kind = "temp"
	KindVolt    Kind = "volt"
	KindAmp     Kind = "amp"
	KindWatt    Kind = "watt"
	KindFanRPM  Kind = "fan_rpm"
	KindPercent Kind = "percent"
	KindCount   Kind = "count"
	KindOther   Kind = "other"
)
