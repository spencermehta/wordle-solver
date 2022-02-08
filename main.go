package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
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

func letterEliminationScore(c byte, pos int, repeat bool, possibleWords []string) int {
  greenElims := len(removeWordsLackingAtPos(c, pos, possibleWords))
  var yellowElims, greyElims int
  if !repeat {
    greyElims = len(removeWordsContaining(c, possibleWords))
    yellowElims = len(removeWordsLacking(c, possibleWords))
  }

  return greenElims + greyElims + yellowElims
}

func contains(c byte, arr []byte) bool {
  var isIn bool
  for _, char := range arr {
    if c == char {
      isIn = true
    }
  }
  return isIn
}

func wordEliminationScore(word string, possibleWords []string) int {
  var score int
  var usedChars []byte
  for i := 0; i < len(word); i++ {
    repeat := contains(word[i], usedChars)
    score += letterEliminationScore(word[i], i, repeat, possibleWords)
    usedChars = append(usedChars, word[i])
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

func guessWord(todaysWord string, allowedWords, possibleWords []string) (int, []string) {
  var numGuesses int
  var guesses []string
  for {
    guess := getBestGuess(allowedWords, possibleWords)
    guesses = append(guesses, guess)
    allowedWords = removeWord(guess, allowedWords)
    numGuesses += 1
    if guess == todaysWord {
      return numGuesses, guesses
    }
    score := getGuessScore(guess, todaysWord)
    possibleWords = eliminateWords(guess, score, possibleWords)
  }
}

func main() {
  allowedWords, err := readLines("wordle-allowed-guesses.txt")
  if err != nil {
    panic(err)
  }
  realAnswers, err := readLines("wordle-answers-alphabetical.txt")
  if err != nil {
    panic(err)
  }
  allowedWords = append(allowedWords, realAnswers...)
  possibleWords := allowedWords

  csvFile, err := os.Create("results.csv")
  csvwriter := csv.NewWriter(csvFile)
  defer csvFile.Close()

  for i, word := range realAnswers {
    fmt.Printf("Word %d/%d (%0.1f%%): %v\n", i, len(realAnswers), float64(i)/float64(len(realAnswers)), word )
    start := time.Now()
    numGuesses, guesses := guessWord(word, allowedWords, possibleWords)
    err = csvwriter.Write(guesses)
    if err != nil {
      panic(err)
    }
    csvwriter.Flush()
    fmt.Printf("Found in %d guesses in %0.2f seconds: %s\n", numGuesses, time.Since(start).Seconds(), guesses)
  }
}
