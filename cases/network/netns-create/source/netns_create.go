package main

import (
    "log"
    "encoding/json"
    "os"
    "runtime"
    netns "../../netns"
)

type netnsResult struct {
     Netns_create map[string]string `json:"Network.Spec.Netns.Netns_create"`
}

func main() {
    // Lock the OS Thread so we don't accidentally switch namespaces
    runtime.LockOSThread()
    defer runtime.UnlockOSThread()

    // Save the current network namespace
    origns, _ := netns.Get()
    defer origns.Close()
 
    result := new(netnsResult)
    result.Netns_create = make(map[string]string)

    // Create a new network namespace
    newns, err := netns.New()
    defer newns.Close()
    if err != nil {
        log.Fatalf("Create network namspace failed: %v\n", err)
        result.Netns_create["res"] = "failed"
    } else {
        result.Netns_create["res"] = "passed"
    }

    jsonStr, err := json.Marshal(result)
    if err != nil {
        log.Fatalf("Convert to json error: %v\n", err)
        return
    }

    foutfile := "netns-create-res.json"
    fout, err := os.Create(foutfile)
    defer fout.Close()

    if err != nil {
        log.Fatal(err)
    } else {
        fout.WriteString(string(jsonStr))
    }

    // Switch back to the original namespace
    netns.Set(origns)
}
