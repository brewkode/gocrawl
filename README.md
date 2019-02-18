# Crawler

This is my attempt at building a minimal crawler using some of golang's primitives. The crawler is visualized to contain multiple pieces. 
- Seeding - the starting point 
- Filtering - dropping urls that are crawled already
- Fetching - doing the actual http fetch, ensuring redirects are followed
- Parsing - parsing the html body and extracting links from it to be sent back for further crawl
- Sitemap building - book-keeping of urls & their outlinks

The solution is built using core primitives of golang - viz., goroutines & channels. Each piece mentioned above is connected to other piece via channels. This design helps us with the following
- identifying functionality and separating them from interaction
- extensible to add more features (for ex: the sitemap builder was a late addition to the system)
- testable
- as all blocks interfaces with others via channels, it creates nice decoupling between these pieces as well
- scale

## Block diagram
```
+---------------+
|  Seeder       |
|               |
+--------+------+
         |
         v
+--------+------+         +-------------+           +---------------+          +--------------------+
|               |         |             |           |               |          |                    |
|    Filter     +--------->   Fetch     +---------->+   Parse       +--------->+  Sitemap Builder   |
|               |         |             |           |               |          |                    |
+--------+------+         +-------------+           +--------+------+          +--------------------+
         ^                                                   |
         |                                                   |
         |                                                   |
         |                                                   |
         |                                                   |
         |                                                   |
         |                                                   |
         +---------------------------------------------------+

```

## Scaleability
- since all these blocks are connected via channels, we can scale individual blocks as we find bottlenecks
- communication via channels allows the blocks to be (potentially) even scaled out if need be

## Design Limitations
- State management in the Filter & Sitemap Builder limits scale. Filter is trivial and is not a big deal. If memory footprint becomes a concern, we could move to some probabilistic data structures too
- Sitemap builder manages state in-memory; addressing large volume of urls and their sitemap is possible by way of using a more appropriate data store

## Build & Run
- `make clean test`
- `make setup`
- `make build`
- For usage `./gocrawl --help`
- To crawl with a default website: `./gocrawl`
- To crawl a different website: `./gocrawl -home=<home_page>`
- The program waits(in a loop) to read `url` on stdin for which the incoming/outgoing urls will be printed to screen.

## TODOs
- use a proper dependency manager
- redirection mapping info isn't preserved internally.
- implement tee'd channel so that the link extractor can use that as its output instead of taking two output channel args
