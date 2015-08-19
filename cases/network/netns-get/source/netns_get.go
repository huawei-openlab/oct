package main

import (
    "log"
    "encoding/json"
    "os"
    "runtime"
    netns "../../netns"
)

type netnsResult struct {
     Netns_get map[string]string `json:"Network.Spec.Netns.Netns_get"`
}

func main() {
    // Lock the OS Thread so we don't accidentally switch namespaces
    runtime.LockOSThread()
    defer runtime.UnlockOSThread()

    // Save the current network namespace
    origns, _ := netns.Get()
    defer origns.Close()
 
    result := new(netnsResult)
    result.Netns_get = make(map[string]string)

    newns, err := netns.New()
    defer newns.Close()
    if err != nil {
        log.Fatal("Create network namespace failed!")
        return
    }

    netns.Set(newns)

    ns, err := netns.Get()
    defer ns.Close()
    if err != nil {
        result.Netns_get["res"] = "failed"
    } else {
        result.Netns_get["res"] = "passed"
    }

    jsonStr, err := json.Marshal(result)
    if err != nil {
        log.Fatalf("Convert to json error: %v\n", err)
        return
    }

    foutfile := "netns-get-res.json"
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
