package main

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	// test
	status := m.Run()

	// exit
	os.Exit(status)
}

func Test(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		baseStr         string
		overlayStr      string
		expectedStr     string
		baseType        interface{}
		overlayType     interface{}
		expectedType    interface{}
		expectedIsError bool
	}{
		{
			name:         `[normal]異なるkeyはmergeされる`,
			baseStr:      `{"a":"ok"}`,
			overlayStr:   `{"b":"ok"}`,
			expectedStr:  `{"a":"ok","b":"ok"}`,
			baseType:     map[string]string{},
			overlayType:  map[string]string{},
			expectedType: map[string]string{},
		},
		{
			name:         `[normal]同じkeyはbaseが優先される`,
			baseStr:      `{"a":"ok"}`,
			overlayStr:   `{"a":"ng"}`,
			expectedStr:  `{"a":"ok"}`,
			baseType:     map[string]string{},
			overlayType:  map[string]string{},
			expectedType: map[string]string{},
		},
		{
			name:         `[normal]overlayが空の場合はbaseがそのまま返る`,
			baseStr:      `{"a":"ok"}`,
			overlayStr:   `{}`,
			expectedStr:  `{"a":"ok"}`,
			baseType:     map[string]string{},
			overlayType:  map[string]string{},
			expectedType: map[string]string{},
		},
		{
			name:         `[normal]baseが空の場合はoverlayがそのまま返る`,
			baseStr:      `{}`,
			overlayStr:   `{"a":"ok"}`,
			expectedStr:  `{"a":"ok"}`,
			baseType:     map[string]string{},
			overlayType:  map[string]string{},
			expectedType: map[string]string{},
		},
		{
			name:         `[normal]map内の異なるkeyはmergeされる`,
			baseStr:      `{"a":{"a1":"ok"}}`,
			overlayStr:   `{"b":{"b1":"ok"}}`,
			expectedStr:  `{"a":{"a1":"ok"},"b":{"b1":"ok"}}`,
			baseType:     map[string]interface{}{},
			overlayType:  map[string]interface{}{},
			expectedType: map[string]interface{}{},
		},
		{
			name:         `[normal]map内の同じkeyはbaseが優先される`,
			baseStr:      `{"a":{"a1":"ok"}}`,
			overlayStr:   `{"a":{"a1":"ng"}}`,
			expectedStr:  `{"a":{"a1":"ok"}}`,
			baseType:     map[string]interface{}{},
			overlayType:  map[string]interface{}{},
			expectedType: map[string]interface{}{},
		},
		{
			name:         `[normal]map内の異なるkeyはmergeされる`,
			baseStr:      `{"a":{"a1":"ok"}}`,
			overlayStr:   `{"a":{"a2":"ok"}}`,
			expectedStr:  `{"a":{"a1":"ok","a2":"ok"}}`,
			baseType:     map[string]interface{}{},
			overlayType:  map[string]interface{}{},
			expectedType: map[string]interface{}{},
		},
		{
			name:         `[normal]array内の異なるvalueはmergeされる`,
			baseStr:      `{"a":["a1"]}`,
			overlayStr:   `{"a":["a2"]}`,
			expectedStr:  `{"a":["a1","a2"]}`,
			baseType:     map[string]interface{}{},
			overlayType:  map[string]interface{}{},
			expectedType: map[string]interface{}{},
		},
		{
			name:         `[normal]array内の同じvalueはmergeされ重複する`,
			baseStr:      `{"a":["a1","a2"]}`,
			overlayStr:   `{"a":["a1"]}`,
			expectedStr:  `{"a":["a1","a2","a1"]}`,
			baseType:     map[string]interface{}{},
			overlayType:  map[string]interface{}{},
			expectedType: map[string]interface{}{},
		},
		{
			name:         `[normal]array内の同じkeyはmergeされ重複する`,
			baseStr:      `["a"]`,
			overlayStr:   `["a"]`,
			expectedStr:  `["a","a"]`,
			baseType:     map[string][]interface{}{},
			overlayType:  map[string][]interface{}{},
			expectedType: map[string][]interface{}{},
		},
		{
			name:         `[normal]多段map内の異なるkeyはmergeされる`,
			baseStr:      `{"a":{"a1":{"a11":"ok"}}}`,
			overlayStr:   `{"a":{"a1":{"a12":"ok"}}}`,
			expectedStr:  `{"a":{"a1":{"a11":"ok","a12":"ok"}}}`,
			baseType:     map[string]interface{}{},
			overlayType:  map[string]interface{}{},
			expectedType: map[string]interface{}{},
		},
		{
			name:         `[normal]array-map内の同じkeyはmergeされ重複する`,
			baseStr:      `{"a":[{"a1":"base01"},{"a1":"base02"}]}`,
			overlayStr:   `{"a":[{"a1":"overlay"}]}`,
			expectedStr:  `{"a":[{"a1":"base01"},{"a1":"base02"},{"a1":"overlay"}]}`,
			baseType:     map[string]interface{}{},
			overlayType:  map[string]interface{}{},
			expectedType: map[string]interface{}{},
		},
		{
			name:         `[normal]array-map内の同じkey-valueはmergeされ重複する`,
			baseStr:      `{"a":[{"a1":"base01"},{"a1":"value"}]}`,
			overlayStr:   `{"a":[{"a1":"value"}]}`,
			expectedStr:  `{"a":[{"a1":"base01"},{"a1":"value"},{"a1":"value"}]}`,
			baseType:     map[string]interface{}{},
			overlayType:  map[string]interface{}{},
			expectedType: map[string]interface{}{},
		},
		{
			name:         `[normal]多段array内の同じkeyはmergeされ重複する`,
			baseStr:      `{"a":[["a11","a12"],["a11","a12"]]}`,
			overlayStr:   `{"a":[["a11"]]}`,
			expectedStr:  `{"a":[["a11","a12"],["a11","a12"],["a11"]]}`,
			baseType:     map[string]interface{}{},
			overlayType:  map[string]interface{}{},
			expectedType: map[string]interface{}{},
		},
		{
			name:            `[abnormal]top-levelの型が違う場合エラーが返る`,
			baseStr:         `{"a":"hoge"}`,
			overlayStr:      `["a"]`,
			baseType:        map[string]interface{}{},
			overlayType:     map[string]interface{}{},
			expectedType:    map[string]interface{}{},
			expectedIsError: true,
		},
		{
			name:            `[abnormal]second-levelの型が違う場合エラーが返る`,
			baseStr:         `{"a":"v"}`,
			overlayStr:      `{"a":1}`,
			baseType:        map[string]interface{}{},
			overlayType:     map[string]interface{}{},
			expectedType:    map[string]interface{}{},
			expectedIsError: true,
		},
		{
			name:            `[abnormal]second-levelの型が違う場合エラーが返る`,
			baseStr:         `{"a":"v"}`,
			overlayStr:      `{"a":{"a1":"v"}}`,
			baseType:        map[string]interface{}{},
			overlayType:     map[string]interface{}{},
			expectedType:    map[string]interface{}{},
			expectedIsError: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			Run(t, tt.baseStr, tt.overlayStr, tt.expectedStr, tt.baseType, tt.expectedType, tt.expectedType, tt.expectedIsError)
		})
	}
}

func Run(t *testing.T, baseStr, overlayStr, expectedStr string, base, overlay, expected interface{}, expectedIsError bool) {
	if err := json.Unmarshal([]byte(baseStr), &base); err != nil {
		panic(err)
	}
	if err := json.Unmarshal([]byte(overlayStr), &overlay); err != nil {
		panic(err)
	}
	if !expectedIsError {
		if err := json.Unmarshal([]byte(expectedStr), &expected); err != nil {
			panic(err)
		}
	}

	// logic
	actual, err := mergeJson(base, overlay)

	// assert
	if !expectedIsError {
		if err != nil {
			t.Fatalf("[test failed] error is occured: %v\n", err)
		}
		equal(t, expected, actual)
	} else {
		if err == nil {
			t.Fatalf("[test failed] error is not occured\n")
		}
	}
}

func equal(t *testing.T, expected, actual interface{}) {
	expectedBytes, err := json.Marshal(expected)
	if err != nil {
		panic(err)
	}
	actualBytes, err := json.Marshal(actual)
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(expectedBytes, actualBytes) {
		t.Fatalf("[test failed] expected: %s, actual: %s\n", string(expectedBytes), string(actualBytes))
	}
}
