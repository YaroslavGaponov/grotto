package masterservice

import (
	"strconv"
	"strings"
)

type FileMetadata struct {
	chunks map[int]string
}

func NewFileMetadata() FileMetadata {
	return FileMetadata{
		chunks: make(map[int]string),
	}
}

func (md *FileMetadata) Load(body []byte) {
	pairs := strings.Split(string(body), "\n")
	for _, pair := range pairs {
		if len(pair) == 0 {
			continue
		}
		idurl := strings.Split(pair, "|")
		id, err := strconv.Atoi(idurl[0])
		if err != nil {
			continue
		}
		md.chunks[id] = idurl[1]
	}
}

func (md *FileMetadata) AddChunk(id int, url string) {
	md.chunks[id] = url
}

func (fileMetadata *FileMetadata) ToByteArray() []byte {
	var sb strings.Builder
	for id, url := range fileMetadata.chunks {
		sb.WriteString(strconv.Itoa(id))
		sb.WriteString("|")
		sb.WriteString(url)
		sb.WriteString("\n")
	}
	return []byte(sb.String())
}
