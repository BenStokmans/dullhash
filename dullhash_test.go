package dullhash

import (
	"encoding/hex"
	"github.com/dgryski/go-onlinestats"
	"gonum.org/v1/gonum/stat"
	"math/big"
	"math/rand"
	"testing"
	"time"
)

const correlationBatchSize = 5000000
var inputs []float64 = nil
var outputs []float64 = nil

func TestSumAdjacentCollisions(t *testing.T) {
	colls := 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 16; i++ {
		data1, data2 := make([]byte, 256 + i), make([]byte, 256 + i)
		if _, err := rand.Read(data1); err != nil {
			t.Fatalf("error while reading random bytes: %v\n", err)
		}
		for j := 0; j < len(data1); j++ {
			copy(data2, data1)
			data2[j]++
			sum1, sum2 := Sum(data1), Sum(data2)
			if hex.EncodeToString(sum1[:]) == hex.EncodeToString(sum2[:]) {
				t.Errorf(
					"hash value is the same for data1 and data2 for bytes:\ndata1: %v\ndata2: %v\nat index %v, data1: %v, data2: %v\n",
					data1,
					data2,
					j,
					data1[j],
					data2[j],
				)
				colls++
			}
		}
	}
	if colls > 0 {
		t.Errorf("total of %v output collisions found\n", colls)
	}
}

func TestSumPerformance(t *testing.T) {
	start := time.Now()
	// generate sample in first test which will be used in the next 2 tests
	generateSample(correlationBatchSize)
	t.Logf("hash time: %v, batch size: %d", time.Since(start), correlationBatchSize)
}

func TestSumCorrelationCoefficient(t *testing.T) {
	pearsons := stat.Correlation(inputs, outputs, nil)

	if pearsons > .001 || pearsons < -.001 {
		t.Errorf("pearsons correlation coefficient of %v is too high/low, expected [-0.001, 0.001]\n", pearsons)
	}
	t.Logf("pearsons correlation: %v, batch size: %v\n", pearsons, correlationBatchSize)
}

func TestSumSpearmanRhoCorrelationCoefficient(t *testing.T) {
	spearmanr, p := onlinestats.Spearman(inputs, outputs)

	if spearmanr > .001 || spearmanr < -.001 {
		t.Errorf("spearmanr correlation coefficient of %v is too high/low, expected [-0.001, 0.001]\n", spearmanr)
	}
	t.Logf("spearmanr correlation: %v, associated p-value: %v, batch size: %v\n", spearmanr, p, correlationBatchSize)
}

func generateSample(batchSize int) {
	rand.Seed(time.Now().UnixNano())
	inputs, outputs = make([]float64, batchSize), make([]float64, batchSize)
	for i := 0; i < len(inputs); i++ {
		inputs[i] = float64(rand.Int63())
	}
	for i := 0; i < batchSize; i++ {
		sum := Sum(big.NewInt(int64(i)).Bytes())
		sumbi := big.NewInt(0).SetBytes(sum[:])
		outputs[i] = float64(sumbi.Div(sumbi, big.NewInt(4)).Int64())
	}
}