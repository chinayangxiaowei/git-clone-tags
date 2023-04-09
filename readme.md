# Instructions
    Breakpoint cloning of very large git repositories by fetching tags one by one.

### Build
```shell
    # Go to the source code directory
    go build -o git_clone_tags
```

### Usage
```
  -end-build
    	calculating the maximum build version (default true)
  -min-major int
        filter minimum major version number
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
64.0.3282.204
65.0.3325.230
66.0.3359.203
86.0.4240.281
88.0.4324.218
90.0.4430.246
96.0.4664.219
100.0.4896.241
102.0.5005.200
106.0.5249.225
108.0.5359.228
```
2. Filter versions using matching expressions
```shell
./git_clone_tags -remote https://chromium.googlesource.com/chromium/src -show-tags -tags "100.0.*.*"
```
3. clone the filtered version
```shell
./git_clone_tags -remote https://chromium.googlesource.com/chromium/src  -tags "4.*.*.*" -min-build 200 -repo chromium/src
```