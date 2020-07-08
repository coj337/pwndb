# pwndb
pwndb is a tool that looks for leaked passwords from a dark web breach database given a user or domain.


### Credit
Tool has been extended from the one made by [PostRequest](https://github.com/postrequest) [here](https://github.com/postrequest/tools/tree/master/pwndb).

### Usage
Any number of users and domains can be specified.
```
pwndb -user foo -user bar
pwndb -domain foo.com -domain bar.com -domain baz.com
```

If at least one user and domain is specfied, all permutations will be checked but it will return *only* results that contain both.  
The below command will check foo@baz.com and bar@baz.com
```
pwndb -user foo -user bar -domain baz.com
```

### Installation
1. Download and install the go tools. [https://golang.org/dl/](https://golang.org/dl/)
2. Find the location of your GOPATH.
    ```
    go env GOPATH
    ```
3. Drop pwndb.go in GOPATH/src/pwndb (create any folders that don't exist)
4. Run the go install command and the binary should be able to run from anywhere.

    ```
    go install pwndb
    pwndb
    ```
### Quick Setup
Download and run the binary for your platform below. (or from [releases](https://github.com/coj337/pwndb/releases))  
[Windows](https://github.com/coj337/pwndb/releases/download/v1.1/pwndb-windows.exe)  
[Linux](https://github.com/coj337/pwndb/releases/download/v1.1/pwndb-linux)  
[OSX](https://github.com/coj337/pwndb/releases/download/v1.1/pwndb-osx)

### Demo
![demo image](https://github.com/coj337/pwndb/blob/master/demo.png)
