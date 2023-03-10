testcity:
	go test ./features/city... -coverprofile=cover.out && go tool cover -html=cover.out

testreview:
	go test ./features/review... -coverprofile=cover.out && go tool cover -html=cover.out

testdiscussion:
	go test ./features/discussion... -coverprofile=cover.out && go tool cover -html=cover.out

testadditional:
	go test ./features/additional... -coverprofile=cover.out && go tool cover -html=cover.out

testuser:
	go test ./features/user... -coverprofile=cover.out && go tool cover -html=cover.out

testclient:
	go test ./features/client... -coverprofile=cover.out && go tool cover -html=cover.out

testpartner:
	go test ./features/partner... -coverprofile=cover.out && go tool cover -html=cover.out

run:
	go run main.go