# Releaser

Releaser is a tool which help you do version release works.

Features:

- [List pull requests in a milestone(version)](#list-pull-requests-in-a-milestone)
- [List release version in a milestone(version)](#list-release-version-in-a-milestone)
- [Check the module version consistency between repos](#check-the-module-version-consistency-between-repos)

## Before start

Compile releaser.

```bash
git clone https://github.com/you06/releaser.git
cd releaser && go build
```

Edit the config file

```toml
# You can generate your token from https://github.com/settings/tokens
# KEEP IT SECRET!
github-token = ""
# Managed repositories
repos = ["pingcap/tidb", "tikv/tikv", "pingcap/pd"]
# Release note repo
release-note-repo = "pingcap/release-note"
# Release note example, {repo} will be replaced by repo name
# eg. https://github.com/you06/release-note-poc
release-note-path = "/{repo}"
# Default release note language pull request
pull-language = "en"
```

## List pull requests in a milestone

- `-config` specify config file.
- `-version` milestone name

```text
./releaser pr-list -config config.toml -version v3.0.9
+------------+------+--------------+--------------------------------+
|    REPO    |  PR  |    AUTHOR    |             TITLE              |
+------------+------+--------------+--------------------------------+
| tikv/tikv  | 6431 | youjiali1995 | deadlock: more solid role      |
|            |      |              | change observer (#6415)        |
| tikv/tikv  | 6430 | BusyJay      | raftstore: speed up conf       |
|            |      |              | change (#6421)                 |
| tikv/tikv  | 6429 | youjiali1995 | lock_manager: update default   |
|            |      |              | config (#6426)                 |
| tikv/tikv  | 6426 | youjiali1995 | lock_manager: update default   |
|            |      |              | config                         |
| tikv/tikv  | 6422 | youjiali1995 | lock_manager: more metrics     |
|            |      |              | (#6392)                        |
| tikv/tikv  | 6421 | BusyJay      | raftstore: speed up conf       |
|            |      |              | change                         |
| tikv/tikv  | 6392 | youjiali1995 | lock_manager: more metrics     |
| tikv/tikv  | 6330 | youjiali1995 | test: fix unstable             |
|            |      |              | test_waiter_manager_wake_up    |
|            |      |              | (#6304)                        |
| tikv/tikv  | 6260 | yiwu-arbug   | Update rocksdb. Update         |
|            |      |              | rust-rocksdb branch for 3.0    |
|            |      |              | branch                         |
| pingcap/pd | 2083 | sre-bot      | feature(label): tolerate       |
|            |      |              | backslash in the label name    |
|            |      |              | (#1595)                        |
| pingcap/pd | 2067 | sre-bot      | statistic: fix the issue of    |
|            |      |              | tombstone for label counter    |
|            |      |              | (#2060)                        |
+------------+------+--------------+--------------------------------+

No milestone repos: pingcap/tidb, pingcap/br, pingcap/tidb-operator
```

## List release version in a milestone

Arguments:

- `-config` specify config file.
- `-version` milestone name

```text
./releaser release-notes -config ./config.toml -version v4.0.0-beta.1
+--------------+-------+-------------+--------------------------------+--------------+----+----+
|     REPO     |  PR   |   AUTHOR    |             TITLE              | RELEASE NOTE | CN | EN |
+--------------+-------+-------------+--------------------------------+--------------+----+----+
| pingcap/tidb | 14971 | bb7133      | util: move `ColumnsToProto`    | x            | x  | x  |
|              |       |             | from parser to util pkg        |              |    |    |
| pingcap/tidb | 14964 | crazycs520  | infoschema: add more metrics   | x            | x  | x  |
|              |       |             | table for diagnose report.     |              |    |    |
| pingcap/tidb | 14954 | AilinKid    | tables: add sequence binlog    | √            | √  | √  |
|              |       |             | support                        |              |    |    |
| pingcap/tidb | 14878 | crazycs520  | *: make CLUSTER_SLOW_QUERY     | √            | x  | x  |
|              |       |             | support query slow log at any  |              |    |    |
|              |       |             | time                           |              |    |    |
| pingcap/tidb | 14876 | wjhuang2016 | *: support index encode/decode | x            | x  | x  |
|              |       |             | for new collation              |              |    |    |
| pingcap/tidb | 14840 | crazycs520  | executor: make SLOW_QUERY      | √            | √  | x  |
|              |       |             | support query slow log at any  |              |    |    |
|              |       |             | time                           |              |    |    |
| pingcap/tidb | 14839 | lonng       | executor: fix the potential    | √            | x  | x  |
|              |       |             | goroutine leak in cluster log  |              |    |    |
|              |       |             | retriever                      |              |    |    |
+--------------+-------+-------------+--------------------------------+--------------+----+----+
```

## Check the module version consistency between repos

For a complex system, there will usually be many units, and they are in different repos, have different dependencies manager files, like `go.mod`, `Cargo.toml`.

Support dependencies file:

- `go.mod`
- `Cargo.toml`

Arguments:

- `-config` specify config file.
- `-version` can be a tag name or branch name

```text
./releaser check-module -config config.toml -version v3.0.9
pingcap/tidb go.mod: https://github.com/pingcap/tidb/blob/v3.0.9/go.mod
tikv/tikv Cargo.toml: https://github.com/tikv/tikv/blob/v3.0.9/Cargo.toml
pingcap/pd go.mod: https://github.com/pingcap/pd/blob/v3.0.9/go.mod
-----------------------
github.com/coreos/pkg | pingcap/tidb: v0.0.0-20180928190104-399ea9e2e55f, pingcap/pd: v0.0.0-20160727233714-3ac0863d7acf
github.com/dustin/go-humanize | pingcap/tidb: v1.0.0, pingcap/pd: v0.0.0-20180421182945-02af3965c54e
github.com/gogo/protobuf | pingcap/tidb: v1.2.0, pingcap/pd: v1.2.1
github.com/golang/protobuf | pingcap/tidb: v1.2.0, pingcap/pd: v1.3.2
github.com/golang/snappy | pingcap/tidb: v0.0.1, pingcap/pd: v0.0.0-20180518054509-2e65f85255db
github.com/google/btree | pingcap/tidb: v0.0.0-20180813153112-4030bb1f1f0c, pingcap/pd: v1.0.0
github.com/gorilla/context | pingcap/tidb: v1.1.1, pingcap/pd: v0.0.0-20160226214623-1ea25387ff6f
github.com/gorilla/mux | pingcap/tidb: v1.6.2, pingcap/pd: v1.6.1
github.com/gorilla/websocket | pingcap/tidb: v1.4.0, pingcap/pd: v1.2.0
github.com/montanaflynn/stats | pingcap/tidb: v0.0.0-20180911141734-db72e6cae808, pingcap/pd: v0.0.0-20151014174947-eeaced052adb
github.com/pingcap/kvproto | pingcap/tidb: v0.0.0-20191106014506-c5d88d699a8d, pingcap/pd: v0.0.0-20190516013202-4cf58ad90b6c
github.com/prometheus/client_golang | pingcap/tidb: v0.9.0, pingcap/pd: v1.0.0
github.com/unrolled/render | pingcap/tidb: v0.0.0-20180914162206-b9786414de4d, pingcap/pd: v0.0.0-20171102162132-65450fb6b2d3
go.etcd.io/etcd | pingcap/tidb: v0.0.0-20190320044326-77d4b742cdbf, pingcap/pd: v0.5.0-alpha.5.0.20191023171146-3cf2f69b5738
go.uber.org/zap | pingcap/tidb: v1.9.1, pingcap/pd: v1.10.0
google.golang.org/grpc | pingcap/tidb: v1.17.0, pingcap/pd: v1.23.1
```
