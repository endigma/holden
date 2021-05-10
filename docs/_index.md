---
Title: README
---

# holden

dynamic markdown documentation written in go

## running

holden will run out of the box and create a config next to the binary

holden also looks for configs in:
- ~/.holden/
- /etc/holden
- ./

with the filename `config.toml`.

## features

- no javascript
- entirely server sided
- reactive formatting on mobile
- easily hackable css and formatting
- dynamically generated sidebar/index
- optional custom pages
- fancy curl rendering with ansi (if you request the raw file)

## update from source control

if you want to host your documentation elsewhere, try integrating [fennec](https://github.com/endigma/fennec) into your workflow, it can accept POST requests from git services and trigger a script to update your production docs.

## custom system pages

- `_sidebar.md` in the docroot will replace the default autogenned sidebar
- `_index.md` will be the index for a folder, if a folder has an index it is clickable in the autogenned sidebar
- `_404.md` in the docroot will replace the default 404 page

## configuration

for configuration instructions refer to the comments in the example config file in `example.toml`

you can also use envars like `HLDN_<KEY>=<VALUE>`

## scripting

if you'd like to get the raw markdown source using `curl`, send the `RawPlease` header with the value `true`, like:

```
curl http://docs.example.com/_index.md -H "RawPlease: true"
```