memfs (pronounced mem-fis) is a in-memory file server. It provides major speed benefits over http.FileServer() for servers that have memory to spare. Package benchmarks suggest 8-25x speed increases over the std lib http.FileServer().

# How do you know if memfs is right for your server?
* Your server desires high throughput and has memory to spare in excess of the target directory's size.
* If you're not sure whether you have memory to spare, run a test with TODO log.Println("Dir Size: ", memfs.DirSize("./YOUR_DIR")) and consider whether your servers can spare that much memory at runtime. Alternatively, you can scale up instance size if it's an issue.
* You just want to serve files faster.
    
# How does memfs work
* memfs supplies a FileServer that points to a directory.
* It reads the contents into memory, skipping any that do not form a valid HTTP URL.
    * File creations, updates, or deletions after the FileServer has been created have no effect. This is usually not an issue because this package's primary use case is to quickly serve static front end builds, such as a React App.
* Files are matched by exact paths. Additionally, "/" matches "index.*", such as index.html or index.production.html in the directory or any subdirectory ("/" finds "/index.html", while "/nested/" finds "/nested/index.html").
    * For variable pattern matching, you'll need to implement middleware appropriate to your use case.
    * Files outside of the target directory cannot be matched.
* Files are served by providing the http.ResponseWriter with the pre-loaded []byte read from the target directory when the FileServer was first created.

# Example of a good use case
* Your server has 8gb memory. At peak times, the server is running at 100% CPU and has 50% memory utilization. Part of its job is to serve a React App that has a total build size (including all assets, such as images) of 1gb. memfs can use 1gb of the 4gb available to speed up file serving.

# Example
Structure:
```
📦example
 ┣ 📂app
 ┃ ┣ 📂dist
 ┃ ┃ ┣ 📂assets
 ┃ ┃ ┃ ┣ 📜favicon.ico
 ┃ ┃ ┃ ┣ 📜invalid url.txt
 ┃ ┃ ┃ ┣ 📜script.js
 ┃ ┃ ┃ ┗ 📜someImage.jpg
 ┃ ┃ ┗ 📜index.html
 ┃ ┣ 📂src
 ┃ ┃ ┣ 📂assets
 ┃ ┃ ┃ ┣ 📜favicon.ico
 ┃ ┃ ┃ ┣ 📜invalid url.txt
 ┃ ┃ ┃ ┣ 📜script.js
 ┃ ┃ ┃ ┗ 📜someImage.jpg
 ┃ ┃ ┣ 📜otherSource.tsx
 ┃ ┃ ┗ 📜source.tsx
 ┃ ┣ 📜index.html
 ┃ ┗ 📜tsconfig.json
 ┗ 📜main.go

// main.go
func main() {
    fs, failed, err := memfs.New("./app/dist")
    if err != nil {
        panic(err)
    }

    // Prints: The following routes failed to load: invalid url.txt
    log.Println("The following routes failed to load: ", failed)

    fs.NotFoundHandler(handler.NotFoundHandler)

    mux := http.NewServeMux()
    mux.Handle("GET /endpoint", handler.EndpointHandler)
    mux.Handle("GET /", fs)

    s := http.Server{
        handler: mux,
    }

    if err := s.ListenAndServe(); err != nil {
        log.Println("Server shutting down: ", err)
    }
}
```

# Version Policy
memfs follows semantic versioning for the documented public API on stable releases.

# Testing
TODO

# Future Directions
If you'd like some additional functionality, feel free to reach out. We can talk about your needs and see whether I can implement them. If one of the topics below could be useful for you, make sure to mention it so that it can become a priority instead of a possibility.
* Allow the file server to account for creates/updates/deletes on the target directory while handling requests.

# Donations
If this library has saved you or your company money through a reduction in server costs, or you appreciate the simple benefits provided by the package, please consider a financial donation to: . Much appreciated!