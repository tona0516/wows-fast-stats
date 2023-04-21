package repo

import (
	"encoding/base64"
	"os"
	"path/filepath"
)

type Screenshot struct {}

func (s *Screenshot) Save(path string, base64Data string) error {
    dir := filepath.Dir(path)
    os.Mkdir(dir, 0755)

    data, err := base64.StdEncoding.DecodeString(base64Data)
    if err != nil {
        return err
    }

    f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)

	return err
}

