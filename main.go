package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
  possibleWords, allowedWords, _ := wordLists()
  patterns := patterns()

  fmt.Println(bestWord(allowedWords, possibleWords, patterns))
}

func wordLists() ([]string, []string, []string) {
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

  return possibleWords, allowedWords, realAnswers
}

func bestWord( allowedWords, possibleWords, patterns []string) string {
  var maxEi float64
  bestWord := ""
  for i, word := range allowedWords {
    fmt.Printf("%d/%d\r\n", i, len(allowedWords))
    ei := wordExpectedInformation(word, patterns, possibleWords)
    if ei > maxEi {
      maxEi = ei
      bestWord = word
    }
  }
  return bestWord
}

func wordExpectedInformation(realWord string, patterns, possibleWords []string) float64 {
  var sum float64
  for _, pattern := range patterns {
    p := patternProbability(realWord, pattern, possibleWords) 
    i := patternInformation(realWord, pattern, possibleWords)
    var ei float64
    if !math.IsInf(i, 0){
      ei = p * i
    } else {
      ei = 0.0
    }
    sum += ei
  }
  return sum
}

func patternInformation(realWord, pattern string, possibleWords []string) float64 {
  return math.Log2(1.0/patternProbability(realWord, pattern, possibleWords))
}

func patternProbability(realWord, pattern string, possibleWords []string) float64 {
  sum := 0
  for _, word := range possibleWords {
    if matchesPattern(realWord, pattern, word) {
      sum++
    }
  }
  return float64(sum)/float64(len(possibleWords))
}

func matchesPattern(realWord, pattern, word string) bool {
  for i, c := range pattern {
    switch c {
    case '0':
      if contains(word, []rune(realWord)[i]) {
        return false
      }
    case '1':
      if !([]rune(realWord)[i] != []rune(word)[i] && contains(word, []rune(realWord)[i])) {
        return false
      }
    case '2':
      if []rune(realWord)[i] != []rune(word)[i] {
        return false
      }
    }
  }
  return true
}

func patterns() []string {
  /*
  3^5 options
  */
  var patterns []string
  matchPossibilities := []rune{'0','1','2'}
  for _, first := range matchPossibilities {
    word := string(first)
    for _, second := range matchPossibilities {
      word += string(second)
      for _, third := range matchPossibilities {
        word += string(third)
        for _, fourth := range matchPossibilities {
          word += string(fourth)
          for _, fifth := range matchPossibilities {
            word += string(fifth)
            patterns = append(patterns, word)
            word = word[:len(word)-1]
          }
            word = word[:len(word)-1]
        }
            word = word[:len(word)-1]
      }
            word = word[:len(word)-1]
    }
            word = word[:len(word)-1]
  }

  return patterns
}

func contains(arr string, c rune) bool {
  var isIn bool
  for _, char := range arr {
    if c == char {
      isIn = true
    }
  }
  return isIn
}

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

