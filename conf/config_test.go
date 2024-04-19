package conf

import (
	"fmt"
	"testing"
)

func TestInitBusinessConfig(t *testing.T) {
	config := initBusinessConfig()
	if config == nil {
		t.Errorf("config is nil")
	}
	fmt.Println(config)
}
func TestConvertConfig(t *testing.T) {
	originConfig := initBusinessConfig()
	config := convertConfig(originConfig)
	if config == nil {
		t.Errorf("config is nil")
	}
	ethPaymaster := config.SupportPaymaster["Ethereum"]

	fmt.Println(ethPaymaster)
}
func TestStrategyConfig(t *testing.T) {
	GetBasicStrategyConfig("Ethereum")

}
