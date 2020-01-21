package main

// GetAllUser returns all users
func (h Handlers) GetAllUser() ([]User, error) {
	var users []User
	results, err := h.DB.Query(`
		SELECT
			id, username, email, lastname, firstname, mobile, need_sms,
			street, zip, city, date_of_birth, date_of_entry,
			COALESCE(date_of_exit, '') as date_of_exit,
			state, credit, credit_date, credit_comment,
			COALESCE(iban, '') as iban,
			COALESCE(bic, '') as bic,
			COALESCE(sepa, '') as sepa,
			COALESCE(additionals, '') as additionals, 
			COALESCE(comment, '') as comment,
			COALESCE(group_comment, '') as group_comment,
			created_at, updated_at,
			COALESCE(last_login, '') as last_login
		FROM users`)
	if err != nil {
		return users, err
	}
	defer results.Close()

	for results.Next() {
		var user User
		err = results.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Lastname,
			&user.Firstname,
			&user.Mobile,
			&user.NeedSMS,
			&user.Street,
			&user.ZIP,
			&user.City,
			&user.DateOfBirth,
			&user.DateOfEntry,
			&user.DateOfExit,
			&user.State,
			&user.Credit,
			&user.CreditDate,
			&user.CreditComment,
			&user.IBAN,
			&user.BIC,
			&user.SEPA,
			&user.Additionals,
			&user.Comment,
			&user.GroupComment,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.LastLogin,
		)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

// UpdateUserData updates user
func (h Handlers) UpdateUserData(user User) error {
	_, err := h.DB.Exec(
		`UPDATE users
	   SET username = ?, email = ?, lastname = ? 		 
		 WHERE id = ?`,
		user.Username,
		user.Email,
		user.Lastname,
		user.ID,
	)
	return err
}

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
