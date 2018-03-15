/*
	Run dictionary attack over decrypting a pgp message.
*/

package main

import (
	"bufio"
	"os"
)

// DictionaryAttack Interface implementing the dictionary attack vector.
type DictionaryAttack interface {
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

// Attack function to call with an implementation of DictionaryAttack
func Attack(attackInterface DictionaryAttack) string {
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

// UserDictionaryAttack interface implmentation for running dictionary attack on user password
type UserDictionaryAttack struct{}

func (userAttack *UserDictionaryAttack) attack(word string) bool {
	return true
}

func main() {

}
