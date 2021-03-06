package opendmm

import (
	"testing"
)

func TestDmm(t *testing.T) {
	assertSearchable(
		t,
		[]string{
			"MIDE-029",
			"mide-029",
			"XV-100",
			"XV-1001",
			"IPZ687",
			"WPE-037",
			"ONEZ_067",
			"MMGH00010",
			"HODV021158",
			"RHE-463",
			"parathd-2000",
			"140c02202",
			"khy159",
			"DVDMS00392",
		},
		dmmEngine,
	)

	assertUnsearchable(
		t,
		[]string{
			"MCB3DBD-25",
			"S2MBD-048",
		},
		dmmEngine,
	)
}
