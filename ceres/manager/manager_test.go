package manager

import (
	"ceres/aql"
	"ceres/config"
	"ceres/freespace"
	"ceres/schema"
	"ceres/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func Test_readData(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()

	blocks := [][]int{{1, 2}, {4, 4}, {6, 11}}
	var expectedError error
	var expectedData []map[string]interface{}
	expectedError = nil
	for _, block := range blocks {
		for idx := block[0]; idx <= block[1]; idx++ {
			var tempInterface map[string]interface{}
			json.Unmarshal([]byte(fmt.Sprintf("{\"foo\":\"bar\",\".id\":\"bar.%d\"}", idx)), &tempInterface)
			expectedData = append(expectedData, tempInterface)
		}
	}

	actual, err := readData("db1", "foo", "bar", blocks)

	if err != expectedError {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, expectedError)
	}
	if !reflect.DeepEqual(actual, expectedData) {
		t.Errorf("Data was incorrect, got: %v, want: %v", actual, expectedData)
	}
}

func Test_readNoFile(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()

	blocks := [][]int{{1, 2}, {4, 4}, {6, 11}}
	var expectedError error
	var expectedData []map[string]interface{}
	expectedError = &os.PathError{}
	expectedData = nil

	actual, err := readData("db1", "foo", "bar12345", blocks)

	if !errors.As(err, &expectedError) {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, expectedError)
	}
	if !reflect.DeepEqual(actual, expectedData) {
		t.Errorf("Data was incorrect, got: %v, want: %v", actual, expectedData)
	}
}

func Test_readBadContents(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()

	blocks := [][]int{{0, 2}, {4, 4}, {6, 11}}
	var expectedError error
	var expectedData []map[string]interface{}
	expectedError = &json.SyntaxError{}
	expectedData = nil

	actual, err := readData("db1", "foo", "bad", blocks)

	if !errors.As(err, &expectedError) {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, expectedError)
	}
	if !reflect.DeepEqual(actual, expectedData) {
		t.Errorf("Data was incorrect, got: %v, want: %v", actual, expectedData)
	}
}

func Test_writeData(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()

	blocks := [][]int{{20, 22}, {24, 25}}
	data := make([]map[string]interface{}, 0)
	expected := "{\"foo\":\"bar\",\".id\":\"bar.0\"}\n{\"foo\":\"bar\",\".id\":\"bar.1\"}\n{\"foo\":\"bar\",\".id\":\"bar.2\"}\n{\"foo\":\"bar\",\".id\":\"bar.3\"}\n{\"foo\":\"bar\",\".id\":\"bar.4\"}\n{\"foo\":\"bar\",\".id\":\"bar.5\"}\n{\"foo\":\"bar\",\".id\":\"bar.6\"}\n{\"foo\":\"bar\",\".id\":\"bar.7\"}\n{\"foo\":\"bar\",\".id\":\"bar.8\"}\n{\"foo\":\"bar\",\".id\":\"bar.9\"}\n{\"foo\":\"bar\",\".id\":\"bar.10\"}\n{\"foo\":\"bar\",\".id\":\"bar.11\"}\n{\"foo\":\"bar\",\".id\":\"bar.12\"}\n{\"foo\":\"bar\",\".id\":\"bar.13\"}\n{\"foo\":\"bar\",\".id\":\"bar.14\"}\n{\"foo\":\"bar\",\".id\":\"bar.15\"}\n{\"foo\":\"bar\",\".id\":\"bar.16\"}\n{\"foo\":\"bar\",\".id\":\"bar.17\"}\n{\"foo\":\"bar\",\".id\":\"bar.18\"}\n{\"foo\":\"bar\",\".id\":\"bar.19\"}\n{\".id\":\"bar.20\",\"foo\":\"bar\"}\n{\".id\":\"bar.21\",\"foo\":\"bar\"}\n{\".id\":\"bar.22\",\"foo\":\"bar\"}\n\n{\".id\":\"bar.24\",\"foo\":\"bar\"}\n{\".id\":\"bar.25\",\"foo\":\"bar\"}\n\n\n\n\n\n\n"
	var expectedError error
	expectedError = nil
	for idx := 20; idx <= 24; idx++ {
		datum := make(map[string]interface{})
		datum["foo"] = "bar"
		data = append(data, datum)
	}

	err := writeData("db1", "foo", "bar", blocks, data)

	dat, _ := os.ReadFile("../../test/.ceres/data/db1/foo/bar")

	if err != expectedError {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, expectedError)
	}
	if string(dat) != expected {
		t.Errorf("Data was incorrect, got: %v, want: %v", string(dat), expected)
	}
}

func Test_writeNoFile(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()

	blocks := [][]int{{1, 2}, {4, 4}, {6, 11}}
	data := make([]map[string]interface{}, 0)
	var expectedError error
	expectedError = &os.PathError{}

	err := writeData("db1", "foo", "bar12345", blocks, data)

	if !errors.As(err, &expectedError) {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, expectedError)
	}
}

func Test_writeBadContents(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()

	blocks := [][]int{{20, 22}, {24, 25}}
	data := make([]map[string]interface{}, 0)
	var expectedError error
	expectedError = &json.SyntaxError{}
	for idx := 20; idx <= 24; idx++ {
		datum := make(map[string]interface{})
		datum["foo"] = math.Inf(1)
		data = append(data, datum)
	}

	err := writeData("db1", "foo", "bar", blocks, data)

	if !errors.As(err, &expectedError) {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, expectedError)
	}
}

func Test_deleteData(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()

	blocks := [][]int{{20, 22}, {24, 25}}
	expected := "{\"foo\":\"bar\",\".id\":\"bar.0\"}\n{\"foo\":\"bar\",\".id\":\"bar.1\"}\n{\"foo\":\"bar\",\".id\":\"bar.2\"}\n{\"foo\":\"bar\",\".id\":\"bar.3\"}\n{\"foo\":\"bar\",\".id\":\"bar.4\"}\n{\"foo\":\"bar\",\".id\":\"bar.5\"}\n{\"foo\":\"bar\",\".id\":\"bar.6\"}\n{\"foo\":\"bar\",\".id\":\"bar.7\"}\n{\"foo\":\"bar\",\".id\":\"bar.8\"}\n{\"foo\":\"bar\",\".id\":\"bar.9\"}\n{\"foo\":\"bar\",\".id\":\"bar.10\"}\n{\"foo\":\"bar\",\".id\":\"bar.11\"}\n{\"foo\":\"bar\",\".id\":\"bar.12\"}\n{\"foo\":\"bar\",\".id\":\"bar.13\"}\n{\"foo\":\"bar\",\".id\":\"bar.14\"}\n{\"foo\":\"bar\",\".id\":\"bar.15\"}\n{\"foo\":\"bar\",\".id\":\"bar.16\"}\n{\"foo\":\"bar\",\".id\":\"bar.17\"}\n{\"foo\":\"bar\",\".id\":\"bar.18\"}\n{\"foo\":\"bar\",\".id\":\"bar.19\"}\n\n\n\n\n\n\n\n\n\n\n\n\n"
	var expectedError error
	expectedError = nil

	err := deleteData("db1", "foo", "bar", blocks)

	dat, _ := os.ReadFile("../../test/.ceres/data/db1/foo/bar")

	if err != expectedError {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, expectedError)
	}
	if string(dat) != expected {
		t.Errorf("Data was incorrect, got: %v, want: %v", string(dat), expected)
	}
}

