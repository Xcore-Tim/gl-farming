package mongoParams

import (
	"fmt"
	"strings"
)

const (
	TestAddress      string = "mongodb://localhost:27017"
	AzureProdAddress string = "mongodb://farming-mongodb:7wnnOjnZgpq4Ruprtqq5qXxsS7ZfCF8LxhHJYIZgzenmAJc3l1ZrFEsT5AuCYXjWtvGAJ6Fdfj0lACDbWPXUiw==@farming-mongodb.mongo.cosmos.azure.com:10255/?ssl=true&retrywrites=false&maxIdleTimeMS=120000&appName=@farming-mongodb@"
	// ProdAddress      string = "mongodb://62.84.112.127:27017"
	ProdAddress string = "mongodb://localhost:27017"
)

func GetConnectionString() string {
	DB_NAME := "gypsyland"
	DB_HOSTS := []string{
		"rc1b-1bdh1alctlvnb8k9.mdb.yandexcloud.net:27018",
	}
	DB_USER := "gipsy-dev"
	DB_PASS := "gipsydev"

	// CACERT := "/etc/ssl/certs/ya-root.crt"
	CACERT := "ya-root.crt"

	url := fmt.Sprintf("mongodb://%s:%s@%s/%s?tls=true&tlsCaFile=%s",
		DB_USER,
		DB_PASS,
		strings.Join(DB_HOSTS, ","),
		DB_NAME,
		CACERT)

	return url
}
