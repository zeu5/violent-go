/*
	Run dictionary attack over decrypting a pgp message.
*/

package main

import (
	"bufio"
	"fmt"
	"os"
)

// DictionaryAttack Interface implementing the dictionary attack vector.
type Attacker interface {
	attack(word string) bool
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func readLines(f *os.File) []string {
	// Read lines of the given file
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

// Attack function to call with an implementation of Attacker
func Attack(attackInterface Attacker) string {
	dictFile, err := os.Open("/usr/share/dict/words")
	check(err)
	defer dictFile.Close()

	var answer string
	done := make(chan string, 10)

	for _, word := range readLines(dictFile) {
		go func(word string, done chan string) {
			if attackInterface.attack(word) {
				done <- word
			} else {
				done <- ""
			}
		}(word, done)

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

// UserAttacker interface implmentation for running dictionary attack on user password
type UserAttacker struct{}

func (userAttack *UserAttacker) attack(word string) bool {
	if word == "Zyzzogeton" {
		return true
	}
	return false
}

// NewUserAttacker constructor for the interface
func NewUserAttacker() *UserAttacker {
	return &UserAttacker{}
}

func main() {
	userAttacker := NewUserAttacker()
	if password := Attack(userAttacker); password != "" {
		fmt.Println("Found password: " + password)
	} else {
		fmt.Println("No password found")
	}
}
