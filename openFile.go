package main

import "os"

/*
 死循环openfile
*/
func main() {
	i := 0
	for i < 100000 {
		size := int64(1024 * 1024)
		path := "/data/common/4MR3XguF2Ff/file.fem"
		f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
		}
		if err := f.Truncate(size); err != nil {
			f.Close()
		}

		f.Close()
		i++
	}
}
