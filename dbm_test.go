package main

import "testing"

func TestNewVersion(t *testing.T) {
	v, err := NewVersion("1.2.3")
	if err != nil {
		t.Errorf("Got error parsing version '1.2.3': %s", err)
	}
	if len(v) != 3 {
		t.Errorf("Bad version length: got %v instead of 3", len(v))
	}
	if v[0] != 1 || v[1] != 2 || v[2] != 3 {
		t.Error("Bad error parsing")
	}
	v, err = NewVersion("1")
	if err != nil {
		t.Errorf("Got error parsing version '1': %s", err)
	}
	if len(v) != 1 {
		t.Errorf("Bad version length: got %v instead of 1", len(v))
	}
}

func TestBadVersion(t *testing.T) {
	_, err := NewVersion("foo")
	if err == nil {
		t.Error("Got no error parsing version 'foo'")
	}
	if err.Error() != "Error parsing version 'foo'" {
		t.Error("Got bad error parsing version 'foo'")
	}
	_, err = NewVersion("")
	if err == nil {
		t.Error("Got no error parsing void version")
	}
}

func TestVersionCompareTo(t *testing.T) {
	v1, _ := NewVersion("0")
	v2, _ := NewVersion("0")
	if v1.CompareTo(v2) != 0 {
		t.Fail()
	}
	v2, _ = NewVersion("0.0")
	if v1.CompareTo(v2) >= 0 {
		t.Fail()
	}
	v2, _ = NewVersion("0.1")
	if v1.CompareTo(v2) >= 0 {
		t.Fail()
	}
	v1, _ = NewVersion("1.2.3")
	v2, _ = NewVersion("2.3.4")
	if v1.CompareTo(v2) >= 0 {
		t.Fail()
	}
	v2, _ = NewVersion("1.2.4")
	if v1.CompareTo(v2) >= 0 {
		t.Fail()
	}
	v2, _ = NewVersion("1.2.2")
	if v1.CompareTo(v2) <= 0 {
		t.Fail()
	}
	v2, _ = NewVersion("1.2.3")
	if v1.CompareTo(v2) != 0 {
		t.Fail()
	}
}
