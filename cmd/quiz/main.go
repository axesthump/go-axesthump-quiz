package main

import (
	"bufio"
	"context"
	"fmt"
	"go-axesthump-quiz/internal/config"
	"os"
	"strings"
	"time"
)

type answers struct {
	right int
	wrong int
}

func (a *answers) String() string {
	return fmt.Sprintf("Your right answers: %d! Your wrong answers: %d\n", a.right, a.wrong)
}

func main() {
	conf, err := config.NewAppConfig()
	if err != nil {
		exit(fmt.Sprintf("Err: %v", err))
	}
	ctx := context.Background()
	var cancel context.CancelFunc

	if conf.TimeLimit != 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(conf.TimeLimit))
		defer cancel()
	}

	userAnswers := make(chan answers)
	go playGame(ctx, userAnswers, conf)
	result := <-userAnswers
	fmt.Println()
	fmt.Printf("%s\n", &result)

}

func playGame(ctx context.Context, ch chan answers, conf *config.AppConfig) {
	sc := bufio.NewScanner(os.Stdin)
	userAnswers := answers{}
	userAnswerChan := make(chan string)

	for i := 0; i < len(conf.QuizData); i++ {
		quizElement := conf.QuizData[i]
		fmt.Print(quizElement.Question)
		go func() {
			sc.Scan()
			userAnswer := sc.Text()
			userAnswer = strings.TrimSpace(userAnswer)
			userAnswerChan <- userAnswer
		}()
		select {
		case userAnswer := <-userAnswerChan:
			if userAnswer == quizElement.Answer {
				userAnswers.right += 1
			} else {
				userAnswers.wrong += 1
			}
		case <-ctx.Done():
			ch <- userAnswers
			return
		}
	}
	ch <- userAnswers
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
