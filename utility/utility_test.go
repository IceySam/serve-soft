package utility_test

import (
	"database/sql"
	"reflect"
	"strings"
	"testing"

	"github.com/IceySam/serve-soft/examples"
	"github.com/IceySam/serve-soft/utility"
)

type Animal struct {
	Id    int
	Niche string
	Color string
	Age   int
}

type Case struct {
	flat int
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		data    interface{}
		wantErr bool
		msg     string
	}{
		{
			name:    "valid",
			data:    examples.Food{Name: "Beans"},
			wantErr: false,
		},
		{
			name:    "missing name",
			data:    examples.Food{},
			wantErr: true,
			msg:     "name is required",
		},
		{
			name:    "incomplete animal",
			data:    Animal{Id: 1, Niche: "Forest", Color: "gray", Age: 2},
			wantErr: false,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			err := utility.Validate(v.data)
			if (err != nil) != v.wantErr {
				t.Errorf("Validate() got = %v, wantErr %v", err, v.wantErr)
				return
			}
			if v.wantErr && !strings.Contains(err.Error(), v.msg) {
				t.Errorf("Validate() error message = %v, msg to contain %v", err, v.msg)
			}
		})
	}
}

func TestToMap(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		wantType reflect.Type
		wantName string
		wantErr  bool
	}{
		{
			name:     "simple struct",
			input:    Animal{Id: 1, Niche: "Forest", Color: "gray", Age: 2},
			wantType: reflect.TypeOf(Animal{}),
			wantName: "Animal",
			wantErr:  false,
		},
		{
			name:     "invalid struct (non-JSON-serializable)",
			input:    Case{flat: 8},
			wantType: reflect.TypeOf(Case{}),
			wantName: "Case",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotType, gotName, err := utility.ToMap(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotType != tt.wantType {
				t.Errorf("ToMap() gotType = %v, want %v", gotType, tt.wantType)
			}
			if gotName != tt.wantName {
				t.Errorf("ToMap() gotName = %v, want %v", gotName, tt.wantName)
			}
		})
	}
}

func BenchmarkToMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utility.ToMap(Animal{Id: 1, Niche: "Forest", Color: "gray", Age: 2})
	}
}

func TestToStruct(t *testing.T) {
	tests := []struct {
		name     string
		inputMap map[string]interface{}
		outPtr   interface{}
		wantErr  bool
	}{
		{
			name:     "simple valid",
			inputMap: map[string]interface{}{"Id": 1, "Niche": "Forest", "Color": "gray", "Age": 2},
			outPtr:   &Animal{},
			wantErr:  false,
		},
		{
			name:     "wrong map field",
			inputMap: map[string]interface{}{"Id": 1, "Niche": 34, "Color": "gray", "Age": 2},
			outPtr:   &Animal{},
			wantErr:  true,
		},
		{
			name:     "wrong ptr",
			inputMap: map[string]interface{}{"Id": 1, "Niche": "Forest"},
			outPtr:   &Case{},
			wantErr:  true,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			err := utility.ToStruct(v.inputMap, v.outPtr)
			if (err != nil) != v.wantErr {
				t.Errorf("ToStruct() error = %v, wantErr %v", err, v.wantErr)
				return
			}

			if !v.wantErr {
				ptr := v.outPtr.(*Animal)
				if reflect.TypeOf(*ptr).NumField() != 4 {
					t.Errorf("ToStruct() error = %v", "field mismatch")
				}
			}

		})
	}
}

func TestToStructArray(t *testing.T) {
	tests := []struct {
		name          string
		inputMapSlice []map[string]interface{}
		outPtr        interface{}
		wantErr       bool
	}{
		{
			name:          "simple valid",
			inputMapSlice: []map[string]interface{}{{"Id": 1, "Niche": "Forest", "Color": "gray", "Age": 2}, {"Id": 2, "Niche": "mountain", "Color": "dark", "Age": 5}},
			outPtr:        &[]Animal{},
			wantErr:       false,
		},
		{
			name:          "wrong map field",
			inputMapSlice: []map[string]interface{}{{"Id": 1, "Niche": 34, "Color": "gray", "Age": 2}},
			outPtr:        &[]Animal{},
			wantErr:       true,
		},
		{
			name:          "wrong ptr",
			inputMapSlice: []map[string]interface{}{{"Id": 1, "Niche": "Forest"}},
			outPtr:        &[]Case{},
			wantErr:       true,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			err := utility.ToStructArray(v.inputMapSlice, v.outPtr)
			if (err != nil) != v.wantErr {
				t.Errorf("ToStructArray() error = %v, wantErr %v", err, v.wantErr)
				return
			}

		})
	}
}

func TestTypeEquals(t *testing.T) {
	tests := []struct {
		name      string
		data      any
		compare   any
		wantEqual bool
	}{
		{
			name:      "same type",
			data:      10,
			compare:   10,
			wantEqual: true,
		},
		{
			name:      "different type",
			data:      "hello",
			compare:   10,
			wantEqual: false,
		},
		// {
		// 	name:      "same underlying type",
		// 	data:      int32(10),
		// 	compare:   10,
		// 	wantEqual: true,
		// },
		{
			name:      "pointer vs non-pointer",
			data:      reflect.Ptr,
			compare:   10,
			wantEqual: false,
		},
		{
			name:      "interface vs concrete type",
			data:      12,
			compare:   interface{}(12),
			wantEqual: true,
		},
		{
			name:      "nil vs nil",
			data:      nil,
			compare:   nil,
			wantEqual: true,
		},
		{
			name:      "nil vs non-nil",
			data:      nil,
			compare:   10,
			wantEqual: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEqual := utility.TypeEquals(tt.data, tt.compare)
			if gotEqual != tt.wantEqual {
				t.Errorf("TypeEquals(%v, %v) = %v, want %v", tt.data, tt.compare, gotEqual, tt.wantEqual)
			}
		})
	}
}

func TestParseAny(t *testing.T) {
	tests := []struct {
		name    string
		input   sql.RawBytes
		want    any
		wantErr bool
	}{
		{
			name:    "valid int",
			input:   []byte("123"),
			want:    int64(123),
			wantErr: false,
		},
		{
			name:    "valid float",
			input:   []byte("3.14"),
			want:    float64(3.14),
			wantErr: false,
		},
		{
			name:    "valid bool (lowercase)",
			input:   []byte("true"),
			want:    true,
			wantErr: false,
		},
		{
			name:    "valid bool (uppercase)",
			input:   []byte("TRUE"),
			want:    true,
			wantErr: false,
		},
		{
			name:    "nil byte slice",
			input:   nil,
			want:    "NULL",
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   []byte(""),
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utility.ParseAny(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseAny() got = %v, want %v", got, tt.want)
			}
		})
	}
}
