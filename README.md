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


Current:

[
  {
    "HostName": "wcp12cr2",
    "AppDirName": "jre",
    "DynLibBinPath": "/oracle/fmw/ohs/oracle_common/jdk/jre/lib/amd64/server/libjvm.so",
    "JavaBinPath": "/oracle/fmw/ohs/oracle_common/jdk/jre/bin/java",
    "JavaCBinPath": "",
    "IsJDK": false,
    "BaseDir": "/oracle/fmw/ohs/oracle_common/jdk/jre",
    "JavaHome": "/oracle/fmw/ohs/oracle_common/jdk/jre",
    "JavaRuntimeName": "Java(TM) SE Runtime Environment",
    "JavaRuntimeVersion": "1.8.0_51-b16",
    "JavaVendor": "Oracle Corporation",
    "JavaVersion": "1.8.0_51",
    "JavaVersionDate": "Java(TM) SE Runtime Environment (build 1.8.0_51-b16)",
    "JavaVMName": "Java HotSpot(TM) 64-Bit Server VM",
    "JavaVMVendor": "Oracle Corporation",
    "JavaVMVersion": "25.51-b03",
    "ProcessRunning": true,
    "ProcessPath": "/oracle/fmw/ohs/oracle_common/jdk/jre/bin/java",
    "CommandLine": "/oracle/fmw/ohs/oracle_common/jdk/jre/bin/java -server -Xms32m -Xmx200m -Dcoherence.home=/oracle/fmw/ohs/wlserver/../coherence -Dbea.home=/oracle/fmw/ohs/wlserver/.. -Dweblogic.RootDirectory=/oracle/fmw/ohs/user_projects/domains/base_domain -Djava.system.class.loader=com.oracle.classloader.weblogic.LaunchClassLoader -Djava.security.policy=/oracle/fmw/ohs/wlserver/server/lib/weblogic.policy -Dweblogic.nodemanager.JavaHome=/oracle/fmw/ohs/oracle_common/jdk/jre weblogic.NodeManager -v"
  },
  {
    "HostName": "wcp12cr2",
    "AppDirName": "jre",
    "DynLibBinPath": "/usr/java/jdk1.8.0_60/jre/lib/amd64/server/libjvm.so",
    "JavaBinPath": "/usr/java/jdk1.8.0_60/jre/bin/java",
    "JavaCBinPath": "",
    "IsJDK": false,
    "BaseDir": "/usr/java/jdk1.8.0_60/jre",
    "JavaHome": "/usr/java/jdk1.8.0_60/jre",
    "JavaRuntimeName": "Java(TM) SE Runtime Environment",
    "JavaRuntimeVersion": "1.8.0_60-b27",
    "JavaVendor": "Oracle Corporation",
    "JavaVersion": "1.8.0_60",
    "JavaVersionDate": "Java(TM) SE Runtime Environment (build 1.8.0_60-b27)",
    "JavaVMName": "Java HotSpot(TM) 64-Bit Server VM",
    "JavaVMVendor": "Oracle Corporation",
    "JavaVMVersion": "25.60-b23",
    "ProcessRunning": false,
    "ProcessPath": "",
    "CommandLine": ""
  },
  {
    "HostName": "wcp12cr2",
    "AppDirName": "jre",
    "DynLibBinPath": "/usr/lib/jvm/java-1.5.0-gcj-1.5.0.0/jre/lib/x86_64/server/libjvm.so",
    "JavaBinPath": "/usr/lib/jvm/java-1.5.0-gcj-1.5.0.0/jre/bin/java",
    "JavaCBinPath": "",
    "IsJDK": false,
    "BaseDir": "/usr/lib/jvm/java-1.5.0-gcj-1.5.0.0/jre",
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
    "AppDirName": "jre",
    "DynLibBinPath": "/usr/lib/jvm/java-1.6.0-openjdk-1.6.0.0.x86_64/jre/lib/amd64/server/libjvm.so",
    "JavaBinPath": "/usr/lib/jvm/java-1.6.0-openjdk-1.6.0.0.x86_64/jre/bin/java",
    "JavaCBinPath": "",
    "IsJDK": false,
    "BaseDir": "/usr/lib/jvm/java-1.6.0-openjdk-1.6.0.0.x86_64/jre",
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
    "AppDirName": "jre",
    "DynLibBinPath": "/usr/lib/jvm/java-1.7.0-openjdk-1.7.0.9.x86_64/jre/lib/amd64/server/libjvm.so",
    "JavaBinPath": "/usr/lib/jvm/java-1.7.0-openjdk-1.7.0.9.x86_64/jre/bin/java",
    "JavaCBinPath": "",
    "IsJDK": false,
    "BaseDir": "/usr/lib/jvm/java-1.7.0-openjdk-1.7.0.9.x86_64/jre",
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
  },
  {
    "HostName": "wcp12cr2",
    "AppDirName": "jre",
    "DynLibBinPath": "/oracle/fmw/ohs/oracle_common/jdk/jre/lib/amd64/server/libjvm.so",
    "JavaBinPath": "/oracle/fmw/ohs/oracle_common/jdk/jre/bin/java",
    "JavaCBinPath": "",
    "IsJDK": false,
    "BaseDir": "/oracle/fmw/ohs/oracle_common/jdk/jre",
    "JavaHome": "/oracle/fmw/ohs/oracle_common/jdk/jre",
    "JavaRuntimeName": "Java(TM) SE Runtime Environment",
    "JavaRuntimeVersion": "1.8.0_51-b16",
    "JavaVendor": "Oracle Corporation",
    "JavaVersion": "1.8.0_51",
    "JavaVersionDate": "Java(TM) SE Runtime Environment (build 1.8.0_51-b16)",
    "JavaVMName": "Java HotSpot(TM) 64-Bit Server VM",
    "JavaVMVendor": "Oracle Corporation",
    "JavaVMVersion": "25.51-b03",
    "ProcessRunning": true,
    "ProcessPath": "/oracle/fmw/ohs/oracle_common/jdk/jre/bin/java",
    "CommandLine": "/oracle/fmw/ohs/oracle_common/jdk/jre/bin/java -server -Xms32m -Xmx200m -Dcoherence.home=/oracle/fmw/ohs/wlserver/../coherence -Dbea.home=/oracle/fmw/ohs/wlserver/.. -Dweblogic.RootDirectory=/oracle/fmw/ohs/user_projects/domains/base_domain -Djava.system.class.loader=com.oracle.classloader.weblogic.LaunchClassLoader -Djava.security.policy=/oracle/fmw/ohs/wlserver/server/lib/weblogic.policy -Dweblogic.nodemanager.JavaHome=/oracle/fmw/ohs/oracle_common/jdk/jre weblogic.NodeManager -v"
  },
  {
    "HostName": "wcp12cr2",
    "AppDirName": "jre",
    "DynLibBinPath": "/oracle/db/ohome/jdk/jre/lib/amd64/server/libjvm.so",
    "JavaBinPath": "/oracle/db/ohome/jdk/jre/bin/java",
    "JavaCBinPath": "",
    "IsJDK": false,
    "BaseDir": "/oracle/db/ohome/jdk/jre",
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
  }
]


