migrate-up:
	migrate -database="mysql://root:retroced@tcp(localhost:3306)/retroced" -source=file://./db/schema up

migrate-down:
	migrate -database="mysql://root:retroced@tcp(localhost:3306)/retroced" -source=file://./db/schema down

migrate-drop:
	migrate -database="mysql://root:retroced@tcp(localhost:3306)/retroced" -source=file://./db/schema drop -f

build:
	go build -o build/server cmd/server/main.go

run:
	build/server
