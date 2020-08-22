cq
---

<p align="center">
  <img width="180px" src="https://github.com/mikemackintosh/cq/blob/master/.github/logo.png?raw=true">
</p>

# Overview
First things first, this project was inspired by [`jq`](https://stedolan.github.io/jq/).

There are plenty of legacy systems that only output or support CSV. However, much tooling in the tech field these days doesn't support automation around CSV. As a result, we need a way to convert or adjust

```sh
cat file.csv | cq '\(.name)'
```

Or convert the CSV to JSON:

```sh
cat file.csv | cq -j
```
