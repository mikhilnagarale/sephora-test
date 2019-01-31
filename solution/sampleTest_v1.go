package main
import ("fmt"
//        "log"
        "os/exec"
        "strings"
//        "io/ioutil"
        "math/rand"
        "time"
        "os"
)


//L at the end indicates
type Node struct{
        parentL  map[string]string
        childL   map[string]string
        //message is for current Node status
        //-1 = Job Not Started //0 = Job Started //1 = Job Completed
        message  []int
}

func check(e error){
        if e !=nil{
                panic(e)
        }
}

func getDirectories(path string)string{
out,err:=exec.Command("ls",path).Output()
check(err)
return string(out)
}

func getFileData(path string)string{
out,err:=exec.Command("cat",path).Output()
check(err)
return string(out)
}

func loadVertices(myVertices *map[string]Node,basePath string,checkFolders []string){
//fmt.Println(checkFolders)
//fmt.Println(basePath)
//fmt.Println((*myVertices))

for i := range checkFolders{
        output := getDirectories(basePath+"/"+checkFolders[i])
        folderPath := basePath+"/"+checkFolders[i]
        lines := strings.Split(output,"\n")
        schema := checkFolders[i]
        for lineNo := range lines{
        lineSplit := strings.Split(lines[lineNo]," ")
        //fmt.Println(lineSplit[0],len(lineSplit),len(lineSplit[0]))
        if len(lineSplit[0])>0{
                //fmt.Println(lineSplit)
                //Adding table Name with Schema in Map myVertices
                if !(strings.Contains(lineSplit[0],".sql")){
                        tblName := schema+"."+lineSplit[0]
                        (*myVertices)[tblName]= Node{parentL: map[string]string{},childL: map[string]string{},message:[]int{-1}}
                        }
                //If the given element is .sql script then add the target table in Map and then add current script as node & then add
                //target table as a parent to script
                //fmt.Println(lineSplit[0])
                if strings.Contains(lineSplit[0],".sql"){
                        tblName := schema+"."+strings.Replace(lineSplit[0],".sql","",1)
                        //fmt.Println(tblName)
                        sqlScript := lineSplit[0]
                        sqlScriptName := schema+"."+sqlScript
                        (*myVertices)[tblName]= Node{parentL: map[string]string{},childL: map[string]string{sqlScriptName:""},message:[]int{-1}}
                //Adding .sql script as a Node & target table as it's parent
                        (*myVertices)[sqlScriptName]= Node{parentL:map[string]string{tblName:""},childL:map[string]string{},message:[]int{-1}}
                //Read the script Line by Line & add the dependencies accordingly : getFileData
                        filePath := folderPath+"/"+lineSplit[0]
                        //fmt.Println(filePath)
                        fileData := getFileData(filePath)
                        //fmt.Println(fileData)
                        fileLines := strings.Split(fileData,"\n")
                        for fileLineNo := range fileLines{
                                fileLineSplits := strings.Split(fileLines[fileLineNo]," ")
                                for splitIndex := range fileLineSplits{
                                        if strings.Contains(fileLineSplits[splitIndex],"raw.")||strings.Contains(fileLineSplits[splitIndex],"final.")||strings.Contains(fileLineSplits[splitIndex],"tmp.")              {
                                                childTbl := strings.Trim(fileLineSplits[splitIndex],"`")
                                                //fmt.Println(sqlScriptName+"<-"+childTbl)
                                                //Add the identified child table to the childL of curerent script
                                                //statement myVertices[lineSplit[8]]  will return the Node. We need to add childTbl to the childL map of this Node
                                                (*myVertices)[sqlScriptName].childL[childTbl]=""
                                                //Add the current script to the parentL of the childTbl. Since We've already added all the tables from rsw schema
                                                //& then added all the target tables for script then all the child must be present in the myVertices.
                                                //fmt.Println(sqlScriptName)
                                                //fmt.Println(myVertices[childTbl])
                                                //fmt.Println(childTbl)
                                                (*myVertices)[childTbl].parentL[sqlScriptName]=""
                                                }
                                        }
                                }
                        }
                }
        }
    }
fmt.Println((*myVertices))
}

func startNode(script string, myVertices *map[string]Node){
sleepSeconds := 1+rand.Intn(10)
if (*myVertices)[script].message[0] == -1{
        fmt.Println("\nStarted the execution/loading of "+script+"\n")
        (*myVertices)[script].message[0] = 0
        time.Sleep(time.Duration(sleepSeconds)*time.Second)
        (*myVertices)[script].message[0] = 1
        fmt.Println("\nCompleted the execution/loading of "+script+"\n")
        }
}

func initiateDependencyExe(myVertices *map[string]Node,inProgress *map[string]string){

totalNoOfVertices := 0
incompleteVertices := 0
//fmt.Println(myVertices)

totalNoOfVertices = len((*myVertices))
incompleteVertices = len((*myVertices))
fmt.Println(incompleteVertices)

//fmt.Println(myVertices)
//startNode("raw.inventory_items",&myVertices)

//Traverse the Map - myVertices & start the leaf nodes / nodes with no child
for key := range (*myVertices){
        //fmt.Println(myVertices[key],len(myVertices[key].childL))
        if len((*myVertices)[key].childL)==0{
        go startNode(key,myVertices)
        (*inProgress)[key] = ""
        }
}

for len((*inProgress))!=0  {
        for inProgressKey := range (*inProgress){
        if ((*myVertices)[inProgressKey].message[0]==1){
                //Check Parent Status
                for parentKey := range (*myVertices)[inProgressKey].parentL{
                        fmt.Println("checking Parent "+parentKey,(*myVertices)[parentKey].message[0])
                        //Check Parent's Child Status
                        noOfChild := len((*myVertices)[parentKey].childL)
                        completedChild := 0
                        for childKey := range (*myVertices)[parentKey].childL{
                                fmt.Println("checking childs "+childKey,(*myVertices)[childKey].message[0])
                                //Parent's action based on Child status
                                if (*myVertices)[childKey].message[0]==1{
                                        completedChild = completedChild + 1
                                }
                        }
                        //Check If all Child for current Parent are completed. If yes then start the current Parent.
                        if completedChild == noOfChild{
                                //Check if Parent Node is already started then only start the parent
                                if (*myVertices)[parentKey].message[0]==-1{
                                go startNode(parentKey,myVertices)
                                (*inProgress)[parentKey] = ""
                                }
                        }
                fmt.Println("\n")
                //Added this delay since the loop is checking next parent before startNode() process starts the previous parent.
                //This was causing issue when two childs had same parent.
                time.Sleep(1*time.Second)
                }
                //remove completed child
                delete((*inProgress),inProgressKey)
                incompleteVertices = incompleteVertices-1
                }
        }
        //fmt.Println(myVertices)
        fmt.Println("In Progress Queue = ")
        fmt.Println((*inProgress))
        fmt.Println("\n")
        fmt.Printf("totalNoOfVertices=%d,incompleteVertices=%d\n",totalNoOfVertices,incompleteVertices)
        time.Sleep(10*time.Second)
        }
}

func main(){

//basePath := "/home/maersk/test_data/git"
basePath := os.Args[1]
checkFolders := []string{"raw","tmp","final"}
//Declaring Map to hold Parent-Child relationship
myVertices := make(map[string]Node)
loadVertices(&myVertices,basePath,checkFolders)
//This Array will have scripts/nodes which are currently running
inProgress := map[string]string{}
initiateDependencyExe(&myVertices,&inProgress)

        fmt.Println(myVertices)
        fmt.Println(inProgress)
        time.Sleep(10*time.Second)
}


