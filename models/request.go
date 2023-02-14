package models

import "time"

type File struct {
	Id         int       `json:"Id"`
	Name       string    `json:"name"`
	Version    float32   `json:"version"`
	DateCreate time.Time `json:"dateCreate"`
}

type CreateFileReq struct {
	Name    string  `form:"name"`
	Version float32 `form:"version"`
	File    []byte  `form:"file"`
}
