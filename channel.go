package main

import (
	"fmt"
	"net/url"
)

type DownloadParams struct {
	host        string `json:"host"`
	Path        string `json:"path"`
	PathRewrite string `json:"path_rewrite"`
	Base        string `json:"base"`
	Bucket      string `json:"bucket"`
	ProjectID   string `json:"project_id"`
	Token       string `json:"token"`
}

func main() {
	urlStr := "https://shanhe-test.yuansuan.cn:21312/api/filemanager/single/download?path=%25%5E%26(_%2B.txt&path_rewrite=undefined&base=.&bucket=common&project_id=4ECeF9o27bu&token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiI0eHZSV1ptODJkTCIsImlhdCI6MTY4NTQyOTExMywiZXhwIjoxNjg1NTE1NTEzfQ.weLpQadoR4lnmPE8ysy_KImY4brx-GZwZepMM89W24ZYdmOS4CITTABSLj6Lvaee2kwgn0l5ZCgIYKs6gqtGhGa4_b3qTxAwGfCd2V09wXee9yRd5mNtoOkOlJ662n2GJD7csrjmiW_cHYjOEEX_sS21RD61oloOGI8xSIV7eP4"

	u, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println("URL parsing error:", err)
		return
	}
	host := u.Host

	params := DownloadParams{
		host:        host,
		Path:        u.Query().Get("path"),
		PathRewrite: u.Query().Get("path_rewrite"),
		Base:        u.Query().Get("base"),
		Bucket:      u.Query().Get("bucket"),
		ProjectID:   u.Query().Get("project_id"),
		Token:       u.Query().Get("token"),
	}

	fmt.Printf("Parsed parameters:\n%+v\n", params)
}
