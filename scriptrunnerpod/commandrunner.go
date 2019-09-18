package main

import (
    "os"
    "os/exec"
    "fmt"
    "time"
    "strconv"
    "strings"

     scriptrunnerv1alpha1 "github.com/MoserMichael/scriptrunner/pkg/apis/scriptrunner/v1alpha1"

    //"context"
    //"sigs.k8s.io/controller-runtime/pkg/client/config"
    //"sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/patch"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/types"
    "k8s.io/apimachinery/pkg/runtime/schema"
    "k8s.io/client-go/dynamic"
    "k8s.io/client-go/rest"

)

type TaskConfig struct {
    //ScriptRunnerObjectName string
    NodeName string
    InstName string

    CmdHistorySize int32
    SerialNumOfRun int32

    PackagesToInstall string
    PipToInstall string
    PythonScript string

    InitWait int32
    RunPeriod int32
    NumRepetitions int
}


func newTaskConfig() (*TaskConfig, error) {

    var initialWait =  os.Getenv("INITIAL_WAIT")
    var runPeriod = os.Getenv("RUN_PERIOD")
    var numRepetitions = os.Getenv("NUM_REPETITIONS")
    var historySize = os.Getenv("HISTORY_SIZE")

    //var scriptRunnerObjectName = os.Getenv("EXTRACT_INFO_OBJECT_NAME") // unique id of instance that will get run results.
    var nodeName = os.Getenv("EXTRACT_INFO_NODE_NAME")
    var instName = os.Getenv("EXTRACT_INFO_NAME")

    var script = os.Getenv("PYTHON_SCRIPT")
    var err error
    var nval int64

    var packagesToInstall = os.Getenv("PACKAGES_TO_INSTALL")
    var pipToInstall = os.Getenv("PIP_TO_INSTALL")

    cfg := TaskConfig{}

    nval, err = strconv.ParseInt( initialWait, 10, 32 )
    if err != nil {
        return nil, fmt.Errorf("Wrong value of INITIAL_WAIT %s", err)
    }
    cfg.InitWait = int32(nval)

    nval, err =  strconv.ParseInt( numRepetitions, 10, 32 )
    if err != nil {
        return nil, fmt.Errorf("Wrong value of NUM_REPETITIONS %s", err)
    }
    cfg.NumRepetitions = int(nval)

    nval, err =  strconv.ParseInt( runPeriod, 10, 32 )
    if err != nil {
        return nil, fmt.Errorf("Wrong value of RUN_PERIOD %s", err)
    }
    cfg.RunPeriod = int32(nval)

    nval, err = strconv.ParseInt( historySize, 10, 32 )
    if err != nil {
        return nil, fmt.Errorf("Wrong value of HISTORY_SIZE %s", err)
    }
    cfg.CmdHistorySize = int32(nval)

    /*
    if scriptRunnerObjectName == "" {
        return nil, fmt.Errorf("missing EXTRACT_INFO_OBJECT_NAME env. var")
    }
    cfg.ScriptRunnerObjectName = scriptRunnerObjectName
    */
    if nodeName == "" {
        return nil, fmt.Errorf("missing EXTRACT_INFO_NODE_NAME env. var")
    }
    cfg.NodeName = nodeName

    if instName == "" {
        return nil, fmt.Errorf("missing EXTRACT_INFO_NAME env. var")
    }
    cfg.InstName = instName

    if script == "" {
        return nil, fmt.Errorf("missing PYTHON_SCRIPT env. var")
    }
    cfg.PythonScript = script
    cfg.SerialNumOfRun = 0
    cfg.PackagesToInstall = packagesToInstall
    cfg.PipToInstall = pipToInstall

    return &cfg, nil
}

// CmdResult result of running a textual command - split into lines
type CmdResult struct {
    // StdOut what command produced as standard output
    StdOut string
    // StdErr what command produced as standard error
    StdErr string
    // Status process exit code
    Status int32
    // TimeStart when process was created
    TimeStart time.Time
    // TimeEnd when proces finished
    TimeEnd   time.Time
    // SerialNum serial number of this run
    SerialNum int32
}

// NewCommand run a command and parse the output into lines (CmdResult)
func NewCommand(cmd []string) (*CmdResult, error) {
        res := CmdResult{}

        res.TimeStart = time.Now()

	clicmd := exec.Command(cmd[0], cmd[1:]...)

        output, err := clicmd.Output()

        res.TimeEnd = time.Now()

	if err != nil {
            exitError, ok := err.(*exec.ExitError)
            if ok {
               res.StdErr = string( exitError.Stderr )
               res.Status = int32( exitError.ExitCode() )
            }
        }

        res.StdOut = string(output)

        fmt.Println("Command:", cmd, "result:", res.StdOut)

        return &res, err
}

