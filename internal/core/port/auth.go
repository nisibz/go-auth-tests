package port

type AuthService interface {
	Register(name, email, password string) (string, error)
	Login(email, password string) (string, error)
}
