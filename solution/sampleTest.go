package main
import ("fmt"
//        "log"
        "os/exec"
        "strings"
//        "io/ioutil"
)


//L at the end indicates
type Node struct{
        name string
        parentL  map[string]string
        childL   map[string]string
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

func main(){

basePath := "/home/maersk/test_data/git"
checkFolders := []string{"raw","tmp","final"}
//Declaring Map to hold Parent-Child relationship
myVertices := make(map[string]Node)
//fmt.Println(myVertices)
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
                        myVertices[tblName]= Node{parentL: map[string]string{},childL: map[string]string{}}
                        }
                //If the given element is .sql script then add the target table in Map and then add current script as node & then add
                //target table as a parent to script
                //fmt.Println(lineSplit[0])
                if strings.Contains(lineSplit[0],".sql"){
                        tblName := schema+"."+strings.Replace(lineSplit[0],".sql","",1)
                        //fmt.Println(tblName)
                        sqlScriptName := lineSplit[0]
                        myVertices[tblName]= Node{parentL: map[string]string{},childL: map[string]string{lineSplit[0]:""}}
                //Adding .sql script as a Node & target table as it's parent
                        myVertices[sqlScriptName]= Node{parentL:map[string]string{tblName:""},childL:map[string]string{}}
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
                                                myVertices[sqlScriptName].childL[childTbl]=""
                                                //Add the current script to the parentL of the childTbl. Since We've already added all the tables from rsw schema
                                                //& then added all the target tables for script then all the child must be present in the myVertices.
                                                //fmt.Println(sqlScriptName)
                                                //fmt.Println(myVertices[childTbl])
                                                //fmt.Println(childTbl)
                                                myVertices[childTbl].parentL[sqlScriptName]=""
                                                }
                                        }
                                }
                        }
                }
        }


//fmt.Println(myVertices)
//      fmt.Println(output)





        }
fmt.Println(myVertices)
}


