package main

// UsersResponse JSON API Spec Wrapper
type UsersResponse struct {
	Users []User `json:"data"`
}

// User holds properties
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// TransactionResponse JSON API Spec Wrapper
type TransactionResponse struct {
	Transactions []Transaction `json:"data"`
}

// UserTransactionResponse JSON API Spec Wrapper
type UserTransactionResponse struct {
	UserTransaction `json:"data"`
}

// DoorCodeResponse JSON API Spec Wrapper
type DoorCodeResponse struct {
	DoorCode `json:"data"`
}

// UserTransaction holds last transaction plus actual balance
type UserTransaction struct {
	Balance      `json:"balance"`
	Transactions []Transaction `json:"transactions"`
}

// Balance holds actual user balance
type Balance float32

// DoorCode holds actual door codes
type DoorCode struct {
	Value     string `json:"doorcode"`
	UpdatedAt string `json:"updated_at"`
	UpdatedBy int    `json:"updated_by"`
}

// Transaction holds properties
// swagger:model transaction
type Transaction struct {
	ID          int     `json:"id"`
	Amount      float32 `json:"amount"`
	CreatedAt   string  `json:"created_at"`
	FirstName   string  `json:"firstname"`
	LastName    string  `json:"lastname"`
	Status      int     `json:"status"`
	Reason      string  `json:"reason"`
	CategoryID  int     `json:"category_id"`
	Type        string  `json:"type"`
	ValidatedBy int     `json:"validated_by"`
}

// TransactionCategoryResponse JSON API Spec Wrapper
type TransactionCategoryResponse struct {
	TransactionCategories []TransactionCategory `json:"data"`
}

// TransactionCategory implements categories of a transaction
type TransactionCategory struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// TransactionRequest implements user and requested transactions
type TransactionRequest struct {
	Transactions []Transaction `json:"transactions"`
	User         `json:"user"`
}