func Test_deleteNoFile(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()

	blocks := [][]int{{20, 22}, {24, 25}}
	var expectedError error
	expectedError = &os.PathError{}

	err := deleteData("db1", "foo", "bar12345", blocks)

	if !errors.As(err, &expectedError) {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, expectedError)
	}
}

func TestRead(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()

	ids := []string{"bar.1", "bar.2", "bar.4", "bar.6", "bar.7", "bar.8", "bar.9", "bar.10", "bar.11"}
	var expectedError error
	expectedData := make([]map[string]interface{}, 0)
	expectedError = nil
	for _, id := range ids {
		parts := strings.Split(id, ".")
		var tempInterface map[string]interface{}
		json.Unmarshal([]byte(fmt.Sprintf("{\"foo\":\"bar\",\".id\":\"bar.%s\"}", parts[1])), &tempInterface)
		expectedData = append(expectedData, tempInterface)
	}

	actual, err := Read("db1", "foo", ids)

	if err != expectedError {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, expectedError)
	}
	if !reflect.DeepEqual(actual, expectedData) {
		t.Errorf("Data was incorrect, got: %v, want: %v", actual, expectedData)
	}
}

func TestReadNoFile(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()

	ids := []string{"bar12345.1", "bar.2", "bar.4", "bar.6", "bar.7", "bar.8", "bar.9", "bar.10", "bar.11"}
	var expectedError error
	expectedData := make([]map[string]interface{}, 0)
	expectedError = &os.PathError{}
	expectedData = nil

	actual, err := Read("db1", "foo", ids)

	if !errors.As(err, &expectedError) {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, expectedError)
	}
	if !reflect.DeepEqual(actual, expectedData) {
		t.Errorf("Data was incorrect, got: %v, want: %v", actual, expectedData)
	}
}

func TestWrite(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	data := make([]map[string]interface{}, 0)
	expected := "{\"foo\":\"bar\",\".id\":\"bar.0\"}\n{\"foo\":\"bar\",\".id\":\"bar.1\"}\n{\"foo\":\"bar\",\".id\":\"bar.2\"}\n{\"foo\":\"bar\",\".id\":\"bar.3\"}\n{\"foo\":\"bar\",\".id\":\"bar.4\"}\n{\"foo\":\"bar\",\".id\":\"bar.5\"}\n{\"foo\":\"bar\",\".id\":\"bar.6\"}\n{\"foo\":\"bar\",\".id\":\"bar.7\"}\n{\"foo\":\"bar\",\".id\":\"bar.8\"}\n{\"foo\":\"bar\",\".id\":\"bar.9\"}\n{\"foo\":\"bar\",\".id\":\"bar.10\"}\n{\"foo\":\"bar\",\".id\":\"bar.11\"}\n{\"foo\":\"bar\",\".id\":\"bar.12\"}\n{\"foo\":\"bar\",\".id\":\"bar.13\"}\n{\"foo\":\"bar\",\".id\":\"bar.14\"}\n{\"foo\":\"bar\",\".id\":\"bar.15\"}\n{\"foo\":\"bar\",\".id\":\"bar.16\"}\n{\"foo\":\"bar\",\".id\":\"bar.17\"}\n{\"foo\":\"bar\",\".id\":\"bar.18\"}\n{\"foo\":\"bar\",\".id\":\"bar.19\"}\n{\".id\":\"bar.20\",\"foo\":\"bar\"}\n{\".id\":\"bar.21\",\"foo\":\"bar\"}\n{\".id\":\"bar.22\",\"foo\":\"bar\"}\n{\".id\":\"bar.23\",\"foo\":\"bar\"}\n{\".id\":\"bar.24\",\"foo\":\"bar\"}\n\n\n\n\n\n\n\n"
	var expectedError error
	expectedError = nil
	for idx := 20; idx <= 24; idx++ {
		datum := make(map[string]interface{})
		datum["foo"] = "bar"
		data = append(data, datum)
	}

	err := Write("db1", "foo", data)

	dat, _ := os.ReadFile("../../test/.ceres/data/db1/foo/bar")

	if err != expectedError {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, expectedError)
	}
	if string(dat) != expected {
		t.Errorf("Data was incorrect, got: %v, want: %v", string(dat), expected)
	}
}

func TestWriteBadContents(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	data := make([]map[string]interface{}, 0)
	var expectedError error
	expectedError = &json.SyntaxError{}
	for idx := 20; idx <= 24; idx++ {
		datum := make(map[string]interface{})
		datum["foo"] = math.Inf(1)
		data = append(data, datum)
	}

	err := Write("db1", "foo", data)

	if !errors.As(err, &expectedError) {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, expectedError)
	}
}

func TestWriteOverflow(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	data := make([]map[string]interface{}, 0)
	var expectedError error
	expectedError = nil
	for idx := 0; idx <= 256; idx++ {
		datum := make(map[string]interface{})
		datum["foo"] = "bar"
		data = append(data, datum)
	}

	err := Write("db1", "foo", data)

	files, _ := ioutil.ReadDir("../../test/.ceres/data/db1/foo")

	if err != expectedError {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, expectedError)
	}
	if len(files) != 14 {
		t.Errorf("Number of files was incorrect, got: %d, want: %d", len(files), 14)
	}
}

