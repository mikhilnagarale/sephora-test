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
        parentL []string
        childL  []string
}

func check(e error){
        if e !=nil{
                panic(e)
        }
}

func getDirectories(path string)string{
out,err:=exec.Command("ls","-lrt",path).Output()

check(err)
return string(out)

}
func main(){

basePath := "/home/maersk/test_data/git"
checkFolders := []string{"raw","tmp","final"}
//Declaring Map to hold Parent-Child relationship
myVertices := map[string]Node{}
fmt.Println(myVertices)
for i := range checkFolders{
        output := getDirectories(basePath+"/"+checkFolders[i])
        lines := strings.Split(output,"\n")
        schema := checkFolders[i]
        for lineNo := range lines{
        lineSplit := strings.Split(lines[lineNo]," ")
        if len(lineSplit)==9{
                //fmt.Println(lineSplit)
                //Adding table Name with Schema in Map myVertices
                if !(strings.Contains(lineSplit[8],".sql")){
                        tblName := schema+"."+lineSplit[8]
                        myVertices[tblName]= Node{parentL:[]string{},childL:[]string{}}
                        }
                //If the given element is .sql script then add the target table in Map and then add current script as node & then add
                //target table as a parent to script
                if strings.Contains(lineSplit[8],".sql"){
                        tblName := schema+"."+strings.Trim(lineSplit[8],".sql")
                        myVertices[tblName]= Node{parentL:[]string{},childL:[]string{lineSplit[8]}}
                //Adding .sql script as a Node & target table as it's parent
                        myVertices[lineSplit[8]]= Node{parentL:[]string{tblName},childL:[]string{}}
                //Read the script Line by Line & add the dependencies accordingly
                        }
                }
        }


//fmt.Println(myVertices)
//      fmt.Println(output)





        }
fmt.Println(myVertices)
}

