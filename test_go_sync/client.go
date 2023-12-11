package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	gosync "github.com/Redundancy/go-sync"
	"github.com/Redundancy/go-sync/blocksources"
	"github.com/Redundancy/go-sync/filechecksum"
	"github.com/Redundancy/go-sync/indexbuilder"
	"io"
	"os"
)

//10.0.6.189
const RsyncIp = "http://localhost:8000/content"
const ReadFileIp = "http://localhost:8001/file/content"

const BLOCK_SIZE = 100

const Local_PATH = "/Users/mac/GolandProjects/go-example/test_go_sync/input.text"

const sink_path = "/Users/mac/GolandProjects/go-example/test_go_sync/output.text"

const LOCAL_VERSION = "The qwik brown fox jumped 0v3r the lazy"

func main() {
	// 机器ip
	LOCAL_URL := fmt.Sprintf("%s", RsyncIp)

	generator := filechecksum.NewFileChecksumGenerator(BLOCK_SIZE)
	stat, err := os.Stat("/Users/mac/Downloads/Fluent_test_Pipe.cas")

	if err != nil {
		return
	}

	// 得到远程的文件.gosync


	_, referenceFileIndex, checksumLookup, err := indexbuilder.BuildIndexFromString(generator, LOCAL_VERSION)

	if err != nil {
		return
	}

	//referened
	fileSize := int64(stat.Size())

	// This would normally be saved in a file
	blockCount := fileSize / BLOCK_SIZE
	if fileSize%BLOCK_SIZE != 0 {
		blockCount++
	}

	fs := &gosync.BasicSummary{
		ChecksumIndex:  referenceFileIndex,
		ChecksumLookup: checksumLookup,
		BlockCount:     uint(blockCount),
		BlockSize:      uint(BLOCK_SIZE),
		FileSize:       fileSize,
	}

	// Need to replace the output and the input
	file, err := os.ReadFile(Local_PATH)
	if err != nil {
		return
	}
	inputFile := bytes.NewReader(file)
	patchedFile := bytes.NewBuffer(nil)

	resolver := blocksources.MakeFileSizedBlockResolver(
		uint64(fs.GetBlockSize()),
		fs.GetFileSize(),
	)

	rsync := &gosync.RSync{
		Input:  inputFile,
		Output: patchedFile,
		Source: blocksources.NewHttpBlockSource(
			LOCAL_URL,
			1,
			resolver,
			&filechecksum.HashVerifier{
				Hash:                md5.New(),
				BlockSize:           fs.GetBlockSize(),
				BlockChecksumGetter: fs,
			},
		),
		Summary: fs,
		OnClose: nil,
	}

	err = rsync.Patch()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	err = rsync.Close()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Patched content: \"%v\"\n", patchedFile.String())

	// Just for inspection
	remoteReferenceSource := rsync.Source.(*blocksources.BlockSourceBase)
	fmt.Printf("Downloaded Bytes: %v\n", remoteReferenceSource.ReadBytes())

}

func readLocal(path string) (io.Reader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	buffer := bytes.NewBuffer(make([]byte, 0))

	_, err = reader.WriteTo(buffer)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}
