package platform

import "testing"

func TestLookupFamilyStub(t *testing.T) {
	f, ok := lookupFamily("Macmini9,1")
	if ok || f != "" {
		t.Fatalf("expected stub lookup to return ok=false and an empty family")
	}
}
