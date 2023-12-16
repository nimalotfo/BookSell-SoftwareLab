module gitlab.com/narm-group/user-service

go 1.18

require (
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.7
	github.com/sirupsen/logrus v1.9.0
	gitlab.com/narm-group/service-api v0.0.0-20230321134530-8c76a0727107
	golang.org/x/crypto v0.7.0
	google.golang.org/grpc v1.53.0
	gorm.io/gorm v1.24.7-0.20230306060331-85eaf9eeda11
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.3.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

require (
	github.com/cenkalti/backoff/v4 v4.2.0
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	gorm.io/driver/postgres v1.5.0
)

replace gitlab.com/narm-group/service-api => ../service-api