If we gather all files


javaPaths := []string{

  // From MyPC
  
  "/home/alin/Documents/android-studio/jbr/bin/java",
  "/home/alin/Documents/android-studio/jbr/bin/javac",
  "/home/alin/Documents/android-studio/jbr/lib/server/libjvm.so", 

  "/usr/bin/java", 
  "/usr/bin/javac", 

  "/usr/lib/jvm/java-19-openjdk-amd64/bin/java", 
  "/usr/lib/jvm/java-19-openjdk-amd64/bin/javac", 
  "/usr/lib/jvm/java-19-openjdk-amd64/lib/server/libjvm.so", 

  "/usr/share/bash-completion/completions/java", 
  "/usr/share/bash-completion/completions/javac",

  "/snap/core/16091/etc/apparmor.d/abstractions/ubuntu-browsers.d/java" ,
  "/snap/core/16091/usr/share/bash-completion/completions/java",
  "/snap/core/16091/usr/share/bash-completion/completions/javac" ,
  "/snap/core/16202/etc/apparmor.d/abstractions/ubuntu-browsers.d/java" ,
  "/snap/core/16202/usr/share/bash-completion/completions/java" ,
  "/snap/core/16202/usr/share/bash-completion/completions/javac" ,
  "/snap/core18/2790/etc/apparmor.d/abstractions/ubuntu-browsers.d/java" ,
  "/snap/core18/2790/usr/share/bash-completion/completions/java",
  "/snap/core18/2790/usr/share/bash-completion/completions/javac" ,
  "/snap/core18/2796/etc/apparmor.d/abstractions/ubuntu-browsers.d/java" ,
  "/snap/core18/2796/usr/share/bash-completion/completions/java" ,
  "/snap/core18/2796/usr/share/bash-completion/completions/javac" ,
  "/snap/core20/1974/etc/apparmor.d/abstractions/ubuntu-browsers.d/java" ,
  "/snap/core20/1974/usr/share/bash-completion/completions/java" ,
  "/snap/core20/1974/usr/share/bash-completion/completions/javac" ,
  "/snap/core20/2015/etc/apparmor.d/abstractions/ubuntu-browsers.d/java" ,
  "/snap/core20/2015/usr/share/bash-completion/completions/java" ,
  "/snap/core20/2015/usr/share/bash-completion/completions/javac" ,
  "/snap/core22/858/etc/apparmor.d/abstractions/ubuntu-browsers.d/java" ,
  "/snap/core22/858/usr/share/bash-completion/completions/java" ,
  "/snap/core22/858/usr/share/bash-completion/completions/javac" ,
  "/snap/core22/864/etc/apparmor.d/abstractions/ubuntu-browsers.d/java" ,
  "/snap/core22/864/usr/share/bash-completion/completions/java" ,
  "/snap/core22/864/usr/share/bash-completion/completions/javac" ,

  "/snap/dbeaver-ce/268/usr/share/dbeaver-ce/jre/bin/java" ,
  "/snap/dbeaver-ce/268/usr/share/dbeaver-ce/jre/lib/server/libjvm.so" ,

  "/snap/dbeaver-ce/270/usr/share/dbeaver-ce/jre/bin/java" ,
  "/snap/dbeaver-ce/270/usr/share/dbeaver-ce/jre/lib/server/libjvm.so" ,

  "/snap/snapd/20092/usr/lib/snapd/apparmor.d/abstractions/ubuntu-browsers.d/java" ,
  "/snap/snapd/20290/usr/lib/snapd/apparmor.d/abstractions/ubuntu-browsers.d/java",

  // From Oracle WebCenter Portal 12c R2

  "/var/lib/alternatives/javac",
  "/usr/java/jdk1.8.0_60/bin/javac",
  "/usr/lib/jvm/java-1.6.0-openjdk-1.6.0.0.x86_64/bin/javac",
  "/usr/lib/jvm/java-1.7.0-openjdk-1.7.0.9.x86_64/bin/javac",
  "/oracle/db/ohome/jdk/bin/javac",


  "/usr/java/jdk1.8.0_60/jre/lib/amd64/server/libjvm.so",
  "/usr/lib/jvm/java-1.6.0-openjdk-1.6.0.0.x86_64/jre/lib/amd64/server/libjvm.so",
  "/usr/lib/jvm/java-1.7.0-openjdk-1.7.0.9.x86_64/jre/lib/amd64/server/libjvm.so",
  "/usr/lib64/gcj-4.4.4/libjvm.so",
  "/oracle/fmw/ohs/oracle_common/jdk/jre/lib/amd64/server/libjvm.so",
  "/oracle/db/ohome/jdk/jre/lib/amd64/server/libjvm.so",

  "/var/lib/alternatives/java",
  "/usr/java/jdk1.8.0_60/jre/bin/java",
  "/usr/java/jdk1.8.0_60/bin/java",
  "/usr/lib/jvm/java-1.6.0-openjdk-1.6.0.0.x86_64/jre/bin/java",
  "/usr/lib/jvm/java-1.6.0-openjdk-1.6.0.0.x86_64/bin/java",
  "/usr/lib/jvm/java-1.7.0-openjdk-1.7.0.9.x86_64/jre/bin/java",
  "/usr/lib/jvm/java-1.7.0-openjdk-1.7.0.9.x86_64/bin/java",
  "/oracle/fmw/ohs/oracle_common/jdk/jre/bin/java",
  "/oracle/db/ohome/jdk/jre/bin/java",
  "/oracle/db/ohome/jdk/bin/java",

}