func TestPatch(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	expected := "{\"foo\":\"bar\",\".id\":\"bar.0\"}\n{\"foo\":\"bar\",\".id\":\"bar.1\"}\n{\"foo\":\"bar\",\".id\":\"bar.2\"}\n{\"foo\":\"bar\",\".id\":\"bar.3\"}\n{\"foo\":\"bar\",\".id\":\"bar.4\"}\n{\"foo\":\"bar\",\".id\":\"bar.5\"}\n{\"foo\":\"bar\",\".id\":\"bar.6\"}\n{\"foo\":\"bar\",\".id\":\"bar.7\"}\n{\"foo\":\"bar\",\".id\":\"bar.8\"}\n{\"foo\":\"bar\",\".id\":\"bar.9\"}\n{\"foo\":\"bar\",\".id\":\"bar.10\"}\n{\"foo\":\"bar\",\".id\":\"bar.11\"}\n{\"foo\":\"bar\",\".id\":\"bar.12\"}\n{\"foo\":\"bar\",\".id\":\"bar.13\"}\n{\"foo\":\"bar\",\".id\":\"bar.14\"}\n{\"foo\":\"bar\",\".id\":\"bar.15\"}\n{\"foo\":\"bar\",\".id\":\"bar.16\"}\n{\"foo\":\"bar\",\".id\":\"bar.17\"}\n{\"foo\":\"bar\",\".id\":\"bar.18\"}\n{\"foo\":\"bar\",\".id\":\"bar.19\"}\n{\".id\":\"bar.20\",\"foo\":\"bar\"}\n{\".id\":\"bar.21\",\"foo\":\"bar\"}\n{\".id\":\"bar.22\",\"foo\":\"bar\"}\n{\".id\":\"bar.23\",\"foo\":\"bar\"}\n{\".id\":\"bar.24\",\"foo\":\"baz\"}\n{\".id\":\"bar.25\",\"foo\":\"bar\"}\n{\".id\":\"bar.26\",\"foo\":\"bar\"}\n{\".id\":\"bar.27\",\"foo\":\"bar\"}\n{\".id\":\"bar.28\",\"foo\":\"bar\"}\n{\".id\":\"bar.29\",\"foo\":\"bar\"}\n{\".id\":\"bar.30\",\"foo\":\"bar\"}\n{\".id\":\"bar.31\",\"foo\":\"bar\"}\n"
	var expectedError error
	expectedError = nil
	datum := map[string]interface{}{"foo": "baz"}
	ids := []string{"bar.24"}

	err := Patch("db1", "foo", ids, datum)

	dat, _ := os.ReadFile("../../test/.ceres/data/db1/foo/bar")

	if err != expectedError {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, expectedError)
	}
	if string(dat) != expected {
		t.Errorf("Data was incorrect, got: %v, want: %v", string(dat), expected)
	}
}

func TestOverWrite(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	expected := "{\"foo\":\"bar\",\".id\":\"bar.0\"}\n{\"foo\":\"bar\",\".id\":\"bar.1\"}\n{\"foo\":\"bar\",\".id\":\"bar.2\"}\n{\"foo\":\"bar\",\".id\":\"bar.3\"}\n{\"foo\":\"bar\",\".id\":\"bar.4\"}\n{\"foo\":\"bar\",\".id\":\"bar.5\"}\n{\"foo\":\"bar\",\".id\":\"bar.6\"}\n{\"foo\":\"bar\",\".id\":\"bar.7\"}\n{\"foo\":\"bar\",\".id\":\"bar.8\"}\n{\"foo\":\"bar\",\".id\":\"bar.9\"}\n{\"foo\":\"bar\",\".id\":\"bar.10\"}\n{\"foo\":\"bar\",\".id\":\"bar.11\"}\n{\"foo\":\"bar\",\".id\":\"bar.12\"}\n{\"foo\":\"bar\",\".id\":\"bar.13\"}\n{\"foo\":\"bar\",\".id\":\"bar.14\"}\n{\"foo\":\"bar\",\".id\":\"bar.15\"}\n{\"foo\":\"bar\",\".id\":\"bar.16\"}\n{\"foo\":\"bar\",\".id\":\"bar.17\"}\n{\"foo\":\"bar\",\".id\":\"bar.18\"}\n{\"foo\":\"bar\",\".id\":\"bar.19\"}\n{\".id\":\"bar.20\",\"foo\":\"bar\"}\n{\".id\":\"bar.21\",\"foo\":\"bar\"}\n{\".id\":\"bar.22\",\"foo\":\"bar\"}\n{\".id\":\"bar.23\",\"foo\":\"bar\"}\n{\".id\":\"bar.24\",\"foo\":\"bar\"}\n{\".id\":\"bar.25\",\"foo\":\"bar\"}\n{\".id\":\"bar.26\",\"foo\":\"bar\"}\n{\".id\":\"bar.27\",\"foo\":\"bar\"}\n{\".id\":\"bar.28\",\"foo\":\"bar\"}\n{\".id\":\"bar.29\",\"foo\":\"bar\"}\n{\".id\":\"bar.30\",\"foo\":\"baz\"}\n{\".id\":\"bar.31\",\"foo\":\"bar\"}\n"
	var expectedError error
	expectedError = nil
	data := []map[string]interface{}{{".id": "bar.24", "foo": "bar"}, {".id": "bar.30", "foo": "baz"}}

	err := OverWrite("db1", "foo", data)

	dat, _ := os.ReadFile("../../test/.ceres/data/db1/foo/bar")

	if err != expectedError {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, expectedError)
	}
	if string(dat) != expected {
		t.Errorf("Data was incorrect, got: %v, want: %v", string(dat), expected)
	}
}

func TestDelete(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	ids := []string{"bar.20", "bar.21", "bar.22", "bar.23", "bar.24", "bar.25", "bar.26", "bar.27", "bar.28", "bar.29", "bar.30", "bar.31"}
	var expectedError error
	expected := "{\"foo\":\"bar\",\".id\":\"bar.0\"}\n{\"foo\":\"bar\",\".id\":\"bar.1\"}\n{\"foo\":\"bar\",\".id\":\"bar.2\"}\n{\"foo\":\"bar\",\".id\":\"bar.3\"}\n{\"foo\":\"bar\",\".id\":\"bar.4\"}\n{\"foo\":\"bar\",\".id\":\"bar.5\"}\n{\"foo\":\"bar\",\".id\":\"bar.6\"}\n{\"foo\":\"bar\",\".id\":\"bar.7\"}\n{\"foo\":\"bar\",\".id\":\"bar.8\"}\n{\"foo\":\"bar\",\".id\":\"bar.9\"}\n{\"foo\":\"bar\",\".id\":\"bar.10\"}\n{\"foo\":\"bar\",\".id\":\"bar.11\"}\n{\"foo\":\"bar\",\".id\":\"bar.12\"}\n{\"foo\":\"bar\",\".id\":\"bar.13\"}\n{\"foo\":\"bar\",\".id\":\"bar.14\"}\n{\"foo\":\"bar\",\".id\":\"bar.15\"}\n{\"foo\":\"bar\",\".id\":\"bar.16\"}\n{\"foo\":\"bar\",\".id\":\"bar.17\"}\n{\"foo\":\"bar\",\".id\":\"bar.18\"}\n{\"foo\":\"bar\",\".id\":\"bar.19\"}\n\n\n\n\n\n\n\n\n\n\n\n\n"
	expectedError = nil

	err := Delete("db1", "foo", ids)

	dat, _ := os.ReadFile("../../test/.ceres/data/db1/foo/bar")

	if err != expectedError {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, expectedError)
	}
	if string(dat) != expected {
		t.Errorf("Data was incorrect, got: %v, want: %v", string(dat), expected)
	}
}

