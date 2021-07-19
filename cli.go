package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/jlaffaye/ftp"
)

//CONSTANTS
var goToDir string
var ipadd string = "demo.wftpserver.com:21"
var username string = "demo"
var password string = "demo"
var file1 string = "Summer"
var file2 string = "Winter"

type jsonread struct {
	GoToDir string `json:"goToDir"`
	Hrsmins []struct {
		Hour   int `json:"hour"`
		Minute int `json:"minute"`
	} `json:"hrsmins"`
}

func ftpconnect(ipadd, username, password, goToDir, filename string) bool {
	//connect to Server
	c, err := ftp.Dial(ipadd, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
		return false
	}
	fmt.Println("Connected to FTP Server")

	//login
	err = c.Login(username, password)
	if err != nil {
		log.Fatal(err)
		return false
	}
	fmt.Println("Access Granted")

	//change folder
	fmt.Println("Changing Directory")
	err = c.ChangeDir(goToDir)
	if err != nil {
		log.Fatal(err)
		return false
	}
	fmt.Println("Directory Changed ")
	// open a file
	fmt.Println("Opening " + filename + ".jpg .....")
	r, err := c.Retr(filename + ".jpg")
	if err != nil {
		log.Fatal(err)
		return false
	}
	fmt.Println(filename + ".jpg now open")
	//read from file
	fmt.Println("Reading......")
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
		return false
	}
	//write data on file
	fmt.Println("Creating local copy: ")
	err = ioutil.WriteFile(filename+"_"+time.Now().Format("2006-01-02 15:04:05")+".jpg", buf, 0644)
	if err != nil {
		log.Fatal(err)
		return false
	}
	fmt.Println("Copy Created")
	r.Close()
	//Exit
	if err := c.Quit(); err != nil {
		log.Fatal(err)
		return false
	}
	fmt.Println("Success")
	fmt.Println("Restarting")
	return true
}

func main() {

	jsonread := jsonread{}
	file, err := ioutil.ReadFile("input.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	err = json.Unmarshal([]byte(file), &jsonread)
	if err != nil {
		log.Fatal(err)
		return
	}
	goToDir = jsonread.GoToDir
	tim := [6]time.Time{
		time.Date(int(time.Now().Year()), time.Now().Month(), time.Now().Day(), jsonread.Hrsmins[0].Hour, jsonread.Hrsmins[0].Minute, 0, 0, time.Local),
		time.Date(int(time.Now().Year()), time.Now().Month(), time.Now().Day(), jsonread.Hrsmins[1].Hour, jsonread.Hrsmins[1].Minute, 0, 0, time.Local),
		time.Date(int(time.Now().Year()), time.Now().Month(), time.Now().Day(), jsonread.Hrsmins[2].Hour, jsonread.Hrsmins[2].Minute, 0, 0, time.Local),
		time.Date(int(time.Now().Year()), time.Now().Month(), time.Now().Day(), jsonread.Hrsmins[3].Hour, jsonread.Hrsmins[3].Minute, 0, 0, time.Local),
		time.Date(int(time.Now().Year()), time.Now().Month(), time.Now().Day(), jsonread.Hrsmins[4].Hour, jsonread.Hrsmins[4].Minute, 0, 0, time.Local),
		time.Date(int(time.Now().Year()), time.Now().Month(), time.Now().Day(), jsonread.Hrsmins[5].Hour, jsonread.Hrsmins[5].Minute, 0, 0, time.Local),
	}
	//checks time diff for next ftp call and calls ftpconnect
	index := 0
	for index < 6 {
		if tim[index].After(time.Now()) {
			diff := time.Until(tim[index])
			fmt.Println("Sleeping for " + diff.String())
			time.Sleep(diff)
		}
		go ftpconnect(ipadd, username, password, goToDir, file1)
		ftpconnect(ipadd, username, password, goToDir, file2)
		index++
	}
	fmt.Println("Files downloaded for today!!")
}
