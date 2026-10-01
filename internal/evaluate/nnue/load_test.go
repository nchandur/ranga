package nnue

import (
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// helper to build a deterministic binary buffer of expected network size
func generateMockNetworkBytes(t *testing.T) ([]byte, int16, int16, int16, int16) {
	t.Helper()
	buf := make([]byte, expectedNetworkBytes)
	pos := 0

	writeI16 := func(val int16) {
		binary.LittleEndian.PutUint16(buf[pos:pos+2], uint16(val))
		pos += 2
	}

	firstWeight := int16(111)
	lastWeight := int16(222)
	writeI16(firstWeight)
	for i := 1; i < 768*HiddenSize-1; i++ {
		writeI16(int16(i % 100))
	}
	writeI16(lastWeight)

	firstBias := int16(-350)
	writeI16(firstBias)
	for i := 1; i < HiddenSize; i++ {
		writeI16(int16(i))
	}

	for i := range 2 * HiddenSize {
		writeI16(int16(i + 10))
	}

	outputBias := int16(-42)
	writeI16(outputBias)

	return buf, firstWeight, lastWeight, firstBias, outputBias
}

func TestLoadNetwork(t *testing.T) {
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
	t.Run("parseNetwork deserializes weights and biases correctly", func(t *testing.T) {
		data, wantFirstWeight, wantLastWeight, wantFirstBias, wantOutputBias := generateMockNetworkBytes(t)

		net, err := parseNetwork(data)
		if err != nil {
			t.Fatalf("parseNetwork returned unexpected error: %v", err)
		}

		if got := net.FeatureWeights[0].Values[0]; got != wantFirstWeight {
			t.Errorf("FeatureWeights[0][0] = %d; want %d", got, wantFirstWeight)
		}
		if got := net.FeatureWeights[767].Values[HiddenSize-1]; got != wantLastWeight {
			t.Errorf("FeatureWeights[767][last] = %d; want %d", got, wantLastWeight)
		}

		if got := net.FeatureBias.Values[0]; got != wantFirstBias {
			t.Errorf("FeatureBias[0] = %d; want %d", got, wantFirstBias)
		}

		if got := net.OutputBias; got != wantOutputBias {
			t.Errorf("OutputBias = %d; want %d", got, wantOutputBias)
		}
	})
	t.Run("parseNetwork ignores extra trailing bytes", func(t *testing.T) {
		data, _, _, _, _ := generateMockNetworkBytes(t)
		paddedData := append(data, []byte{0xFF, 0xEE, 0xDD}...)

		net, err := parseNetwork(paddedData)
		if err != nil {
			t.Fatalf("expected parseNetwork to accept extra trailing bytes, got: %v", err)
		}
		if net == nil {
			t.Fatalf("expected non-nil network")
		}
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
		data, _, _, _, wantOutputBias := generateMockNetworkBytes(t)

		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test_net.nnue")
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}

		net, err := LoadNetwork(filePath)
		if err != nil {
			t.Fatalf("LoadNetwork failed: %v", err)
		}

		if net.OutputBias != wantOutputBias {
			t.Errorf("loaded OutputBias = %d; want %d", net.OutputBias, wantOutputBias)
		}
	})
}