func TestDeleteNoFile(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()

	ids := []string{"bar12345.20", "bar.21", "bar.22", "bar.23", "bar.24", "bar.25", "bar.26", "bar.27", "bar.28", "bar.29", "bar.30", "bar.31"}
	var expectedError error
	expectedError = &os.PathError{}

	err := Delete("db1", "foo", ids)

	if !errors.As(err, &expectedError) {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, expectedError)
	}
}

func TestDoAnd(t *testing.T) {
	A := []string{"foo", "bar", "hello", "world"}
	B := []string{"baz", "bar", "Hello", "world"}
	expected := []string{"bar", "world"}
	output := doAnd(A, B)

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Output was incorrect, got: %v, want: %v", output, expected)
	}
}

func TestDoOr(t *testing.T) {
	A := []string{"foo", "bar", "hello", "world"}
	B := []string{"baz", "bar", "Hello", "world"}
	expected := []string{"Hello", "bar", "baz", "foo", "hello", "world"}
	output := doOr(A, B)

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Output was incorrect, got: %v, want: %v", output, expected)
	}
}

func TestDoNot(t *testing.T) {
	A := []string{"foo", "bar", "hello", "world"}
	B := []string{"bar", "hello", "baz"}
	expected := []string{"foo", "world"}
	output := doNot(A, B)

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Output was incorrect, got: %v, want: %v", output, expected)
	}
}

func TestDoXor(t *testing.T) {
	A := []string{"foo", "bar", "hello", "world"}
	B := []string{"bar", "hello", "baz"}
	expected := []string{"baz", "foo", "world"}
	output := doXor(A, B)

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Output was incorrect, got: %v, want: %v", output, expected)
	}
}

func TestDoOrderASC(t *testing.T) {
	expectedOutput := []map[string]interface{}{{"foo": "a"}, {"foo": "b"}, {"foo": "c"}}
	input := []map[string]interface{}{{"foo": "c"}, {"foo": "b"}, {"foo": "a"}}
	output := doOrderASC(input, "foo")
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("Output was incorrect, got: %v, want: %v", output, expectedOutput)
	}
}

func TestDoOrderDSC(t *testing.T) {
	input := []map[string]interface{}{{"foo": "a"}, {"foo": "b"}, {"foo": "c"}}
	expectedOutput := []map[string]interface{}{{"foo": "c"}, {"foo": "b"}, {"foo": "a"}}
	output := doOrderDSC(input, "foo")
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("Output was incorrect, got: %v, want: %v", output, expectedOutput)
	}
}

func TestProcessFilterBool(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	CreateDatabase("filter")
	CreateCollection("filter", "test", map[string]interface{}{"foo": "BOOL"})
	data := []map[string]interface{}{{"foo": true}, {"foo": true}, {"foo": true}, {"foo": false}, {"foo": false}, {"foo": false}, {"foo": false}, {"foo": false}}
	Write("filter", "test", data)

	files, _ := filePathWalkDir("../../test/.ceres/data/filter/test")
	id_prefix := filepath.Base(files[0])

	expectedData := []string{id_prefix + ".0", id_prefix + ".1", id_prefix + ".2"}

	nodeL := aql.Node{Value: "foo"}
	nodeR := aql.Node{Value: "true"}
	nodeC := aql.Node{Value: "=", Left: &nodeL, Right: &nodeR}

	ids, err := ProcessFilter("filter", "test", nodeC)

	if !reflect.DeepEqual(ids, expectedData) {
		t.Errorf("IDs were incorrect, got: %v, want: %v", ids, expectedData)
	}
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", ids, "<nil>")
	}
	COLDELlection("filter", "test")
	DeleteDatabase("filter")
}

func TestProcessFilterFloat(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	CreateDatabase("filter")
	CreateCollection("filter", "test", map[string]interface{}{"foo": "FLOAT"})
	data := []map[string]interface{}{{"foo": 1.0}, {"foo": 2.0}, {"foo": 2.0}, {"foo": 3.0}, {"foo": 4.0}, {"foo": 5.0}, {"foo": 6.0}, {"foo": 7.0}}
	Write("filter", "test", data)

	files, _ := filePathWalkDir("../../test/.ceres/data/filter/test")
	id_prefix := filepath.Base(files[0])

	expectedData := []string{id_prefix + ".0", id_prefix + ".1", id_prefix + ".2"}

	nodeL := aql.Node{Value: "foo"}
	nodeR := aql.Node{Value: "3.0"}
	nodeC := aql.Node{Value: "<", Left: &nodeL, Right: &nodeR}

	ids, err := ProcessFilter("filter", "test", nodeC)

	if !reflect.DeepEqual(ids, expectedData) {
		t.Errorf("IDs were incorrect, got: %v, want: %v", ids, expectedData)
	}
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", ids, "<nil>")
	}
	COLDELlection("filter", "test")
	DeleteDatabase("filter")
}

func TestProcessFilterString(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	CreateDatabase("filter")
	CreateCollection("filter", "test", map[string]interface{}{"foo": "STRING"})
	data := []map[string]interface{}{{"foo": "1"}, {"foo": "2"}, {"foo": "2"}, {"foo": "3"}, {"foo": "4"}, {"foo": "5"}, {"foo": "6"}, {"foo": "7"}}
	Write("filter", "test", data)

	files, _ := filePathWalkDir("../../test/.ceres/data/filter/test")
	id_prefix := filepath.Base(files[0])

	expectedData := []string{id_prefix + ".0", id_prefix + ".1", id_prefix + ".2"}

	nodeL := aql.Node{Value: "foo"}
	nodeR := aql.Node{Value: "3"}
	nodeC := aql.Node{Value: "<", Left: &nodeL, Right: &nodeR}

	ids, err := ProcessFilter("filter", "test", nodeC)

	if !reflect.DeepEqual(ids, expectedData) {
		t.Errorf("IDs were incorrect, got: %v, want: %v", ids, expectedData)
	}
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", ids, "<nil>")
	}
	COLDELlection("filter", "test")
	DeleteDatabase("filter")
}

