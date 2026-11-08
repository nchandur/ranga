package nnue

const (
	HiddenSize        int   = 64
	Scale             int32 = 400
	QA                int16 = 255
	QB                int16 = 64
	InputBuckets      int   = 1
	OutputBuckets     int   = 8
	FeaturesPerBucket int   = 768
	InputSize         int   = InputBuckets * FeaturesPerBucket
)

var inputBucketLayout = [32]int{
	0, 0, 0, 0,
	0, 0, 0, 0,
	0, 0, 0, 0,
	0, 0, 0, 0,
	0, 0, 0, 0,
	0, 0, 0, 0,
	0, 0, 0, 0,
	0, 0, 0, 0,
}

var bucketTable [2][64]uint8
