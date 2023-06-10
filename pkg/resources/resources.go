package resources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/lensesio/tableprinter"
	"gopkg.in/yaml.v3"
)

var supportedTypes = map[string]struct{}{
	"vpc": {},
	"eks": {},
}

type AccountResources map[string][]Resource
type Resource struct {
	Type    string
	Account string
	Region  string
	ID      string
}

func Read(resourceStr, fileName, fileFormat string) (AccountResources, error) {
	resources, err := parseResourceString(resourceStr)
	if err != nil {
		return nil, err
	}
	res, err := readResourceFile(fileName, fileFormat)
	if err != nil {
		return nil, err
	}
	result := make(map[string][]Resource)
	for _, r := range append(resources, res...) {
		if _, ok := result[r.Account]; !ok {
			result[r.Account] = make([]Resource, 0)
		}
		result[r.Account] = append(result[r.Account], r)
	}
	return result, nil
}

type line struct {
	Index   int    `header:"#"`
	Account string `header:"Account"`
	Region  string `header:"Region"`
	Type    string `header:"Type"`
	ID      string `header:"ID/Name"`
}

func (ar AccountResources) String() string {
	var lines []line
	counter := 0
	for account, res := range ar {
		for _, r := range res {
			counter++
			lines = append(lines, line{
				Index:   counter,
				Account: account,
				Region:  r.Region,
				Type:    r.Type,
				ID:      r.ID,
			})
		}
	}
	buf := bytes.NewBufferString("")
	tableprinter.Print(buf, lines)
	return buf.String()
}

func parseResourceString(resourceStr string) ([]Resource, error) {
	if resourceStr == "" {
		return nil, nil
	}
	var resources []Resource
	for _, r := range strings.Split(resourceStr, ",") {
		parts := strings.Split(r, ":")
		account, region := "default", "default"
		var rType, id string
		switch len(parts) {
		case 2:
			rType, id = parts[0], parts[1]
		case 3:
			rType, region, id = parts[0], parts[1], parts[2]
		case 4:
			rType, account, region, id = parts[0], parts[1], parts[2], parts[3]
		default:
			return nil, fmt.Errorf("invalid resource: %s", r)
		}
		rType = strings.ToLower(rType)
		if _, ok := supportedTypes[rType]; !ok {
			return nil, fmt.Errorf("unsupported resource type: %s [%s]", rType, r)
		}
		resources = append(resources, Resource{
			Type:    rType,
			Account: account,
			Region:  region,
			ID:      id,
		})
	}
	return resources, nil
}

var reader = map[string]func(in []byte, out interface{}) error{
	"json": json.Unmarshal,
	"yaml": yaml.Unmarshal,
}

func readResourceFile(name, format string) ([]Resource, error) {
	if name == "" {
		return nil, nil
	}
	file, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	var content struct {
		Resources []string `json:"resources" yaml:"resources"`
	}
	if err := reader[format](file, &content); err != nil {
		return nil, err
	}
	var resources []Resource
	for _, r := range content.Resources {
		res, err := parseResourceString(r)
		if err != nil {
			return nil, err
		}
		resources = append(resources, res...)
	}
	return resources, nil
}
