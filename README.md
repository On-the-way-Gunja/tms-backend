# Prototype
- Test server: http://gcp.kimdictor.kr
- API document in test server : http://gcp.kimdictor.kr/docs/index.html

## How to build
```bash
go get -u github.com/swaggo/swag/cmd/swag
git clone https://github.com/On-the-way-Gunja/Prototype.git
cd ./Prototype
go env #copy your $GOPATH
<$GOPATH>/bin/swag init go build
nano keys.txt #insert access key through json format
sudo ./proto*

```
