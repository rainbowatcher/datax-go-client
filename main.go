package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

var (
	dataxHome           string
	classpath           string
	logbackFile         string
	defaultJVM          string
	defaultPropertyConf string
	childProcess        *os.Process
	options             Options
	RET_STATE           = map[string]int{
		"KILL":  143,
		"FAIL":  -1,
		"OK":    0,
		"RUN":   1,
		"RETRY": 2,
	}
)

const (
	DataxVersion      = "DATAX-OPENSOURCE-3.0"
	RemoteDebugConfig = "-Xdebug -Xrunjdwp:transport=dt_socket,server=y,address=9999"
)

type Options struct {
	jvmParameters string
	jobid         string
	mode          string
	params        string
	remoteDebug   bool
	loglevel      string
	reader        string
	writer        string
}

func init() {
	// 获取当前文件的绝对路径
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	// 获取当前文件所在的目录
	exeDir := filepath.Dir(exePath)
	dataxHome = filepath.Dir(exeDir)
	if isWindows() {
		cmd := exec.Command("CMD", "/C", "chcp", "65001")
		err := cmd.Run()
		if err != nil {
			fmt.Println("[Error]=> " + fmt.Sprint(err))
		}
		classpath = fmt.Sprintf("%s/lib/*", dataxHome)
	} else {
		classpath = dataxHome + "/lib/*:."
	}
	logbackFile = dataxHome + "/conf/logback.xml"
	defaultJVM = "-Xms1g -Xmx1g -XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=" + dataxHome + "/log"
	defaultPropertyConf = "-Dfile.encoding=UTF-8 -Dlogback.statusListenerClass=ch.qos.logback.core.status.NopStatusListener -Djava.security.egd=file:///dev/urandom -Ddatax.home=" +
		dataxHome + " -Dlogback.configurationFile=" + logbackFile
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func registerSignal() {
	if !isWindows() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

		go func() {
			sig := <-signals
			suicide(sig)
		}()
	}
}

func suicide(signum os.Signal) {
	fmt.Fprintf(os.Stderr, "[Error] DataX receive unexpected signal %d, starts to suicide.\n", signum)

	if childProcess != nil {
		childProcess.Signal(syscall.SIGQUIT)
		time.Sleep(1 * time.Second)
		childProcess.Kill()
	}
	fmt.Fprintln(os.Stderr, "DataX Process was killed ! you did ?")
	os.Exit(1)
}

func getLocalIp() []string {
	var ipStrs []string
	ips, err := net.LookupIP(getOutboundIP())
	if err != nil {
		panic(err)
	}
	for _, ip := range ips {
		ipStrs = append(ipStrs, ip.String())
	}
	return ipStrs
}

func isUrl(path string) bool {
	if path == "" {
		return false
	}
	match, _ := regexp.MatchString("^http[s]?://\\S+\\w*", strings.ToLower(path))
	return match
}

func generateJobConfigTemplate(reader string, writer string) {
	readerRef := fmt.Sprintf("Please refer to the %s document:\n     https://github.com/alibaba/DataX/blob/master/%s/doc/%s.md \n", reader, reader, reader)
	writerRef := fmt.Sprintf("Please refer to the %s document:\n     https://github.com/alibaba/DataX/blob/master/%s/doc/%s.md \n ", writer, writer, writer)
	fmt.Println(readerRef)
	fmt.Println(writerRef)
	jobGuid := "Please save the following configuration as a json file and  use\n     python {DATAX_HOME}/bin/datax.py {JSON_FILE_NAME}.json \nto run the job.\n"
	fmt.Println(jobGuid)

	// Define job template
	jobTemplate := map[string]any{
		"job": map[string]any{
			"setting": map[string]any{
				"speed": map[string]any{
					"channel": "",
				},
			},
			"content": []map[string]any{
				{
					"reader": map[string]any{},
					"writer": map[string]any{},
				},
			},
		},
	}

	// Set reader and writer templates
	readerTemplatePath := fmt.Sprintf("%s/plugin/reader/%s/plugin_job_template.json", dataxHome, reader)
	writerTemplatePath := fmt.Sprintf("%s/plugin/writer/%s/plugin_job_template.json", dataxHome, writer)
	readerPar, err := readPluginTemplate(readerTemplatePath)
	if err != nil {
		fmt.Printf("Read reader[%s] template error: can't find file %s\n", reader, readerTemplatePath)
	}
	writerPar, err := readPluginTemplate(writerTemplatePath)
	if err != nil {
		fmt.Printf("Read writer[%s] template error: can't find file %s\n", writer, writerTemplatePath)
	}

	jobTemplate["job"].(map[string]any)["content"].([]map[string]any)[0]["reader"] = readerPar
	jobTemplate["job"].(map[string]any)["content"].([]map[string]any)[0]["writer"] = writerPar

	// Print the job template as JSON
	jobJSON, _ := json.MarshalIndent(jobTemplate, "", "    ")
	fmt.Println(string(jobJSON))
}

