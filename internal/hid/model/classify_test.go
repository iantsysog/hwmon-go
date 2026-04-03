package model

import (
	"testing"

	"github.com/iantsysog/hwmon-go/internal/test"
)

func TestClassifySensorsUsagePage(t *testing.T) {
	kind, unit, ok := classify(ElementMeta{UsagePage: 0x20, Usage: 0x33, Name: "Temperature"})
	test.True(t, ok)
	test.Eq(t, KindTemp, kind)
	test.Eq(t, "°C", unit)

	kind, unit, ok = classify(ElementMeta{UsagePage: 0x20, Usage: 0x26, Name: "Voltage"})
	test.True(t, ok)
	test.Eq(t, KindVolt, kind)
	test.Eq(t, "V", unit)
}

func TestNormalizeTdevSp78(t *testing.T) {
	v := normalizeValue(KindTemp, "PMU tdev1", 6400.0)
	test.Eq(t, 25.0, v)
	v = normalizeValue(KindTemp, "PMU tdev1", 80.0)
	test.Eq(t, 80.0, v)
}
