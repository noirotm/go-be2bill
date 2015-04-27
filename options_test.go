package be2bill

import (
	"net/url"
	"reflect"
	"testing"
)

func TestOptionsSortedKeys(t *testing.T) {
	o := Options{
		"z": "test",
		"m": "vale",
		"a": "echo",
	}

	keys := o.sortedKeys()
	expected := []string{"a", "m", "z"}

	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("Got %v, expected %v", keys, expected)
	}
}

func TestOptionsUrlValuesSimple(t *testing.T) {
	o := Options{
		"z": "test",
		"m": "vale",
		"a": "echo",
	}

	params := o.urlValues()
	expected := url.Values{
		"z": {"test"},
		"m": {"vale"},
		"a": {"echo"},
	}

	if !reflect.DeepEqual(params, expected) {
		t.Errorf("Got %v, expected %v", params, expected)
	}
}

func TestOptionsUrlValuesRecursive(t *testing.T) {
	o := Options{
		"a": "echo",
		"p": Options{
			"z": "subopt1",
			"y": "subopt2",
		},
	}

	params := o.urlValues()
	expected := url.Values{
		"a":    {"echo"},
		"p[z]": {"subopt1"},
		"p[y]": {"subopt2"},
	}

	if !reflect.DeepEqual(params, expected) {
		t.Errorf("Got %v, expected %v", params, expected)
	}
}
