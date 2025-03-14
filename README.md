## How to Use
https://github.com/finos/git-proxy
https://git-proxy.finos.org/docs/installation
1. Compile and run the proxy:

```
npm install @finos/git-proxy
git-proxy -c proxy.config.json
```

2. Configure Git to use your proxy:
```bash
git config --global http.proxy http://localhost:8080
# OR at the repo level (replace https://github.com with http://localhost:8080)
git remote set-url origin http://localhost:8000/mathieubellon/tsripe.com.git
# OR for SSH
# Use a tool like socat to redirect SSH traffic through your proxy
```

## Customization Options

1. **Intercepting Logic**: Modify the condition in the `if strings.Contains(data, "git-receive-pack")` block to add your custom interception logic.

2. **Blocking Pushes**: To block pushes, you can replace the forwarding code with a Git protocol error response.

3. **Filtering**: Add content filtering by modifying the buffer before forwarding it.

4. **Authentication**: Add authentication checks before allowing pushes.

Note: For production use, you would need to handle SSH protocols properly, implement proper timeout handling, and add error recovery mechanisms.
