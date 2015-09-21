package session

import (
	"bytes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/gob"
)

func init() {
	gob.Register([]interface{}{})
	gob.Register(map[int]interface{}{})
	gob.Register(map[string]interface{}{})
	gob.Register(map[interface{}]interface{}{})
	gob.Register(map[string]string{})
	gob.Register(map[int]string{})
	gob.Register(map[int]int{})
	gob.Register(map[int]int64{})
}

func encodeGod(data map[interface{}]interface{}) ([]byte, error) {
	for _, v := range data {
		gob.Register(v)
	}
	buf := bytes.NewBuffer(nil)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(data)
	if err != nil {
		return []byte(""), err
	}
	return buf.Bytes(), nil
}

func decodeGob(data []byte) (map[interface{}]interface{}, error) {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	var out map[interface{}]interface{}
	err := decoder.Decode(&out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func encode(value []byte) []byte {
	encoded := make([]byte, base64.URLEncoding.EncodedLen(len(value)))
	base64.URLEncoding.Encode(encoded, value)
	return encoded
}

func decode(value []byte) ([]byte, error) {
	decoded := make([]byte, base64.URLEncoding.DecodedLen(len(value)))
	b, err := base64.URLEncoding.Decode(decoded, value)
	if err != nil {
		return nil, err
	}
	return decoded[:b], nil
}

func encodeCookie(block cipher.Block, hashKey, name string, value map[interface{}]interface{}) (string, error) {
	var err error
	var b []byte

}
