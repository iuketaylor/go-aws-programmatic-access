package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/eiannone/keyboard"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	WormholePath       string `envconfig:"WORMHOLE_PATH"`
	AwsCredentialsPath string `envconfig:"AWS_CREDENTIALS_PATH"`
	DevId              string `envconfig:"DEV_ID"`
	ToolingId          string `envconfig:"TOOLING_ID"`
	ProdId             string `envconfig:"PROD_ID"`
	IcatId             string `envconfig:"ICAT_ID"`
	InnovationId       string `envconfig:"INNOVATION_ID"`
}

func main() {
	cfg, configError := getConfig()
	if configError != nil {
		log.Fatal(configError)
	}
	options := []string{"ibl-dev", "ibl-tooling", "ibl-prod", "icat", "innovation"}
	selectedIndex := 0

	openKeyboardError := keyboard.Open()
	if openKeyboardError != nil {
		panic(openKeyboardError)
	}
	defer keyboard.Close()

	printOptions(options, selectedIndex)

	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyArrowDown {
			selectedIndex = (selectedIndex + 1) % len(options)
		} else if key == keyboard.KeyArrowUp {
			if selectedIndex == 0 {
				selectedIndex = len(options) - 1
			} else {
				selectedIndex = (selectedIndex - 1) % len(options)
			}
		} else if key == keyboard.KeyEnter {
			break
		} else if key == keyboard.KeyEsc {
			return
		}

		// Move cursor up to the first option
		fmt.Printf("\033[%dA", len(options))
		// Print options again with updated selection
		printOptions(options, selectedIndex)
	}

	// Execute shell script based on the selected option
	switch options[selectedIndex] {
	case "ibl-dev":
		runScript(cfg.DevId, options[selectedIndex], *cfg)
	case "ibl-tooling":
		runScript(cfg.ToolingId, options[selectedIndex], *cfg)
	case "ibl-prod":
		runScript(cfg.ProdId, options[selectedIndex], *cfg)
	case "icat":
		runScript(cfg.IcatId, options[selectedIndex], *cfg)
	case "innovation":
		runScript(cfg.InnovationId, options[selectedIndex], *cfg)
	default:
		fmt.Println("No script configured for this option")
	}
}

func getConfig() (*Config, error) {
	loadConfigError := godotenv.Load()
	if loadConfigError != nil {
		return nil, fmt.Errorf("Error loading .env file: %w", loadConfigError)
	}

	var cfg Config
	loadEnvVariablesError := envconfig.Process("", &cfg)
	if loadEnvVariablesError != nil {
		return nil, fmt.Errorf("Error loading environment variables: %s\n", loadEnvVariablesError)
	}

	return &cfg, nil
}

func runScript(accountId string, accountName string, cfg Config) {
	cmd := exec.Command(cfg.WormholePath, accountId)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = writeToFile(out, cfg.AwsCredentialsPath)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	err = exec.Command("bash", "-c", "source /Users/taylol04/.aws_credentials").Run()
	if err != nil {
		fmt.Println("Error sourcing file:", err)
		return
	}

	fmt.Printf("Successfully wormholed into %s. \n", accountName)
}

func printOptions(options []string, selectedIndex int) {
	for i, option := range options {
		if i == selectedIndex {
			fmt.Printf("\x1b[7m%s\x1b[0m\n", option)
		} else {
			fmt.Println(option)
		}
	}
}

func writeToFile(data []byte, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
