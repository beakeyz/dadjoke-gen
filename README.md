# Dadjoke-gen

## Eh, why?

to make you suffer with me on my journy to learn golang. make sure to flame this garbage code as much as possible, so I can improve this piece of garbare in the future.

## Okay, but what does this do?

well the plan is to have a system that can handle data in the form of json. it should have full cli intregration and useability. in the future, I might add some kind of reddit post stealer, so I dont have to come up with dumb jokes to add...
since I've been only focused on the actual functionality of the system, it cant really do anything cool yet (because there is no command line madness lmao)
since it does work now, and jokes can be saved, removed, and edited, I should probably start work on the fun part now.
also, because I hate Tensorflow, I am going to TRY to write something along the lines of an AI to determine if a joke is funny enough to be put into the database, but I am not big-brained enough to do that yet. (if YOU are big-brained enough, likeeeee, hmu?)

## Well thats cool, but how do I use it?

Clone the mf and you are basicaly done. There is a makefile to make running and building easier.
The makefile is simple as hell and prob redundant but who really cares right?

To build:
```bash
make build-server
```
To run:
```bash
make run-server
```
to run the client, use the same thing but with the 'client' suffix instead of 'server'
Since the client was built on an older version of the server, its broken at this time. I will either remove it (and continue with the flutter client) or revive it and make it cool 😎
## Some side-notes (just cuz)

I know the format of this project does not cohere to the standard golang specifications, buuuut idrgaf. I will probably fix it sometime later, but for now you are going to have to deal with my way of structure
