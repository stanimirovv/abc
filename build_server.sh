cd server
echo "Running the automated tests..."
/usr/local/go/bin/go test -coverprofile=c.out
#echo "Generating HTML report..."
#/usr/local/go/bin/go tool cover -html=c.out
echo "Removing the test coverprofile file..."
rm c.out
echo "Building the executable..."
/usr/local/go/bin/go build
echo "Done!"
