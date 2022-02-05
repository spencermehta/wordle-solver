package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readLines(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }

  return lines, scanner.Err()

}

func getGuessScore(guess, todaysWord string) string {
  var score string
  for i := 0; i < len(guess); i++ {
    if guess[i] == todaysWord[i] {
      score += "2"
    } else if strings.Contains(todaysWord, string(guess[i])) {
      score += "1"
    } else {
      score += "0"
    }
  }
  return score
}

func removeWordsContaining(c byte, possibleWords []string) []string {
  var newPossibleWords []string
  for _, word := range possibleWords {
    if !strings.Contains(word, string(c)) {
      newPossibleWords = append(newPossibleWords, word)
    }
  }
  return newPossibleWords
}

func removeWordsLacking(c byte, possibleWords []string) []string {
  var newPossibleWords []string
  for _, word := range possibleWords {
    if strings.Contains(word, string(c)) {
      newPossibleWords = append(newPossibleWords, word)
    }
  }
  return newPossibleWords
}

func removeWordsLackingAtPos(c byte, i int, possibleWords []string) []string {
  var newPossibleWords []string
  for _, word := range possibleWords {
    if word[i] == c {
      newPossibleWords = append(newPossibleWords, word)
    }
  }
  return newPossibleWords
}

func eliminateWords(guess, score string, possibleWords []string) []string {
  var newPossibleWords []string
  for i := 0; i < len(guess); i++ {
    if score[i] == "2"[0] {
      newPossibleWords = removeWordsLackingAtPos(guess[i], i, possibleWords)
    } else if score[i] == "1"[0] {
      newPossibleWords = removeWordsLacking(guess[i], possibleWords)
    } else {
      newPossibleWords = removeWordsContaining(guess[i], possibleWords)
    }
  }

  var isIn bool
  var index int
  for i, w := range newPossibleWords {
    if w == guess {
      isIn = true
      index = i
    }
  }
  if isIn {
    newPossibleWords = append(newPossibleWords[0:index], newPossibleWords[index+1:]...)
  }
  return newPossibleWords
}

func main() {
  allowedWords, err := readLines("wordle-allowed-guesses.txt")
  possibleWords := allowedWords
  if err != nil {
    panic(err)
  }

  fmt.Print("Enter today's word:")
  var todaysWord string
  fmt.Scan(&todaysWord)

  var guesses int
  for {
    guess := possibleWords[0]
    guesses += 1
    fmt.Printf("Guess #%d: %v, %v, %d words left\n", guesses, guess, todaysWord, len(possibleWords))
    if guess == todaysWord {
      fmt.Println("found!")
      return
    }
    score := getGuessScore(guess, todaysWord)
    possibleWords = eliminateWords(guess, score, possibleWords)
    fmt.Println(len(possibleWords))
  }
}
