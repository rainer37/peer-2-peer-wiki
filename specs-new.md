# P2P Wikipedia Specifications

People who install the app will be able to create, view, modify and delete Wikipedia-style
articles. The articles will be stored on other computers (peers) where the app is
installed but each user will be able to access the full collection of articles by
performing lookups.

To achieve this, the peers will be structured using the P-Grid DHT.
To allow new peers to join the network, a traditional web server will be used to
provide IP addresses of (some) existing peers.

Communication among the peers will be performed using TCP for P-Grid operations and
RPC for article operations. TCP will also be used by nodes to communicate with the
web server.
[NOTE]: # (We may want to use message passing instead of RPCs since that is better supported by P-Grid)

## P-Grid Summary
- Peers use a common binary trie structure to organize their routing tables.

### System Parameters
- The maximum path length parameter `maxl` is a _global_ parameter that prevents
  over-specialization of peers. We will estimate it locally by setting a threshold
  on the number of articles that must be stored under a given prefix. Partitioning
  of the key space will not happen until the number of articles with a certain
  prefix is greater than the threshold.

### Node Operations
#### Start

#### Stop

#### Join
A new node joins the P-Grid by connecting to an existing peer (whose address
was obtained from the web service) and together they decide how to split the search
space that the existing peer is responsible for. This is done recursively by the
`exchange` function which will, based on the system parameters, situate the joining
node in the P-Grid.

#### Leave

#### Lookup
Get all the articles (including replicas) with the specified title (key search)
or has a title in the specified range (range/substring search).

For key searches, a breadth-first search is shown to have the highest success rate
for finding the specified article for a fixed number of messages.

For range searches, there are two algorithms:
- `minmax` where the range search is executed sequentially, and
- `shower` where the range search is executed in parallel
[TODO]: # (Provide the latency bounds)
[NOTE]: # (The shower algorithm will return results as soon as they are found)



## Articles
### Properties
```
title: string
content: string

// indicates whether the article should be displayed
active: boolean

// logical clock value?
timestamp: int
```

### Operations
#### Insert
Create a new article in the P-Grid if no existing articles already have the title
specified.

The peer that receives the command will forward it to all peers with the same path
(i.e. its replicas) using an epidemic algorithm.

#### Update
Modify the content of an existing article in the P-Grid, setting the timestamp to
the latest logical time.

The peer that receives the command will forward it to all peers with the same path
(i.e. its replicas) using an epidemic algorithm.

#### Delete
Sets an article's active field to false and updates the timestamp to the latest
logical time.

The peer that receives the command will forward it to all peers with the same path
(i.e. its replicas) using an epidemic algorithm.

#### Search
Supports direct key search, substring search and range queries on article titles.

This is calls P-Grid's `lookup(key)` function which returns all matching articles
from all alive peers in the range (this will include article replicas). Only the
latest distinct articles will be returned.
