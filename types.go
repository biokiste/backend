package main

// UsersResponse JSON API Spec Wrapper
type UsersResponse struct {
	Users []User `json:"data"`
}

// UserResponse JSON API Spec Wrapper
type UserResponse struct {
	User `json:"data"`
}

// Auth0Bearer represents token object
type Auth0Bearer struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// Auth0User represents Auth0 User data
type Auth0User struct {
	Connection string `json:"connection"`
	UserID     string `json:"user_id,omitempty"`
	Email      string `json:"email"`
	Password   string `json:"password,omitempty"`
	LastLogin  string `json:"last_login,omitempty"`
}

// User holds properties
type User struct {
	ID            int     `json:"id"`
	UserID        string  `json:"user_id"`
	Username      string  `json:"username,omitempty"`
	Email         string  `json:"email,omitempty"`
	Lastname      string  `json:"lastname,omitempty"`
	Firstname     string  `json:"firstname,omitempty"`
	Mobile        string  `json:"mobile,omitempty"`
	NeedSMS       int     `json:"need_sms,omitempty"`
	Phone         string  `json:"phone,omitempty"`
	Street        string  `json:"street,omitempty"`
	ZIP           string  `json:"zip,omitempty"`
	City          string  `json:"city,omitempty"`
	DateOfBirth   string  `json:"date_of_birth,omitempty"`
	DateOfEntry   string  `json:"date_of_entry,omitempty"`
	DateOfExit    string  `json:"date_of_exit,omitempty"`
	State         int     `json:"state,omitempty"`
	Credit        float32 `json:"credit,omitempty"`
	CreditDate    string  `json:"credit_date,omitempty"`
	CreditComment string  `json:"credit_comment,omitempty"`
	IBAN          string  `json:"iban,omitempty"`
	BIC           string  `json:"bic,omitempty"`
	SEPA          string  `json:"sepa,omitempty"`
	RememberToken string  `json:"remember_token,omitempty"`
	Additionals   string  `json:"additionals,omitempty"`
	Comment       string  `json:"comment,omitempty"`
	GroupComment  string  `json:"group_comment,omitempty"`
	CreatedAt     string  `json:"created_at,omitempty"`
	UpdatedAt     string  `json:"updated_at,omitempty"`
	LastLogin     string  `json:"last_login,omitempty"`
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

// TransactionStateResponse JSON API Spec Wrapper
type TransactionStateResponse struct {
	TransactionStates []TransactionState `json:"data"`
}

// TransactionState implements states of a transaction
type TransactionState struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

// TransactionRequest implements user and requested transactions
type TransactionRequest struct {
	Transactions []Transaction `json:"transactions"`
	User         `json:"user"`
}

// GroupType represents type of member group
type GroupType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GroupTypesRequest lists all types of groups
type GroupTypesRequest struct {
	GroupTypes []GroupType `json:"data"`
}

// GroupUserEntry implements group with user entry
type GroupUserEntry struct {
	GroupID    int `json:"group_id"`
	UserID     int `json:"user_id"`
	PositionID int `json:"position_id"`
}

// Group implements list of user ids and group leader ids
type Group struct {
	ID        int   `json:"id"`
	UserIDs   []int `json:"user_ids"`
	LeaderIDs []int `json:"leader_ids"`
}

// GroupRequest implements list of groups
type GroupRequest struct {
	Groups []Group `json:"data"`
}
