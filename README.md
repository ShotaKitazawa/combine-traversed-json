# traverse-and-merge-json

```
$ tree .
.
└── hoge
     ├── fuga
     │   └── metadata.json
     └── metadata.json
```

```
$ cat hoge/fuga/metadata.json
{
  "a1": "ok",
  "a2": "ok",
  "b1": [
    "a",
    "b"
  ],
  "c1": {
    "c11": "ok",
    "c12": "ok"
  }
}
```

```
$ cat hoge/metadata.json
{
  "a1": "invalid",
  "b1": [
    "a",
    "ok"
  ],
  "c1": {
    "c11": "invalid",
    "c13": "ok"
  }
}
```

```
$ traverse-and-merge-json hoge/fuga metadata.json
{"a1":"ok","a2":"ok","b1":["a","b","ok"],"c1":{"c11":"ok","c12":"ok","c13":"ok"}}
```
