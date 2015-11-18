package hostendvalidate

import (
	"fmt"
	"os"
	"strings"

	"github.com/huawei-openlab/oct/utils"
)

func ContainerOutputHandler(output string) error {
	if err := namespaceHostHandler(output); err != nil {
		return err
	}
	return nil
}

func namespaceHostHandler(output string) error {
	nsout := utils.GetBetweenStr(output, "[namespace_output_start]", "[namespace_output_end]")
	if strings.EqualFold(nsout, "") {
		return nil
	}
	for _, ns := range strings.Split(nsout, "\n") {
		if !strings.EqualFold(ns, "") {
			linkc := strings.Split(ns, ",")[0]
			nsname := strings.Split(linkc, ":")[0]
			path := strings.Split(ns, ",")[1]
			if !strings.EqualFold(path, "") {
				linkh, _ := os.Readlink(path)
				if !strings.EqualFold(linkh, linkc) {
					return fmt.Errorf("%v namespace expected: %v, actual: %v ", nsname, linkh, linkc)

				}
			}
			if strings.EqualFold(path, "") {
				linkh, _ := os.Readlink("/proc/1/ns/" + nsname)
				if strings.EqualFold(linkh, linkc) {
					return fmt.Errorf("namespace %v path is empty, but namespace inside and outside container is the same", nsname)
				}
			}
		}
	}
	return nil
}