If not windows gather only java/javac files in the bin/server folder

[

  /home/alin/Documents/android-studio/jbr/bin/java 
  /home/alin/Documents/android-studio/jbr/bin/javac 
  /home/alin/Documents/android-studio/jbr/lib/server/libjvm.so 

  /usr/bin/java 
  /usr/bin/javac 

  /usr/lib/jvm/java-19-openjdk-amd64/bin/java 
  /usr/lib/jvm/java-19-openjdk-amd64/bin/javac 
  /usr/lib/jvm/java-19-openjdk-amd64/lib/server/libjvm.so 

  /snap/dbeaver-ce/268/usr/share/dbeaver-ce/jre/bin/java 
  /snap/dbeaver-ce/268/usr/share/dbeaver-ce/jre/lib/server/libjvm.so 

  /snap/dbeaver-ce/270/usr/share/dbeaver-ce/jre/bin/java 
  /snap/dbeaver-ce/270/usr/share/dbeaver-ce/jre/lib/server/libjvm.so

  /oracle/fmw/ohs/oracle_common/jdk/jre/lib/amd64/server/libjvm.so
  /oracle/fmw/ohs/oracle_common/jdk/jre/bin/java

  /usr/lib/jvm/java-1.6.0-openjdk-1.6.0.0.x86_64/jre/lib/amd64/server/libjvm.so
  /usr/lib/jvm/java-1.6.0-openjdk-1.6.0.0.x86_64/jre/bin/java
  /usr/lib/jvm/java-1.6.0-openjdk-1.6.0.0.x86_64/bin/javac


]


