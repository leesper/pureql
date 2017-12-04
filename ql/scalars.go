package ql

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

// built-in scalars
var (
	Int     = NewScalar(resultInt, inputInt)
	Float   = NewScalar(resultFloat, inputFloat)
	String  = NewScalar(resultString, inputString)
	Boolean = NewScalar(resultBoolean, inputBoolean)
	ID      = NewScalar(resultString, inputID)
)

func resultInt(value interface{}) interface{} {
	switch value := value.(type) {
	case int:
		if value > math.MaxInt32 || value < math.MinInt32 {
			return errors.New("integer out of range(-2^31~2^31)")
		}
		return value
	case float32:
		if value < float32(math.MinInt32) || value > float32(math.MaxInt32) {
			return errors.New("integer out of range(-2^31~2^31)")
		}
		return int(value)
	case float64:
		if value < float64(math.MinInt32) || value > float64(math.MaxInt32) {
			return errors.New("integer out of range(-2^31~2^31)")
		}
		return int(value)
	case string:
		val, err := strconv.ParseFloat(value, 0)
		if err != nil {
			return err
		}
		return resultInt(val)
	case bool:
		if value {
			return 1
		}
		return 0
	case *int:
		return *value
	case *float32:
		return resultInt(*value)
	case *float64:
		return resultInt(*value)
	case *string:
		return resultInt(*value)
	case *bool:
		return resultInt(*value)
	default:
		return fmt.Errorf("cannot coerce %T", value)
	}
}

func inputInt(value interface{}) interface{} {
	switch value := value.(type) {
	case int:
		if value > math.MaxInt32 || value < math.MinInt32 {
			return errors.New("integer out of range(-2^31~2^31)")
		}
		return value
	case *int:
		return *value
	default:
		return fmt.Errorf("cannot coerce %T", value)
	}
}

func resultFloat(value interface{}) interface{} {
	switch value := value.(type) {
	case int:
		return float64(value)
	case float32:
		return float64(value)
	case float64:
		return value
	case string:
		val, err := strconv.ParseFloat(value, 0)
		if err != nil {
			return err
		}
		return val
	case bool:
		if value {
			return float64(1)
		}
		return float64(0)
	case *int:
		return float64(*value)
	case *float32:
		return float64(*value)
	case *float64:
		return *value
	case *string:
		return resultFloat(*value)
	case *bool:
		return resultFloat(*value)
	default:
		return fmt.Errorf("cannot coerce %T", value)
	}
}

func inputFloat(value interface{}) interface{} {
	switch value := value.(type) {
	case int:
		return float64(value)
	case float32:
		return float64(value)
	case float64:
		return value
	case *int:
		return float64(*value)
	case *float32:
		return float64(*value)
	case *float64:
		return float64(*value)
	default:
		return fmt.Errorf("cannot coerce %T", value)
	}
}

func resultString(value interface{}) interface{} {
	if val, ok := value.(*string); ok {
		return *val
	}
	return fmt.Sprintf("%v", value)
}

func inputString(value interface{}) interface{} {
	switch value := value.(type) {
	case *string:
		return *value
	case string:
		return value
	default:
		return fmt.Errorf("cannot coerce %T", value)
	}
}

func resultBoolean(value interface{}) interface{} {
	switch value := value.(type) {
	case int:
		if value == 0 {
			return false
		}
		return true
	case float32:
		if value != 0 {
			return true
		}
		return false
	case float64:
		if value != 0 {
			return true
		}
		return false
	case string:
		if value == "false" || value == "" {
			return false
		}
		return true
	case bool:
		return value
	case *int:
		return resultBoolean(*value)
	case *float32:
		return resultBoolean(*value)
	case *float64:
		return resultBoolean(*value)
	case *string:
		return resultBoolean(*value)
	case *bool:
		return resultBoolean(*value)
	default:
		return fmt.Errorf("cannot coerce %T", value)
	}
}

func inputBoolean(value interface{}) interface{} {
	switch value := value.(type) {
	case bool:
		return value
	case *bool:
		return *value
	default:
		return fmt.Errorf("cannot coerce %T", value)
	}
}

func inputID(value interface{}) interface{} {
	switch value := value.(type) {
	case int:
		return resultString(value)
	case string:
		return resultString(value)
	case *int:
		return resultString(*value)
	case *string:
		return resultString(*value)
	default:
		return fmt.Errorf("cannot coerce %T", value)
	}
}
