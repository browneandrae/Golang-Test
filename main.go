package main

import "os"
import "bufio"
import "fmt"
import "strings"
import "flag"
import "Memory"


func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
	var fileBool bool
	args := os.Args
	argCount := len(args)
	flag.BoolVar( &fileBool, "file", false, "Use to signal that inputting  text file to the database - input the filename (No Spaces) next") 
	interMode := flag.Bool("i", false, "Prefer interactive mode")
	helpTime := flag.Bool("help", false, "Shows how to use the program")
    flag.Parse()

	if fileBool == true {
		processFile(os.Args[2])
	} else if argCount == 1 || *helpTime == true{
		showHelp()
	} else if *interMode == true {
		runInteractive()
	}
	
}

/**********************
Routing for various functions used in interactive and file mode
**********************/
func decipherCommand(text string, memory *Memory.Database, interactive bool ) {
	/**********************
	Declare what's needed
	**********************/
	var command string 
	var key string
	var value string
	words := strings.Split(text, " ")
	wordsLength := len(words)
	if wordsLength > 0 {
		command = strings.ToUpper(strings.Trim(words[0], " "))

	}
	if wordsLength > 1 {
	
		key = strings.ToUpper(strings.Trim(words[1], " "))
	
	}
	if wordsLength > 2 {
		value = strings.ToUpper(strings.Trim(words[2], " "))
	}
	

	if strings.Contains("BEGIN", command){
		doBegin(memory)
	}else if strings.Contains("GET", command){
		get := doGet( key,memory)
		if !interactive {
			fmt.Print(" "+get)
		}else {
			fmt.Println(get)
		}
	} else if strings.Contains("SET", command){
		doSet( key,value,memory)
	} else if strings.Contains("UNSET", command) {
		doUnset(key,value,memory)
	}else if strings.Contains("ROLLBACK", command) {
		doRollBack(memory)
		if interactive {
			fmt.Println("")
		}
	} else if strings.Contains("NUMEQUALTO", command) {
		numequal := doNumEqualTo(key,memory)
		fmt.Print(" ")
		fmt.Print(numequal)
	} else if strings.Contains("END", command) {
		if !interactive {
		fmt.Println("")
		}
		doEnd()
	}else if strings.Contains("COMMIT", command) {
		doCommit(memory, interactive)
	}
	if !interactive {
		fmt.Println("")
	}

	
}

/**********************
 Open a new transaction block. Transaction blocks can be nested;
 a BEGIN can be issued inside of an existing block.
**********************/
func doBegin( m *Memory.Database) {
	m.Begin()
}

/**********************
– Close all open transaction blocks, permanently applying the changes made in them.
 Print nothing if successful, or print NO TRANSACTION if no transaction is in progress.
**********************/
func doCommit( m *Memory.Database, i bool) {
	m.Commit(i)
}

/**********************
END – Exit the program. Your program will
 always receive this as its last command.
**********************/
func doEnd() {
	os.Exit(0)
}

/**********************
Print out the value of the variable name, 
or NULL if that variable is not set.
**********************/
func doGet(key string, m *Memory.Database) string{
	value := m.Get(key)
	if value == "" {
		value = "NULL"
	}
	return value
}

/**********************
Print out the number of variables that are currently set to value.
 If no variables equal that value, print 0.
**********************/
func doNumEqualTo(key string, m *Memory.Database) int {
	return m.NumEqualTo(key) 
}

/**********************
Undo all of the commands issued in the most recent transaction block,
 and close the block. Print nothing if successful, or print NO TRANSACTION 
 if no transaction is in progress.
**********************/
func doRollBack( m *Memory.Database) {
	m.Rollback() 
}

/**********************
Set the variable name to the value value. 
Neither variable names nor values will contain spaces.
**********************/
func doSet(key string,value string, m *Memory.Database) {
	m.Set(key,value)
}

/**********************
Unset the variable name, making it just
 like that variable was never set.
**********************/
func doUnset(key string,value string, m *Memory.Database) {
	m.Unset(key, value)
}

/**********************
Called if flag -file is set
**********************/
func processFile(filePath string) {
	/**********************
	Get pwd, write unabbreviated here
	**********************/
	pwd, err := os.Getwd()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
	
    /**********************
	Read file with scanner
	**********************/
	readLine(pwd+"/"+filePath)
	
}

/**********************
Read Lines of a file
**********************/
func readLine(path string) {
  inFile, _ := os.Open(path)
  defer inFile.Close()
  scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines) 
	memory := new(Memory.Database)

  for scanner.Scan() {
    fmt.Print(scanner.Text())
	/**********************
	Do routing of commands to various functions
	**********************/
	decipherCommand(scanner.Text(), memory, false) 
  }
}

/**********************
Run Interactively
**********************/
func runInteractive() {
	/**********************
	Open Interaction
	**********************/
	scanner := bufio.NewScanner(os.Stdin)
	memory := new(Memory.Database)

	for scanner.Scan() {
		decipherCommand(scanner.Text(), memory, true) 
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

/**********************
Help documentation for command
**********************/
func showHelp() {
	fmt.Println("")
	fmt.Println("T")
	fmt.Println("")
	fmt.Println("Compilation instructions - To run the program without compiling , use the command 'go run [filename.go]' in a go enabled environment to use the database program. If you need to compile our program first install go onto your local or remote UNIX machine. Copy the files to a directory and make sure you have added the command go to your path. Then run go build main.go. Once compiled you may deploy this program anywhere, cd into the directory and run './[filename] followed by commands -f [filename.txt] or -i.")
	fmt.Println("")
	fmt.Println("-file - To input a file, use the -file flag after command line entry -  go [filename].go. The file must be in the same directory as the program. Each transaction should be on a new line, each word or letter should be separated by one space. Variable can either be a letters or string.")
	fmt.Println("")
	fmt.Println("-i - To use this program interactively, enter flag -i after command line entry -  go [filename].go. If using compile file change go [filename].go to ./[filename] then use the same command. Commands interactively must be entered one by one. Press enter to submit a command.")
	fmt.Println("")
	fmt.Println("Use END to exit the program in interactive or file mode")
	fmt.Println("")
	fmt.Println("View Readme file for command usage")	
	fmt.Println("")
	fmt.Println("This program is not 100% in memory. Due to golang's dependency on pointers it was challenging to copy arrays into one another and other related methods to created backups that could be rolled back to, so saving the data in a file became an alternate solution.")	
	fmt.Println("")
	fmt.Println("Lastly, files rollbacknum and rollbackstorage should be created and be writeable by the user running the go program in the shell, if you get an error file cannot be created, please create these files")

}

