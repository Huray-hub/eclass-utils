package auth

type Credentials struct {
	Username string
	Password string
}

// UsernameEmpty method checks if Username is empty
func (crd Credentials) UsernameEmpty() bool {
	return crd.Username == ""
}

// PasswordEmpty method checks if Password is empty
func (crd Credentials) PasswordEmpty() bool {
	return crd.Password == ""
}

// ClearCredentials method clears both Username & Password values
func (crd *Credentials) ClearCredentials() {
	crd.Username = ""
	crd.Password = ""
}
