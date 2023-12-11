package main

import (
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
)

type Interval struct {
	Start int
	End   int
}

func main() {
	totalSize := 80 // 总区间大小
	intervalChan := make(chan Interval)
	doneChan := make(chan struct{})

	go processIntervals(totalSize, intervalChan, doneChan)

	// 模拟动态发送子区间
	go sendIntervals(intervalChan)

	// 等待通道处理完毕
	<-doneChan

	fmt.Println("Program finished")
}

func sendIntervals(intervalChan chan<- Interval) {
	// 模拟动态发送子区间
	intervals := []Interval{
		{Start: 1, End: 1},
		{Start: 1, End: 4},
		{Start: 8, End: 8},
	}

	for _, interval := range intervals {
		intervalChan <- interval
	}

	close(intervalChan)
}

func processIntervals(totalSize int, intervalChan <-chan Interval, doneChan chan<- struct{}) {
	missingIntervals := make([]Interval, 0)
	intervalMap := treemap.NewWithIntComparator()

	for interval := range intervalChan {
		intervalMap.Put(interval.Start, interval.End)
	}

	missingIntervals = findMissingIntervals(totalSize, intervalMap)
	//mergedIntervals := mergeIntervals(missingIntervals)

	printIntervals("Missing Intervals:", missingIntervals)
	//printIntervals("Merged Intervals:", mergedIntervals)

	doneChan <- struct{}{}
}

func findMissingIntervals(totalSize int, intervalMap *treemap.Map) []Interval {
	missingIntervals := make([]Interval, 0)

	iterator := intervalMap.Iterator()
	iterator.Begin()

	start := 1
	for iterator.Next() {
		end := iterator.Key().(int) - 1
		if start <= end {
			missingIntervals = append(missingIntervals, Interval{Start: start, End: end})
		}
		start = iterator.Value().(int) + 1
	}

	if start <= totalSize {
		missingIntervals = append(missingIntervals, Interval{Start: start, End: totalSize})
	}

	return missingIntervals
}

func mergeIntervals(intervals []Interval) []Interval {
	if len(intervals) == 0 {
		return intervals
	}

	mergedIntervals := make([]Interval, 0)
	currentInterval := intervals[0]

	for i := 1; i < len(intervals); i++ {
		interval := intervals[i]
		if interval.Start == currentInterval.End+1 {
			currentInterval.End = interval.End
		} else {
			mergedIntervals = append(mergedIntervals, currentInterval)
			currentInterval = interval
		}
	}

	mergedIntervals = append(mergedIntervals, currentInterval)

	return mergedIntervals
}

func printIntervals(title string, intervals []Interval) {
	fmt.Println(title)
	for _, interval := range intervals {
		fmt.Printf("[%d, %d]\n", interval.Start, interval.End)
	}
}
