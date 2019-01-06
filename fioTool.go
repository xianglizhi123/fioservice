package main
import(
	"os"
	"os/exec"
	"bufio"
	"strings"
	"io"
	"fmt"
	"strconv"
	"encoding/json"
)
//report response 
type Result struct{
    Results []Value `json:"eachJob"`
    RunStatus RunStatus `json:"overAll"`
}
// block of the report,may be read, write
type Block struct{
    BlockType string 
    Lines []string
}
// this struct store the parameters of the report except total run status part
type Value struct{
    Read_io string `json:"read_io"`
    Read_bw string `json:"read_bw"`
    Read_iops string `json:"read_iops"`
    Read_runt string `json:"read_runt"`
    Read_clat_min string `json:"read_clat_min"`
    Read_clat_max string `json:"read_clat_max"`
    Read_clat_avg string `json:"read_clat_avg"`
    Read_clat_stdev string `json:"read_clat_stdev"`
    Read_lat_min string `json:"read_lat_min"`
    Read_lat_max string `json:"read_lat_max"`
    Read_lat_avg string `json:"read_lat_avg"`
    Read_lat_stdev string `json:"read_lat_stdev"`
    Read_bw_details_min string `json:"read_bw_details_min"`
    Read_bw_details_max string `json:"read_bw_details_max"`
    Read_bw_details_per string `json:"read_bw_details_per"`
    Read_bw_details_avg string `json:"read_bw_details_avg"`
    Read_bw_details_stdev string `json:"read_bw_details_stdev"`
    Write_io string `json:"write_io"`
    Write_bw string `json:"write_bw"`
    Write_iops string `json:"write_iops"`
    Write_runt string `json:"write_runt"`
    Write_clat_min string `json:"write_clat_min"`
    Write_clat_max string `json:"write_clat_max"`
    Write_clat_avg string `json:"write_clat_avg"`
    Write_clat_stdev string `json:"write_clat_stdev"`
    Write_lat_min string `json:"write_lat_min"`
    Write_lat_max string `json:"write_lat_max"`
    Write_lat_avg string `json:"write_lat_avg"`
    Write_lat_stdev string `json:"write_lat_stdev"`
    Write_bw_details_min string `json:"write_bw_details_min"`
    Write_bw_details_max string `json:"write_bw_details_max"`
    Write_bw_details_per string `json:"write_bw_details_per"`
    Write_bw_details_avg string `json:"write_bw_details_avg"`
    Write_bw_details_stdev string `json:"write_bw_details_stdev"`
}
//this struct store the runstatus of the report
type RunStatus struct{
    Read_io string `json:"total_read_io"`
    Read_aggrb string `json:"total_read_aggrb"`
    Read_minb string `json:"total_read_minb"`
    Read_maxb string `json:"total_read_maxb"`
    Read_mint string `json:"total_read_mint"`
    Read_maxt string `json:"total_read_maxt"`
    Write_io  string `json:"total_write_io"`
    Write_aggrb string `json:"total_write_aggrb"`
    Write_minb string `json:"total_write_minb"`
    Write_maxb string `json:"total_write_maxb"`
    Write_mint string `json:"total_write_mint"`
    Write_maxt string `json:"total_write_maxt"`
    Disk_stats_ios string `json:"disk_stats_ios"`
    Disk_stats_merge string `json:"disk_stats_merge"`
    Disk_stats_ticks string `json:"disk_stats_ticks"`
    Disk_stats_in_queue string `json:"disk_stats_in_queue"`
    Disk_stats_util string `json:"disk_stats_util"`
}
func RetriveWriteBw(str string, resp *Value){
    str=str[0:len(str)-1]
    var temp string
    for i:=0;i!=len(str);i++{
        if int(str[i])!=32{
            temp=temp+str[i:i+1]
        }
    }
    temp=strings.Split(temp,":")[1]
    infors:=strings.Split(temp,",")
    resp.Write_bw_details_min=strings.Split(infors[0],"=")[1]+"KB/s"
    resp.Write_bw_details_max=strings.Split(infors[1],"=")[1]+"KB/s"
    resp.Write_bw_details_per=strings.Split(infors[2],"=")[1]
    resp.Write_bw_details_avg=strings.Split(infors[3],"=")[1]
    resp.Write_bw_details_stdev=strings.Split(infors[4],"=")[1]
}
func RetriveReadBw(str string, resp *Value){
    str=str[0:len(str)-1]
    var temp string
    for i:=0;i!=len(str);i++{
        if int(str[i])!=32{
            temp=temp+str[i:i+1]
        }
    }
    temp=strings.Split(temp,":")[1]
    infors:=strings.Split(temp,",")
    resp.Read_bw_details_min=strings.Split(infors[0],"=")[1]+"KB/s"
    resp.Read_bw_details_max=strings.Split(infors[1],"=")[1]+"KB/s"
    resp.Read_bw_details_per=strings.Split(infors[2],"=")[1]
    resp.Read_bw_details_avg=strings.Split(infors[3],"=")[1]
    resp.Read_bw_details_stdev=strings.Split(infors[4],"=")[1]
}
func RetriveReadClat(str string, resp *Value){
    str=str[0:len(str)-1]
    var temp string
    for i:=0;i!=len(str);i++{
        if int(str[i])!=32{
            temp=temp+str[i:i+1]
        }
    }
    unit:=GetUnit(temp)
    temp=strings.Split(temp,":")[1]
    infors:=strings.Split(temp,",")
    resp.Read_clat_min=strings.Split(infors[0],"=")[1]+unit
    resp.Read_clat_max=strings.Split(infors[1],"=")[1]+unit
    resp.Read_clat_avg=strings.Split(infors[2],"=")[1]+unit
    resp.Read_clat_stdev=strings.Split(infors[3],"=")[1]
}
func RetriveReadLat(str string, resp* Value){
    //fmt.Printf(str)
    str=str[0:len(str)-1]
    var temp string
    for i:=0;i!=len(str);i++{
        if int(str[i])!=32{
            temp=temp+str[i:i+1]
        }
    }
    unit:=GetUnit(temp)
    temp=strings.Split(temp,":")[1]
    infors:=strings.Split(temp,",")
    resp.Read_lat_min=strings.Split(infors[0],"=")[1]+unit
    resp.Read_lat_max=strings.Split(infors[1],"=")[1]+unit
    resp.Read_lat_avg=strings.Split(infors[2],"=")[1]+unit
    resp.Read_lat_stdev=strings.Split(infors[3],"=")[1]
}
func RetriveWriteClat(str string, resp *Value){
    str=str[0:len(str)-1]
    var temp string
    for i:=0;i!=len(str);i++{
        if int(str[i])!=32{
            temp=temp+str[i:i+1]
        }
    }
    unit:=GetUnit(temp)
    temp=strings.Split(temp,":")[1]
    infors:=strings.Split(temp,",")
    resp.Write_clat_min=strings.Split(infors[0],"=")[1]+unit
    resp.Write_clat_max=strings.Split(infors[1],"=")[1]+unit
    resp.Write_clat_avg=strings.Split(infors[2],"=")[1]+unit
    resp.Write_clat_stdev=strings.Split(infors[3],"=")[1]
}
func RetriveWriteLat(str string, resp* Value){
	str=str[0:len(str)-1]
    var temp string
    for i:=0;i!=len(str);i++{
        if int(str[i])!=32{
            temp=temp+str[i:i+1]
        }
	}
	unit:=GetUnit(temp)
	temp=strings.Split(temp,":")[1]
    infors:=strings.Split(temp,",")
    resp.Write_lat_min=strings.Split(infors[0],"=")[1]+unit
    resp.Write_lat_max=strings.Split(infors[1],"=")[1]+unit
    resp.Write_lat_avg=strings.Split(infors[2],"=")[1]+unit
    resp.Write_lat_stdev=strings.Split(infors[3],"=")[1]
}
func GetUnit(str string) string{
    var start int
    var end int
    for i:=0;i!=len(str);i++{
        if str[i:i+1]=="("{
            //fmt.Println(str[i:i+1])
            start=i
        }
        if str[i:i+1]==")"{
            end=i
            break
        }
    }
    return str[start+1:end]
}
// retrive first line, io=,bw=...
func RetriveFirstLine(str string,resp *Value){
   str=str[2:len(str)-1]
   var temp string
   for i:=0;i!=len(str);i++{
       if int(str[i])!=32{
           temp=temp+str[i:i+1]
       }
   }
   mode:=strings.Split(temp,":")[0]
   infor:=strings.Split(temp,":")[1]
   //fmt.Printf("mode is %s,infor is %s\n",mode,infor)
   attributes:=strings.Split(infor,",")
   io:=strings.Split(attributes[0],"=")[1]
   bw:=strings.Split(attributes[1],"=")[1]
   iops:=strings.Split(attributes[2],"=")[1]
   runt:=strings.Split(attributes[3],"=")[1]
   //fmt.Printf("io=%s,bw=%s,iops=%s,runt=%s\n",io,bw,iops,runt)
   if mode=="read"{
       resp.Read_io=io
       resp.Read_bw=bw
       resp.Read_iops=iops
       resp.Read_runt=runt
   }else if mode=="write"{
       resp.Write_io=io
       resp.Write_bw=bw
       resp.Write_iops=iops
       resp.Write_runt=runt
   }
}
func RemoveBlockSpaceAndEmptyLine(block *Block){
     var str string
     str=""
     for i:=0;i!=len(block.Lines);i++{
        for j:=0;j!=len(block.Lines[i]);j++{
	   if block.Lines[i][j:j+1]!=""{
	      str=str+block.Lines[i][j:j+1]
	    }
	}
	block.Lines[i]=str
	str=""
     }
     var res []string
     for i:=0;i!=len(block.Lines);i++{
        if block.Lines[i]!="\n"{
	res=append(res,block.Lines[i])
	}
     }
     block.Lines=res    
}
func GetRunStatus(block Block,runStatus *RunStatus, fioMode string){
    if fioMode=="randread"{
        fmt.Println(block.Lines)
        numLines:=len(block.Lines)
        readLine:=strings.Split(block.Lines[1],":")[1]
        readLine=readLine[0:len(readLine)-1]
	readDetails:=strings.Split(readLine,",")
	runStatus.Read_io=strings.Split(readDetails[0],"=")[1]
        runStatus.Read_aggrb=strings.Split(readDetails[1],"=")[1]
	runStatus.Read_minb=strings.Split(readDetails[2],"=")[1]
	runStatus.Read_maxb=strings.Split(readDetails[3],"=")[1]
	runStatus.Read_mint=strings.Split(readDetails[4],"=")[1]
	runStatus.Read_maxt=strings.Split(readDetails[5],"=")[1]
	if numLines>2{
        diskLine:=strings.Split(block.Lines[3],":")[1]
        diskLine=diskLine[0:len(diskLine)-1]
	fmt.Printf("diskLine is %s\n",diskLine)
        diskDetails:=strings.Split(diskLine,",")
        runStatus.Disk_stats_ios=strings.Split(diskDetails[0],"=")[1]
        runStatus.Disk_stats_merge=strings.Split(diskDetails[1],"=")[1]
        runStatus.Disk_stats_ticks=strings.Split(diskDetails[2],"=")[1]
        runStatus.Disk_stats_in_queue=strings.Split(diskDetails[3],"=")[1]
        runStatus.Disk_stats_util=strings.Split(diskDetails[4],"=")[1]
	}
    }else if fioMode=="randwrite"{
        fmt.Println(block.Lines)
        numLines:=len(block.Lines)
	writeLine:=strings.Split(block.Lines[1],":")[1]
	writeLine=writeLine[0:len(writeLine)-1]
        writeDetails:=strings.Split(writeLine,",")
	runStatus.Write_io=strings.Split(writeDetails[0],"=")[1]
	runStatus.Write_aggrb=strings.Split(writeDetails[1],"=")[1]
        runStatus.Write_minb=strings.Split(writeDetails[2],"=")[1]
	runStatus.Write_maxb=strings.Split(writeDetails[3],"=")[1]
	runStatus.Write_mint=strings.Split(writeDetails[4],"=")[1]
	runStatus.Write_maxt=strings.Split(writeDetails[5],"=")[1]
	if numLines>2{					
	diskLine:=strings.Split(block.Lines[3],":")[1]
	diskLine=diskLine[0:len(diskLine)-1]
	fmt.Printf("diskLine is %s\n",diskLine)
        diskDetails:=strings.Split(diskLine,",")
        runStatus.Disk_stats_ios=strings.Split(diskDetails[0],"=")[1]
        runStatus.Disk_stats_merge=strings.Split(diskDetails[1],"=")[1]
        runStatus.Disk_stats_ticks=strings.Split(diskDetails[2],"=")[1]
        runStatus.Disk_stats_in_queue=strings.Split(diskDetails[3],"=")[1]
        runStatus.Disk_stats_util=strings.Split(diskDetails[4],"=")[1]
	}
    }else if fioMode=="randrw"{
        numLines:=len(block.Lines)
	readLine:=strings.Split(block.Lines[1],":")[1]
	readLine=readLine[0:len(readLine)-1]
        readDetails:=strings.Split(readLine,",")
        runStatus.Read_io=strings.Split(readDetails[0],"=")[1]
        runStatus.Read_aggrb=strings.Split(readDetails[1],"=")[1]
        runStatus.Read_minb=strings.Split(readDetails[2],"=")[1]
        runStatus.Read_maxb=strings.Split(readDetails[3],"=")[1]
        runStatus.Read_mint=strings.Split(readDetails[4],"=")[1]
        runStatus.Read_maxt=strings.Split(readDetails[5],"=")[1]
	writeLine:=strings.Split(block.Lines[2],":")[1]
	writeLine=writeLine[0:len(writeLine)-1]
        writeDetails:=strings.Split(writeLine,",")
        runStatus.Write_io=strings.Split(writeDetails[0],"=")[1]
        runStatus.Write_aggrb=strings.Split(writeDetails[1],"=")[1]
        runStatus.Write_minb=strings.Split(writeDetails[2],"=")[1]
        runStatus.Write_maxb=strings.Split(writeDetails[3],"=")[1]
        runStatus.Write_mint=strings.Split(writeDetails[4],"=")[1]
        runStatus.Write_maxt=strings.Split(writeDetails[5],"=")[1]
	if numLines>4{
	diskLine:=strings.Split(block.Lines[4],":")[1]
	diskLine=diskLine[0:len(diskLine)-1]
	fmt.Printf("diskLine is %s\n",diskLine)
        diskDetails:=strings.Split(diskLine,",")
        runStatus.Disk_stats_ios=strings.Split(diskDetails[0],"=")[1]
        runStatus.Disk_stats_merge=strings.Split(diskDetails[1],"=")[1]
        runStatus.Disk_stats_ticks=strings.Split(diskDetails[2],"=")[1]
        runStatus.Disk_stats_in_queue=strings.Split(diskDetails[3],"=")[1]
        runStatus.Disk_stats_util=strings.Split(diskDetails[4],"=")[1]
	}
    }
}
func BuildResult(blocks []Block,fioMode string) Result{
    BlockValue:=make([]Value,len(blocks)-1)
    if fioMode=="randread"{
       for i:=0;i!=len(blocks)-1;i++{
		 RetriveFirstLine(blocks[i].Lines[1],&BlockValue[i])
		 fmt.Printf("read_io=%s,read_bw=%s,read_iops=%s,read_runt=%s\n",BlockValue[i].Read_io,BlockValue[i].Read_bw,BlockValue[i].Read_iops,BlockValue[i].Read_runt)
		 RetriveReadBw(blocks[i].Lines[10],&BlockValue[i])
		 fmt.Printf("read_bw_details_min=%s,read_bw_details_max=%s,read_bw_details_per=%s,read_bw_details_avg=%s,read_bw_details_stdev=%s\n",BlockValue[i].Read_bw_details_min,BlockValue[i].Read_bw_details_max,BlockValue[i].Read_bw_details_avg,BlockValue[i].Read_bw_details_per,BlockValue[i].Read_bw_details_stdev)
		 RetriveReadClat(blocks[i].Lines[2],&BlockValue[i])
		 fmt.Printf("Read_clat_min=%s,Read_clat_max=%s,Read_clat_avg=%s,Read_clat_stdev=%s\n",BlockValue[i].Read_clat_min,BlockValue[i].Read_clat_max,BlockValue[i].Read_clat_avg,BlockValue[i].Read_clat_stdev)
		 RetriveReadLat(blocks[i].Lines[3],&BlockValue[i])
		 fmt.Printf("Read_lat_min=%s,Read_lat_max=%s,Read_lat_avg=%s,Read_lat_stdev=%s\n",BlockValue[i].Read_lat_min,BlockValue[i].Read_lat_max,BlockValue[i].Read_lat_avg,BlockValue[i].Read_lat_stdev)
       }
    }else if fioMode=="randwrite"{
        for i:=0;i!=len(blocks)-1;i++{
			RetriveFirstLine(blocks[i].Lines[1],&BlockValue[i])
			fmt.Printf("read_io=%s,read_bw=%s,read_iops=%s,read_runt=%s\n",BlockValue[i].Write_io,BlockValue[i].Write_bw,BlockValue[i].Write_iops,BlockValue[i].Write_runt)
			RetriveWriteBw(blocks[i].Lines[10],&BlockValue[i])
			fmt.Printf("read_bw_details_min=%s,read_bw_details_max=%s,read_bw_details_per=%s,read_bw_details_avg=%s,read_bw_details_stdev=%s\n",BlockValue[i].Write_bw_details_min,BlockValue[i].Write_bw_details_max,BlockValue[i].Write_bw_details_avg,BlockValue[i].Write_bw_details_per,BlockValue[i].Write_bw_details_stdev)
			RetriveWriteClat(blocks[i].Lines[2],&BlockValue[i])
			fmt.Printf("Read_clat_min=%s,Read_clat_max=%s,Read_clat_avg=%s,Read_clat_stdev=%s\n",BlockValue[i].Write_clat_min,BlockValue[i].Write_clat_max,BlockValue[i].Write_clat_avg,BlockValue[i].Write_clat_stdev)
			RetriveWriteLat(blocks[i].Lines[3],&BlockValue[i])
			fmt.Printf("Read_lat_min=%s,Read_lat_max=%s,Read_lat_avg=%s,Read_lat_stdev=%s\n",BlockValue[i].Write_lat_min,BlockValue[i].Write_lat_max,BlockValue[i].Write_lat_avg,BlockValue[i].Write_lat_stdev)
        }
    }else if fioMode=="randrw"{
        for i:=0;i!=len(blocks)-1;i++{
			RetriveFirstLine(blocks[i].Lines[1],&BlockValue[i])
			fmt.Printf("read_io=%s,read_bw=%s,read_iops=%s,read_runt=%s\n",BlockValue[i].Read_io,BlockValue[i].Read_bw,BlockValue[i].Read_iops,BlockValue[i].Read_runt)
			RetriveReadClat(blocks[i].Lines[2],&BlockValue[i])
			fmt.Printf("Read_clat_min=%s,Read_clat_max=%s,Read_clat_avg=%s,Read_clat_stdev=%s\n",BlockValue[i].Read_clat_min,BlockValue[i].Read_clat_max,BlockValue[i].Read_clat_avg,BlockValue[i].Read_clat_stdev)
			RetriveReadLat(blocks[i].Lines[3],&BlockValue[i])
			fmt.Printf("Read_lat_min=%s,Read_lat_max=%s,Read_lat_avg=%s,Read_lat_stdev=%s\n",BlockValue[i].Read_lat_min,BlockValue[i].Read_lat_max,BlockValue[i].Read_lat_avg,BlockValue[i].Read_lat_stdev)
			RetriveReadBw(blocks[i].Lines[10],&BlockValue[i])
			fmt.Printf("read_bw_details_min=%s,read_bw_details_max=%s,read_bw_details_per=%s,read_bw_details_avg=%s,read_bw_details_stdev=%s\n",BlockValue[i].Read_bw_details_min,BlockValue[i].Read_bw_details_max,BlockValue[i].Read_bw_details_avg,BlockValue[i].Read_bw_details_per,BlockValue[i].Read_bw_details_stdev)
			RetriveFirstLine(blocks[i].Lines[11],&BlockValue[i])
			fmt.Printf("read_io=%s,read_bw=%s,read_iops=%s,read_runt=%s\n",BlockValue[i].Write_io,BlockValue[i].Write_bw,BlockValue[i].Write_iops,BlockValue[i].Write_runt)
			RetriveWriteClat(blocks[i].Lines[12],&BlockValue[i])
			fmt.Printf("read_bw_details_min=%s,read_bw_details_max=%s,read_bw_details_per=%s,read_bw_details_avg=%s,read_bw_details_stdev=%s\n",BlockValue[i].Write_bw_details_min,BlockValue[i].Write_bw_details_max,BlockValue[i].Write_bw_details_avg,BlockValue[i].Write_bw_details_per,BlockValue[i].Write_bw_details_stdev)
			RetriveWriteLat(blocks[i].Lines[13],&BlockValue[i])
			fmt.Printf("Read_lat_min=%s,Read_lat_max=%s,Read_lat_avg=%s,Read_lat_stdev=%s\n",BlockValue[i].Write_lat_min,BlockValue[i].Write_lat_max,BlockValue[i].Write_lat_avg,BlockValue[i].Write_lat_stdev)
			RetriveWriteBw(blocks[i].Lines[20],&BlockValue[i])
			fmt.Printf("read_bw_details_min=%s,read_bw_details_max=%s,read_bw_details_per=%s,read_bw_details_avg=%s,read_bw_details_stdev=%s\n",BlockValue[i].Write_bw_details_min,BlockValue[i].Write_bw_details_max,BlockValue[i].Write_bw_details_avg,BlockValue[i].Write_bw_details_per,BlockValue[i].Write_bw_details_stdev)
          }
    }
    var runStatus RunStatus
    RemoveBlockSpaceAndEmptyLine(&blocks[len(blocks)-1])
    fmt.Println("hello 1")
    fmt.Println(blocks[len(blocks)-1].Lines)
    fmt.Println("hello 2")
    GetRunStatus(blocks[len(blocks)-1],&runStatus,fioMode)
    var result Result
    result.Results=BlockValue
    result.RunStatus=runStatus
    return result
}
func GetUsefulBlocks(rawStr[]string,projectName string,fioMode string,numJobs string) []Block{
	BlockBeginChars:=projectName+": (groupid="
	var begin int
	for i:=0;i!=len(rawStr);i++{
	  if len(rawStr[i])>=len(BlockBeginChars){
		  if BlockBeginChars==rawStr[i][0:len(BlockBeginChars)]{
			  begin=i
			  break
		  }
	  }
	}
	strs:=rawStr[begin:]
	//fmt.Println(strs)
	var strs2 []string
	//delte null line from the report
	for i:=0;i!=len(strs);i++{
		if strs[i]!="\n"{
		strs2=append(strs2,strs[i])
		}
	}
	strs=strs2
	numjobs,err:=strconv.Atoi(numJobs)
	if err!=nil{
		numjobs=1
	}
	var AllJobsInforStart int
	blocks:=make([]Block,numjobs+1)
	var pre int
	var post int
	pre=0
	if fioMode=="randread"||fioMode=="randwrite"{
	 for i:=0;i!=numjobs;i++{
		 blocks[i].BlockType=fioMode
		 for j:=pre+1;j!=len(strs);j++{
		  if BlockBeginChars==strs[j][0:len(BlockBeginChars)]||strs[j][0:len("Run status group")]=="Run status group"{
			 post=j
			 break
		   }
		 }
		 blocks[i].Lines=strs[pre:post]
		 pre=post
	 }
	}else{
	 for i:=0;i!=numjobs;i++{
		 blocks[i].BlockType=fioMode
		 for j:=pre+1;j!=len(strs);j++{
		   if BlockBeginChars==strs[j][0:len(BlockBeginChars)]||strs[j][0:len("Run status group")]=="Run status group"{
			  post=j
			  break
			}
		  }
		  blocks[i].Lines=strs[pre:post]
		  pre=post
		 
	 }
	}
   for i:=0;i!=len(strs);i++{
	   if len(strs[i])>len("Run status group")&&strs[i][0:len("Run status group")]=="Run status group"{
		 AllJobsInforStart=i 
		 break
	   }
   } 
   blocks[numjobs].BlockType="AllJobsStatus"
   blocks[numjobs].Lines=strs[AllJobsInforStart:]
   for i:=0;i!=len(blocks[numjobs].Lines);i++{
	   var temp string
	   for j:=0;j!=len(blocks[numjobs].Lines[i]);j++{
		   if int(blocks[numjobs].Lines[i][j])!=32{
			   temp=temp+blocks[numjobs].Lines[i][j:j+1]
		   }
	   }
	   blocks[numjobs].Lines[i]=temp
   }
   return blocks
}
func main(){
 //fmt.Printf("start fioAgent process\n")
 var cmds []string
 cmds=os.Args[1:len(os.Args)]
 fmt.Println(cmds)
 cmd:=exec.Command("fio",cmds...)
 stdout,_ :=cmd.StdoutPipe()
 cmd.Start()
 reader:=bufio.NewReader(stdout)
 var resp []string
  for {
    line,err2:=reader.ReadString('\n')
      if err2!=nil ||io.EOF ==err2{
      break
      }
      //resp=resp+line
      resp=append(resp,line)
  }
  fmt.Println("fio result")
  fmt.Println(resp)
  fmt.Println("fio ends")
  var fioMode string
  var numJobs string
  var projectName string
  for i:=0;i!=len(cmds);i++{
     if cmds[i][0:4]=="-rw="{
     fioMode=cmds[i][4:len(cmds[i])]
     }
     if len(cmds[i])>8&&cmds[i][0:8]=="-numjobs"{
     numJobs=cmds[i][9:len(cmds[i])]
     }
     if len(cmds[i])>6&&cmds[i][0:5]=="-name"{
     projectName=cmds[i][6:len(cmds[i])]
     }
  }
  fmt.Printf("fioMode=%s,numJobs=%s,projectName=%s\n",fioMode,numJobs,projectName)
  blocks:=GetUsefulBlocks(resp,projectName,fioMode,numJobs)
  fmt.Println(blocks)
  res:=BuildResult(blocks,fioMode)
  b,_:=json.Marshal(res)
  fmt.Println("report")
  fmt.Println(b)
  fmt.Println("report end")
  var reportName string
  reportName="/go/src/fioProject/report/"+strconv.Itoa(os.Getpid())+"-"+projectName
  f,_:=os.Create(reportName)
  defer f.Close()
  f.Write(b)
  w:=bufio.NewWriter(f)
  w.Flush()
  cmd.Wait()
}