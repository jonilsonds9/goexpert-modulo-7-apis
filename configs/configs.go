package configs

var cfg *conf

type conf struct {
	DBDriver      string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	WebServerPort string
	JwtSecret     string
	JwtExpiresIn  int
}

func LoadConfig(path string) {

}
