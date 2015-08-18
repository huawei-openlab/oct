package main

import (
    "log"
    "encoding/json"
    "os"
    "runtime"
    netns "../../netns"
)

type netnsResult struct {
     Netns_del map[string]string `json:"Network.Spec.Netns.Netns_del"`
}

func main() {
    // Lock the OS Thread so we don't accidentally switch namespaces
    runtime.LockOSThread()
    defer runtime.UnlockOSThread()

    // Save the current network namespace
    origns, _ := netns.Get()
    defer origns.Close()
 
    result := new(netnsResult)
    result.Netns_del = make(map[string]string)

    // Create a new network namespace
    newns, err := netns.New()
    defer newns.Close()
    if err != nil {
        log.Fatalf("Create network namspace failed: %v\n", err)
        return
    }

    newns.Close()
    if newns.IsOpen() {
        log.Fatal("Delete network namespace failed!")
        result.Netns_del["res"] = "failed"
    } else {
        result.Netns_del["res"] = "passed"
    }

    jsonStr, err := json.Marshal(result)
    if err != nil {
        log.Fatalf("Convert to json error: %v\n", err)
        return
    }

    foutfile := "netns-del-res.json"
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
