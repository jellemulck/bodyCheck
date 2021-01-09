# bodyCheck

A simple Go tool that fetches a URL and checks the body for content.

# Install

```
# go get github.com/jellemulck/bodyCheck
```

# usage example

```
bodyCheck -file /tmp/file_with_urls.txt -content root:x -threads 10 -path ../../etc/passwd
```

# description:

```
- file : add an input file with domain URLs
- content : the content to search for
- threads : number of concurrent times the program runs
- path : add a path after the domain url
```