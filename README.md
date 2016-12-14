# P2P Wikipedia

## Installation
To install p2pwiki, clone the repo into your golang working directory and run
```
> go install github.com/nickbradley/p2pwiki
```

## Usage
### Running the chord server
To start the first node in the chord ring:
```
> p2pwiki 127.0.0.1:2222 server start create
```

To join an existing chord ring:
```
> p2pwiki 127.0.0.1:3333 server start join 127.0.0.1:2222
```

### Interacting with an article
All article commands have the form
```
p2pwiki <ip:port client> article <ip:port local server> <op> <title> <opargs>
```


Before viewing/editing an article, it must first be pulled from a peer and stored in the local cache:
```
> p2pwiki 127.0.0.1:3334 article 127.0.0.1:3333 pull A1

Article A1 has been pulled successfully from 127.0.0.1:3333.
```

To view a local article:
```
> p2pwiki 127.0.0.1:3334 article 127.0.0.1:3333 view A1

A1
----
B
```

Local articles can be edited using:
```
> p2pwiki 127.0.0.1:3334 article 127.0.0.1:3333 insert A1 2 "second sentence"
Insert into A1 succeed...

> p2pwiki 127.0.0.1:3334 article 127.0.0.1:3333 view A1

A1
---
B
second sentence

> p2pwiki 127.0.0.1:3334 article 127.0.0.1:3333 delete A1 1

A1
---
second sentence
```

To share your changes you must push your edits to the peer server:
```
> p2pwiki 127.0.0.1:3334 article 127.0.0.1:3333 push A1

Article A1 has sent to 127.0.0.1:3333.
```
Notice that the CRDT will automagically merge changes from multiple pushes.

You can also discard edits by running:
```
> p2pwiki 127.0.0.1:3334 article 127.0.0.1:3333 discard A1
```

## Project structure
```
parent_dir
 |-- bin
 |-- pkg
 |-- src
     |-- github.com
         |-- nickbradley
             |-- p2pwiki
                 |-- article
                     |-- article.go
                     |-- treedoc.go
                 |-- chord
                     |-- chord.go
                 |-- p2pwiki.go
```
