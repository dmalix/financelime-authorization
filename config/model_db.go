package config

type DB struct {
	Connect struct {
		Host     string
		Port     int
		SslMode  string
		DbName   string
		User     string
		Password string
	}
	Migration struct {
		DropFile   string
		CreateFile string
		InsertFile string
	}
}
