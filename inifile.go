package goini

import "os"
import "io"
import "bufio"
import "fmt"
import "strings"
import "strconv"
import "sync"

type IniFile struct {
	fileHandle *os.File
	aryLines   []string
	fileName   string
	lock       sync.Mutex
}

func (f *IniFile) Reload() bool {

	f.lock.Lock()

	defer f.lock.Unlock()

	var err error
	f.fileHandle, err = os.Open(f.fileName)
	if err != nil {
		return false
	}

	defer f.fileHandle.Close()

	reader := bufio.NewReader(f.fileHandle)

	if reader == nil {
		return false
	}

	for {
		text, _, err := reader.ReadLine()
		if err != nil || io.EOF == err {
			break
		}
		f.aryLines = append(f.aryLines, string(text))
	}

	return true
}

func (f *IniFile) Init(fileName string) bool {
	var err error
	f.fileHandle, err = os.Open(fileName)
	if err != nil {
		return false
	}

	defer f.fileHandle.Close()

	reader := bufio.NewReader(f.fileHandle)

	if reader == nil {
		return false
	}

	for {
		text, _, err := reader.ReadLine()
		if err != nil || io.EOF == err {
			break
		}
		f.aryLines = append(f.aryLines, string(text))
	}

	return true
}

func (f *IniFile) ReadString(section string, key string, defaultVal string) string {

	f.lock.Lock()
	defer f.lock.Unlock()

	section = "[" + section + "]"
	var bFindSec bool = false
	var strValue string = defaultVal
	for _, value := range f.aryLines {
		if value != "" && value[0] == '#' {
			continue
		}

		if value == section {
			bFindSec = true
			continue
		}

		if bFindSec {
			if value != "" && value[0] == '[' {
				break
			}

			index := strings.Index(value, "=")
			if index == -1 {
				continue
			} else {
				item := strings.Split(value, "=")
				if len(item) < 2 {
					continue
				} else {
					if item[0] == key {
						return item[1]
					} else {
						continue
					}
				}
			}
		}
	}

	return strValue
}

func (f *IniFile) ReadInt(section string, key string, defaultVal int) int {

	f.lock.Lock()
	defer f.lock.Unlock()

	section = "[" + section + "]"
	var bFindSec bool = false
	var intValue int = defaultVal
	for _, value := range f.aryLines {
		if value != "" && value[0] == '#' {
			continue
		}

		if value == section {
			bFindSec = true
			continue
		}

		if bFindSec {
			if value != "" && value[0] == '[' {
				break
			}

			index := strings.Index(value, "=")
			if index == -1 {
				continue
			} else {
				item := strings.Split(value, "=")
				if len(item) < 2 {
					continue
				} else {
					if item[0] == key {
						convertVal, err := strconv.Atoi(item[1])
						if err != nil {
							return intValue
						} else {
							return convertVal
						}
					} else {
						continue
					}
				}
			}
		}
	}

	return intValue
}
