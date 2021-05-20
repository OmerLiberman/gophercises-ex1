package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type Problem struct {
	question string
	answer string
}

func ReadCsv(path *string) [][]string {
	file, err := os.Open(*path)
	if err != nil {
		fmt.Println("Some informative failure message")
		os.Exit(1)
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		fmt.Println("Some informative failure message")
		os.Exit(1)
	}

	return lines
}

func ParseLines(lines [][]string) []Problem {
	ret := make([]Problem, len(lines))
	for i, line := range lines {
		ret[i] = Problem{
			question: strings.TrimSpace(line[0]),
			answer: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

func RunGame(problems []Problem, timeLimit *int) (int, int) {
	correct := 0
	incorrect := 0

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i + 1, p.question)

		// This is an anonymous function which writes each time into
		// the channel the word read from the user.
		answerCh := make(chan string)
		go func() {
			var givenAnswer string
			fmt.Scanf("%s\n", &givenAnswer)
			answerCh <- givenAnswer
		} ()

		// select lets the goroutine wait on multiple communication operations.
		select {
			case <-timer.C:  // End everything when timer ticks off.
				fmt.Printf("Final score: correct %d, incorrect: %d", correct, incorrect)
				return correct, incorrect
			case answer := <-answerCh:
				if answer == p.answer {
					fmt.Println("Good answer !")
					correct++
				} else {
					fmt.Println("Wrong !")
					incorrect++
				}
		}
	}

	return correct, incorrect
}

func main() {
	// Param can be given via cmd line.
	csvFilePath := flag.String("csv", "problems.csv", "A csv file to be used")
	timeLimit := flag.Int("limit", 10, "Limit is in seconds")
	flag.Parse()

	lines := ReadCsv(csvFilePath)
	problems := ParseLines(lines)
	correct, incorrect := RunGame(problems, timeLimit)

	fmt.Printf("Correct answers: %d, Incorrect: %d", correct, incorrect)
}