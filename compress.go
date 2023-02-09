package main

/*
type SizeableWriter struct {
	Count int
}

func (w *SizeableWriter) Write(data []byte) (int, error) {
	w.Count += len(data)
	return len(data), nil
}

var (
	filePath = []string{
		"/Users/mac/Desktop/imei.20190824.01.simac",
		"/Users/mac/Downloads/STARCCM_test_Blade@00300.sim",
		"/Users/mac/Desktop/e_10_no_solid_steady_mesh_trim@07500.sim",
	}

	GZipLevel = []int{
		1, // 目前程序中使用的压缩算法
		//-2, 1, //2, 3, 4, 5, 6, 7, 8, 9,

	}
)

func GZip(r io.Reader, lvl int) (int, error) {
	var w SizeableWriter

	gw, err := gzip.NewWriterLevel(&w, gzip.BestSpeed)
	if err != nil {
		return 0, err
	}

	if _, err = io.Copy(gw, r); err != nil {
		return 0, err
	}

	return w.Count, nil
}

func ZStd(r io.Reader) (int, error) {
	var w SizeableWriter

	e, err := zstd.NewWriter(&w)
	if err != nil {
		return 0, err
	}

	if _, err = io.Copy(e, r); err != nil {
		return 0, err
	}

	return w.Count, nil
}

func WithReset(r io.ReadSeeker, f func(io.Reader) error) error {
	_, err := r.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	return f(r)
}

func TimeSpend(r io.Reader) (time.Duration, error) {
	start := time.Now()

	_, err := io.Copy(io.Discard, r)
	if err != nil {
		return 0, err
	}

	return time.Now().Sub(start), nil
}

func main() {
	fp, err := os.Open(filePath[0])
	if err != nil {
		panic(err)
	}
	defer func() { _ = fp.Close() }()

	//ioSpend, err := TimeSpend(fp)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("io spend: %s\n", ioSpend)

	fi, err := fp.Stat()
	if err != nil {
		panic(err)
	}
	fileSize := float64(fi.Size())

	err = WithReset(fp, func(r io.Reader) error {
		before := time.Now()

		c, err := ZStd(r)
		if err != nil {
			return err
		}

		fmt.Printf("zstd, before = %d, after = %d, ratio = %.3f, spend = %s\n",
			fi.Size(), c, float64(c)/fileSize, time.Now().Sub(before))
		return nil
	})

	for _, lvl := range GZipLevel {
		err = WithReset(fp, func(r io.Reader) error {
			before := time.Now()

			c, err := GZip(r, lvl)
			if err != nil {
				return err
			}

			fmt.Printf("gzip with level %d, before = %d, after = %d, ratio = %.3f, spend = %s\n",
				lvl, fi.Size(), c, float64(c)/fileSize, time.Now().Sub(before))
			return nil
		})
	}

}*/
