package convert

import "strconv"

func GetInterfaceToInt(t1 interface{}) int64 {
	var t2 int64
	switch t1.(type) {
	case uint:
		t2 = int64(t1.(uint))
		break
	case int8:
		t2 = int64(t1.(int8))
		break
	case uint8:
		t2 = int64(t1.(uint8))
		break
	case int16:
		t2 = int64(t1.(int16))
		break
	case uint16:
		t2 = int64(t1.(uint16))
		break
	case int32:
		t2 = int64(t1.(int32))
		break
	case uint32:
		t2 = int64(t1.(uint32))
		break
	case int64:
		t2 = int64(t1.(int64))
		break
	case uint64:
		t2 = int64(t1.(uint64))
		break
	case float32:
		t2 = int64(t1.(float32))
		break
	case float64:
		t2 = int64(t1.(float64))
		break
	case string:
		t2, _ = strconv.ParseInt(t1.(string), 10, 64)
		break
	default:
		t2 = t1.(int64)
		break
	}
	return t2
}
