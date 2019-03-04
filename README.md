# personal-kv

Really simple (and dump) key-value store accessible via a simple API. I'm using it to serialize data between iOS Shortcuts actions.

## Usage

Set a key:

```
POST /

{
  "action": "set",
  "key": "foo",
  "val": "bar"
}
```

Get a key:

```
POST /

{
  "action": "get",
  "key": "foo"
}
```
