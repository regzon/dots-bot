# Dots Game Bot

This is an implementation of a minimax algorithm that
serves as a bot for the Dots game.

## How to run

Firsty, build the docker container:

```sh
docker build -t dots .
```

Then, run it:

```sh
# Replace WIDTH, HEIGHT and DEPTH with suitable parameters
# I usually use 8, 8 and 5 respectively
docker run -it --rm dots dots-game WIDTH HEIGHT DEPTH
```

## Comments

Unfortunately, I haven't documented code enough and haven't
replaced some brute-force arguments with more optimized ones.  
Probably, those will come in a future update.

If you'd like to change players, you can check the file 
`internal/game/game.go`
