package go_configuration_client

import (
	"errors"
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
	"encoding/json"
	"reflect"
	"github.com/newm4n/go-utility"
)

type Configuration struct {
	Name            string           `json:"name,omitempty"`
	Profiles        []string         `json:"profiles,omitempty"`
	Label           string           `json:"label,omitempty"`
	Version         string           `json:"version,omitempty"`
	State           string           `json:"state,omitempty"`
	PropertySources []PropertySource `json:"propertySources,omitempty"`
}

func (cnf *Configuration) ContainsKey(key string) bool {
	for _, s := range cnf.PropertySources {
		if s.Source[key] != "" {
			return true
		}
	}
	return false
}

func (cnf *Configuration) Get(key string) interface{} {
	for _, s := range cnf.PropertySources {
		if s.Source[key] != "" {
			log.Println("Found configuration %s [%v]", key, reflect.TypeOf(s.Source[key]))
			return s.Source[key]
		}
	}
	return ""
}

func (cnf *Configuration) GetRequired(key string) (interface{}, error) {
	for _, s := range cnf.PropertySources {
		if s.Source[key] != "" {
			log.Println("Found configuration %s [%v]", key, reflect.TypeOf(s.Source[key]))
			return s.Source[key], nil
		}
	}
	return "", errors.New(fmt.Sprintf("required key '%s' not exist", key))
}

func (cnf *Configuration) GetDefaulted(key string, defaultValue interface{}) interface{} {
	for _, s := range cnf.PropertySources {
		if s.Source[key] != "" {
			log.Println("Found configuration %s [%v]", key, reflect.TypeOf(s.Source[key]))
			return s.Source[key]
		}
	}
	return defaultValue
}

type PropertySource struct {
	Name   string                 `json:"name,omitempty"`
	Source map[string]interface{} `json:"source,omitempty"`
}

type ConfigurationClient struct {
	client *http.Client `json:"-"`
	Host   string
	Port   int
	SecurePort int
}

func NewConfigurationClient(host string, port, securePort int) *ConfigurationClient {
	return &ConfigurationClient{
		client: go_utility.GetDefaultHttpClient(true),
		Host: host,
		Port: port,
		SecurePort: securePort,
	}
}

func (cs *ConfigurationClient) GetConfiguration(name, profile string) (*Configuration, error) {
	var scheme string
	var port int
	if cs.SecurePort > 0 {
		scheme = "https"
		port = cs.SecurePort
	} else if cs.Port > 0 {
		scheme = "http"
		port = cs.Port
	} else {
		return nil, errors.New("both port or secureport is 0, dont know how to connect to config server")
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s://%s:%d/%s/%s", scheme, cs.Host, port, name, profile), nil)
	if err != nil {
		return nil, err
	}
	resp, err := cs.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if err != nil {
			return nil, err
		}
		return nil, errors.New(fmt.Sprintf("unexpected response code %d : %s", resp.StatusCode, string(body)))
	}
	log.Printf("config body : %s", string(body))
	c := &Configuration{}
	if err = json.Unmarshal(body, c); err != nil {
		return nil, err
	}
	return c, nil
}
