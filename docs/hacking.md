# Hacking and Customization

## Themeing

we have two main css files:

- `vars.css`, which stores variables with colors that are used by the main css file and is extremely easy to edit.
- `style.css`, which actually styles the website.

we only support changing `vars.css`, if you change anything in `style.css` you're on your own.

## Changing template

in order to change the template you need to edit `assets/static/page.html`, but this again is something we do not support. you're on your own.

### template variables

```go
type Page struct {
	Prefix          string // the prefix defined in the config
	Contents        string // html rendered from markdown source
	Meta            map[string]interface{} // anything you define in the yaml header, look below
	SidebarContents string // html sidebar content
	Raw             string // raw path of the file being served
}
```

meta uses the yaml header for the document, for example if your markdown file is:

```md
---
Title = Hello, world!
Author = Me
---

# Hello, world!
```

`{{ .Meta.Title }}` will be "Hello, world!" and `{{ .Meta.Author }}` will be "Me".