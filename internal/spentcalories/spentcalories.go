package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

const (
	activityRunning = "Бег"
	activityWalking = "Ходьба"
)

func parseTraining(data string) (int, string, time.Duration, error) {
	stringsSlice := strings.Split(data, ",")
	if len(stringsSlice) != 3 {
		return 0, "", 0, fmt.Errorf("data не соотвествует формату")
	}

	stepsCount, err := strconv.Atoi(stringsSlice[0])
	if err != nil {
		return 0, "", 0, err
	}
	if stepsCount <= 0 {
		return 0, "", 0, fmt.Errorf("Шаги ноль")
	}

	duration, err := time.ParseDuration(stringsSlice[2])
	if err != nil {
		return 0, "", 0, err
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("Продолжительность ноль")
	}

	return stepsCount, stringsSlice[1], duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := float64(height) * stepLengthCoefficient
	return stepLength * float64(steps) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	return distance(steps, height) / float64(duration.Hours())
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var spentCalories float64
	switch activityType {
	case activityWalking:
		spentCalories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case activityRunning:
		spentCalories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	distance := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activityType, duration.Hours(), distance, speed, spentCalories), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("некорректные данные")
	}

	return (weight * meanSpeed(steps, height, duration) * duration.Minutes()) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("некорректные данные")
	}

	calories := (meanSpeed(steps, height, duration) * weight * duration.Minutes()) / minInH
	return calories * walkingCaloriesCoefficient, nil
}
