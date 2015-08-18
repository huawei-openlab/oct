package main

import (
    "fmt"
    "log"
    "encoding/json"
    "os"
    "runtime"
    netns "../../netns"
)

type netnsResult struct {
     netns_Create map[string]string `json:"Network.Spec.Netns.Create"`
}

func main() {
    // Lock the OS Thread so we don't accidentally switch namespaces
    runtime.LockOSThread()
    defer runtime.UnlockOSThread()

    // Save the current network namespace
    origns, _ := netns.Get()
    defer origns.Close()
 
    result := new(netnsResult)
    result.netns_Create = make(map[string]string)

    // Create a new network namespace
    newns, err := netns.New()
    defer newns.Close()
    if err != nil {
        log.Fatalf("Create network namspace failed: %v\n", err)
        result.netns_Create["false"] = "failed"
        fmt.Println(result)
    } else {
        fmt.Println("Create network namespace successfully!")
        result.netns_Create["false"] = "passed"
        fmt.Println(result)
    }

    jsonStr, err := json.Marshal(result)
    if err != nil {
        log.Fatalf("Convert to json error: %v\n", err)
        return
    }
    fmt.Println(jsonStr)

    foutfile := "netns-create-res.json"
    fout, err := os.Create(foutfile)
    if err != nil {
        log.Fatal(err)
    } else {
        fmt.Println(jsonStr)
        fout.WriteString(string(jsonStr))
    }

    // Switch back to the original namespace
    netns.Set(origns)
}
