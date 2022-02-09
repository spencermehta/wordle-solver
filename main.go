package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
  possibleWords, allowedWords, realAnswers := wordLists()
  patterns := patterns()

  for i, answer := range realAnswers {
    fmt.Printf("Word %d/%d\t", i+1, len(realAnswers))
    guesses, numGuesses := completeWordle(answer, allowedWords, possibleWords, patterns)
    fmt.Printf("Completed in %d guesses: %s\n", numGuesses, guesses)
  }
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

func completeWordle(realWord string, allowedWords, possibleWords, patterns []string) ([]string, int) {
  guess := "tares" // hard coded first step since it takes forever to run
  allowedWords = removeWord(guess, allowedWords)
  var guesses []string
  guesses = append(guesses, guess)
  numGuesses := 1
  for {
    patt := generatePattern(realWord, guess)
    possibleWords = eliminateWords(guess, patt, possibleWords)
    guess = bestWord(allowedWords, possibleWords, patterns)
    numGuesses += 1
    allowedWords = removeWord(guess, allowedWords)
    guesses = append(guesses, guess)
    if guess == realWord {
      break;
    }
  }
  return guesses, numGuesses
}

func bestWord(allowedWords, possibleWords, patterns []string) string {
  var maxEi float64 = -1
  bestWord := ""
  for _, word := range allowedWords {
    ei := wordExpectedInformation(word, patterns, possibleWords)
    if ei > maxEi {
      maxEi = ei
      bestWord = word
    }
  }
  for _, word := range possibleWords {
    ei := wordExpectedInformation(word, patterns, possibleWords)
    if ei >= maxEi {
      maxEi = ei
      bestWord = word
    }
  }
  return bestWord
}

func wordExpectedInformation(guess string, patterns, possibleWords []string) float64 {
  var sum float64
  for _, pattern := range patterns {
    p := patternProbability(guess, pattern, possibleWords) 
    i := patternInformation(guess, pattern, possibleWords)
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

func patternInformation(guess, pattern string, possibleWords []string) float64 {
  return math.Log2(1.0/patternProbability(guess, pattern, possibleWords))
}

func patternProbability(guess, pattern string, possibleWords []string) float64 {
  sum := 0
  for _, word := range possibleWords {
    if matchesPattern(guess, pattern, word) {
      sum++
    }
  }
  return float64(sum)/float64(len(possibleWords))
}

func eliminateWords(guess, pattern string, possibleWords []string) []string {
  var newPossiblewords []string
  for _, word := range possibleWords {
    if matchesPattern(guess, pattern, word) {
      newPossiblewords = append(newPossiblewords, word)
    } else {
    }
  }
  return newPossiblewords
}

func matchesPattern(guess, pattern, word string) bool {
  for i, c := range pattern {
    switch c {
    case '0':
      if contains(word, []rune(guess)[i]) {
        return false
      }
    case '1':
      if !([]rune(guess)[i] != []rune(word)[i] && contains(word, []rune(guess)[i])) {
        return false
      }
    case '2':
      if []rune(guess)[i] != []rune(word)[i] {
        return false
      }
    }
  }
  return true
}

func generatePattern(realWord, guess string) string {
  pattern := ""
  for i, c := range guess {
    switch {
    case c == []rune(realWord)[i]:
      pattern += string('2')
    case c != []rune(realWord)[i] && contains(realWord, c):
      pattern += string('1')
    case !contains(realWord, c):
      pattern += string('0')
    default:
      pattern += string('!')
    }
  }
  return pattern
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

