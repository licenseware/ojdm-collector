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





[oracle@wcp12cr2 /]$ find / -type f -name java 2>/dev/null

/var/lib/alternatives/java
/usr/java/jdk1.8.0_60/jre/bin/java
/usr/java/jdk1.8.0_60/bin/java
/usr/lib/jvm/java-1.6.0-openjdk-1.6.0.0.x86_64/jre/bin/java
/usr/lib/jvm/java-1.6.0-openjdk-1.6.0.0.x86_64/bin/java
/usr/lib/jvm/java-1.7.0-openjdk-1.7.0.9.x86_64/jre/bin/java
/usr/lib/jvm/java-1.7.0-openjdk-1.7.0.9.x86_64/bin/java
/oracle/fmw/ohs/oracle_common/jdk/jre/bin/java
/oracle/db/ohome/jdk/jre/bin/java
/oracle/db/ohome/jdk/bin/java


[oracle@wcp12cr2 /]$ find / -type f -name libjvm.so 2>/dev/null

/usr/java/jdk1.8.0_60/jre/lib/amd64/server/libjvm.so
/usr/lib/jvm/java-1.6.0-openjdk-1.6.0.0.x86_64/jre/lib/amd64/server/libjvm.so
/usr/lib/jvm/java-1.7.0-openjdk-1.7.0.9.x86_64/jre/lib/amd64/server/libjvm.so
/usr/lib64/gcj-4.4.4/libjvm.so
/oracle/fmw/ohs/oracle_common/jdk/jre/lib/amd64/server/libjvm.so
/oracle/db/ohome/jdk/jre/lib/amd64/server/libjvm.so


[oracle@wcp12cr2 /]$ find / -type f -name javac 2>/dev/null

/var/lib/alternatives/javac
/usr/java/jdk1.8.0_60/bin/javac
/usr/lib/jvm/java-1.6.0-openjdk-1.6.0.0.x86_64/bin/javac
/usr/lib/jvm/java-1.7.0-openjdk-1.7.0.9.x86_64/bin/javac
/oracle/db/ohome/jdk/bin/javac



[oracle@wcp12cr2 oracle]$ find /oracle -type f -name libjvm.so
/oracle/fmw/ohs/oracle_common/jdk/jre/lib/amd64/server/libjvm.so
/oracle/db/ohome/jdk/jre/lib/amd64/server/libjvm.so



Output from Oracle VM:



[oracle@wcp12cr2 ojdm-collector_0.0.19_linux_amd64]$ ./ojdm-collector

Licenseware OJDM Collector - Gather all java info in one place
Found libjvm.so in path /oracle/db/ohome/jdk/jre/lib/amd64/server/libjvm.so
Found libjvm.so in path /usr/lib/jvm/java-1.5.0-gcj-1.5.0.0/jre/lib/x86_64/server/libjvm.so
Found libjvm.so in path /usr/lib/jvm/java-1.6.0-openjdk-1.6.0.0.x86_64/jre/lib/amd64/server/libjvm.so
Found libjvm.so in path /usr/lib/jvm/java-1.7.0-openjdk-1.7.0.9.x86_64/jre/lib/amd64/server/libjvm.so
Encountered error lstat /snap: no such file or directory

All running processes

[
  {
    "Name": "java",
    "ProcDir": "/oracle/fmw/ohs/oracle_common/jdk/jre/bin/java",
    "CommandLine": "/oracle/fmw/ohs/oracle_common/jdk/jre/bin/java -server -Xms32m -Xmx200m -Dcoherence.home=/oracle/fmw/ohs/wlserver/../coherence -Dbea.home=/oracle/fmw/ohs/wlserver/.. -Dweblogic.RootDirectory=/oracle/fmw/ohs/user_projects/domains/base_domain -Djava.system.class.loader=com.oracle.classloader.weblogic.LaunchClassLoader -Djava.security.policy=/oracle/fmw/ohs/wlserver/server/lib/weblogic.policy -Dweblogic.nodemanager.JavaHome=/oracle/fmw/ohs/oracle_common/jdk/jre weblogic.NodeManager -v"
  },
]

Java Info with Running Processes:

