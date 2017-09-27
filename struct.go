package stream

type User struct {
	Name  string
	Email string
}

func (u User) String() string {
	return u.Name + " <" + u.Email + ">"
}
