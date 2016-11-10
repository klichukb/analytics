server:
	go run wrapper.go models.go api.go db.go server.go

clients:
	go run wrapper.go models.go client.go