Get java/javac binary paths based on libjvm.so path

[
  /home/alin/Documents/android-studio/jbr/lib/server/libjvm.so 
  /usr/lib/jvm/java-19-openjdk-amd64/lib/server/libjvm.so 
  /snap/dbeaver-ce/268/usr/share/dbeaver-ce/jre/lib/server/libjvm.so 
  /snap/dbeaver-ce/270/usr/share/dbeaver-ce/jre/lib/server/libjvm.so  
]


jbasePath /home/alin/Documents/android-studio/jbr
javaBinPath /home/alin/Documents/android-studio/jbr/bin/java
javaCBinPath /home/alin/Documents/android-studio/jbr/bin/javac


jbasePath /usr/lib/jvm/java-19-openjdk-amd64
javaBinPath /usr/lib/jvm/java-19-openjdk-amd64/bin/java
javaCBinPath /usr/lib/jvm/java-19-openjdk-amd64/bin/javac


jbasePath /snap/dbeaver-ce/268/usr/share/dbeaver-ce/jre
javaBinPath /snap/dbeaver-ce/268/usr/share/dbeaver-ce/jre/bin/java
javaCBinPath 


jbasePath /snap/dbeaver-ce/270/usr/share/dbeaver-ce/jre
javaBinPath /snap/dbeaver-ce/270/usr/share/dbeaver-ce/jre/bin/java
javaCBinPath 


On the Oracle VM


jbasePath /usr/java/jdk1.8.0_60/jre
javaBinPath /usr/java/jdk1.8.0_60/bin/java
javaCBinPath /usr/java/jdk1.8.0_60/bin/javac


jbasePath /usr/lib/jvm/java-1.6.0-openjdk-1.6.0.0.x86_64/jre
javaBinPath /usr/lib/jvm/java-1.6.0-openjdk-1.6.0.0.x86_64/bin/java
javaCBinPath /usr/lib/jvm/java-1.6.0-openjdk-1.6.0.0.x86_64/bin/javac


jbasePath /usr/lib/jvm/java-1.7.0-openjdk-1.7.0.9.x86_64/jre
javaBinPath /usr/lib/jvm/java-1.7.0-openjdk-1.7.0.9.x86_64/bin/java
javaCBinPath /usr/lib/jvm/java-1.7.0-openjdk-1.7.0.9.x86_64/bin/javac


jbasePath /
javaBinPath /bin/java
javaCBinPath /bin/javac


jbasePath /oracle/fmw/ohs/oracle_common/jdk/jre
javaBinPath /oracle/fmw/ohs/oracle_common/jdk/bin/java
javaCBinPath /oracle/fmw/ohs/oracle_common/jdk/bin/javac


jbasePath /oracle/db/ohome/jdk/jre
javaBinPath /oracle/db/ohome/jdk/bin/java
javaCBinPath /oracle/db/ohome/jdk/bin/javac