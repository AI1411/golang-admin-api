package util

import (
	"os"

	"github.com/AI1411/golang-admin-api/util/errors"
)

func CheckDir(fileDir string) error {
	_, err := os.Stat(fileDir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(fileDir, 0o777)
		if err != nil {
			return errors.NewInternalServerError("failed to create directory", err)
		}
		return errors.New("ディレクトリを作成しました。")
	}

	if err != nil {
		return err
	}
	return nil
}