func readPluginTemplate(plugin string) (map[string]any, error) {
	// Open plugin file for reading
	f, err := os.Open(plugin)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Load JSON data from file
	data := make(map[string]any)
	err = json.NewDecoder(f).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func parseArgs() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `Usage: %s [options] job-url-or-path

Options:
  -h, --help            show this help message and exit

  Product Env Options:
    Normal user use these options to set jvm parameters, job runtime mode
    etc. Make sure these options can be used in Product Env.

    -j, --jvm=<jvm parameters>
                        Set jvm parameters if necessary.
    --jobid=<job unique id>
                        Set job unique id when running by Distribute/Local
                        Mode.
    -m, --mode=<job runtime mode>
                        Set job runtime mode such as: standalone, local,
                        distribute. Default mode is standalone.
    -p, --params=<parameter used in job config>
                        Set job parameter, eg: the source tableName you want
                        to set it by command, then you can use like this:
                        -p"-DtableName=your-table-name", if you have mutiple
                        parameters: -p"-DtableName=your-table-name
                        -DcolumnName=your-column-name".Note: you should config
                        in you job tableName with ${tableName}.
    -r, --reader=<parameter used in view job config[reader] template>
                        View job config[reader] template, eg:
                        mysqlreader,streamreader
    -w, --writer=<parameter used in view job config[writer] template>
                        View job config[writer] template, eg:
                        mysqlwriter,streamwriter

  Develop/Debug Options:
    Developer use these options to trace more details of DataX.

    -d, --debug         Set to remote debug mode.
    --loglevel=<log level>
                        Set log level such as: debug, info, all etc.%s`, os.Args[0], "\n")
	}
	flag.StringVar(&options.jvmParameters, "jvm", defaultJVM, "Set jvm parameters if necessary.")
	flag.StringVar(&options.jvmParameters, "j", defaultJVM, "Set jvm parameters if necessary.")
	flag.StringVar(&options.jobid, "jobid", "-1", "Set job unique id when running by Distribute/Local Mode.")
	flag.StringVar(&options.mode, "mode", "standalone", "Set job runtime mode such as: standalone, local, distribute. Default mode is standalone.")
	flag.StringVar(&options.mode, "m", "standalone", "Set job runtime mode such as: standalone, local, distribute. Default mode is standalone.")
	flag.StringVar(&options.params, "params", "", "Set job parameter, eg: the source tableName you want to set it by command, then you can use like this: -p\"-DtableName=your-table-name\", if you have mutiple parameters: -p\"-DtableName=your-table-name -DcolumnName=your-column-name\".Note: you should config in you job tableName with ${tableName}.")
	flag.StringVar(&options.params, "p", "", "Set job parameter, eg: the source tableName you want to set it by command, then you can use like this: -p\"-DtableName=your-table-name\", if you have mutiple parameters: -p\"-DtableName=your-table-name -DcolumnName=your-column-name\".Note: you should config in you job tableName with ${tableName}.")
	flag.StringVar(&options.reader, "reader", "", "")
	flag.StringVar(&options.reader, "r", "", "")
	flag.StringVar(&options.writer, "writer", "", "")
	flag.StringVar(&options.writer, "w", "", "")
	flag.BoolVar(&options.remoteDebug, "debug", false, "Set to remote debug mode.")
	flag.StringVar(&options.loglevel, "loglevel", "info", "Set log level such as: debug, info, all etc.")
	flag.Parse()
}

func buildStartCommand(options Options, args []string) string {
	commandMap := make(map[string]string)
	var tempJVMCommand string
	if options.jvmParameters != "" {
		tempJVMCommand = tempJVMCommand + " " + options.jvmParameters
	}

	if options.remoteDebug {
		tempJVMCommand = tempJVMCommand + " " + RemoteDebugConfig
		fmt.Println("local ip: ", getLocalIp())
	}

	if options.loglevel != "" {
		tempJVMCommand = tempJVMCommand + " " + fmt.Sprintf("-Dloglevel=%s", options.loglevel)
	}

	// jobResource 可能是 URL，也可能是本地文件路径（相对,绝对）
	jobResource := args[0]
	if !isUrl(jobResource) {
		jobResource, _ = os.Getwd()
		jobResource = filepath.Join(jobResource, args[0])
		if strings.HasPrefix(strings.ToLower(jobResource), "file://") {
			jobResource = jobResource[len("file://"):]
		}
	}

	jobParams := fmt.Sprintf("-Dlog.file.name=%s",
		strings.ReplaceAll(strings.ReplaceAll(jobResource, "/", "_"), ".", "_")[:20])
	if options.params != "" {
		jobParams = jobParams + " " + options.params
	}

	if options.jobid != "" {
		commandMap["jobid"] = options.jobid
	}

	commandMap["jvm"] = tempJVMCommand
	commandMap["params"] = jobParams
	commandMap["job"] = jobResource

	err := godotenv.Load()
	if err != nil {
		fmt.Println(err.Error())
	}

	javaHome := os.Getenv("JAVA_HOME")
	if javaHome == "" {
		javaHome = "java"
	} else {
		javaHome = javaHome + "/bin/java"
	}
	return fmt.Sprintf("%s -server %s %s -classpath %s %s com.alibaba.datax.core.Engine -mode %s -jobid %s -job %s",
		javaHome,
		tempJVMCommand,
		defaultPropertyConf,
		classpath,
		jobParams,
		options.mode,
		options.jobid,
		jobResource,
	)
}

func printCopyright() {
	fmt.Printf("\nDataX (%s), From Alibaba !\n", DataxVersion)
	fmt.Printf("Copyright (C) 2010-2017, Alibaba Group. All Rights Reserved.\n\n")
	os.Stdout.Sync()
}

func main() {
	printCopyright()
	parseArgs()
	// fmt.Println(options)
	if options.reader != "" && options.writer != "" {
		generateJobConfigTemplate(options.reader, options.writer)
		os.Exit(RET_STATE["OK"])
	}
	// fmt.Println(os.Args[len(os.Args)-1:])
	job := os.Args[len(os.Args)-1:]
	if len(job) != 1 {
		flag.Usage()
		os.Exit(RET_STATE["FAIL"])
	}
	startCommand := buildStartCommand(options, job)
	var cmd *exec.Cmd
	if isWindows() {
		cmd = exec.Command("cmd", "/C", startCommand)
	} else {
		cmd = exec.Command("sh", "-c", startCommand)
	}

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] Fail to start DataX: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Fprintln(os.Stderr, "DataX is finished !")
		os.Exit(0)
	}
}
