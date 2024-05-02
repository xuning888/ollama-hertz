wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["sec-ch-ua"] = "\"Chromium\";v=\"124\", \"Google Chrome\";v=\"124\", \"Not-A.Brand\";v=\"99\""
wrk.headers["sec-ch-ua-platform"] = "\"macOS\""
wrk.headers["Referer"] = "http://localhost:8080/"
wrk.headers["sec-ch-ua-mobile"] = "?0"
wrk.headers["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
wrk.body   = '{"content":"time.Second * time.Duration(1000) 表示什么？","llmTimeoutSecond":30,"userId":"1111","maxWindows":30,"llmModel":"llama3"}'

