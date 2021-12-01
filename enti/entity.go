package enti

type Registration struct {
	Id        int    `pg:",pk" `
	Firstname string `validate:"nonzero"`
	Lastname  string `validate:"nonzero"`
	Username  string `sql:",unique" validate:"nonzero"`
	Password  string `validate:"nonzero"`
	Email     string `sql:",unique" validate:"nonzero"`
	Otp       string
}
type LoginInfo struct {
	Username string `validate:"nonzero"`
	Password string `validate:"nonzero"`
}
type ResetPassword struct {
	Username string `validate:"nonzero"`
}
