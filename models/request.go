package models

import "time"

type File struct {
	Id         int       `json:"Id"`
	Name       string    `json:"name"`
	Version    string    `json:"version"`
	DateCreate time.Time `json:"dateCreate"`
}

type CreateFileReq struct {
	Name    string `form:"name"`
	Version string `form:"version"`
	File    []byte `form:"file"`
}
