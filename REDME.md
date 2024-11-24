Create migration file
``` bash
migrate create -ext sql -dir migrations/ -seq ${NAME_OF_MIGRATION}
```

Run migration
``` bash
migrate -path migrations/ -database ${DATABASE} -verbose up
```

Rollback migration
``` bash
migrate -path migrations/ -database ${DATABASE} -verbose down
```

Run Application with nodemon
``` bash
nodemon --exec go run cmd/main.go --signal SIGTERM
```