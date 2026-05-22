package repository1

// In file `dbentity.go` we declare structures for mapping into db tables
// Scan query objects into these structures
// Use them in lower case so they cannot be exported
// 
// type <table_name>DB struct {
//     <some_field> <some_type> `db:"<some_field"`
// }

type userDB struct {
	ID           string    `db:"id"`
	Email        string    `db:"email"`
	Role         string    `db:"role"`
	PasswordHash string    `db:"password_hash"`
	RegisteredAt time.Time `db:"registered_at"`
}