package shared

import (
	"encoding/json"
	"log"
	"os"
)

// JSONPretty
func JSONPretty(d interface{}) {
	if val, ok := d.(string); ok {
		JSONUnmarshal([]byte(val), &d)
	}

	bb, _ := json.MarshalIndent(d, " ", " ")
	log.Println(string(bb))
}

// JSONMarshal
func JSONMarshal(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

// JSONMarshalByte
func JSONMarshalByte(obj interface{}) []byte {
	b, _ := json.Marshal(obj)
	return b
}

// JSONMarshalStr
func JSONMarshalStr(obj interface{}) string {
	b, _ := json.Marshal(obj)
	return string(b)
}

// JSONUnmarshal
func JSONUnmarshal(b []byte, destination interface{}) error {
	return json.Unmarshal(b, &destination)
}

// JSONConvert
func JSONConvert(source interface{}, destination interface{}) error {
	b, err := json.Marshal(source)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, &destination); err != nil {
		return err
	}

	return nil
}

// JSONFromFile
func JSONFromFile(path2file string, destination interface{}) error {
	b, err := os.ReadFile(path2file)
	if err != nil {
		return err
	}

	return JSONUnmarshal(b, &destination)
}
