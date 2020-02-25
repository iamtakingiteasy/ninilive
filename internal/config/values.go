package config

// Default config values
var Default = Values{
	DB: ValuesDB{
		URL:       "postgresql://ninilive:ninilive@localhost:5432/ninilive?sslmode=disable",
		HardLimit: 1000,
	},
	Upload: ValuesUpload{
		Dir:      "upload",
		MaxBytes: 1024 * 1024 * 100,
		Suffixes: []string{".jpg", ".png", ".gif", ".webm", ".mp4"},
	},
	HTTP: ValuesHTTP{
		Listen: ":8080",
	},
}

// Values root configuration
type Values struct {
	DB       ValuesDB
	Security ValuesSecurity
	Upload   ValuesUpload
	HTTP     ValuesHTTP
}

// ValuesDB database configuration
type ValuesDB struct {
	URL       string
	HardLimit uint64
}

// ValuesSecurity security configuration
type ValuesSecurity struct {
	PasswordSalt string
}

// ValuesUpload upload configuration
type ValuesUpload struct {
	Dir      string
	MaxBytes uint64
	Suffixes []string
}

// ValuesHTTP http server configuration
type ValuesHTTP struct {
	Listen string
}
