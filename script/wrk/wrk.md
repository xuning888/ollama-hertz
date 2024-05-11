
# stream
```shell
wrk -t4 -c4 -d10s --timeout=10m -s chat_stream.lua http://localhost:8080/api/v1/chat/stream
```