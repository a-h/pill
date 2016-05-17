go test ./httpservice/main/ -coverprofile .main.out
go tool cover -html=.main.out -o coverage_main.html
go test ./dataaccess/ -coverprofile .dataaccess.out
go tool cover -html=.dataaccess.out -o coverage_dataaccess.html
