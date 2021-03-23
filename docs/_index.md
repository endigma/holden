---
Title: README
---

# holden

dynamic markdown documentation written in go

## running

in order to run holden you have to pass a config file through the command line, or put a `config.toml` file where the binary is located. to use the shipped config run:
```
./holden example.toml
```

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

for configuration instruction refer to the comments in the example config file in `example.toml`

## scripting

if you'd like to get the raw markdown source using `curl`, send the `RawPlease` header with the value `true`, like:

```
curl http://docs.example.com/_index.md -H "RawPlease: true"
```