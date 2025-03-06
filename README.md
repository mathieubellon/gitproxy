## How to Use

1. Compile and run the proxy:
```
go build -o gitproxy
./gitproxy
```

2. Configure Git to use your proxy:
```bash
git config --global http.proxy http://localhost:8080
# OR for SSH
# Use a tool like socat to redirect SSH traffic through your proxy
```

## Customization Options

1. **Intercepting Logic**: Modify the condition in the `if strings.Contains(data, "git-receive-pack")` block to add your custom interception logic.

2. **Blocking Pushes**: To block pushes, you can replace the forwarding code with a Git protocol error response.

3. **Filtering**: Add content filtering by modifying the buffer before forwarding it.

4. **Authentication**: Add authentication checks before allowing pushes.

Note: For production use, you would need to handle SSH protocols properly, implement proper timeout handling, and add error recovery mechanisms.
