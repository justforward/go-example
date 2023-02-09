package main

import "fmt"

var (
	filePaths = []string{"/Users/mac/Downloads/STARCCM_test_Blade@00300.sim",
		"/Users/mac/Desktop/e_10_no_solid_steady_mesh_trim@07500.sim",
	}

	GZipLevels = []int{
		-2, //1, //2, 3, 4, 5, 6, 7, 8, 9,
	}
)

type Request struct {
	//path   string `form:"path"`
	Offset int64 `form:"offset"`
	Length int64 `form:"length"`
}

func main() {
	i := 1
	s := []string{"A", "B", "C"}

	// 多重赋值 先计算左边的表达式，然后计算右边的值进行赋值
	i, s[i-1] = 2, "Z"
	fmt.Printf("s: %v \n", s)
}