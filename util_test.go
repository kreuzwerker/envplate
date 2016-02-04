package envplate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testBackupFile = "test/filemode.txt.bak"
	testFile       = "test/filemode.txt"
)

func TestBackup(t *testing.T) {

	defer _delete(t, testBackupFile)

	assert := assert.New(t)

	modeBackup, err := filemode(testBackupFile)
	assert.Error(err)

	modeFile, err := filemode(testFile)
	assert.NoError(err)

	assert.Equal("hello world\n", _read(t, testFile))

	assert.NoError(createBackup(testFile))

	modeBackup, err = filemode(testBackupFile)
	assert.NoError(err)

	modeFile, err = filemode(testFile)
	assert.NoError(err)

	assert.Equal(modeBackup.String(), modeFile.String())
	assert.Equal(_read(t, testFile), _read(t, testBackupFile))

}

func TestBackupErrors(t *testing.T) {

	assert := assert.New(t)

	_tmpdir(t, func(dir string) {

		err := createBackup(dir)
		assert.Error(err)

	})

}

func TestFilemode(t *testing.T) {

	assert := assert.New(t)

	mode, err := filemode(testFile)

	assert.NoError(err)
	assert.Equal("-rwxrw-r-x", mode.String())

	mode, err = filemode("test/nofile.txt")
	assert.Error(err)

}
