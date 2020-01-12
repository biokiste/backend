package main

// UsersResponse JSON API Spec Wrapper
type UsersResponse struct {
	Users []User `json:"data"`
}

// User holds properties
type User struct {
	ID            int     `json:"id"`
	Username      string  `json:"username"`
	Email         string  `json:"email"`
	Lastname      string  `json:"lastname"`
	Firstname     string  `json:"firstname"`
	Mobile        string  `json:"mobile"`
	NeedSMS       int     `json:"need_sms"`
	Phone         string  `json:"phone,omitempty"`
	Street        string  `json:"street"`
	ZIP           string  `json:"zip"`
	City          string  `json:"city"`
	DateOfBirth   string  `json:"date_of_birth"`
	DateOfEntry   string  `json:"date_of_entry"`
	DateOfExit    string  `json:"date_of_exit"`
	State         int     `json:"state"`
	Credit        float32 `json:"credit"`
	CreditDate    string  `json:"credit_date"`
	CreditComment string  `json:"credit_comment,omitempty"`
	IBAN          string  `json:"iban,omitempty"`
	BIC           string  `json:"bic,omitempty"`
	SEPA          string  `json:"sepa,omitempty"`
	RememberToken string  `json:"remember_token,omitempty"`
	Additionals   string  `json:"additionals,omitempty"`
	Comment       string  `json:"comment,omitempty"`
	GroupComment  string  `json:"group_comment,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	LastLogin     string  `json:"last_login"`
	Password      string  `json:"password,omitempty"`
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
