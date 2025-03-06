package main

import (
    "fmt"
    "io"
    "log"
    "net"
    "strings"
    "time"
)

const (
    proxyAddr  = ":8080"        // Address where the proxy server will listen
    targetAddr = "github.com:22" // Default remote Git server (can be configured)
)

func main() {
    fmt.Printf("Starting Git push proxy server on %s -> %s\n", proxyAddr, targetAddr)

    // Create TCP listener
    listener, err := net.Listen("tcp", proxyAddr)
    if err != nil {
        log.Fatalf("Failed to listen on %s: %v", proxyAddr, err)
    }
    defer listener.Close()

    for {
        // Accept client connections
        clientConn, err := listener.Accept()
        if err != nil {
            log.Printf("Failed to accept connection: %v", err)
            continue
        }

        // Handle each connection in a goroutine
        go handleConnection(clientConn)
    }
}

func handleConnection(clientConn net.Conn) {
    defer clientConn.Close()

    // Connect to the target Git server
    targetConn, err := net.Dial("tcp", targetAddr)
    if err != nil {
        log.Printf("Failed to connect to target %s: %v", targetAddr, err)
        return
    }
    defer targetConn.Close()

    // Create channels to signal when proxy operations are done
    clientDone := make(chan bool)
    targetDone := make(chan bool)

    // Forward data from client to target with inspection
    go func() {
        defer close(clientDone)

        buffer := make([]byte, 4096)
        isPushDetected := false

        for {
            clientConn.SetReadDeadline(time.Now().Add(5 * time.Second))
            n, err := clientConn.Read(buffer)
            if err != nil {
                if err != io.EOF && !strings.Contains(err.Error(), "timeout") {
                    log.Printf("Error reading from client: %v", err)
                }
                break
            }

            // Look for Git push command signatures in the buffer
            data := string(buffer[:n])
            if strings.Contains(data, "git-receive-pack") ||
               strings.Contains(data, "push") {
                isPushDetected = true
                log.Printf("Git push detected!")

                // You can add custom logic here - e.g., validate the push,
                // modify content, log details, etc.

                // Example: Log push details
                log.Printf("Push details: %s", data)

                // Example: If you want to reject the push
                // You could return a Git protocol error here instead of forwarding
                // For now, we just let it through but log it
            }

            // Forward the data to the target
            _, err = targetConn.Write(buffer[:n])
            if err != nil {
                log.Printf("Error writing to target: %v", err)
                break
            }
        }

        if isPushDetected {
            log.Printf("Finished handling Git push")
        }
    }()

    // Forward data from target to client (unmodified)
    go func() {
        defer close(targetDone)

        buffer := make([]byte, 4096)
        for {
            targetConn.SetReadDeadline(time.Now().Add(5 * time.Second))
            n, err := targetConn.Read(buffer)
            if err != nil {
                if err != io.EOF && !strings.Contains(err.Error(), "timeout") {
                    log.Printf("Error reading from target: %v", err)
                }
                break
            }

            _, err = clientConn.Write(buffer[:n])
            if err != nil {
                log.Printf("Error writing to client: %v", err)
                break
            }
        }
    }()

    // Wait for both directions to complete
    <-clientDone
    <-targetDone
}
