# P2P Wikipedia

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
Before viewing/editing an article, it must first be pulled from a peer and stored
in the local cache:
```
> p2pwiki 127.0.0.1:2222 article pull beer

Article beer has been pulled successfully from 127.0.0.1:3333.
```

To view a cached article:
```
> p2pwiki 127.0.0.1:2222 article view beer

Beer
----
Beer is delicious.
There are many types of beer.
```

Cached articles can be edited using:
```
> p2pwiki 127.0.0.1:2222 article insert 3 "Beer has been around for a long time."

Beer
---
Beer is delicious.
There are many types of beer.
Beer has been around for a long time.


> p2pwiki 127.0.0.1:2222 article delete 1

Beer
---
There are many types of beer.
Beer has been around for a long time.
```

To share your changes you must push your edits to the peer server:
```
> p2pwiki 127.0.0.1:2222 article push beer

Article beer has sent to 127.0.0.1:3333.
```
Notice that the CRDT will automagically merge changes from multiple pushes.

## Project structure
You'll need the following directory structure (create it in some parent directory):

```
parent_dir
 |-- bin
 |-- pkg
 |-- src
     |-- github.com
         |-- nickbradley
```

Then `cd` into `src/github.com/nickbradley` and run `git clone https://github.com/nickbradley/p2pwiki.git`. Finally, create a directory
`vendor` in `p2pwiki`. After cloning, the directory should look like:

```
parent_dir
 |-- bin
 |-- pkg
 |-- src
     |-- github.com
         |-- nickbradley
             |-- p2pwiki
                 |-- vendor
                 |-- article.go
                 |-- p2pwiki.go
                 |-- peer.go
```
