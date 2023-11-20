package ojdmcollector

import (
	"reflect"
	"testing"
)

func Test_full_version_settings(t *testing.T) {

	expectedFullJInfo := JavaInfoRunningProcs{
		JavaHome:           "/home/alin/Documents/android-studio/jbr",
		JavaRuntimeName:    "OpenJDK Runtime Environment",
		JavaRuntimeVersion: "17.0.6+0-17.0.6b802.4-9586694",
		JavaVendor:         "JetBrains s.r.o.",
		JavaVersion:        "17.0.6",
		JavaVersionDate:    "2023-01-17",
		JavaVMName:         "OpenJDK 64-Bit Server VM",
		JavaVMVendor:       "JetBrains s.r.o.",
		JavaVMVersion:      "17.0.6+0-17.0.6b802.4-9586694",
	}

	fullVersionSettings := `
	VM settings:
    Max. Heap Size (Estimated): 3.85G
    Using VM: OpenJDK 64-Bit Server VM

Property settings:
    file.encoding = UTF-8
    file.separator = /
    java.class.path = 
    java.class.version = 61.0
    java.home = /home/alin/Documents/android-studio/jbr
    java.io.tmpdir = /tmp
    java.library.path = /usr/java/packages/lib
        /usr/lib64
        /lib64
        /lib
        /usr/lib
    java.runtime.name = OpenJDK Runtime Environment
    java.runtime.version = 17.0.6+0-17.0.6b802.4-9586694
    java.specification.name = Java Platform API Specification
    java.specification.vendor = Oracle Corporation
    java.specification.version = 17
    java.vendor = JetBrains s.r.o.
    java.vendor.url = https://openjdk.java.net/
    java.vendor.url.bug = https://bugreport.java.com/bugreport/
    java.version = 17.0.6
    java.version.date = 2023-01-17
    java.vm.compressedOopsMode = Zero based
    java.vm.info = mixed mode
    java.vm.name = OpenJDK 64-Bit Server VM
    java.vm.specification.name = Java Virtual Machine Specification
    java.vm.specification.vendor = Oracle Corporation
    java.vm.specification.version = 17
    java.vm.vendor = JetBrains s.r.o.
    java.vm.version = 17.0.6+0-17.0.6b802.4-9586694
    jbr.virtualization.information = No virtualization detected
    jdk.debug = release
    line.separator = \n 
    native.encoding = UTF-8
    os.arch = amd64
    os.name = Linux
    os.version = 6.2.0-36-generic
    path.separator = :
    sun.arch.data.model = 64
    sun.boot.library.path = /home/alin/Documents/android-studio/jbr/lib
    sun.cpu.endian = little
    sun.io.unicode.encoding = UnicodeLittle
    sun.java.launcher = SUN_STANDARD
    sun.jnu.encoding = UTF-8
    sun.management.compiler = HotSpot 64-Bit Tiered Compilers
    sun.stderr.encoding = UTF-8
    sun.stdout.encoding = UTF-8
    user.country = US
    user.dir = /home/alin/Documents/licenseware/ojdm-collector
    user.home = /home/alin
    user.language = en
    user.name = alin
	`

	t.Run(
		"Full version settings with -XshowSettings:all -version",
		func(t *testing.T) {
			jinfo := extractInfoFromFullVersionSettings(fullVersionSettings)
			if !reflect.DeepEqual(jinfo, expectedFullJInfo) {
				t.Errorf("extractInfoFromFullVersionSettings() = \n %+v,\n expected =\n %+v", jinfo, expectedFullJInfo)
			}
		},
	)

}

func Test_partial_version_settings(t *testing.T) {

	expectedPartialJInfo := JavaInfoRunningProcs{
		JavaHome:           "/home/alin/Documents/android-studio/jbr",
		JavaRuntimeName:    "OpenJDK Runtime Environment",
		JavaRuntimeVersion: "17.0.6+0-17.0.6b802.4-9586694",
		JavaVersion:        "17.0.6",
		JavaVersionDate:    "2023-01-17",
		JavaVMName:         "OpenJDK 64-Bit Server VM",
	}

	partialVersionSettings := `
/home/alin/Documents/android-studio/jbr/bin/java
openjdk version "17.0.6" 2023-01-17
OpenJDK Runtime Environment (build 17.0.6+0-17.0.6b802.4-9586694)
OpenJDK 64-Bit Server VM (build 17.0.6+0-17.0.6b802.4-9586694, mixed mode)
`

	t.Run(
		"Partial version settings with -version",
		func(t *testing.T) {
			jinfo := extractInfoFromFullVersionSettings(partialVersionSettings)
			if !reflect.DeepEqual(jinfo, expectedPartialJInfo) {
				t.Errorf("extractInfoFromFullVersionSettings() = \n %+v,\n expected =\n %+v", jinfo, expectedPartialJInfo)
			}
		},
	)

}