type CmdContext struct {
    cfg *TaskConfig
    //cl *client.Client
    cl *dynamic.NamespaceableResourceInterface
    emptyInstance *scriptrunnerv1alpha1.ScriptRunner
    currentInstance *scriptrunnerv1alpha1.ScriptRunner
}

//func makeContext(cfg *TaskConfig, cl *client.Client) (* CmdContext) {
func makeContext(cfg *TaskConfig, cl *dynamic.NamespaceableResourceInterface) (* CmdContext) {

    currentInstance := &scriptrunnerv1alpha1.ScriptRunner{}

    runnerStatus :=  &scriptrunnerv1alpha1.CommandRunnerStatus{}

    runnerStatus.CmdRunHistory =  make([]*scriptrunnerv1alpha1.CommandStatus, cfg.CmdHistorySize)

    currentInstance.Status.CommandStatus =  make(map[string]*scriptrunnerv1alpha1.CommandRunnerStatus)

    currentInstance.Status.CommandStatus[ cfg.NodeName ] = runnerStatus

    return  &CmdContext{ cfg, cl, &scriptrunnerv1alpha1.ScriptRunner{}, currentInstance }
}

func LogCommandResults( ctx* CmdContext, res *CmdResult )  {
    fmt.Println("Command result: status=", res.Status, " stdout=", res.StdOut, " stderr=", res.StdErr, " serialOfRun=", res.SerialNum, " timeStart=", res.TimeStart, " timeEnd=", res.TimeEnd)

    cfg := (*ctx).cfg

    curInstance := (*ctx).currentInstance
    commandStatus := &((*curInstance).Status)

    index := int( (*cfg).SerialNumOfRun % (*cfg).CmdHistorySize )
    (*commandStatus).CommandStatus[cfg.NodeName].CmdRunHistory[ index ] =  &scriptrunnerv1alpha1.CommandStatus{ int64( (*cfg).SerialNumOfRun), int(res.Status), res.StdOut, res.StdErr, res.TimeStart.String(), res.TimeEnd.String() }
    cfg.SerialNumOfRun ++

    cl := (*ctx).cl

    resCl, ok := (*cl).(dynamic.NamespaceableResourceInterface)
    if !ok {
        fmt.Println("not of type NamespaceableResourceInterface")
        return
    }

    // what to do with default namespace?
    resClD := resCl.Namespace("default")

    // Patch(name string, pt types.PatchType, data []byte, options metav1.UpdateOptions, subresources ...string) (*unstructured.Unstructured, error)
    patchObjs, errr := patch.NewJSONPatch( (*ctx).emptyInstance, (*ctx).currentInstance )
    if errr != nil {
        fmt.Println("Can't create json patch", errr)
        return
    }

    for _, po := range patchObjs {
        patchData, err :=  po.MarshalJSON()
        if err != nil {
            fmt.Println("Can't marshal json (2)", err)
            return
        }

        fmt.Println("attempt patch of name: ", (*cfg).InstName, " patch: ", string(patchData))
        //_, err = resCl.Patch( (*cfg).InstName, types.JSONPatchType, patchData, metav1.UpdateOptions{} )
        //_, err = resClD.Patch( (*cfg).InstName, types.JSONPatchType, patchData, metav1.UpdateOptions{} )
        _, err = resClD.Patch( (*cfg).InstName, types.MergePatchType, patchData, metav1.UpdateOptions{} )
        if err != nil {
            fmt.Println("Can't patch object ", (*cfg).InstName, " error: ", err)
        } else {
            fmt.Println("patch succeeded object ", (*cfg).InstName)
        }
    }
}

func RunOneTask( ctx* CmdContext ) ()  {
    var scriptFileName = "/tmp/runscript"
    var res *CmdResult = nil

    cfg := (*ctx).cfg

    file, err := os.Create(scriptFileName)
    if err != nil {
        fmt.Println("Can't create script file err=", err)
        return
    }
    file.WriteString( cfg.PythonScript )
    file.Close()

    cmd := []string{ "python3", scriptFileName }
    file.Write( []byte( cfg.PythonScript ) )
    file.Close()
    res, errr := NewCommand( cmd )

    if errr != nil {
        fmt.Println(errr,"Error during command execution")
    }

    if res != nil {
        LogCommandResults(ctx, res)
    }
}

