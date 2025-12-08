package util

import (
	"errors"
	"fmt"
	"log"
	"math"

	"reflect"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Converts a string array to an any array.
func ConvertStringArrayToAnyArray(input []string) []any {
	args := make([]any, len(input))
	for i, v := range input {
		args[i] = v
	}
	return args
}

// Repeats an input for a given amount.
func RepeatInput[T any | string | int](input T, amount int) []T {
	arr := make([]T, amount)
	for i := 0; i < amount; i++ {
		arr[i] = input
	}
	return arr
}

// Converts a struct to an any array.
func ConvertStructToAnyArray(entry interface{}) ([]any, error) {
	val := reflect.ValueOf(entry)
	if val.Kind() != reflect.Struct {
		return nil, errors.New("input is not of type struct")
	}

	numFields := val.NumField()
	slice := make([]any, numFields)

	for i := 0; i < numFields; i++ {
		field := val.Field(i).Interface()
		slice[i] = field
	}

	return slice, nil
}

// Converts a struct's fields pointers to an any array.
func ConvertStructPointerFieldsToPointerAnyArray(entry any) ([]any, error) {
	ptrValue := reflect.ValueOf(entry)

	if ptrValue.Kind() != reflect.Ptr || ptrValue.Elem().Kind() != reflect.Struct {
		return nil, errors.New("input is not of type pointer of struct")
	}

	structType := ptrValue.Elem().Type()
	structValue := ptrValue.Elem()

	fieldPointers := make([]any, structType.NumField())

	for i := 0; i < structType.NumField(); i++ {
		fieldValue := structValue.Field(i).Addr().Interface()
		fieldPointers[i] = fieldValue
	}

	return fieldPointers, nil
}

// ConvertMapToInterface converts a map to an interface.
func ConvertMapToInterface(input map[string]interface{}, outStruct interface{}) error {
	val := reflect.ValueOf(outStruct).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldName := typ.Field(i).Tag.Get("cass")
		if fieldName == "" {
			continue // skip fields without the cass tag
		}

		mapValue, found := input[fieldName]
		if !found {
			continue // Skip fields not present in the map
		}

		field := val.Field(i)
		value := reflect.ValueOf(mapValue)

		if value.Type().ConvertibleTo(field.Type()) {
			field.Set(value.Convert(field.Type()))
		} else {
			return fmt.Errorf("type mismatch for field: %s", fieldName)
		}
	}

	return nil
}

func ConvertStringPointerToPgTypeString(input *string) pgtype.Text {
	if input == nil {
		return pgtype.Text{String: "", Valid: false}
	}
	return pgtype.Text{String: *input, Valid: true}
}

func ConvertFloat32PointerToPgTypeFloat8(input *float32) pgtype.Float8 {
	if input == nil {
		return pgtype.Float8{Float64: 0, Valid: false}
	}
	return pgtype.Float8{Float64: float64(*input), Valid: true}
}

func ConvertIntPointerToPgTypeTimestamptz(input *int) pgtype.Timestamptz {
	if input == nil {
		return pgtype.Timestamptz{Time: time.Now(), Valid: false}
	}
	return pgtype.Timestamptz{Time: time.Unix(int64(*input), 0), Valid: true}
}

var floatType = reflect.TypeOf(float64(0))

func ConvertInterfaceToFloat64(unk interface{}) (float64, error) {
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(floatType) {
		return math.NaN(), fmt.Errorf("cannot convert %v to float64", v.Type())
	}
	fv := v.Convert(floatType)
	return fv.Float(), nil
}

func ConvertStringToFloat(input string) (float64, error) {
	f, err := strconv.ParseFloat(input, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return 0, err
	}
	return f, nil
}

// Convert float64 to pgtype.Numeric
func ToPgNumeric(val float64) pgtype.Numeric {
	n := pgtype.Numeric{}
	valStr := fmt.Sprintf("%.2f", val) // forces format like 890.00
	err := n.Scan(valStr)
	if err != nil {
		log.Printf("ToPgNumeric scan failed: %v", err)
		return pgtype.Numeric{Valid: false}
	}
	return n
}

// Convert int32 to pgtype.Int4
func ToPgInt4(val int32) pgtype.Int4 {
	return pgtype.Int4{Valid: true, Int32: val}
}

// Convert bool to pgtype.Bool
func ToPgBool(val bool) pgtype.Bool {
	return pgtype.Bool{Valid: true, Bool: val}
}

// // Convert db.Status to db.NullStatus
// func ToNullStatus(val string) db.NullStatus {
// 	return db.NullStatus{Status: db.Status(val), Valid: true}
// }

// ParseOptionalInt parses a string to *int. Returns nil if empty or invalid.
func ParseOptionalInt(s string) *int {
	if s == "" {
		return nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &v
}

func PgTimeToString(ts pgtype.Timestamptz) string {
	if ts.Valid {
		return ts.Time.Format("2006-01-02T15:04:05Z07:00")
	}
	return ""
}


func NumericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	v, err := n.Value()
	if err != nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return val
	case string:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0
		}
		return f
	default:
		return 0
	}
} 