func TestProcessFilterGT(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	CreateDatabase("filter")
	CreateCollection("filter", "test", map[string]interface{}{"foo": "INT"})
	data := []map[string]interface{}{{"foo": 1}, {"foo": 2}, {"foo": 2}, {"foo": 3}, {"foo": 4}, {"foo": 5}, {"foo": 6}, {"foo": 7}}
	Write("filter", "test", data)

	files, _ := filePathWalkDir("../../test/.ceres/data/filter/test")
	id_prefix := filepath.Base(files[0])

	expectedData := []string{id_prefix + ".3", id_prefix + ".4", id_prefix + ".5", id_prefix + ".6", id_prefix + ".7"}

	nodeL := aql.Node{Value: "foo"}
	nodeR := aql.Node{Value: "2"}
	nodeC := aql.Node{Value: ">", Left: &nodeL, Right: &nodeR}

	ids, err := ProcessFilter("filter", "test", nodeC)

	if !reflect.DeepEqual(ids, expectedData) {
		t.Errorf("IDs were incorrect, got: %v, want: %v", ids, expectedData)
	}
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", ids, "<nil>")
	}
	COLDELlection("filter", "test")
	DeleteDatabase("filter")
}

func TestProcessFilterGE(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	CreateDatabase("filter")
	CreateCollection("filter", "test", map[string]interface{}{"foo": "INT"})
	data := []map[string]interface{}{{"foo": 1}, {"foo": 2}, {"foo": 2}, {"foo": 3}, {"foo": 4}, {"foo": 5}, {"foo": 6}, {"foo": 7}}
	Write("filter", "test", data)

	files, _ := filePathWalkDir("../../test/.ceres/data/filter/test")
	id_prefix := filepath.Base(files[0])

	expectedData := []string{id_prefix + ".1", id_prefix + ".2", id_prefix + ".3", id_prefix + ".4", id_prefix + ".5", id_prefix + ".6", id_prefix + ".7"}

	nodeL := aql.Node{Value: "foo"}
	nodeR := aql.Node{Value: "2"}
	nodeC := aql.Node{Value: ">=", Left: &nodeL, Right: &nodeR}

	ids, err := ProcessFilter("filter", "test", nodeC)

	if !reflect.DeepEqual(ids, expectedData) {
		t.Errorf("IDs were incorrect, got: %v, want: %v", ids, expectedData)
	}
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", ids, "<nil>")
	}
	COLDELlection("filter", "test")
	DeleteDatabase("filter")
}

func TestProcessFilterLT(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	CreateDatabase("filter")
	CreateCollection("filter", "test", map[string]interface{}{"foo": "INT"})
	data := []map[string]interface{}{{"foo": 1}, {"foo": 2}, {"foo": 2}, {"foo": 3}, {"foo": 4}, {"foo": 5}, {"foo": 6}, {"foo": 7}}
	Write("filter", "test", data)

	files, _ := filePathWalkDir("../../test/.ceres/data/filter/test")
	id_prefix := filepath.Base(files[0])

	expectedData := []string{id_prefix + ".0"}

	nodeL := aql.Node{Value: "foo"}
	nodeR := aql.Node{Value: "2"}
	nodeC := aql.Node{Value: "<", Left: &nodeL, Right: &nodeR}

	ids, err := ProcessFilter("filter", "test", nodeC)

	if !reflect.DeepEqual(ids, expectedData) {
		t.Errorf("IDs were incorrect, got: %v, want: %v", ids, expectedData)
	}
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", ids, "<nil>")
	}
	COLDELlection("filter", "test")
	DeleteDatabase("filter")
}

func TestProcessFilterLE(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	CreateDatabase("filter")
	CreateCollection("filter", "test", map[string]interface{}{"foo": "INT"})
	data := []map[string]interface{}{{"foo": 1}, {"foo": 2}, {"foo": 2}, {"foo": 3}, {"foo": 4}, {"foo": 5}, {"foo": 6}, {"foo": 7}}
	Write("filter", "test", data)

	files, _ := filePathWalkDir("../../test/.ceres/data/filter/test")
	id_prefix := filepath.Base(files[0])

	expectedData := []string{id_prefix + ".0", id_prefix + ".1", id_prefix + ".2"}

	nodeL := aql.Node{Value: "foo"}
	nodeR := aql.Node{Value: "2"}
	nodeC := aql.Node{Value: "<=", Left: &nodeL, Right: &nodeR}

	ids, err := ProcessFilter("filter", "test", nodeC)

	if !reflect.DeepEqual(ids, expectedData) {
		t.Errorf("IDs were incorrect, got: %v, want: %v", ids, expectedData)
	}
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", ids, "<nil>")
	}
	COLDELlection("filter", "test")
	DeleteDatabase("filter")
}

func TestProcessFilterEQ(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	CreateDatabase("filter")
	CreateCollection("filter", "test", map[string]interface{}{"foo": "INT"})
	data := []map[string]interface{}{{"foo": 1}, {"foo": 2}, {"foo": 2}, {"foo": 3}, {"foo": 4}, {"foo": 5}, {"foo": 6}, {"foo": 7}}
	Write("filter", "test", data)

	files, _ := filePathWalkDir("../../test/.ceres/data/filter/test")
	id_prefix := filepath.Base(files[0])

	expectedData := []string{id_prefix + ".1", id_prefix + ".2"}

	nodeL := aql.Node{Value: "foo"}
	nodeR := aql.Node{Value: "2"}
	nodeC := aql.Node{Value: "=", Left: &nodeL, Right: &nodeR}

	ids, err := ProcessFilter("filter", "test", nodeC)

	if !reflect.DeepEqual(ids, expectedData) {
		t.Errorf("IDs were incorrect, got: %v, want: %v", ids, expectedData)
	}
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", ids, "<nil>")
	}
	COLDELlection("filter", "test")
	DeleteDatabase("filter")
}

