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
  var newPossibleWords []string = possibleWords 
  for i := 0; i < len(guess); i++ {
    if score[i] == "2"[0] {
      newPossibleWords = removeWordsLackingAtPos(guess[i], i, newPossibleWords)
    } else if score[i] == "1"[0] {
      newPossibleWords = removeWordsLacking(guess[i], newPossibleWords)
    } else {
      newPossibleWords = removeWordsContaining(guess[i], newPossibleWords)
    }
  }
  newPossibleWords = removeWord(guess, newPossibleWords)
  return newPossibleWords
}

func removeWord(word string, words []string) []string {
  var isIn bool
  var index int
  for i, w := range words {
    if w == word {
      isIn = true
      index = i
    }
  }
  if isIn {
    words = append(words[0:index], words[index+1:]...)
  }
  return words
}

func letterEliminationScore(c byte, pos int, possibleWords []string) int {
  greenElims := len(removeWordsLackingAtPos(c, pos, possibleWords))
  greyElims := len(removeWordsContaining(c, possibleWords))
  yellowElims := len(removeWordsLacking(c, possibleWords))

  return greenElims + greyElims + yellowElims
}

func wordEliminationScore(word string, possibleWords []string) int {
  var score int
  for i := 0; i < len(word); i++ {
    score += letterEliminationScore(word[i], i, possibleWords)
  }
  return score
}

func min(arr []int) (int, int) {
  var mi, m int
  for i, v := range arr {
    if v < m {
      mi = i
      v = m
    }
  }
  return mi, m
}

func max(arr []int) (int, int) {
  var mi, m int
  for i, v := range arr {
    if v > m {
      mi = i
      v = m
    }
  }
  return mi, m
}

func getBestGuess(allowedWords, possibleWords []string) string {
  var eliminationScores []int
  var possEliminationScores []int
  for _, word := range allowedWords {
    eliminationScores = append(eliminationScores, wordEliminationScore(word, possibleWords))
  }
  for _, word := range possibleWords {
    possEliminationScores = append(possEliminationScores, wordEliminationScore(word, possibleWords))
  }
  mi, m := max(eliminationScores)
  mip, mp := max(possEliminationScores)
  if m > mp {
    return allowedWords[mi]
  }
  return possibleWords[mip]
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
    guess := getBestGuess(allowedWords, possibleWords)
    allowedWords = removeWord(guess, allowedWords)
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
