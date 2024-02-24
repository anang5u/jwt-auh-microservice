package shared

import "os"

func FileWriteToFile(data interface{}, filename string) error {
	file, err := JSONMarshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, file, 0644)
}
