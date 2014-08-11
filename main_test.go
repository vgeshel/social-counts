package main

import (
	"fmt"
	"encoding/json"
	"testing"
)

func ExampleJson() {
	var jsonBlob = []byte(`
		{"Name": "Platypus", "Order": "Monotremata"}
	`)
	var x map[string]interface{}

	err := json.Unmarshal(jsonBlob, &x)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%s", x["Name"])
	// Output: Platypus
}


func TestGPlus(t *testing.T) {
	cnt, err := GetGPlus("http://no.such.domain.org/foo/bar?baz")

	if err != nil {
		t.Errorf("G+ request failed: %+v", err)
		return
	}

	if cnt != 0 {
		t.Errorf("G+ returned unexpected value: %d", cnt)
	}

	goodUrl := "https://developers.google.com/speed/pagespeed/service/faq"
	cnt, err = GetGPlus(goodUrl)

	if err != nil {
		t.Errorf("G+ request failed: %+v", err)
		return
	}

	if cnt <= 0 {
		t.Errorf("G+ returned unexpected value: %d", cnt)
	}

	t.Logf("G+ returned %d for %s", cnt, goodUrl)
}

func TestStumbleUpon(t *testing.T) {
	cnt, err := GetStumbleUpon("http://no.such.domain.org/foo/bar?baz")

	if err != nil {
		t.Errorf("SU request failed: %+v", err)
		return
	}

	if cnt != 0 {
		t.Errorf("SU returned unexpected value: %d", cnt)
	}

	goodUrl := "http://www.stopfake.org"
	cnt, err = GetStumbleUpon(goodUrl)

	if err != nil {
		t.Errorf("SU request failed: %+v", err)
		return
	}

	if cnt <= 0 {
		t.Errorf("SU returned unexpected value: %d", cnt)
	}

	t.Logf("SU returned %d for %s", cnt, goodUrl)
}

func TestPinterest(t *testing.T) {
	cnt, err := GetPinterest("http://no.such.domain.org/foo/bar?baz")

	if err != nil {
		t.Errorf("Pinterest request failed: %+v", err)
		return
	}

	if cnt != 0 {
		t.Errorf("Pinterest returned unexpected value: %d", cnt)
	}

	goodUrl := "http://www.stopfake.org"
	cnt, err = GetStumbleUpon(goodUrl)

	if err != nil {
		t.Errorf("Pinterest request failed: %+v", err)
		return
	}

	if cnt <= 0 {
		t.Errorf("Pinterest returned unexpected value: %d", cnt)
	}

	t.Logf("Pinterest returned %d for %s", cnt, goodUrl)
}

func TestVK(t *testing.T) {
	cnt, err := GetVK("http://no.such.domain.org/foo/bar?baz")

	if err != nil {
		t.Errorf("VK request failed: %+v", err)
		return
	}

	if cnt != 0 {
		t.Errorf("VK returned unexpected value: %d", cnt)
	}

	goodUrl := "http://www.stopfake.org"
	cnt, err = GetVK(goodUrl)

	if err != nil {
		t.Errorf("VK request failed: %+v", err)
		return
	}

	if cnt <= 650 {
		t.Errorf("VK returned unexpected value: %d", cnt)
	}

	t.Logf("VK returned %d for %s", cnt, goodUrl)
}
