cd server
echo "Running the automated tests..."
$GOROOT/bin/go test -coverprofile=c.out -logtostderr=true
#echo "Generating HTML report..."
#/usr/local/go/bin/go tool cover -html=c.out
echo "Removing the test coverprofile file..."
rm c.out
echo "Building the executable..."
$GOROOT/bin/go build
echo "Done!"
