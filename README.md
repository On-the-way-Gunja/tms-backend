# OTW api backend
## How to build
```bash
go get -u github.com/swaggo/swag/cmd/swag
git clone https://github.com/On-the-way-Gunja/tms-backend.git
cd ./tms-backend
go env #copy your $GOPATH
<$GOPATH>/bin/swag init && go build
nano keys.txt #insert access key through json format
sudo ./proto* #execute built binary

```
