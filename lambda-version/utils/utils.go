package utils

import (
	"stori-technical-challenge/lambda-version/email"
	"stori-technical-challenge/lambda-version/transactions"
)

// ConvertTransactionSummaryToEmailSummary converts a map of transactions.Summary to a map of email.Summary.
func ConvertTransactionSummaryToEmailSummary(input map[string]transactions.Summary) map[string]email.Summary {
	converted := make(map[string]email.Summary)
	for key, value := range input {
		converted[key] = email.Summary{
			NumTransactions: value.NumTransactions,
			AvgCredit:       value.AvgCredit,
			AvgDebit:        value.AvgDebit,
		}
	}
	return converted
}
