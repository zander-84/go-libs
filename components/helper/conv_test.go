package helper

import (
	"testing"
)

func TestConv_ShouldStoU(t *testing.T) {
	uInt := defaultConv.ShouldStoU("123123123")

	if uInt != 123123123 {
		t.Errorf("not eq 123123123, uint is %d", uInt)
	}

	uInt32 := defaultConv.ShouldStoU32("4294967295")

	if uInt32 != 4294967295 {
		t.Errorf("not eq 4294967295, uInt32 is %d", uInt32)
	}

	uInt64 := defaultConv.ShouldStoU64("4294967296")

	if uInt64 != 4294967296 {
		t.Errorf("not eq 4294967296, uInt64 is %d", uInt64)
	}
}
