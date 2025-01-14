Create migration file
``` bash
migrate create -ext sql -dir migrations/ -seq ${NAME_OF_MIGRATION}
```

Run migration
``` bash
migrate -path migrations/ -database ${DATABASE} -verbose up

# examle
migrate -path migrations/ -database "mysql://root:123@tcp(127.0.0.1:3306)/go_restaurant" -verbose up
```

Rollback migration
``` bash
migrate -path migrations/ -database ${DATABASE} -verbose down
```

Run Application
``` bash
go run cmd/main.go
```

Run Application with nodemon
``` bash
nodemon --exec go run cmd/main.go --signal SIGTERM
```

Edit dependency injection
``` bash
wire ./pkg/di
```
