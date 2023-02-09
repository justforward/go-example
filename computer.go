package main

import (
	"fmt"
	"math"
	"os"
)
// sumHead :{62404 62403 16 7233} 一共分为62404个分块，
func main() {
	// 根据文件的大小得到分块信息
	filePath := "/Users/mac/Downloads/Q3.sim"
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	stat, err := f.Stat()
	if err != nil {
		return
	}

	sumHead := SumSizesSqroot(stat.Size())

	fmt.Printf("sumHead :%v", sumHead)
}

type SumHead struct {
	// “number of blocks” (openrsync)
	// “how many chunks” (rsync) 分块大小
	ChecksumCount int32

	// “block length in the file” (openrsync)
	// maximum (1 << 29) for older rsync, (1 << 17) for newer  分块长度
	BlockLength int32

	// “long checksum length” (openrsync) 计算出的rollinghash长度
	ChecksumLength int32

	// “terminal (remainder) block length” (openrsync)
	// RemainderLength is flength % BlockLength 文件分块之后 剩余的长度
	RemainderLength int32

	//Sums []SumBuf
}

const blockSize = 700 // rsync/rsync.h

// Corresponds to rsync/generator.c:sum_sizes_sqroot
func SumSizesSqroot(contentLen int64) SumHead {
	// * The block size is a rounded square root of file length.

	// 	The block size algorithm plays a crucial role in the protocol efficiency. In general, the block size is the rounded square root of the total file size. The minimum block size, however, is 700 B. Otherwise, the square root computation is simply sqrt(3) followed by ceil(3)

	// For reasons unknown, the square root result is rounded up to the nearest multiple of eight.

	// TODO: round this
	blockLength := int32(math.Sqrt(float64(contentLen)))
	if blockLength < blockSize {
		blockLength = blockSize
	}

	// * The checksum size is determined according to:
	// *     blocksum_bits = BLOCKSUM_EXP + 2*log2(file_len) - log2(block_len)
	// * provided by Donovan Baarda which gives a probability of rsync
	// * algorithm corrupting data and falling back using the whole md4
	// * checksums.
	const checksumLength = 16 // TODO?

	return SumHead{
		ChecksumCount:   int32((contentLen + (int64(blockLength) - 1)) / int64(blockLength)),
		RemainderLength: int32(contentLen % int64(blockLength)),
		BlockLength:     blockLength,
		ChecksumLength:  checksumLength,
	}
}
