package stream

// User represents a user.
type User struct {
	// Name is a name of a user.
	Name string
	// Email is an email of a user.
	Email string
}

// String implements fmt.Stringer interface.
func (u User) String() string {
	return u.Name + " <" + u.Email + ">"
}

// Entity is an implementation of Aggregator interface.
type Entity struct {
	deps []Source
}
