package configs

import "os"

func GetJwtSecret() string {
	return os.Getenv("SECRET")
}
