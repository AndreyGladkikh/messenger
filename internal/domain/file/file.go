package file

import "bytes"

type File struct {
	data *bytes.Buffer
	dataHash []byte
	url string
}
