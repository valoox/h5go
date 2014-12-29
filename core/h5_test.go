package core

import "testing"

// Tries to initialise the HDF5 lib
func TestInit(t *testing.T) {
	if err := Init(); err != nil {
		t.Fatalf("Error initialising the lib: %s", err)
	}
}

// Tests the version
func TestVersion(t *testing.T) {
	M, m, r, err := GetVersion()
	if err != nil {
		t.Fatalf("Error getting lib version: %s", err)
	}
	t.Logf("Version: %v.%v.%v", M, m, r)
}

// Tests the dynamic setting of limits
func TestLimits(t *testing.T) {
	t.Logf("%s", Limits)
	Limits.Array(1*MB, 8*MB)
	if err := Limits.Set(); err != nil {
		t.Fatalf("Error setting limits: %s", err)
	}
	if Limits.Current()[1][1] != 8*MB {
		t.Fatalf("Limit not set properly")
	}
	t.Logf("Updated %s", Limits)
}
