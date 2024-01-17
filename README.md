# Licenseware ODJM Collector

## Introduction
This program is designed to search for Java installations on the machine where it's running, collect information about the machine and the Java installations and generate a CSV report. It uses the jinfo and jps binaries found in most Java JDKs to connect to the JVM and collect data on running processes. If the binaries are not found, only information on installed Java instances is collected. There is no dependency to a certain JDK vendor or version to run this program.

## How to download
Go to `Releases` page on the right side bar
![Releases Page](image.png)

Then using `Assets` select the right package for your operating system and cpu architecture,
naming convention is `ojdm-collector-version_os_architecture`

![Assets](image-1.png)
## How to use
On Windows:
`ojdm-collector.exe`

On Linux / MacOS:
`./ojdm-collector`

Running the tool on Linux may require to first make it executable by running 
`chmod +x ojdm-collector`

The report.csv file will be generated in the location from which the program was executed

### Optional arguments
    -output-path string
            Optional: Path to csv report. (default "report.csv")
    
    -search-paths string
            Optional: List of paths separated by comma where to search for java info.

    $ ojdm-collector -output-path=/path/to/csvreport.csv
    $ ojdm-collector -search-paths=/home,/oracle,/opt
    $ ojdm-collector -search-paths=/home,/usr,/opt -output-path=/path/to/csvreport.csv

