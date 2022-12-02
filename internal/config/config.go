package config

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Problem struct {
	Question string
	Answer   string
}

type AppConfig struct {
	TimeLimit int
	QuizData  []Problem
}

func NewAppConfig() (*AppConfig, error) {
	fileName := flag.String("csv", "problems.csv", "csv file name in dir csv_files")
	timeLimit := flag.Int("t", 0, "limit time for answer")
	shuffleProblems := flag.Bool("r", false, "shuffleProblems")
	flag.Parse()

	file, err := os.OpenFile(fmt.Sprintf("../../internal/csv/%s", *fileName), os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}

	quizData, err := parseQuizData(file)
	if err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	if *shuffleProblems {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(quizData), func(i, j int) {
			quizData[i], quizData[j] = quizData[j], quizData[i]
		})
	}

	return &AppConfig{
		TimeLimit: *timeLimit,
		QuizData:  quizData,
	}, nil
}

func parseQuizData(file *os.File) ([]Problem, error) {
	r := csv.NewReader(file)
	data, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	quizElements := make([]Problem, len(data))
	for i, quizElement := range data {
		quizElements[i].Question = quizElement[0] + "="
		quizElements[i].Answer = quizElement[1]
	}
	return quizElements, nil
}
