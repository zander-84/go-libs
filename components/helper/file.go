package helper

import (
	"os"
	"strings"
)

type File struct {
	time *Time
}

func NewFile(timeZone string) *File {
	this := new(File)
	this.time = NewTime(timeZone)
	return this
}

func (this *File) GetDayPrefixPath(dir string, filename string) string {
	return strings.TrimRight(dir, "/") + "/" + this.time.FormatDayHyphen() + "-" + filename
}

func (this *File) GetPrefixPath(dir string, filename string) string {
	return strings.TrimRight(dir, "/") + "/" + filename
}

func (this *File) OpenOrCreate(path string, fileName string) (*os.File, error) {
	file, err := os.OpenFile(this.GetPrefixPath(path, fileName), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (this *File) OpenOrCreateWithAction(path string, fileName string, f func(f *os.File)) error {
	file, err := this.OpenOrCreate(path, fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	f(file)
	return nil
}
