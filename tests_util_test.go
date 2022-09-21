package gowalker

import (
	"math"
	"os"
	"testing"
)

func TestConvertDataToString(t *testing.T) {
	var i int = 22
	if convertDataToString(i) != "22" {
		t.Error("could not convert int")
	}
	var i64 int64 = 22
	if convertDataToString(i64) != "22" {
		t.Error("could not convert int64")
	}
	var f32 float32 = 22.5
	if convertDataToString(f32) != "22.5" {
		t.Error("could not convert float32")
	}
	if convertDataToString(nil) != "null" {
		t.Error("could not convert nil value")
	}
}
func TestConvertStringToSameType(t *testing.T) {
	var i int = 22
	if res, _ := convertStringToSameType(i, "55"); res != 55 {
		t.Error("could not convert int")
	}
	var i64 int64 = 22
	if res, _ := convertStringToSameType(i64, "37"); res != 37 {
		t.Error("could not convert int64")
	}
	var f32 float32 = 22.5
	if res, _ := convertStringToSameType(f32, "50.7"); math.Floor(res.(float64)*100)/100 != 50.7 {
		t.Error("could not convert float32")
	}
	var f64 float64 = 22.5
	if res, _ := convertStringToSameType(f64, "50.7"); math.Floor(res.(float64)*100)/100 != 50.7 {
		t.Error("could not convert float32")
	}
	if res, _ := convertStringToSameType(nil, "bananas"); res != "bananas" {
		t.Error("could not convert nil value")
	}
}

func TestLoadTemplatesFromDisk(t *testing.T) {
	_ = os.Mkdir("test_data", os.ModePerm)
	_ = os.WriteFile("test_data/templ1.templ", []byte("foobar1"), os.ModePerm)
	_ = os.WriteFile("test_data/templ2.templ", []byte("foobar2"), os.ModePerm)
	_ = os.WriteFile("test_data/templ3.templ", []byte("foobar3"), os.ModePerm)
	defer func() {
		_ = os.RemoveAll("test_data")
	}()
	t1, subs, _ := LoadTemplatesFromDisk("test_data/templ1.templ")
	if t1 != "foobar1" {
		t.Error("could not load main template")
	}
	if subs["templ2"] != "foobar2" {
		t.Error("could not load sub template")
	}
	if subs["templ3"] != "foobar3" {
		t.Error("could not load sub template")
	}
}
