package service

import "time"

type FileNamer struct{}

func NewFileNamer() *FileNamer {
	return &FileNamer{}
}

func (*FileNamer) SetName() string {
	t := time.Now()
	date := t.Format("02.01.2006-15-04")

	return "rtfconv" + date + ".xlsx"
}