func TestProcessFilterNE(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	CreateDatabase("filter")
	CreateCollection("filter", "test", map[string]interface{}{"foo": "INT"})
	data := []map[string]interface{}{{"foo": 1}, {"foo": 2}, {"foo": 2}, {"foo": 3}, {"foo": 4}, {"foo": 5}, {"foo": 6}, {"foo": 7}}
	Write("filter", "test", data)

	files, _ := filePathWalkDir("../../test/.ceres/data/filter/test")
	id_prefix := filepath.Base(files[0])

	expectedData := []string{id_prefix + ".0", id_prefix + ".3", id_prefix + ".4", id_prefix + ".5", id_prefix + ".6", id_prefix + ".7"}

	nodeL := aql.Node{Value: "foo"}
	nodeR := aql.Node{Value: "2"}
	nodeC := aql.Node{Value: "!=", Left: &nodeL, Right: &nodeR}

	ids, err := ProcessFilter("filter", "test", nodeC)

	if !reflect.DeepEqual(ids, expectedData) {
		t.Errorf("IDs were incorrect, got: %v, want: %v", ids, expectedData)
	}
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", ids, "<nil>")
	}
	COLDELlection("filter", "test")
	DeleteDatabase("filter")
}

func TestProcessFilterAND(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	CreateDatabase("filter")
	CreateCollection("filter", "test", map[string]interface{}{"foo": "INT"})
	data := []map[string]interface{}{{"foo": 1}, {"foo": 2}, {"foo": 2}, {"foo": 3}, {"foo": 4}, {"foo": 5}, {"foo": 6}, {"foo": 7}}
	Write("filter", "test", data)

	files, _ := filePathWalkDir("../../test/.ceres/data/filter/test")
	id_prefix := filepath.Base(files[0])

	expectedData := []string{id_prefix + ".3"}

	nodeL1 := aql.Node{Value: "foo"}
	nodeR1 := aql.Node{Value: "2"}
	nodeC1 := aql.Node{Value: ">", Left: &nodeL1, Right: &nodeR1}
	nodeL2 := aql.Node{Value: "foo"}
	nodeR2 := aql.Node{Value: "4"}
	nodeC2 := aql.Node{Value: "<", Left: &nodeL2, Right: &nodeR2}
	nodeC3 := aql.Node{Value: "AND", Left: &nodeC1, Right: &nodeC2}

	ids, err := ProcessFilter("filter", "test", nodeC3)

	if !reflect.DeepEqual(ids, expectedData) {
		t.Errorf("IDs were incorrect, got: %v, want: %v", ids, expectedData)
	}
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", ids, "<nil>")
	}
	COLDELlection("filter", "test")
	DeleteDatabase("filter")
}

func TestProcessFilterOR(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	CreateDatabase("filter")
	CreateCollection("filter", "test", map[string]interface{}{"foo": "INT"})
	data := []map[string]interface{}{{"foo": 1}, {"foo": 2}, {"foo": 2}, {"foo": 3}, {"foo": 4}, {"foo": 5}, {"foo": 6}, {"foo": 7}}
	Write("filter", "test", data)

	files, _ := filePathWalkDir("../../test/.ceres/data/filter/test")
	id_prefix := filepath.Base(files[0])

	expectedData := []string{id_prefix + ".0", id_prefix + ".1", id_prefix + ".2", id_prefix + ".3", id_prefix + ".4", id_prefix + ".5", id_prefix + ".6", id_prefix + ".7"}

	nodeL1 := aql.Node{Value: "foo"}
	nodeR1 := aql.Node{Value: "2"}
	nodeC1 := aql.Node{Value: ">", Left: &nodeL1, Right: &nodeR1}
	nodeL2 := aql.Node{Value: "foo"}
	nodeR2 := aql.Node{Value: "4"}
	nodeC2 := aql.Node{Value: "<", Left: &nodeL2, Right: &nodeR2}
	nodeC3 := aql.Node{Value: "OR", Left: &nodeC1, Right: &nodeC2}

	ids, err := ProcessFilter("filter", "test", nodeC3)

	if !reflect.DeepEqual(ids, expectedData) {
		t.Errorf("IDs were incorrect, got: %v, want: %v", ids, expectedData)
	}
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", ids, "<nil>")
	}
	COLDELlection("filter", "test")
	DeleteDatabase("filter")
}

func TestProcessFilterXOR(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	CreateDatabase("filter")
	CreateCollection("filter", "test", map[string]interface{}{"foo": "INT"})
	data := []map[string]interface{}{{"foo": 1}, {"foo": 2}, {"foo": 2}, {"foo": 3}, {"foo": 4}, {"foo": 5}, {"foo": 6}, {"foo": 7}}
	Write("filter", "test", data)

	files, _ := filePathWalkDir("../../test/.ceres/data/filter/test")
	id_prefix := filepath.Base(files[0])

	expectedData := []string{id_prefix + ".0", id_prefix + ".1", id_prefix + ".2", id_prefix + ".4", id_prefix + ".5", id_prefix + ".6", id_prefix + ".7"}

	nodeL1 := aql.Node{Value: "foo"}
	nodeR1 := aql.Node{Value: "2"}
	nodeC1 := aql.Node{Value: ">", Left: &nodeL1, Right: &nodeR1}
	nodeL2 := aql.Node{Value: "foo"}
	nodeR2 := aql.Node{Value: "4"}
	nodeC2 := aql.Node{Value: "<", Left: &nodeL2, Right: &nodeR2}
	nodeC3 := aql.Node{Value: "XOR", Left: &nodeC1, Right: &nodeC2}

	ids, err := ProcessFilter("filter", "test", nodeC3)

	if !reflect.DeepEqual(ids, expectedData) {
		t.Errorf("IDs were incorrect, got: %v, want: %v", ids, expectedData)
	}
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", ids, "<nil>")
	}
	COLDELlection("filter", "test")
	DeleteDatabase("filter")
}

func TestProcessFilterNOT(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	CreateDatabase("filter")
	CreateCollection("filter", "test", map[string]interface{}{"foo": "INT"})
	data := []map[string]interface{}{{"foo": 1}, {"foo": 2}, {"foo": 2}, {"foo": 3}, {"foo": 4}, {"foo": 5}, {"foo": 6}, {"foo": 7}}
	Write("filter", "test", data)

	files, _ := filePathWalkDir("../../test/.ceres/data/filter/test")
	id_prefix := filepath.Base(files[0])

	expectedData := []string{id_prefix + ".4", id_prefix + ".5", id_prefix + ".6", id_prefix + ".7"}

	nodeL2 := aql.Node{Value: "foo"}
	nodeR2 := aql.Node{Value: "4"}
	nodeC2 := aql.Node{Value: "<", Left: &nodeL2, Right: &nodeR2}
	nodeC3 := aql.Node{Value: "NOT", Right: &nodeC2}

	ids, err := ProcessFilter("filter", "test", nodeC3)

	if !reflect.DeepEqual(ids, expectedData) {
		t.Errorf("IDs were incorrect, got: %v, want: %v", ids, expectedData)
	}
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", ids, "<nil>")
	}
	COLDELlection("filter", "test")
	DeleteDatabase("filter")
}

