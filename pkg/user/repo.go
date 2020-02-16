package user

// Repository : User Repository to perform CRUD operations
type Repository interface {
	// BeforeSave(*User) error // TBD later : Not sure if this is needed
	CreateUser(*User) (*User, error)
	UpdateUser(*User) (*User, error)
	DeleteUser(uint64) (int64, error)
	GetUserByID(uint64) (*User, error)
	GetUserByUsername(string) (*User, error)
}
