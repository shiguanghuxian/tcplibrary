/*
 * @Author: 时光弧线
 * @Date: 2017-12-30 11:55:06
 * @Last Modified by:   时光弧线
 * @Last Modified time: 2017-12-30 11:55:06
 */
package tcplibrary

import (
	"encoding/binary"
	"math"
)

func ByteToBool(i byte) bool {
	if i == 1 {
		return true
	}
	return false
}

func BytesToUint16(data []byte) uint16 {
	return binary.LittleEndian.Uint16(data)
}

func BytesToUint32(data []byte) uint32 {
	return binary.LittleEndian.Uint32(data)
}

func BytesToUint64(data []byte) uint64 {
	return binary.LittleEndian.Uint64(data)
}

func BytesToInt16(data []byte) int16 {
	return int16(BytesToUint16(data))
}

func BytesToInt32(data []byte) int32 {
	return int32(BytesToUint16(data))
}

func BytesToInt64(data []byte) int64 {
	return int64(BytesToUint64(data))
}

func BytesToInt(data []byte) int {
	switch len(data) {
	case 2:
		return int(BytesToUint16(data))
	case 4:
		return int(BytesToUint32(data))
	case 8:
		return int(BytesToUint64(data))
	}
	return 0
}

//IntToBytes 整形转换成byte数组
func IntToBytes(data interface{}) []byte {
	var buf []byte
	switch data.(type) {
	case int:
		buf = make([]byte, 8)
		target, _ := data.(int)
		binary.LittleEndian.PutUint64(buf, uint64(target))
	case int16:
		buf = make([]byte, 2)
		target, _ := data.(int16)
		binary.LittleEndian.PutUint16(buf, uint16(target))
	case int32:
		buf = make([]byte, 4)
		target, _ := data.(int32)
		binary.LittleEndian.PutUint32(buf, uint32(target))
	case int64:
		buf = make([]byte, 8)
		target, _ := data.(int64)
		binary.LittleEndian.PutUint64(buf, uint64(target))
	case uint:
		buf = make([]byte, 8)
		target, _ := data.(uint)
		binary.LittleEndian.PutUint64(buf, uint64(target))
	case uint16:
		buf = make([]byte, 2)
		target, _ := data.(uint16)
		binary.LittleEndian.PutUint16(buf, target)
	case uint32:
		buf = make([]byte, 4)
		target, _ := data.(uint32)
		binary.LittleEndian.PutUint32(buf, target)
	case uint64:
		buf = make([]byte, 8)
		target, _ := data.(uint64)
		binary.LittleEndian.PutUint64(buf, target)
	}
	return buf
}

func Float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func Float32bytes(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}
