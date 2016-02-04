package envplate

import (
	"fmt"
	"io"
	"os"
)

func createBackup(file string) error {

	source, err := os.Open(file)

	if err != nil {
		return err
	}

	defer source.Close()

	target, err := os.Create(fmt.Sprintf("%s.bak", file))

	if err != nil {
		return err
	}

	defer target.Close()

	if _, err := io.Copy(target, source); err != nil {
		return err
	} else if mode, err := filemode(source.Name()); err != nil {
		return err
	} else if err := os.Chmod(target.Name(), mode); err != nil {
		return err
	}

	return nil

}

func filemode(file string) (os.FileMode, error) {

	fileinfo, err := os.Stat(file)

	if err != nil {
		return 0, fmt.Errorf("cannot stat '%s': %v", file, err)
	}

	return fileinfo.Mode(), nil

}
