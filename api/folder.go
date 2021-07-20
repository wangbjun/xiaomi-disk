package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const (
	TypeFolder = "folder"
	TypeFile   = "file"
)

type File struct {
	Sha1          string      `json:"sha1"`
	ModifyTime    int64       `json:"modifyTime"`
	CreateTime    int64       `json:"createTime"`
	PrivacyStatus string      `json:"privacyStatus"`
	Name          string      `json:"name"`
	Size          int64       `json:"size"`
	Id            string      `json:"id"`
	Type          string      `json:"type"`
	Revision      string      `json:"revision"`
	ParentId      interface{} `json:"parentId"`
	Status        string      `json:"status"`
}

type Msg struct {
	Result    string `json:"result"`
	Retryable bool   `json:"retryable"`
	Code      int    `json:"code"`
	Data      struct {
		HasMore  bool    `json:"hasMore"`
		AllCount int     `json:"allCount"`
		Records  []*File `json:"records"`
	}
}

// GetFoldersById 获取目录下的文件
func (api *Api) GetFoldersById(id string) ([]*File, error) {
	result, err := api.Get(fmt.Sprintf(GetFolders, time.Now().Unix(), id, 1, 100))
	if err != nil {
		return nil, err
	}
	msg := &Msg{}
	err = json.Unmarshal(result, msg)
	if err != nil {
		return nil, err
	}
	if msg.Result == "ok" {
		return msg.Data.Records, nil
	} else {
		return nil, errors.New("获取目录列表失败")
	}
}
