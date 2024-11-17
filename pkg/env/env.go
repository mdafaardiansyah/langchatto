package env

import "github.com/joho/godotenv"

var Env map[string]string

// GetEnv retrieves the value of a specific environment variable, or a default value if the variable does not exist.
func GetEnv(key, def string) string {
	if val, ok := Env[key]; ok {
		return val
	}
	return def
}

// SetupEnvFile loads the .env file and populates the Env map with the key-value pairs in the file.
// If the file does not exist or there is an error reading the file, the function panics.
func SetupEnvFile() {
	envFile := ".env"
	var err error
	Env, err = godotenv.Read(envFile)
	if err != nil {
		panic(err)
	}

}
