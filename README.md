# conf8n

Package conf8n is here to simplify the way you work with your config files. It helps to avoid type-casting-hell in your projects. Now you can load, read and traverse even complicated hierarchical configs a clear & simple way.

## Installation
```
go get github.com/safronizator/conf8n
```

## Use
```go
import (
	"github.com/safronizator/conf8n"
	"fmt"
)

func main() {
	// myconf.yaml:
	// db:
	//   host: 192.168.0.1
	//   account:
	//     login: db_user
	//     password: xxxxxxx
	//   tables: ["posts", "comments", "likes"]
	conf, err := conf8n.NewConfigFromFile("myconf.yaml")
	if err != nil {
		panic(err)
	}
	db := somedb.Connect(
		conf.Get("db.host").DefString("127.0.0.1"), // "192.168.0.1"
		conf.Get("db.port").DefInt(3306))           // 3306
	accConf := conf.Get("db.account").Config()
	login := accConf.Get("login")
	pwd := accConf.Get("password")
	if !login.IsSet() || !pwd.IsSet() {
	  	panic("invalid account data")
	}
	db.Auth(login.String(), password.String())
	fmt.Println("Authorized tables:")
	for i := conf.Get("db.tables").Iterate(); !i.Finished(); i.Next() {
	  	fmt.Println(i.Value())
	}
}
```

## Docs

[![GoDoc](https://godoc.org/github.com/safronizator/conf8n?status.svg)](https://godoc.org/github.com/safronizator/conf8n)

## @todo
- docs & examples
- tests
- values setting
- config data writers
- support for async use
