package util

import (
	"github.com/AI1411/golang-admin-api/util/errors"
	"os"
)

func CheckDir(fileDir string) error {
	_, err := os.Stat(fileDir)
	if os.IsNotExist(err) {
		os.MkdirAll(fileDir, 0777)
		return errors.New("ディレクトリを作成しました。")
	}

	if err != nil {
		return err
	}
	return nil
}
