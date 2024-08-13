package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/veyanrech/homeWebCamera/imagecapture/utils"
)

type Config interface {
	GetString(key string) string
	GetSliceOfStrings(key string) []string
	GetInt(key string) int
	Set(key string, value interface{})
}

type config struct {
	Data map[string]interface{}
}

func NewConfig() Config {
	res := &config{}

	res.Data = make(map[string]interface{})

	filename := "config.json"
	switch opsys := utils.GetOS(); opsys {
	default:
		filename = "." + string(os.PathSeparator) + "macos.config.json"
	case "windows":
		filename = "." + string(os.PathSeparator) + "win.config.json"
	}

	err := res.loadConfigFromFile(filename)

	if err != nil {
		return nil
	}

	return res

}

func (c *config) GetString(key string) string {
	v, ok := c.Data[key]
	if !ok {
		return ""
	}

	switch valueType := v.(type) {
	case string:
		return valueType
	default:
		return fmt.Sprintf("%v", valueType)
	}
}

func (c *config) GetSliceOfStrings(key string) []string {
	v, ok := c.Data[key]
	if !ok {
		return nil
	}

	switch valueType := v.(type) {
	case []interface{}:

		var res []string
		for _, v := range valueType {
			switch vtype := v.(type) {
			case string:
				res = append(res, vtype)
			}
		}

		return res
	default:
		return nil
	}
}

func (c *config) GetInt(key string) int {
	v, ok := c.Data[key]
	if !ok {
		return 0
	}

	switch valueType := v.(type) {
	case int:
		return valueType
	case float64:
		return int(valueType)
	default:
		return 0
	}
}

func (c *config) Set(key string, value interface{}) {
	c.Data[key] = value
}

func (c *config) loadConfigFromFile(filePath string) error {
	filereader, err := os.Open(filePath)
	if err != nil {
		return err
	}

	//read file
	fileContent, err := io.ReadAll(filereader)
	if err != nil {
		return err
	}
	//parse as json
	err = json.Unmarshal(fileContent, &c.Data)
	if err != nil {
		return err
	}

	return nil
}
