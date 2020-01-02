# go-wiki

## About

Following [Writing Web Applications](https://golang.org/doc/articles/wiki/) to learn some go

## Build and run

This is a very naively written appllication and as such must be run from the top level folder in the repo (i.e. where this readme is) so it can pick up the `tmpl` and `data` folders.

```bash
go build
./wiki
```

To verify go to `localhost:9090` either in browser or using `curl`/`wget`/etc.

## Operations

To view a page go to `/view/<page name>`.  
To edit a page go to `/edit/<page name>` or click the `edit` link on the `view` page.

If a page does not already exist you will be taken to the edit page to create one.
