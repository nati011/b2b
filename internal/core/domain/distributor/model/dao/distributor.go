package auth

/*
should compose user to grow it
forexample; if you want to Create specific class of user

// AdminInfo struct to hold additional information for admin users

	type AdminInfo struct {
		Permissions []string
		AccessLevel int
	}

// CustomerInfo struct for customer-specific details

	type CustomerInfo struct {
		CustomerID string
		Rewards     int
	}

// SpecificTypeOfUser struct that contains additional user information

	type SpecificTypeOfUser struct {
		*User
		Admin  *AdminInfo
		Customer *CustomerInfo
	}
*/
type User struct {
	Email    string
	FullName string
	Username string
}
