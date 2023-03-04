package authentication

type Credentials struct {
	Username string
	Password string
}

func (crd Credentials) UsernameEmpty() bool {
	return crd.Username == ""
}

func (crd Credentials) PasswordEmpty() bool {
	return crd.Password == ""
}
