package util

import (
	"fmt"
	"time"
)

// GenerateCode generates a standardized code matching the exact format:
// PREFIX + yyMMdd + 6-digit zero-padded sequence.
// Example: CST260803000001
func GenerateCode(prefix string, sequence int64) string {
	dateStr := time.Now().Format("060102")
	return fmt.Sprintf("%s%s%06d", prefix, dateStr, sequence)
}

// GenerateCustomerCode generates a code with prefix "CST".
func GenerateCustomerCode(sequence int64) string {
	return GenerateCode("CST", sequence)
}

// GenerateInvoiceNumber generates an invoice number with prefix "INV".
func GenerateInvoiceNumber(sequence int64) string {
	return GenerateCode("INV", sequence)
}

// GenerateOrderNumber generates an order number with prefix "ORD".
func GenerateOrderNumber(sequence int64) string {
	return GenerateCode("ORD", sequence)
}

// GenerateServiceCode generates a service code with prefix "SRV".
func GenerateServiceCode(sequence int64) string {
	return GenerateCode("SRV", sequence)
}
