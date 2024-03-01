
## Comiple Ordinary file to go_code and  so the file can be comiled into go binary file 。

#### install 

````
go install github.com/seeadoog/go-rescompile 
````

#### Usage: 
run cmd 
````
go-rescompile pages_src page

````

this command will compile director 'parges_src' to a go package 'page' .
all files is defined in index.go 
```
package webpages
var Files = map[string][]byte{
"index.html":a_69fef4ec6d0bdf042b910bd3f5a084d5,
}
```