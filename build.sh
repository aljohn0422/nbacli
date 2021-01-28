env GOOS=windows GOARCH=amd64 go build -o nba 
zip windows-amd64.zip nba
env GOOS=darwin GOARCH=amd64 go build -o nba 
zip mac-amd64.zip nba
env GOOS=linux GOARCH=amd64 go build -o nba 
zip linux-amd64.zip nba