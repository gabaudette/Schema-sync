package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	ModelExt     string
	InterfaceExt string
	EnumExt      string
	DtoExt       string
	ApiDir       string
	MobileAppDir string
	AppDirs      []string
	IgnoredDirs  []string
}

func GetConfig() Config {
	_, err := os.Stat("config.yml")

	if os.IsNotExist(err) {
		// fmt.Println("Configuration file does not exist")

		//fmt.Println("Creating new configuration file...")

		file, err := os.Create("./config.yml")

		if err != nil {
			//fmt.Println("Error creating configuration file:", err)
			log.Fatal(err)
		}

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				//fmt.Println("Error closing configuration file:", err)
				log.Fatal(err)
			}
		}(file)

		_, err = file.WriteString("model_ext: .model\n")

		if err != nil {
			//fmt.Println("Error writing to configuration file:", err)
			log.Fatal(err)
		}

		_, err = file.WriteString("interface_ext: .interface\n")

		if err != nil {
			//fmt.Println("Error writing to configuration file:", err)
			log.Fatal(err)
		}

		_, err = file.WriteString("enum_ext: .enum\n")

		if err != nil {
			//fmt.Println("Error writing to configuration file:", err)
			log.Fatal(err)
		}

		_, err = file.WriteString("dto_ext: .dto\n")

		if err != nil {
			//fmt.Println("Error writing to configuration file:", err)
			log.Fatal(err)
		}

		_, err = file.WriteString("api_dir: ../api/src\n")

		if err != nil {
			//fmt.Println("Error writing to configuration file:", err)
			log.Fatal(err)
		}

		_, err = file.WriteString("mobile_app_dir: ../mobile\n")

		if err != nil {
			//fmt.Println("Error writing to configuration file:", err)
			log.Fatal(err)
		}

		_, err = file.WriteString("app_dirs: [../dashboard]")

		if err != nil {
			//fmt.Println("Error writing to configuration file:", err)
			log.Fatal(err)
		}

		_, err = file.WriteString("ignored_dirs: [node_modules, dist, build, .git]")

		if err != nil {
			//fmt.Println("Error writing to configuration file:", err)
			log.Fatal(err)
		}

		return Config{
			ModelExt:     ".model",
			InterfaceExt: ".interface",
			EnumExt:      ".enum",
			DtoExt:       ".dto",
			ApiDir:       "../api/src",
			MobileAppDir: "../mobile",
			AppDirs:      []string{"../dashboard"},
			IgnoredDirs:  []string{"node_modules", "dist", "build", ".git", "migrations", "data-migrations", "config"},
		}
	}

	//fmt.Println("Reading configuration file...")

	file, err := os.Open("config.yml")
	if err != nil {
		//fmt.Println("Error reading configuration file:", err)
		log.Fatal(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			//fmt.Println("Error closing configuration file:", err)
			log.Fatal(err)
		}
	}(file)

	readFile, err := os.ReadFile("config.yml")
	if err != nil {
		return Config{}
	}

	config := Config{}

	for _, line := range strings.Split(string(readFile), "\n") {
		configValues := strings.Split(line, ": ")
		if len(configValues) != 2 {
			continue
		}

		switch configValues[0] {
		case "model_ext":
			config.ModelExt = configValues[1]
		case "interface_ext":
			config.InterfaceExt = configValues[1]
		case "enum_ext":
			config.EnumExt = configValues[1]
		case "dto_ext":
			config.DtoExt = configValues[1]
		case "api_dir":
			config.ApiDir = configValues[1]
		case "mobile_app_dir":
			config.MobileAppDir = configValues[1]
		case "app_dirs":
			config.AppDirs = strings.Split(configValues[1], ",")
			for i, appDir := range config.AppDirs {
				config.AppDirs[i] = strings.TrimSpace(appDir)
			}
		case "ignored_dirs":
			config.IgnoredDirs = strings.Split(configValues[1], ",")
			for i, ignoredDir := range config.IgnoredDirs {
				config.IgnoredDirs[i] = strings.TrimSpace(ignoredDir)
			}
		default:
			continue
		}
	}

	return config
}
