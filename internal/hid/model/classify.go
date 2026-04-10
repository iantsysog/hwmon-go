package model

import "strings"

type ElementMeta struct {
	Name      string
	UsagePage uint32
	Usage     uint32
}

type usageKey uint64

func makeUsageKey(usagePage, usage uint32) usageKey {
	return usageKey(uint64(usagePage)<<32 | uint64(usage))
}

type classifyRule struct {
	kind Kind
	unit string
}

var sensorRules = map[usageKey]classifyRule{
	makeUsageKey(0x20, 0x33): {kind: KindTemp, unit: "°C"},
	makeUsageKey(0x20, 0x15): {kind: KindTemp, unit: "°C"},
	makeUsageKey(0x20, 0x26): {kind: KindVolt, unit: "V"},
	makeUsageKey(0x20, 0x22): {kind: KindAmp, unit: "A"},
	makeUsageKey(0x20, 0x23): {kind: KindWatt, unit: "W"},
}

func classify(meta ElementMeta) (Kind, string, bool) {
	if r, ok := sensorRules[makeUsageKey(meta.UsagePage, meta.Usage)]; ok {
		return r.kind, r.unit, true
	}

	lower := strings.ToLower(meta.Name)
	switch {
	case strings.Contains(lower, "temp") || strings.Contains(lower, "temperature"):
		return KindTemp, "°C", true
	case strings.Contains(lower, "volt") || strings.Contains(lower, "voltage"):
		return KindVolt, "V", true
	case strings.Contains(lower, "amp") || strings.Contains(lower, "current"):
		return KindAmp, "A", true
	case strings.Contains(lower, "watt") || strings.Contains(lower, "power"):
		return KindWatt, "W", true
	case strings.Contains(lower, "fan") && (strings.Contains(lower, "rpm") || strings.Contains(lower, "speed")):
		return KindFanRPM, "RPM", true
	case strings.Contains(lower, "percent") || strings.Contains(lower, "%"):
		return KindPercent, "%", true
	}

	return KindOther, "", false
}

func Classify(meta ElementMeta) (Kind, string, bool) { return classify(meta) }

func normalizeValue(kind Kind, name string, val float64) float64 {
	if kind != KindTemp {
		return val
	}

	lower := strings.ToLower(name)
	idx := strings.Index(lower, "tdev")
	if idx < 0 || idx+4 >= len(lower) {
		return val
	}
	next := lower[idx+4]
	if next >= '1' && next <= '9' && val > 130.0 {
		return val / 256.0
	}
	return val
}

func NormalizeValue(kind Kind, name string, val float64) float64 {
	return normalizeValue(kind, name, val)
}
