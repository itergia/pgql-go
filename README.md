# Property Graph Query Language Parser for Go

This is a skeleton implementation of the Property Graph Query Language (PGQL).

## Background

PGQL is a language to query [graphs](https://en.wikipedia.org/wiki/Graph_theory), returning properties and metadata about it.
See [the PGQL homepage](https://pgql-lang.org/) for more information.
The language is inspired by [SQL](https://en.wikipedia.org/wiki/SQL) and [openCypher](https://opencypher.org/), created by [NEO4J](https://neo4j.com/).
It was created by [Oracle](https://www.oracle.com/) for its graph database.
Oracle maintains an open-source [parser in Java](https://github.com/oracle/pgql-lang).
Both openCypher and PGQL are being considered in [GQL](https://www.gqlstandards.org/existing-languages), an attempt to standardize a graph query language under ISO's SQL.

## Scope and Status

This code is a research project, and you should not expect it to be stable and useful in its current form.
We built this to study PGQL as a contender for the query language in Itergia Core, a real-time graph database.
As such, we don't know if this will ever be more than a parser and AST definition.
Reach out to us if you are interested in using it (and don't want to fork it.)

There is an accompanying [blog post](https://tommie.github.io/a/2022/12/pgql-go).

## Example

A query:

```pgql
SELECT a.number AS a,
       b.number AS b,
       COUNT(e) AS pathLength,
       ARRAY_AGG(e.amount) AS amounts
  FROM MATCH ANY SHORTEST (a:Account) -[e:transaction]->* (b:Account)
 WHERE a.number = 10039 AND b.number = 2090
```

Two modifications:

```pgql
INSERT EDGE e BETWEEN x AND y
UPDATE y SET ( y.a = 12 )
  FROM MATCH (x), MATCH (y)
 WHERE id(x) = 1 AND id(y) = 2
```

## Building

```console
$ go mod download
$ go generate ./...
$ go test ./...
```

## Compliance

This code is expected to be compliant with examples in [PGQL version 1.5](https://pgql-lang.org/spec/1.5/).

There are a few odds and ends:

* The grammar can parse the examples, at the expense of not being compliant with the documented grammar.
  E.g. the keywords `SHORTEST` and `CHEAPEST` are missing in `TOP k` productions.
* Multiple statements can be parsed.
* Statements end with a semicolon.
* The `.*` token is divided into `.` and `*` for symmetry with other property accesses.
* An `IS Label` is allowed where only `':' Label` is allowed in the documented grammar.
* Subqueries cannot contain modification queries, where the documented grammar allows it.

## License

For this entire repository, except as noted in individual directories and files:

Copyright 2022 Itergia AB

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
