package api

import (
	"os"
	"ssync/internal/config"
	"ssync/internal/files"
	"ssync/internal/tsgenerator"
)

func ListFiles() {
	conf := config.GetConfig()
	_, paths, _ := readApi()

	tsgenerator.ProcessModelFiles(paths, conf.ModelExt)
	tsgenerator.ProcessEnumFiles(paths, conf.EnumExt)
	tsgenerator.ListInterfaceFiles(paths, conf.InterfaceExt)
}

func ResetTypeFile() {
	_ = os.Remove("types.d.ts")
}

func readApi() ([]os.DirEntry, []string, error) {
	conf := config.GetConfig()

	var directories = []string{conf.ApiDir, conf.MobileAppDir}

	_, paths, err := helpers.ReadDir(directories)
	if err != nil {
		return nil, nil, err
	}

	return nil, paths, nil
}
