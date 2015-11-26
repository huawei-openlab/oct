package hooks

import (
	"fmt"
	"os"
	"strings"

	"github.com/huawei-openlab/oct/utils"
)

func NamespacePostStart(output string) error {
	nsout := utils.GetBetweenStr(output, "[namespace_output_start]", "[namespace_output_end]")
	if strings.EqualFold(nsout, "") {
		return nil
	}
	for _, ns := range strings.Split(nsout, "\n") {
		if !strings.EqualFold(ns, "") {
			if len(strings.Split(ns, ",")) != 2 {
				break
			}
			linkc := strings.Split(ns, ",")[0]
			if len(strings.Split(linkc, ",")) != 2 {
				break
			}
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

func HooksValidatePostStart(output string) error {
	// Poststart Hook Validate
	hout := utils.GetBetweenStr(output, "[poststop_hookvalidate_output_start]", "[poststop_hookvalidate_output_end]")
	if !strings.EqualFold(hout, "folder poststophook is not exsist inside container") {
		return nil
	}
	if !utils.DirExist("./rootfs/poststophook") {
		return fmt.Errorf("Poststop Hook validation failed")
	}
	// remove extra folders
	if _, err := os.Stat("./rootfs/prestarthook"); err == nil {
		os.Remove("./rootfs/prestarthook")
	}
	if _, err := os.Stat("./rootfs/poststophook"); err == nil {
		os.Remove("./rootfs/poststophook")
	}
	return nil
}