func Run( ctx* CmdContext ) (error) {

    cfg := (*ctx).cfg

    /*
    cl := (*ctx).cl
    _, err := findScriptRunner( cfg.ScriptRunnerObjectName, cl )
    if err != nil {
            fmt.Println("Can't find the object for ", cfg.ScriptRunnerObjectName)
            return err
    }
    */

    if cfg.PackagesToInstall != "" {

        fmt.Println("Installing packages: ", cfg.PackagesToInstall)

        cmd := []string{ "sudo", "dnf", "install", "-y" }

        f := func(c rune) bool { return c == ' ' }
	pkgList := strings.FieldsFunc(cfg.PackagesToInstall, f)

        for _,val := range pkgList {
            cmd = append(cmd, val)
        }

        res, err := NewCommand( cmd )
        if err != nil {
            fmt.Println("Errors during intallation of packages. error=", err, " stdout=", res.StdOut, "stderr=", res.StdErr, " status=", res.Status )
        } else {
            fmt.Println("Installation of packages succeeded. stdout=", res.StdOut, "stderr=", res.StdErr, "status=", res.Status )
        }
    }

    if cfg.PipToInstall != "" {

        fmt.Println("Installing pip packages: ", cfg.PipToInstall)

        cmd := []string{ "sudo", "/root/pipinstall.sh" }

        f := func(c rune) bool { return c == ' ' }
	pkgList := strings.FieldsFunc(cfg.PipToInstall, f)

        for _,val := range pkgList {
            cmd = append(cmd, val)
        }

        res, err := NewCommand( cmd )
        if err != nil {
            fmt.Println("Errors during intallation of pip packages. error=", err, " stdout=", res.StdOut, "stderr=", res.StdErr, " status=", res.Status )
        } else {
            fmt.Println("Installation of pip packages succeeded. stdout=", res.StdOut, "stderr=", res.StdErr, "status=", res.Status )
        }
    }



    if cfg.InitWait != 0 {
        fmt.Println("initial wait: ", cfg.InitWait)
        time.Sleep( time.Duration( cfg.InitWait ) * time.Millisecond )
    }

    if cfg.NumRepetitions == 0 {
        RunOneTask(ctx)
    } else {
        for i := 0; i < cfg.NumRepetitions; i++ {

           RunOneTask(ctx)

           if i != (cfg.NumRepetitions-1) && cfg.RunPeriod !=0  {
               time.Sleep( time.Duration( cfg.RunPeriod ) * time.Millisecond )
           }
        }
    }
    return nil
}

/*
func makeClientGo() (*client.Client, error) {

    cfg, err := config.GetConfig()
    if  err!=nil {
        fmt.Println("Can't get configuration", err)
    }
    c, errr := client.New(cfg, client.Options{})
    if errr != nil {
        fmt.Println("Can't create client", errr)
        return nil, errr
    }

    return &c, errr
}
*/

func makeClientGoUnstructured() (*dynamic.NamespaceableResourceInterface, error) {

    config, err := rest.InClusterConfig()
    if err != nil {
       fmt.Println("Can't get config (InClusterConfig) err", err)
       return nil, err
    }

    /*
    c, err := kubernetes.NewForConfig(config)
    if err != nil {
        fmt.Println("Can't get k8s client ", err)
        return nil, err
    }
    */

    oGVR := schema.GroupVersionResource{
		Group:    "scriptrunner.github.com",
		Version:  "v1alpha1",
		Resource: "scriptrunners", //ScriptRunner",
	}

    dynClient, errClient := dynamic.NewForConfig(config)
    if errClient != nil {
        fmt.Printf("Error received creating client %v", errClient)
        return nil, errClient
    }

    client := dynClient.Resource(oGVR)

    return &client , nil
}

/*
func findScriptRunner(objectName string, cl *client.Client) ( * scriptrunnerv1alpha1.ScriptRunner, error) {

    opts := &client.ListOptions{}
    opts.SetLabelSelector("ScriptRunnerObjInstance=" + objectName)

    list := &scriptrunnerv1alpha1.ScriptRunnerList{}

    err := (*cl).List( context.TODO(), opts, list )
    if err != nil {
        fmt.Println("Failed to get ScriptRunner object instance ", objectName)
        return nil, err
    }
    if len(list.Items) == 0 {
        return nil, fmt.Errorf("no object found")
    }

    if len(list.Items) != 1 {
        return nil, fmt.Errorf("too many objects found")
    }

    return &list.Items[0], nil
}
*/


func main() {

    taskCfg, err := newTaskConfig()
    if err != nil {
        fmt.Println("Can't get task config from environment err", err)
        return
    }

    client, errr := makeClientGoUnstructured()
    if errr != nil {
        fmt.Println(errr,"Can't make go client err")
        return
    }

    ctx := makeContext(taskCfg, client)

    // run the main loop
    Run( ctx )

    // todo: should exit from main shutdown the container instance?
    // as crd is owner of this pod it will be cleaned up once the crd deployment is deleted; So guess not.
}
