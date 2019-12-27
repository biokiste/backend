package main

// GetBalance returns actual user balance
func (h Handlers) GetBalance(id int) (Balance, error) {
	var userBalance Balance
	if err := h.DB.QueryRow(
		`SELECT SUM(amount)
		 FROM transactions
		 WHERE user_id = ?`, id).Scan(&userBalance); err != nil {
		return userBalance, err
	}
	return userBalance, nil
}
