export GOOS=linux
export GOARCH=amd64
go build -o bootstrap
if [ -f ./function.zip ]; then
    rm ./function.zip
fi
zip function.zip bootstrap