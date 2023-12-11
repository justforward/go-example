package fs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type NetworkFs struct {
	endPoint string
	client   *http.Client
}

func NewNetworkFs(endPoint string) *NetworkFs {

	return &NetworkFs{
		endPoint: endPoint,
		client:   http.DefaultClient,
	}
}
func (n *NetworkFs) Ls(path string, pageOffset, pageSize int64) ([]os.FileInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (n *NetworkFs) Stat(path string) (os.FileInfo, error) {
	// 构建请求
	query := url.Values{}

	//请求参数
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s://%s%s?%s", b.scheme, b.host, "/api/filemanager/stat", query.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("boxfs:stat build req err:%w", err)
	}

	resp, err := n.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("boxfs:stat do err:%w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("boxfs:stat do code:%d", resp.StatusCode)
	}

	var info *Response

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("boxfs:stat io read err:%w", err)
	}
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("boxfs:stat io Unmarshal err:%w", err)
	}

	return &fs_wrap.FileInfos{FileName: info.Data.Name,
		FileSize:    info.Data.Size,
		FileModTime: time.Unix(info.Data.ModTime, 0),
		FileIsDir:   info.Data.IsDir,
	}, nil

}

func (n *NetworkFs) Download(path string, beginOffset int64, endOffset int64) (io.ReadCloser, error) {
	//TODO implement me
	panic("implement me")
}
