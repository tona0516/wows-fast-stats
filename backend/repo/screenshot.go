package repo

import (
	"encoding/base64"
	"os"
)

type Screenshot struct {
    FileName string}

func (s *Screenshot) Save(base64Data string) error {
    os.Mkdir("screenshot", 0755)

    data, err := base64.StdEncoding.DecodeString(base64Data)
    if err != nil {
        return err
    }

    f, err := os.Create("screenshot/" + s.FileName)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)

	return err
}

