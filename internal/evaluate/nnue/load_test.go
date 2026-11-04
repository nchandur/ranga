package nnue

import (
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

const totalFeatureWeights = InputSize * HiddenSize

func mockFeatureWeight(flat int) int16 {
	switch flat {
	case 0:
		return 111
	case totalFeatureWeights - 1:
		return 222
	}
	return int16(flat % 100)
}

// deterministic expected values for the bucketed output layer
func mockOutputWeight(bucket, i int) int16 {
	return int16(bucket*1000 + i%1000)
}

func mockOutputBias(bucket int) int16 {
	return int16(-(bucket + 1) * 7)
}

// helper to build a deterministic binary buffer of expected network size
func generateMockNetworkBytes(t *testing.T) ([]byte, int16, int16, int16) {
	t.Helper()
	buf := make([]byte, expectedNetworkBytes)
	pos := 0

	writeI16 := func(val int16) {
		binary.LittleEndian.PutUint16(buf[pos:pos+2], uint16(val))
		pos += 2
	}

	for i := range totalFeatureWeights {
		writeI16(mockFeatureWeight(i))
	}

	firstBias := int16(-350)
	writeI16(firstBias)
	for i := 1; i < HiddenSize; i++ {
		writeI16(int16(i))
	}

	for i := range 2 * HiddenSize {
		for b := range OutputBuckets {
			writeI16(mockOutputWeight(b, i))
		}
	}
	// one output bias per bucket
	for b := range OutputBuckets {
		writeI16(mockOutputBias(b))
	}

	if pos != expectedNetworkBytes {
		t.Fatalf("mock network wrote %d bytes; want %d", pos, expectedNetworkBytes)
	}

	return buf, mockFeatureWeight(0), mockFeatureWeight(totalFeatureWeights - 1), firstBias
}

// checks representative feature weights of every input bucket, including its boundaries
func checkFeatureLayer(t *testing.T, net *Network) {
	t.Helper()
	for b := range InputBuckets {
		for _, f := range []int{0, 1, 383, 384, FeaturesPerBucket - 1} {
			for _, v := range []int{0, 1, HiddenSize - 1} {
				idx := b*FeaturesPerBucket + f
				want := mockFeatureWeight(idx*HiddenSize + v)
				if got := net.FeatureWeights[idx].Values[v]; got != want {
					t.Errorf("FeatureWeights[%d][%d] (bucket %d, feature %d) = %d; want %d", idx, v, b, f, got, want)
				}
			}
		}
	}
}

// checks every bucket's bias and a few representative weights
func checkOutputLayer(t *testing.T, net *Network) {
	t.Helper()
	for b := range OutputBuckets {
		// first weight, first weight of the "them" half, and last weight of the bucket
		for _, i := range []int{0, HiddenSize, 2*HiddenSize - 1} {
			if got, want := net.OutputWeights[b][i], mockOutputWeight(b, i); got != want {
				t.Errorf("OutputWeights[%d][%d] = %d; want %d", b, i, got, want)
			}
		}
		if got, want := net.OutputBias[b], mockOutputBias(b); got != want {
			t.Errorf("OutputBias[%d] = %d; want %d", b, got, want)
		}
	}
}

func TestLoadNetwork(t *testing.T) {
	t.Run("expected size accounts for every input bucket", func(t *testing.T) {
		want := (InputBuckets*768*HiddenSize + HiddenSize + OutputBuckets*2*HiddenSize + OutputBuckets) * 2
		if expectedNetworkBytes != want {
			t.Errorf("expectedNetworkBytes = %d; want %d", expectedNetworkBytes, want)
		}
		if InputSize != InputBuckets*768 {
			t.Errorf("InputSize = %d; want %d", InputSize, InputBuckets*768)
		}
	})

	t.Run("parseNetwork rejects truncated byte slices", func(t *testing.T) {
		shortData := make([]byte, expectedNetworkBytes-1)
		net, err := parseNetwork(shortData)

		if err == nil {
			t.Fatalf("expected error for data smaller than %d bytes, got nil", expectedNetworkBytes)
		}
		if net != nil {
			t.Errorf("expected nil network on failure, got %v", net)
		}
	})

	t.Run("parseNetwork rejects a net with a different input bucket count", func(t *testing.T) {
		oversized := make([]byte, expectedNetworkBytes+768*HiddenSize*2)
		net, err := parseNetwork(oversized)

		if err == nil {
			t.Fatalf("expected error for oversized data (bucket count mismatch), got nil")
		}
		if net != nil {
			t.Errorf("expected nil network on failure, got %v", net)
		}
	})

	t.Run("parseNetwork deserializes weights and biases correctly", func(t *testing.T) {
		data, wantFirstWeight, wantLastWeight, wantFirstBias := generateMockNetworkBytes(t)

		net, err := parseNetwork(data)
		if err != nil {
			t.Fatalf("parseNetwork returned unexpected error: %v", err)
		}

		if got := net.FeatureWeights[0].Values[0]; got != wantFirstWeight {
			t.Errorf("FeatureWeights[0][0] = %d; want %d", got, wantFirstWeight)
		}
		if got := net.FeatureWeights[InputSize-1].Values[HiddenSize-1]; got != wantLastWeight {
			t.Errorf("FeatureWeights[%d][last] = %d; want %d", InputSize-1, got, wantLastWeight)
		}
		if got := net.FeatureBias.Values[0]; got != wantFirstBias {
			t.Errorf("FeatureBias[0] = %d; want %d", got, wantFirstBias)
		}

		checkFeatureLayer(t, net)
		checkOutputLayer(t, net)
	})

	t.Run("parseNetwork ignores extra trailing bytes", func(t *testing.T) {
		data, _, _, _ := generateMockNetworkBytes(t)
		paddedData := append(data, []byte{0xFF, 0xEE, 0xDD}...)

		net, err := parseNetwork(paddedData)
		if err != nil {
			t.Fatalf("expected parseNetwork to accept extra trailing bytes, got: %v", err)
		}
		if net == nil {
			t.Fatalf("expected non-nil network")
		}
		checkOutputLayer(t, net)
	})

	t.Run("parseNetwork accepts padding up to 63 bytes but not 64", func(t *testing.T) {
		data, _, _, _ := generateMockNetworkBytes(t)

		if _, err := parseNetwork(append(append([]byte{}, data...), make([]byte, 63)...)); err != nil {
			t.Errorf("63 bytes of padding should be accepted, got: %v", err)
		}
		if _, err := parseNetwork(append(append([]byte{}, data...), make([]byte, 64)...)); err == nil {
			t.Errorf("64 bytes of extra data should be rejected")
		}
	})

	t.Run("parseNetwork accepts a 64-byte aligned (bullet-padded) file", func(t *testing.T) {
		data, _, _, _ := generateMockNetworkBytes(t)
		pad := (64 - len(data)%64) % 64
		net, err := parseNetwork(append(data, make([]byte, pad)...))
		if err != nil {
			t.Fatalf("expected padded network to load, got: %v", err)
		}
		checkFeatureLayer(t, net)
		checkOutputLayer(t, net)
	})

	t.Run("LoadNetwork returns error when file does not exist", func(t *testing.T) {
		_, err := LoadNetwork("non_existent_network_file.nnue")
		if err == nil {
			t.Fatalf("expected error loading non-existent file, got nil")
		}
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("expected wrapped os.ErrNotExist, got %v", err)
		}
	})

	t.Run("LoadNetwork reads and unpacks valid network file from disk", func(t *testing.T) {
		data, _, _, _ := generateMockNetworkBytes(t)

		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test_net.nnue")
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}

		net, err := LoadNetwork(filePath)
		if err != nil {
			t.Fatalf("LoadNetwork failed: %v", err)
		}

		checkFeatureLayer(t, net)
		checkOutputLayer(t, net)
	})
}
