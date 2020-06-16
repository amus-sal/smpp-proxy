go build
echo "Build is done >>>>>>>>>>>>>>>>>>>"
export SERVER_ADDRESS=127.0.0.1:2010
export PROXY_ADDRESS=127.0.0.1:1234
./smpp_proxy