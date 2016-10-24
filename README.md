# P2P Wikipedia

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
