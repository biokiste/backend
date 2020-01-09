package main

// UpdateTransaction updates user payment
func (h Handlers) UpdateTransaction(transaction TransactionRequest) error {
	_, err := h.DB.Exec(
		`UPDATE transactions
	   SET status = ?, validated_by = ?, reason = ? 		 
		 WHERE id = ?`,
		transaction.Transactions[0].Status,
		transaction.Transactions[0].ValidatedBy,
		transaction.Transactions[0].Reason,
		transaction.Transactions[0].ID,
	)
	return err
}
