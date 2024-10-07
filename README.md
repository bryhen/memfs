memfs (pronounced mem-fis) is a in-memory file server. It's provides major speed benefits over http.FileServer() in servers that have excess memory and desire high throughput.

# How do you know if memfs is right for your server?

# How does memfs work
 * memfs.New() supplies a FileServer that points to a directory. It reads the contents into memory, skipping any that cannot form a valid HTTP URL.
 * 
# Example
fs := 