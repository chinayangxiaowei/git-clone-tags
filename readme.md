# Instructions
    Breakpoint cloning of very large git repositories by fetching tags one by one.

### Usage
```
  -end-build
    	calculating the maximum build version (default true)
  -min-build int
    	filter minimum build version number
  -remote string
    	remote url (needed)
  -repo string
    	repository path
  -show-tags
    	display matching tags, but do not clone to local repo
  -tags string
    	tags matching string
```
### Use Cases

1. display the final version of each master release and >= min-build
```shell
./git_clone_tags -remote https://chromium.googlesource.com/chromium/src -show-tags -min-major 60  -min-build 200 -end-build -repo chromium/src
```
2. Filter versions using matching expressions
```shell
./git_clone_tags -remote https://chromium.googlesource.com/chromium/src -show-tags -tags 100.0.*.*
```
3. clone the filtered version
```shell
./git_clone_tags -remote https://chromium.googlesource.com/chromium/src  -tags 4.*.*.* -min-build 200 -repo chromium/src
```