[
  {
    "HostName": "wcp12cr2",
    "DynLibBinPath": "/oracle/db/ohome/jdk/jre/lib/amd64/server/libjvm.so",
    "JavaBinPath": "/oracle/db/ohome/jdk/jre/bin/java",
    "JavaCBinPath": "/oracle/db/ohome/jdk/bin/javac",
    "IsJDK": true,
    "JavaHome": "/oracle/db/ohome/jdk/jre",
    "JavaRuntimeName": "Java(TM) SE Runtime Environment",
    "JavaRuntimeVersion": "1.6.0_37-b06",
    "JavaVendor": "",
    "JavaVersion": "1.6.0_37",
    "JavaVersionDate": "Java(TM) SE Runtime Environment (build 1.6.0_37-b06)",
    "JavaVMName": "Java HotSpot(TM) 64-Bit Server VM",
    "JavaVMVendor": "",
    "JavaVMVersion": "",
    "ProcessRunning": false,
    "ProcessPath": "",
    "CommandLine": ""
  },
  {
    "HostName": "wcp12cr2",
    "DynLibBinPath": "/usr/lib/jvm/java-1.5.0-gcj-1.5.0.0/jre/lib/x86_64/server/libjvm.so",
    "JavaBinPath": "/usr/lib/jvm/java-1.5.0-gcj-1.5.0.0/jre/bin/java",
    "JavaCBinPath": "",
    "IsJDK": false,
    "JavaHome": "",
    "JavaRuntimeName": "",
    "JavaRuntimeVersion": "",
    "JavaVendor": "",
    "JavaVersion": "1.5.0",
    "JavaVersionDate": "gij (GNU libgcj) version 4.4.7 20120313 (Red Hat 4.4.7-16)",
    "JavaVMName": "",
    "JavaVMVendor": "",
    "JavaVMVersion": "",
    "ProcessRunning": false,
    "ProcessPath": "",
    "CommandLine": ""
  },
  {
    "HostName": "wcp12cr2",
    "DynLibBinPath": "/usr/lib/jvm/java-1.6.0-openjdk-1.6.0.0.x86_64/jre/lib/amd64/server/libjvm.so",
    "JavaBinPath": "/usr/lib/jvm/java-1.6.0-openjdk-1.6.0.0.x86_64/jre/bin/java",
    "JavaCBinPath": "/usr/lib/jvm/java-1.6.0-openjdk-1.6.0.0.x86_64/bin/javac",
    "IsJDK": true,
    "JavaHome": "/usr/lib/jvm/java-1.6.0-openjdk-1.6.0.0.x86_64/jre",
    "JavaRuntimeName": "OpenJDK Runtime Environment",
    "JavaRuntimeVersion": "",
    "JavaVendor": "",
    "JavaVersion": "1.6.0_24",
    "JavaVersionDate": "OpenJDK Runtime Environment (IcedTea6 1.11.5) (rhel-1.50.1.11.5.0.1.el6_3-x86_64)",
    "JavaVMName": "OpenJDK 64-Bit Server VM",
    "JavaVMVendor": "",
    "JavaVMVersion": "",
    "ProcessRunning": false,
    "ProcessPath": "",
    "CommandLine": ""
  },
  {
    "HostName": "wcp12cr2",
    "DynLibBinPath": "/usr/lib/jvm/java-1.7.0-openjdk-1.7.0.9.x86_64/jre/lib/amd64/server/libjvm.so",
    "JavaBinPath": "/usr/lib/jvm/java-1.7.0-openjdk-1.7.0.9.x86_64/jre/bin/java",
    "JavaCBinPath": "/usr/lib/jvm/java-1.7.0-openjdk-1.7.0.9.x86_64/bin/javac",
    "IsJDK": true,
    "JavaHome": "/usr/lib/jvm/java-1.7.0-openjdk-1.7.0.9.x86_64/jre",
    "JavaRuntimeName": "OpenJDK Runtime Environment",
    "JavaRuntimeVersion": "1.7.0_09-icedtea-mockbuild_2013_01_16_11_20-b00",
    "JavaVendor": "Oracle Corporation",
    "JavaVersion": "1.7.0_09-icedtea",
    "JavaVersionDate": "OpenJDK Runtime Environment (rhel-2.3.4.1.0.1.el6_3-x86_64)",
    "JavaVMName": "OpenJDK 64-Bit Server VM",
    "JavaVMVendor": "Oracle Corporation",
    "JavaVMVersion": "23.2-b09",
    "ProcessRunning": false,
    "ProcessPath": "",
    "CommandLine": ""
  }
]
Creating csv report...
Done!
