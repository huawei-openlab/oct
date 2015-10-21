// +build v0.1.1

//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package linuxresources

import (
	"github.com/huawei-openlab/oct/tools/runtimeValidator/adaptor"
	"github.com/huawei-openlab/oct/tools/runtimeValidator/manager"
	"github.com/opencontainers/specs"
)

func TestMemoryLimit() string {
	var testResourceseMemory specs.Resources = specs.Resources{
		Memory: specs.Memory{
			Limit:       204800,
			Reservation: 0,
			Swap:        0,
			Kernel:      0,
			Swappiness:  -1,
		},
	}
	linuxspec, linuxruntimespec := setResources(testResourceseMemory)
	failinfo := "Memory Limit"
	c := make(chan bool)
	go func() {
		testResources(&linuxspec, &linuxruntimespec)
		close(c)
	}()
	result, err := checkConfigurationFromHost("memory", "memory.limit_in_bytes", "204800", failinfo)
	<-c
	var testResult manager.TestResult
	testResult.Set("TestMemoryLimit", testResourceseMemory.Memory, err, result)
	adaptor.DeleteRun()
	return testResult.Marshal()
}

func TestCpuQuota() string {
	var testResourceCPU specs.Resources = specs.Resources{
		CPU: specs.CPU{
			Shares:          0,
			Quota:           20000,
			Period:          0,
			RealtimeRuntime: 0,
			RealtimePeriod:  0,
			Cpus:            "",
			Mems:            "",
		},
	}
	linuxspec, linuxruntimespec := setResources(testResourceCPU)
	failinfo := "CPU Quota"
	c := make(chan bool)
	go func() {
		testResources(&linuxspec, &linuxruntimespec)
		close(c)
	}()
	result, err := checkConfigurationFromHost("cpu", "cpu.cfs_quota_us", "20000", failinfo)
	<-c
	var testResult manager.TestResult
	testResult.Set("TestMemoryLimit", testResourceCPU.CPU, err, result)
	adaptor.DeleteRun()
	return testResult.Marshal()
}

func TestBlockIOWeight() string {
	var testResourceBlockIO specs.Resources = specs.Resources{
		BlockIO: specs.BlockIO{
			Weight:                  300,
			WeightDevice:            nil,
			ThrottleReadBpsDevice:   nil,
			ThrottleWriteBpsDevice:  nil,
			ThrottleReadIOPSDevice:  nil,
			ThrottleWriteIOPSDevice: nil,
		},
	}
	linuxspec, linuxruntimespec := setResources(testResourceBlockIO)
	failinfo := "BlockIO Weight"
	c := make(chan bool)
	go func() {
		testResources(&linuxspec, &linuxruntimespec)
		close(c)
	}()
	result, err := checkConfigurationFromHost("blkio", "blkio.weight", "300", failinfo)
	<-c
	var testResult manager.TestResult
	testResult.Set("TestBlockIOWeight", testResourceBlockIO.BlockIO, err, result)
	adaptor.DeleteRun()
	return testResult.Marshal()
}

func TestHugepageLimit() string {
	var testResourcehugtlb specs.Resources = specs.Resources{
		HugepageLimits: []specs.HugepageLimit{
			{
				Pagesize: "2MB",
				Limit:    409600,
			},
		},
	}
	linuxspec, linuxruntimespec := setResources(testResourcehugtlb)
	failinfo := "Hugepage Limit"
	c := make(chan bool)
	go func() {
		testResources(&linuxspec, &linuxruntimespec)
		close(c)
	}()
	result, err := checkConfigurationFromHost("hugetlb", "hugetlb."+testResourcehugtlb.HugepageLimits[0].Pagesize+".limit_in_bytes", "409600", failinfo)
	<-c
	var testResult manager.TestResult
	testResult.Set("TestHugepageLimit", testResourcehugtlb.HugepageLimits, err, result)
	adaptor.DeleteRun()
	return testResult.Marshal()
}
