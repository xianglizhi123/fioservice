package main
import(
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "os/exec"
    "io/ioutil"
    "fmt"
	"os"
    "strconv"	
)
type CheckStatusReturn struct{
	Status string `json:"status"`
}
var ch chan int
type FioAgentPid struct{
	Pid int `json:"fioPid"`
}
type ReportName struct{
	Pid int `json:"fioPid"`
	Task string `json:"task"`
}
type ReportResponse struct{
	Outcome interface{} `json:"fioOutcome"`
	Code int `json:"code"`
	Message string `json:"message"`
}
//fio tool parameters, need these paramters to invoke fio
type FioParameters struct{
	Filename  string `json:"filename,omitempty"`
	Direct    string `json:"direct"`
	Ioengine  string `json:"ioengine"`
	Iodepth   string `json:"iodepth"`
	Rw        string `json:"rw,omitempty"`
	Rwmixwrite string `json:"rwmixwrite"`
	Bs         string `json:"bs,omitempty"`
	Size       string `json:"size"`
	Numjobs    string `json:"numjobs"`
	Runtime    string `json:"runtime"`
	Name       string `json:"name,omitempty"`
}
func SetParameters(fio FioParameters) []string{
	var resp []string
	resp=append(resp,"-filename="+fio.Filename)
	if fio.Direct!=""{
	resp=append(resp,"-direct="+fio.Direct)
	}
	if fio.Ioengine!=""{
	resp=append(resp,"-ioengine="+fio.Ioengine)
	}
	if fio.Iodepth!=""{
	resp=append(resp,"-iodepth="+fio.Iodepth)
	}
	resp=append(resp,"-rw="+fio.Rw)
	if fio.Rw=="randrw"{
	resp=append(resp,"-rwmixwrite="+fio.Rwmixwrite)
	}
	if fio.Bs!="" {
	resp=append(resp,"-bs="+fio.Bs)
	}
	if fio.Size!=""{
	resp=append(resp,"-size="+fio.Size)
	}
	if fio.Numjobs!=""{
	resp=append(resp,"-numjobs="+fio.Numjobs)
	}
	if fio.Runtime!=""{
	resp=append(resp,"-runtime="+fio.Runtime)
	}
	resp=append(resp,"-name="+fio.Name)
	return resp
}
func HandleGetReport(w http.ResponseWriter,r *http.Request){
	fmt.Println("inside get report")
	var reportName ReportName
	var reportResponse ReportResponse
	_=json.NewDecoder(r.Body).Decode(&reportName)
	var cmds []string
	fileName:="/go/src/fioProject/report/"+strconv.Itoa(reportName.Pid)+"-"+reportName.Task
	cmds=append(cmds,"-p")
    cmds=append(cmds,strconv.Itoa(reportName.Pid))
	cmd,_:=exec.Command("ps",cmds...).Output()
	if len(string(cmd))>30{
	   reportResponse.Code=400
	   reportResponse.Message="the task is running, check the report later"
	   json.NewEncoder(w).Encode(reportResponse)
	}else{
		plan,err:=ioutil.ReadFile(fileName)
		if err!=nil{
			reportResponse.Code=400
			reportResponse.Message="no such report"
			json.NewEncoder(w).Encode(reportResponse)
		}else{
			reportResponse.Code=200
			reportResponse.Message="success"
			json.Unmarshal(plan,&reportResponse.Outcome)
			json.NewEncoder(w).Encode(reportResponse)
		}
	}
}
func HandleCheckTask(w http.ResponseWriter, r * http.Request){
	fmt.Println("inside check task")
	var fioAgentPid FioAgentPid
	_=json.NewDecoder(r.Body).Decode(&fioAgentPid)
	var cmds []string
	cmds=append(cmds,"-p")
    cmds=append(cmds,strconv.Itoa(fioAgentPid.Pid))
	cmd,_:=exec.Command("ps",cmds...).Output()
	var checkStatusReturn CheckStatusReturn
	if len(string(cmd))>30{
       checkStatusReturn.Status="Running"
	}else{
       checkStatusReturn.Status="Finished"
	}
	json.NewEncoder(w).Encode(checkStatusReturn)
}
func TestServer(w http.ResponseWriter, r * http.Request){
   	json.NewEncoder(w).Encode("hello hello")
}
func HandleFioRequest(w http.ResponseWriter, r * http.Request){
	fmt.Println("inside fio request")
	var fioParameters FioParameters
	_=json.NewDecoder(r.Body).Decode(&fioParameters)
	cmds:=SetParameters(fioParameters)
	cmd:=exec.Command("./fioTool/fioTool",cmds...)
	//json.NewEncoder(w).Encode("hello")
	cmd.Stdout=os.Stdout
	cmd.Start()
	pid:=cmd.Process.Pid
	//fmt.Printf("fioTool pid is %d\n",pid)
	var res FioAgentPid
	res.Pid=pid
	ch:=make(chan error,1)
	go func(){
		ch<-cmd.Wait()
	}()
	//cmd.Wait()
	json.NewEncoder(w).Encode(res)
}
func main() {
    router := mux.NewRouter()
	router.HandleFunc("/ExecuteFio",HandleFioRequest).Methods("POST")
	router.HandleFunc("/CheckStatus",HandleCheckTask).Methods("POST")
	router.HandleFunc("/GetReport",HandleGetReport).Methods("POST")
	router.HandleFunc("/TestServer",TestServer).Methods("POST")
    log.Fatal(http.ListenAndServe(":8000", router))
}