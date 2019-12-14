package byte2struct

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"
)

// ReadBinaryWithStart convert byte array to struct
func ReadBinaryWithStart(data []byte, s interface{}) error {
	dataLen := len(data)
	var structType reflect.Type
	var structValue reflect.Value
	if value, ok := s.(reflect.Value); ok {
		structValue = value
		structType = value.Type()
	} else if reflect.TypeOf(s).Kind() == reflect.Ptr {
		structType = reflect.TypeOf(s).Elem()
		structValue = reflect.ValueOf(s).Elem()
	} else {
		return fmt.Errorf("unsupport type")
	}
	for index := 0; index < structType.NumField(); index++ {
		fieldType := structType.Field(index)
		fieldValue := structValue.Field(index)
		fmt.Printf("field type %+v field value %+v", fieldType, fieldValue)
		start := fieldType.Tag.Get("start")
		if start == "" {
			return fmt.Errorf("tag miss start")
		}
		startInt, _ := strconv.Atoi(start)
		endInt := startInt + int(fieldType.Type.Size())
		if startInt > dataLen || endInt > dataLen {
			return fmt.Errorf("start out of range")
		}
		selectData := data[startInt:endInt]
		switch fieldType.Type.Kind() {
		case reflect.Uint32:
			ret := binary.LittleEndian.Uint32(selectData)
			structValue.Field(index).Set(reflect.ValueOf(ret))
		case reflect.Int32:
			ret := binary.LittleEndian.Uint32(selectData)
			structValue.Field(index).Set(reflect.ValueOf(int32(ret)))
		case reflect.Int64:
			ret := binary.LittleEndian.Uint64(selectData)
			structValue.Field(index).Set(reflect.ValueOf(int64(ret)))
		case reflect.Uint64:
			ret := binary.LittleEndian.Uint64(selectData)
			structValue.Field(index).Set(reflect.ValueOf(ret))
		case reflect.Struct:
			err := ReadBinaryWithStart(data[startInt:], fieldValue)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupport type %s", fieldType.Type.Kind())
		}
	}
	return nil
}
