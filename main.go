package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

var (
	appName    string
	appVersion string
)

func init() {
	testing.Init()

	var showVersion bool
	flag.BoolVar(&showVersion, "v", false, "show application version")
	flag.BoolVar(&showVersion, "version", false, "show application version")
	flag.Parse()

	if showVersion {
		fmt.Printf("%s %s\n", appName, appVersion)
		os.Exit(0)
	}
}

func main() {
	// parse flag
	var n int
	var err error
	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("invalid arguments")
	}
	dir := args[0]
	jsonName := args[1]
	if len(args) >= 3 {
		n, err = strconv.Atoi(args[2])
		if err != nil {
			log.Fatal("invalid arguments")
		}
	}
	// n: 何階層辿るか を決定
	absPath, err := filepath.Abs(dir)
	if err != nil {
		log.Fatal(err)
	}
	pathes := len(strings.Split(filepath.ToSlash(absPath), "/"))
	if n <= 0 || pathes < n {
		n = pathes
	}

	// merge した結果を格納
	result := ReadFileAndMergeJson(absPath, jsonName, n)
	output, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}

	// output
	if result == nil {
		fmt.Println("{}")
	} else {
		fmt.Println(string(output))
	}
}

func ReadFileAndMergeJson(dir, jsonName string, n int) interface{} {
	var result interface{}
	for i := 0; i < n; i++ {

		// .metadata.json の読み込み
		jsonFilePath := filepath.Join([]string{dir, jsonName}...)
		dir = filepath.Dir(filepath.Dir(jsonFilePath))
		raw, err := ioutil.ReadFile(jsonFilePath)
		if err != nil {
			debugf("%s :file not found\n", jsonFilePath)
			continue
		}
		var m interface{}
		if err := json.Unmarshal(raw, &m); err != nil {
			debugf("Error: Unmarshal JSON\n%v\n", err.Error())
			continue
		}

		// 先程読み込んだ json を今までの json へマージした変数 "result" を取得
		if result == nil || reflect.ValueOf(result).IsNil() {
			result = m
		} else if reflect.TypeOf(result).Kind() != reflect.TypeOf(m).Kind() {
			err = fmt.Errorf("invalid type")
		}
		switch m.(type) {
		case map[string]interface{}:
			result, err = mergeJson(result.(map[string]interface{}), m.(map[string]interface{}))
		case []interface{}:
			result, err = mergeJson(result.([]interface{}), m.([]interface{}))
		default:
			err = fmt.Errorf("unsupported type")
		}
		if err != nil {
			debugf("Error\n%v\n", err.Error())
			continue
		}
	}
	return result
}

func mergeJson(base, overlay interface{}) (interface{}, error) {
	if reflect.TypeOf(base).Kind() != reflect.TypeOf(overlay).Kind() {
		return nil, fmt.Errorf("invalid type")
	}

	switch base.(type) {
	case map[string]interface{}:
		baseMap := base.(map[string]interface{})
		overlayMap := overlay.(map[string]interface{})
		result, err := deepCopyMap(baseMap)
		if err != nil {
			return nil, err
		}
		for key, val := range overlayMap {
			if _, ok := baseMap[key]; !ok {
				result[key] = val
			} else {
				if reflect.TypeOf(val).Kind() != reflect.TypeOf(baseMap[key]).Kind() {
					return nil, fmt.Errorf("no much type")
				}
				switch val.(type) {
				case map[string]interface{}:
					r, err := mergeJson(baseMap[key].(map[string]interface{}), overlayMap[key].(map[string]interface{}))
					if err != nil {
						return nil, err
					}
					result[key] = r
				case []interface{}:
					r, err := mergeJson(baseMap[key].([]interface{}), overlayMap[key].([]interface{}))
					if err != nil {
						return nil, err
					}
					result[key] = r
				default:
					// pass
				}
			}
		}
		return result, nil

	case []interface{}:
		baseSlice := base.([]interface{})
		overlaySlice := overlay.([]interface{})
		result, err := deepCopySlice(baseSlice)
		if err != nil {
			return nil, err
		}
		return append(result, overlaySlice...), nil

	default:
		return nil, fmt.Errorf("unsupported type")
	}
}

func contains(s []interface{}, e interface{}) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func deepCopyMap(base map[string]interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(base)
	if err != nil {
		return nil, err
	}
	o := make(map[string]interface{})
	if err := json.Unmarshal(b, &o); err != nil {
		return nil, err
	}
	return o, nil
}

func deepCopySlice(base []interface{}) ([]interface{}, error) {
	b, err := json.Marshal(base)
	if err != nil {
		return nil, err
	}
	var o []interface{}
	if err := json.Unmarshal(b, &o); err != nil {
		return nil, err
	}
	return o, nil
}

func debugf(format string, a ...interface{}) {
	// for debug
	if false {
		fmt.Printf(format, a...)
	}
}
