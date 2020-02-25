package env

import (
	"regexp"

	"github.com/iamtakingiteasy/ninilive/internal/config"
)

type loader struct {
}

func (*loader) Load() (config.Values, error) {
	return config.Values{
		DB: config.ValuesDB{
			URL:       String("DB_URL", config.Default.DB.URL),
			HardLimit: Uint64("DB_HARD_LIMIT", config.Default.DB.HardLimit),
		},
		Security: config.ValuesSecurity{
			PasswordSalt: String("PASSWORD_SALT", config.Default.Security.PasswordSalt),
		},
		Upload: config.ValuesUpload{
			Dir:      String("UPLOAD_DIR", config.Default.Upload.Dir),
			MaxBytes: Uint64("UPLOAD_MAX_BYTES", config.Default.Upload.MaxBytes),
			Suffixes: StringArray("UPLOAD_SUFFIXES", regexp.MustCompile(`\s+`), config.Default.Upload.Suffixes),
		},
		HTTP: config.ValuesHTTP{
			Listen: String("HTTP_LISTEN", config.Default.HTTP.Listen),
		},
	}, nil
}
