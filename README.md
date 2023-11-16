# OJDM Collector

Find Java installed versions and the software asociated with installed Java (cross-platform).

Download the binary from Releases or clone the repo and run `make build` (you will find the binary in the `dist` folder).

```shell

$ ojdm-collector --help

OJDMCollector - Utility to find JVMs/JDKs report their versions and related running processes

Usage:
    $ ojdm-collector
    $ ojdm-collector -csv=/path/to/csvreport.csv

  -csv string
        Path to csv report. (default "report.csv")
```

*If you provide a custom csv path make sure path exists.


OJDM Collector sample output:

```json
{
    "HostName": "ClimenteA",
    "JavaHome": "/home/alin/Documents/android-studio/jbr",
    "IsJDK": true,
    "DynLibBinPath": "/home/alin/Documents/android-studio/jbr/lib/server/libjvm.so",
    "JavaVersion": "17.0.6",
    "JavaRuntimeName": "OpenJDK Runtime Environment",
    "JavaVendor": "JetBrains s.r.o.",
    "JavaRuntimeVersion": "17.0.6+0-17.0.6b802.4-9586694",
    "JavaVMName": "OpenJDK 64-Bit Server VM",
    "JavaVMVendor": "JetBrains s.r.o.",
    "JavaVMVersion": "17.0.6+0-17.0.6b802.4-9586694",
    "JavaVersionDate": "2023-01-17",
    "AppDirName": "android-studio",
    "JavaBinPath": "/home/alin/Documents/android-studio/jbr/bin/java",
    "JavaCBinPath": "/home/alin/Documents/android-studio/jbr/bin/javac",
    "BaseDir": "/home/alin/Documents/android-studio",
    "ProcessRunning": true,
    "CommandLine": "/home/alin/Documents/android-studio/jbr/bin/java -classpath etc"
}
```

Equivalent JDowser Output:

```json

[
  {
    "host": "ClimenteA",
    "java_home": "/home/alin/Documents/android-studio/jbr",
    "is_jdk": true,
    "libjvm": "/home/alin/Documents/android-studio/jbr/lib/server/libjvm.so",
    "libjvm_hash": "111047e0da0d8a31e576957d31725f88",
    "version_info": {
      "java_version": "17.0.6",
      "runtime_name": "OpenJDK Runtime Environment",
      "java_runtime_vendor": "JetBrains s.r.o.",
      "java_runtime_version": "17.0.6+0-17.0.6b802.4-9586694",
      "java_vm_name": "OpenJDK 64-Bit Server VM",
      "java_vm_vendor": "JetBrains s.r.o.",
      "java_vm_version": "17.0.6+0-17.0.6b802.4-9586694"
    },
    "running_instances": 0
  }
]

```
