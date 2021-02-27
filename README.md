# vkbotserver
[![Go Report Card](https://goreportcard.com/badge/github.com/Sepuka/vkbotserver)](https://goreportcard.com/report/github.com/Sepuka/vkbotserver)

There is a server-bot for vk.com

## Usage
1. Add server as a dependency in your Gopkg.toml

```
[[constraint]]
  name = "github.com/sepuka/vkbotserver"
  version = "v0.0.2"
```
  
2. Instance server

```
var handlers := []Executor{
    message.NewConfirmation(cfg.Server)
}

var handlerMap = make(message.HandlerMap, len(handlers))

for _, cmd := range handlers {
    msgName = cmd.(fmt.Stringer).String()
    handlerMap[msgName] = cmd.(message.Executor)
}

var simpleHandler = func (handler message.Executor, req *domain.Request, resp http.ResponseWriter) error {
    return handler.Exec(req, resp)
}

var server = server.NewSocketServer(cfg.Server, handlerMap, simpleHandler), nil
```

3. Run and listen

```
server.Listen()
```

## Nginx settings

```
        location ~^/your_postfix/ {
                include         fastcgi_params;
                fastcgi_pass    unix:/var/run/SOME_PATH/server.sock;
                access_log      /var/log/nginx/APP_access.log;
                error_log       /var/log/nginx/APP_error.log;
        }
```
