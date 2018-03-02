package db

var (
	username = "sta"
	password = "tweedslug"
	hostname = "localhost"
	dbname   = "sta"
)

func conn() (db *sql.DB, err error) {
	return sql.Open("postgres", fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable", username, password, hostname, dbname))
}
