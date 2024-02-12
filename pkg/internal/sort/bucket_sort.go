package sort

import (
	"math"
	"sort"
	"sync"
	"sync/atomic"
)

type SortJob[T any] struct {
	BucketIndex     func(T) uint    // The index of the bucket this item belongs to
	Less            func(T, T) bool // The items are sorted according to this less function
	Items           []T             // The items to be sorted
	NumberOfBuckets int             // The number of buckets to be used
}

// Sorts the given items in place using parallel bucket sort
func BucketSort[T any](job SortJob[T], threads int) []T {
	bucketCollection, bucketFill := fillBuckets(job, threads)
	merge(job, bucketFill, bucketCollection, threads)
	return job.Items
}

// Inserts morton pairs into the specified number of buckets
// Each thread uses a separate slice of buckets to avoid the need for synchronized access
// Return:
// buckets: [threads][numberOfBuckets]bucket => one slice of buckets for each thread
// bucketFill: holds how many pairs have been inserted into the corresponding bucket
func fillBuckets[T any](job SortJob[T], threads int) (buckets [][][]T, bucketFill []int32) {
	batchSize := int(math.Ceil(float64(len(job.Items)) / float64(threads)))
	bucketCollection := make([][][]T, 0, threads)
	bucketEntries := make([]int32, job.NumberOfBuckets)
	wg := sync.WaitGroup{}

	// Each thread inserts an equal amount of pairs into its seperate slice of buckets
	for i := 0; i < threads; i++ {
		start := i * batchSize
		if start >= len(job.Items) {
			break
		}

		end := int(math.Min(float64(start+batchSize), float64(len(job.Items))))
		bucketCollection = append(bucketCollection, make([][]T, job.NumberOfBuckets))
		wg.Add(1)
		go func(input []T, threadNumber int) {
			for _, item := range input {
				index := job.BucketIndex(item)
				bucketCollection[threadNumber][index] = append(bucketCollection[threadNumber][index], item)
				atomic.AddInt32(&bucketEntries[index], 1)
			}
			wg.Done()
		}(job.Items[start:end], i)
	}
	wg.Wait()
	return bucketCollection, bucketEntries
}

type mergeJob[T any] struct {
	index   int
	buckets [][]T
	out     []T
}

func merge[T any](sortJob SortJob[T], bucketEntries []int32, bucketCollection [][][]T, threads int) {
	// Start workers, each worker inserts pairs into the given interval of the out slice and sorts it
	jobs := make(chan mergeJob[T], threads)
	wg := sync.WaitGroup{}
	wg.Add(threads)
	for i := 0; i < threads; i++ {
		go func() {
			for job := range jobs {
				mergeBuckets(job.buckets, job.out, sortJob.Less)
			}
			wg.Done()
		}()
	}

	// Feed jobs to workers,
	// Bucket fills are used to determine the corresponding interval in the output slice
	// This method is used to avoid allocating a output slice as this would be quite expensive
	start := 0
	for i := 0; i < sortJob.NumberOfBuckets; i++ {
		end := start + int(bucketEntries[i])
		job := mergeJob[T]{
			index: i,
			out:   sortJob.Items[start:end],
		}
		for _, buck := range bucketCollection {
			job.buckets = append(job.buckets, buck[i])
		}
		jobs <- job
		start += int(bucketEntries[i])
	}
	close(jobs)
	wg.Wait()
}

// Merges n buckets with the same bucket index
func mergeBuckets[T any](buckets [][]T, out []T, compare func(T, T) bool) {
	index := 0
	for _, bucket := range buckets {
		for _, pair := range bucket {
			out[index] = pair
			index++
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return compare(out[i], out[j])
	})
}
