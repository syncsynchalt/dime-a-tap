package disklog

import (
	"fmt"
	"os"
	"time"
)

func DumpPacket(directory, label string, data []byte) error {
	if directory == "" {
		return nil
	}
	filename := fmt.Sprintf("%s/%s.%s", directory, label, time.Now().Format("20060102150405.000"))
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("unable to log packet: %s", err.Error())
	}
	defer file.Close()
	file.Write(data)
	return nil
}
