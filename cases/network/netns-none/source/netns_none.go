package main

import (
    "log"
    "encoding/json"
    "os"
    "runtime"
    netns "../../netns"
)

type netnsResult struct {
     Netns_none map[string]string `json:"Network.Spec.Netns.Netns_none"`
}

func main() {
    // Lock the OS Thread so we don't accidentally switch namespaces
    runtime.LockOSThread()
    defer runtime.UnlockOSThread()

    // Save the current network namespace
    origns, _ := netns.Get()
    defer origns.Close()
 
    result := new(netnsResult)
    result.Netns_none = make(map[string]string)

    ns := netns.None()
    if ns.IsOpen() {
        log.Fatal("None ns is open")
        result.Netns_none["res"] = "failed"
    } else {
        result.Netns_none["res"] = "passed"
    }

    jsonStr, err := json.Marshal(result)
    if err != nil {
        log.Fatalf("Convert to json error: %v\n", err)
        return
    }

    foutfile := "netns-none-res.json"
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
