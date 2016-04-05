cd server
echo "Running the automated tests..."
/usr/local/go/bin/go test -cover
echo "Building the executable..."
/usr/local/go/bin/go build
echo "Done!"