## Collected data
| HostName | DynLibBinPath                                                       | JavaBinPath                                                   | JavaCBinPath                              | IsJDK | JavaHome                                         | JavaRuntimeName                 | JavaRuntimeVersion | JavaVendor         | JavaVersion | JavaVersionDate | JavaVMName                        | JavaVMVendor       | JavaVMVersion    | ProcessPath                                                                                                                       | ProcessRunning | CommandLine                                                                         | HostLogicalProcessors |
|----------|---------------------------------------------------------------------|---------------------------------------------------------------|-------------------------------------------|-------|--------------------------------------------------|---------------------------------|--------------------|--------------------|-------------|-----------------|-----------------------------------|--------------------|------------------|-----------------------------------------------------------------------------------------------------------------------------------|----------------|-------------------------------------------------------------------------------------|-----------------------|
| Reactor1 | C:\Program Files\Java\jre-9.0.4\bin\server\jvm.dll                  | C:\Program Files\Java\jre-9.0.4\bin\java.exe                  |                                           | false | C:/Program Files/Java/jre-9.0.4                  | Java(TM) SE Runtime Environment | 9.0.4+11           | Oracle Corporation | 9.0.4       |                 | Java HotSpot(TM) 64-Bit Server VM | Oracle Corporation | 9.0.4+11         | C:/Users/Administrator/Documents/ojdm-collector/apache-tinkerpop-gremlin-console-3.7.1-bin/apache-tinkerpop-gremlin-console-3.7.1 | true           | org.apache.tinkerpop.gremlin.console.Console -Xms32m -Xmx512m -Djline.terminal=none | 24                    |
| Reactor1 | C:\Program Files\Java\jdk-9.0.4\bin\server\jvm.dll                  | C:\Program Files\Java\jdk-9.0.4\bin\java.exe                  | C:/Program Files/Java/jdk-9.0.4/bin/javac | true  | C:/Program Files/Java/jdk-9.0.4                  | Java(TM) SE Runtime Environment | 9.0.4+11           | Oracle Corporation | 9.0.4       |                 | Java HotSpot(TM) 64-Bit Server VM | Oracle Corporation | 9.0.4+11         |                                                                                                                                   | false          |                                                                                     | 24                    |
| Reactor1 | C:\Program Files\Java\jdk-21\bin\server\jvm.dll                     | C:\Program Files\Java\jdk-21\bin\java.exe                     | C:/Program Files/Java/jdk-21/bin/javac    | true  | C:/Program Files/Java/jdk-21                     | Java(TM) SE Runtime Environment | 21.0.1+12-LTS-29   | Oracle Corporation | 21.0.1      | 2023-10-17      | Java HotSpot(TM) 64-Bit Server VM | Oracle Corporation | 21.0.1+12-LTS-29 |                                                                                                                                   | false          |                                                                                     | 24                    |
| Reactor1 | C:\Program Files\Zulu\zulu-21\bin\server\jvm.dll                    | C:\Program Files\Zulu\zulu-21\bin\java.exe                    | C:/Program Files/Zulu/zulu-21/bin/javac   | true  | C:/Program Files/Zulu/zulu-21                    | OpenJDK Runtime Environment     | 21.0.1+12-LTS      | Azul Systems, Inc. | 21.0.1      | 2023-10-17      | OpenJDK 64-Bit Server VM          | Azul Systems, Inc. | 21.0.1+12-LTS    |                                                                                                                                   | false          |                                                                                     | 24                    |
| Reactor1 | C:\Program Files\Java\jdk-19\bin\server\jvm.dll                     | C:\Program Files\Java\jdk-19\bin\java.exe                     | C:/Program Files/Java/jdk-19/bin/javac    | true  | C:/Program Files/Java/jdk-19                     | Java(TM) SE Runtime Environment | 19.0.2+7-44        | Oracle Corporation | 19.0.2      | 2023-01-17      | Java HotSpot(TM) 64-Bit Server VM | Oracle Corporation | 19.0.2+7-44      |                                                                                                                                   | false          |                                                                                     | 24                    |
| Reactor1 | C:\Users\Administrator\AppData\Local\DBeaver\jre\bin\server\jvm.dll | C:\Users\Administrator\AppData\Local\DBeaver\jre\bin\java.exe |                                           | false | C:/Users/Administrator/AppData/Local/DBeaver/jre | OpenJDK Runtime Environment     | 17.0.6+10          | Eclipse Adoptium   | 17.0.6      | 2023-01-17      | OpenJDK 64-Bit Server VM          | Eclipse Adoptium   | 17.0.6+10        |                                                                                                                                   | false          |                                                                                     | 24                    |
| Reactor1 | C:\Program Files\Java\jdk1.8.0_202\jre\bin\server\jvm.dll           | C:\Program Files\Java\jdk1.8.0_202\jre\bin\java.exe           |                                           | false | C:/Program Files/Java/jdk1.8.0_202/jre           | Java(TM) SE Runtime Environment | 1.8.0_202-b08      | Oracle Corporation | 1.8.0_202   |                 | Java HotSpot(TM) 64-Bit Server VM | Oracle Corporation | 25.202-b08       |                                                                                                                                   | false          |                                                                                     | 24                    |
| Reactor1 | C:\Program Files\Java\jre1.8.0_202\bin\server\jvm.dll               | C:\Program Files\Java\jre1.8.0_202\bin\java.exe               |                                           | false | C:/Program Files/Java/jre1.8.0_202               | Java(TM) SE Runtime Environment | 1.8.0_202-b08      | Oracle Corporation | 1.8.0_202   |                 | Java HotSpot(TM) 64-Bit Server VM | Oracle Corporation | 25.202-b08       |                                                                                                                                   | false          |                                                                                     | 24                    |

## Searched paths
By default the program will search for Java installations in several specific locations depending on the operating system. Additional paths can be provided by using the -search-paths parameter when runnning the program. 

On Windows: 
* LocalAppData (from environment variables)
* C:\\Program Files
* C:\\Program Files (x86)

On Linux:
* /home
* /usr/bin
* /usr/local
* /usr/lib
* /usr/share
* /opt
* /snap
* /oracle
* /bin

On MacOs:
* /Applications

## Troubleshooting
If no running processes are identified, it may be because the jinfo and jps utilities could not be found on any of the discovered java installations. The easiest way to fix this is to place an OpenJDK in any of the default search paths or to include the location of the OpenJDK in the additional search paths. 
