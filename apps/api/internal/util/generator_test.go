package util_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/nuryanto2121/laundryos/apps/api/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestGenerateCode(t *testing.T) {
	dateStr := time.Now().Format("060102")

	tests := []struct {
		name         string
		prefix       string
		sequence     int64
		expected     string
		expectedLen  int
	}{
		{
			name:        "sequence 1 with custom prefix",
			prefix:      "TST",
			sequence:    1,
			expected:    fmt.Sprintf("TST%s000001", dateStr),
			expectedLen: 3 + 6 + 6,
		},
		{
			name:        "sequence 42 zero padded",
			prefix:      "TST",
			sequence:    42,
			expected:    fmt.Sprintf("TST%s000042", dateStr),
			expectedLen: 15,
		},
		{
			name:        "sequence 999999 boundary",
			prefix:      "TST",
			sequence:    999999,
			expected:    fmt.Sprintf("TST%s999999", dateStr),
			expectedLen: 15,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := util.GenerateCode(tc.prefix, tc.sequence)
			assert.Equal(t, tc.expected, got)
			assert.Equal(t, tc.expectedLen, len(got))
			assert.True(t, len(got) >= 15)
		})
	}
}

func TestExportedGenerators(t *testing.T) {
	dateStr := time.Now().Format("060102")
	sequence := int64(123)

	t.Run("GenerateCustomerCode", func(t *testing.T) {
		got := util.GenerateCustomerCode(sequence)
		expected := fmt.Sprintf("CST%s000123", dateStr)
		assert.Equal(t, expected, got)
		assert.Equal(t, "CST", got[:3])
	})

	t.Run("GenerateInvoiceNumber", func(t *testing.T) {
		got := util.GenerateInvoiceNumber(sequence)
		expected := fmt.Sprintf("INV%s000123", dateStr)
		assert.Equal(t, expected, got)
		assert.Equal(t, "INV", got[:3])
	})

	t.Run("GenerateOrderNumber", func(t *testing.T) {
		got := util.GenerateOrderNumber(sequence)
		expected := fmt.Sprintf("ORD%s000123", dateStr)
		assert.Equal(t, expected, got)
		assert.Equal(t, "ORD", got[:3])
	})

	t.Run("GenerateServiceCode", func(t *testing.T) {
		got := util.GenerateServiceCode(sequence)
		expected := fmt.Sprintf("SRV%s000123", dateStr)
		assert.Equal(t, expected, got)
		assert.Equal(t, "SRV", got[:3])
	})
}
