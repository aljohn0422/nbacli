# Features
- Auto refresh in play by play mode.
- Show boxscore at timeout (for players on court) and period change (full box).
- Timezone default to system time.

# Installation
## With binary

Download prebuilt binaries from Github Releases. Unzip the binary and run `./nba`

## Build from source
    git clone https://github.com/aljohn0422/nbacli/
    go get github.com/manifoldco/promptui
    go build -o nba .
    ./nba

## Build with Dockerfile
    docker build -t nbacli .
    docker run -it nbacli