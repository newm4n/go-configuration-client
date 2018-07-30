package go_configuration_client

import "testing"

func TestConfigurationClient_GetConfiguration(t *testing.T) {
	server := NewConfigurationClient("localhost", 8888, 0)
	config, err := server.GetConfiguration("discovery", "master")
	if err != nil {
		t.Fatal(err)
	} else {
		if config.ContainsKey("spring.application.name") == false {
			t.Fatal("Key not found")
		}
		if val, err := config.GetRequired("spring.application.name"); err != nil {
			t.Fatal(err)
		} else {
			if val != "discovery" {
				t.Fatal("Incorrect value")
			}
		}
	}
}