func TestProcessActionCreateDB(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	action := aql.Action{Type: "DBADD", Identifier: "action"}
	_, err := ProcessAction(action, []string{})
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, "<nil>")
	}
	filePaths, _ := ioutil.ReadDir("../../test/.ceres/data")
	files := make([]string, 0)
	for _, file := range filePaths {
		files = append(files, file.Name())
	}
	contained := utils.Contains(files, "action")
	if contained != true {
		t.Errorf("Contained was incorrect, got: %v, want: %v", contained, true)
		t.Errorf("Files: %v", files)
	}

	DeleteDatabase("action")
}

func TestProcessActionCreateCOL(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	action := aql.Action{Type: "COLADD", Identifier: "action.test", Data: []map[string]interface{}{{"item": "STRING", "price": "INT"}}}
	_, err := ProcessAction(action, []string{})
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, "<nil>")
	}
	filePaths, _ := ioutil.ReadDir("../../test/.ceres/data/action")
	files := make([]string, 0)
	for _, file := range filePaths {
		files = append(files, file.Name())
	}
	contained := utils.Contains(files, "test")
	if contained != true {
		t.Errorf("Contained was incorrect, got: %v, want: %v", contained, true)
		t.Errorf("Files: %v", files)
	}

	DeleteDatabase("action")
}

func TestProcessActionDBDEL(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	CreateDatabase("action")

	action := aql.Action{Type: "DBDEL", Identifier: "action"}
	_, err := ProcessAction(action, []string{})
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, "<nil>")
	}
	filePaths, _ := ioutil.ReadDir("../../test/.ceres/data")
	files := make([]string, 0)
	for _, file := range filePaths {
		files = append(files, file.Name())
	}
	contained := utils.Contains(files, "action")
	if contained != false {
		t.Errorf("Contained was incorrect, got: %v, want: %v", contained, false)
		t.Errorf("Files: %v", files)
	}
}

func TestProcessActionCOLDEL(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	CreateDatabase("action")
	CreateCollection("action", "test", map[string]interface{}{"item": "STRING", "price": "INT"})

	action := aql.Action{Type: "COLDEL", Identifier: "action.test"}
	_, err := ProcessAction(action, []string{})
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, "<nil>")
	}
	filePaths, _ := ioutil.ReadDir("../../test/.ceres/data")
	files := make([]string, 0)
	for _, file := range filePaths {
		files = append(files, file.Name())
	}
	contained := utils.Contains(files, "action")
	if contained != true {
		t.Errorf("Contained was incorrect, got: %v, want: %v", contained, true)
		t.Errorf("Files: %v", files)
	}

	DeleteDatabase("action")
}

func TestProcessActionGET(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	CreateDatabase("action")
	CreateCollection("action", "test", map[string]interface{}{"item": "STRING", "price": "INT"})

	inputData := []map[string]interface{}{{"item": "bolt", "price": 5}, {"item": "screw", "price": 3}, {"item": "nail", "price": 2}, {"item": "nut", "price": 10}}
	Write("action", "test", inputData)
	nodeL := aql.Node{Value: "price"}
	nodeR := aql.Node{Value: "2"}
	nodeC := aql.Node{Value: ">", Left: &nodeL, Right: &nodeR}

	expectedData := []map[string]interface{}{{"price": 3}}
	action := aql.Action{Type: "GET", Identifier: "action.test", Fields: []string{"price"}, Limit: 1, Filter: nodeC, Order: "price", OrderDir: "ASC"}
	data, err := ProcessAction(action, []string{})
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, "<nil>")
	}
	if !reflect.DeepEqual(data, expectedData) != true {
		t.Errorf("Data was incorrect, got: %v, want: %v", data, expectedData)
	}

	expectedData = []map[string]interface{}{{"price": 10}}
	action = aql.Action{Type: "GET", Identifier: "action.test", Fields: []string{"price"}, Limit: 1, Filter: nodeC, Order: "price", OrderDir: "DSC"}
	data, err = ProcessAction(action, []string{})
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, "<nil>")
	}
	if !reflect.DeepEqual(data, expectedData) != true {
		t.Errorf("Data was incorrect, got: %v, want: %v", data, expectedData)
	}

	expectedData = []map[string]interface{}{{"price": 2}}
	action = aql.Action{Type: "GET", Identifier: "action.test", Fields: []string{"price"}, Limit: 1, Order: "price", OrderDir: "ASC"}
	data, err = ProcessAction(action, []string{})
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, "<nil>")
	}
	if !reflect.DeepEqual(data, expectedData) != true {
		t.Errorf("Data was incorrect, got: %v, want: %v", data, expectedData)
	}

	expectedData = []map[string]interface{}{{"price": 10}}
	action = aql.Action{Type: "GET", Identifier: "action.test", Fields: []string{"price"}, Limit: 1, Order: "price", OrderDir: "DSC"}
	data, err = ProcessAction(action, []string{})
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, "<nil>")
	}
	if !reflect.DeepEqual(data, expectedData) != true {
		t.Errorf("Data was incorrect, got: %v, want: %v", data, expectedData)
	}

	DeleteDatabase("action")
}

func TestProcessActionPOST(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	CreateDatabase("action")
	CreateCollection("action", "test", map[string]interface{}{"item": "STRING", "price": "INT"})

	inputData := []map[string]interface{}{{"item": "bolt", "price": 5}, {"item": "screw", "price": 3}, {"item": "nail", "price": 2}, {"item": "nut", "price": 10}}
	action := aql.Action{Type: "POST", Identifier: "action.test", Data: inputData}
	_, err := ProcessAction(action, []string{})
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, "<nil>")
	}

	expectedData := []map[string]interface{}{{"price": 2}}
	action = aql.Action{Type: "GET", Identifier: "action.test", Fields: []string{"price"}, Limit: 1, Order: "price", OrderDir: "ASC"}
	data, err := ProcessAction(action, []string{})
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, "<nil>")
	}
	if !reflect.DeepEqual(data, expectedData) != true {
		t.Errorf("Data was incorrect, got: %v, want: %v", data, expectedData)
	}

	DeleteDatabase("action")
}

