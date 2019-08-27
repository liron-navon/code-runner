# code runner

A simple experiment to execute files in a secure context with docker through an API, will probably look at other directions since the container spinup + execution takes  ~1-2  seconds.

## How to use:

call `make start`

It will expose one route: `POST <url>/exec` which eccept two arguments, `code` and `language`, and return the output lines:

Currently support 4 languages: `java`, `python`, `go`, `node`

Examples:

**Python3**
```
POST http://localhost:8080/exec

Request:
{
	"code": "print("hello from python")",
	"language": "python3"
}

Response:
{
    "output": ["hello from python"]
}
```

**Java**
```
POST http://localhost:8080/exec

Request:
{
	"code": `
    public class HelloWorld {
        public static void main(String[] args) {
            System.out.println("Hello, from the world of java!");
            System.out.println("Another linefrom the world of java!");
        }
    }
    `,
	"language": "python3"
}

Response:
{
    "output": ["Hello, from the world of java!", "Another linefrom the world of java!"]
}
```

**Node**
```
POST http://localhost:8080/exec

Request:
{
	"code": "console.log("hello from node.js")",
	"language": "node"
}

Response:
{
    "output": ["hello from node.js"]
}
```

**Go**
```
POST http://localhost:8080/exec

Request:
{
	"code": `
    package main

    import "fmt"

    func main() {
        fmt.Println("hello gophers")
    }
    `,
	"language": "go"
}

Response:
{
    "output": ["hello gophers"]
}
```