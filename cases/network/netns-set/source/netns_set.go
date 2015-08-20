package main

import (
    "log"
    "encoding/json"
    "os"
    "runtime"
    netns "../../netns"
)

type netnsResult struct {
     Netns_set map[string]string `json:"Network.Spec.Netns.Netns_set"`
}

func main() {
    
    // Lock the OS Thread so we don't accidentally switch namespaces
    runtime.LockOSThread()
    defer runtime.UnlockOSThread()

    // Save the current network namespace
    origins, _ := netns.Get()
    defer origins.Close()
 
    result := new(netnsResult)
    result.Netns_set = make(map[string]string)

    newns, err := netns.New()
    defer newns.Close()
    if err != nil || origins.Equal(newns) {
        log.Fatal("Create network namespace failed!")
        result.Netns_set["res"] = "failed"
        goto output
    }

    if err = netns.Set(origins); err != nil {
        log.Fatal("Network namespace set function failed!")
        result.Netns_set["res"] = "failed"
        goto output
    } 

    newns.Close()
    if newns.IsOpen() {
        log.Fatal("Network namespace close function failed!")
        result.Netns_set["res"] = "failed"
        goto output
    }

    newns, err = netns.Get()
    if err != nil {
        log.Fatal("Network namespace get function failed!")
        result.Netns_set["res"] = "failed"
        goto output
    }

    if !newns.Equal(origins){
        log.Fatal("Network namespace reset failed")
        result.Netns_set["res"] = "passed"
        goto output
    }

    result.Netns_set["res"] = "passed"

    output: jsonStr, err := json.Marshal(result)
            if err != nil {
                 log.Fatalf("Convert to json error: %v\n", err)
                 return
            }

            foutfile := "netns-set-res.json"
            fout, err := os.Create(foutfile)
            defer fout.Close()

            if err != nil {
            log.Fatal(err)
            } else {
                fout.WriteString(string(jsonStr))
            }
}