func TestProcessActionPUT(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	CreateDatabase("action")
	CreateCollection("action", "test", map[string]interface{}{"item": "STRING", "price": "INT"})

	inputData := []map[string]interface{}{{"item": "bolt", "price": 5}, {"item": "screw", "price": 3}, {"item": "nail", "price": 2}, {"item": "nut", "price": 10}}
	action := aql.Action{Type: "POST", Identifier: "action.test", Data: inputData}
	_, err := ProcessAction(action, []string{})
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, "<nil>")
	}

	action = aql.Action{Type: "GET", Identifier: "action.test", Limit: 1, Order: "price", OrderDir: "ASC"}
	data, _ := ProcessAction(action, []string{})

	data[0]["price"] = 20
	action = aql.Action{Type: "PUT", Identifier: "action.test", Data: data}
	_, err = ProcessAction(action, []string{})
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, "<nil>")
	}

	expectedData := []map[string]interface{}{{"price": 3}}
	action = aql.Action{Type: "GET", Identifier: "action.test", Fields: []string{"price"}, Limit: 1, Order: "price", OrderDir: "ASC"}
	data, err = ProcessAction(action, []string{})
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, "<nil>")
	}
	if !reflect.DeepEqual(data, expectedData) != true {
		t.Errorf("Data was incorrect, got: %v, want: %v", data, expectedData)
	}

	DeleteDatabase("action")
}

func TestProcessActionPATCH(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	CreateDatabase("action")
	CreateCollection("action", "test", map[string]interface{}{"item": "STRING", "price": "INT"})

	inputData := []map[string]interface{}{{"item": "bolt", "price": 5}, {"item": "screw", "price": 3}, {"item": "nail", "price": 2}, {"item": "nut", "price": 10}}
	action := aql.Action{Type: "POST", Identifier: "action.test", Data: inputData}
	_, err := ProcessAction(action, []string{})
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, "<nil>")
	}

	action = aql.Action{Type: "GET", Identifier: "action.test", Limit: 1, Fields: []string{".id"}, Order: "price", OrderDir: "ASC"}
	data, _ := ProcessAction(action, []string{})
	ids := []string{data[0][".id"].(string)}

	patchData := []map[string]interface{}{{"price": 20}}
	action = aql.Action{Type: "PATCH", Identifier: "action.test", IDs: ids, Data: patchData}
	_, err = ProcessAction(action, []string{})
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, "<nil>")
	}

	expectedData := []map[string]interface{}{{"price": 3}}
	action = aql.Action{Type: "GET", Identifier: "action.test", Fields: []string{"price"}, Limit: 1, Order: "price", OrderDir: "ASC"}
	data, err = ProcessAction(action, []string{})
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, "<nil>")
	}
	if !reflect.DeepEqual(data, expectedData) != true {
		t.Errorf("Data was incorrect, got: %v, want: %v", data, expectedData)
	}

	patchData = []map[string]interface{}{{"price": 2}}
	action = aql.Action{Type: "PATCH", Identifier: "action.test", IDs: []string{"-"}, Data: patchData}
	_, err = ProcessAction(action, ids)
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, "<nil>")
	}

	expectedData = []map[string]interface{}{{"price": 2}}
	action = aql.Action{Type: "GET", Identifier: "action.test", Fields: []string{"price"}, Limit: 1, Order: "price", OrderDir: "ASC"}
	data, err = ProcessAction(action, []string{})
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, "<nil>")
	}
	if !reflect.DeepEqual(data, expectedData) != true {
		t.Errorf("Data was incorrect, got: %v, want: %v", data, expectedData)
	}

	DeleteDatabase("action")
}

func TestProcessActionDELETE(t *testing.T) {
	os.Setenv("CERES_CONFIG", "../../test/.ceres/config/config.json")
	config.ReadConfigFile()
	freespace.LoadFreeSpace()
	schema.LoadSchema()

	CreateDatabase("action")
	CreateCollection("action", "test", map[string]interface{}{"item": "STRING", "price": "INT"})

	inputData := []map[string]interface{}{{"item": "bolt", "price": 5}, {"item": "screw", "price": 3}, {"item": "nail", "price": 2}, {"item": "nut", "price": 10}}
	action := aql.Action{Type: "POST", Identifier: "action.test", Data: inputData}
	_, err := ProcessAction(action, []string{})
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, "<nil>")
	}

	action = aql.Action{Type: "GET", Identifier: "action.test", Limit: 1, Fields: []string{".id"}, Order: "price", OrderDir: "ASC"}
	data, _ := ProcessAction(action, []string{})
	ids := []string{data[0][".id"].(string)}

	action = aql.Action{Type: "DELETE", Identifier: "action.test", IDs: ids}
	_, err = ProcessAction(action, []string{})
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, "<nil>")
	}

	expectedData := []map[string]interface{}{{"price": 3}}
	action = aql.Action{Type: "GET", Identifier: "action.test", Fields: []string{"price"}, Limit: 1, Order: "price", OrderDir: "ASC"}
	data, err = ProcessAction(action, []string{})
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, "<nil>")
	}
	if !reflect.DeepEqual(data, expectedData) != true {
		t.Errorf("Data was incorrect, got: %v, want: %v", data, expectedData)
	}

	action = aql.Action{Type: "GET", Identifier: "action.test", Limit: 1, Fields: []string{".id"}, Order: "price", OrderDir: "ASC"}
	data, _ = ProcessAction(action, []string{})
	ids = []string{data[0][".id"].(string)}

	action = aql.Action{Type: "DELETE", Identifier: "action.test", IDs: []string{"-"}}
	_, err = ProcessAction(action, ids)
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, "<nil>")
	}

	expectedData = []map[string]interface{}{{"price": 5}}
	action = aql.Action{Type: "GET", Identifier: "action.test", Fields: []string{"price"}, Limit: 1, Order: "price", OrderDir: "ASC"}
	data, err = ProcessAction(action, []string{})
	if err != nil {
		t.Errorf("Error was incorrect, got: %v, want: %v", err, "<nil>")
	}
	if !reflect.DeepEqual(data, expectedData) != true {
		t.Errorf("Data was incorrect, got: %v, want: %v", data, expectedData)
	}

	DeleteDatabase("action")
}
