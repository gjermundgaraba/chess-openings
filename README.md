# Chess Openings

The main purpose of this chess cli (it might not be a CLI forever) is to provide a way to store
chess openings with as many variations as you want as you learn them. Importantly, it will also
allow storing your own comments on every move, so you can remember the ideas behind the moves and openings.

Some useful queries:

```Cypher
MATCH (o:Position)-[r:IN]->(p:Opening)
RETURN o, p, r
```

Delete everything
```Cypher
MATCH (n)
DETACH DELETE n
```