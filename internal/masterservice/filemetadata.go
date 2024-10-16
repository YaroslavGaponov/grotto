package masterservice

import (
	"strconv"
	"strings"
)

type FileMetadata struct {
	chunks map[int][]string
}

func NewFileMetadata() FileMetadata {
	return FileMetadata{
		chunks: make(map[int][]string),
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
		if _, found := md.chunks[id]; found {
			md.chunks[id] = append(md.chunks[id], idurl[1])
		} else {
			md.chunks[id] = []string{idurl[1]}
		}
	}
}

func (md *FileMetadata) AddChunk(id int, url string) {
	if _, found := md.chunks[id]; found {
		md.chunks[id] = append(md.chunks[id], url)
	} else {
		md.chunks[id] = []string{url}
	}
}

func (fileMetadata *FileMetadata) ToByteArray() []byte {
	var sb strings.Builder
	for id, urls := range fileMetadata.chunks {
		for _, url := range urls {
			sb.WriteString(strconv.Itoa(id))
			sb.WriteString("|")
			sb.WriteString(url)
			sb.WriteString("\n")
		}
	}
	return []byte(sb.String())
}
