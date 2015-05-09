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

func TestOptionsFlatten(t *testing.T) {
	testCases := []struct{ options, expected Options }{
		{
			Options{
				"z": "test",
				"m": "vale",
				"a": "echo",
			},
			Options{
				"z": "test",
				"m": "vale",
				"a": "echo",
			},
		},
		{
			Options{
				"a": "echo",
				"p": Options{
					"z": "subopt1",
					"y": "subopt2",
				},
			},
			Options{
				"a":    "echo",
				"p[z]": "subopt1",
				"p[y]": "subopt2",
			},
		},
		{
			Options{
				"a": "echo",
				"p": Options{
					"z": "subopt1",
					"x": Options{
						"2015-09-12": "3000",
						"2015-10-12": "1000",
						"2015-11-12": "1500",
					},
					"y": "subopt2",
				},
			},
			Options{
				"a":                "echo",
				"p[z]":             "subopt1",
				"p[x][2015-09-12]": "3000",
				"p[x][2015-10-12]": "1000",
				"p[x][2015-11-12]": "1500",
				"p[y]":             "subopt2",
			},
		},
	}

	for _, test := range testCases {
		result := test.options.flatten()

		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("  want %+v\n", test.expected)
			t.Errorf("  got %+v", result)
		}
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

	o = Options{
		"a": "echo",
		"p": Options{
			"z": "subopt1",
			"x": Options{
				"2015-09-12": "3000",
				"2015-10-12": "1000",
				"2015-11-12": "1500",
			},
			"y": "subopt2",
		},
	}

	params = o.urlValues()
	expected = url.Values{
		"a":                {"echo"},
		"p[z]":             {"subopt1"},
		"p[x][2015-09-12]": {"3000"},
		"p[x][2015-10-12]": {"1000"},
		"p[x][2015-11-12]": {"1500"},
		"p[y]":             {"subopt2"},
	}

	if !reflect.DeepEqual(params, expected) {
		t.Errorf("Got %v, expected %v", params, expected)
	}
}
