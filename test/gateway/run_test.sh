# TODO: 
#
# Compile code
cd ../.. && go install
rm -rf ./gen-go
frugal --gen go type_test.frugal && frugal --gen gateway type_test.frugal

# Run server (capture pid to kill later?)
# go run server/main.go &

# Run proxy
# go run proxy/main.go &

# make HTTP requests and verify responses

