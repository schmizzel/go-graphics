package sort_test

import (
	"runtime"
	"testing"

	"github.com/schmizzel/go-graphics/pkg/internal/sort"
	"github.com/stretchr/testify/assert"
)

func TestBucketSort(t *testing.T) {
	job := sort.SortJob[int]{
		Less: func(a, b int) bool { return a < b },
		BucketIndex: func(item int) uint {
			return uint(item / 2)
		},
		NumberOfBuckets: 4,
		Items:           []int{1, 4, 3, 6, 4, 2, 5, 4, 2, 7},
	}

	sort.BucketSort(job, runtime.GOMAXPROCS(0))
	assert.Equal(t, []int{1, 2, 2, 3, 4, 4, 4, 5, 6, 7}, job.Items)
}
