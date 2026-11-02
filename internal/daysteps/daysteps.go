package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	stringsSlice := strings.Split(data, ",")
	if len(stringsSlice) != 2 {
		return 0, 0, fmt.Errorf("data не соотвествует формату")
	}

	stepsCount, err := strconv.Atoi(stringsSlice[0])
	if err != nil {
		return 0, 0, err
	}
	if stepsCount <= 0 {
		return 0, 0, fmt.Errorf("Количество шагов должно быть больше нуля")
	}

	duration, err := time.ParseDuration(stringsSlice[1])
	if err != nil {
		return 0, 0, err
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("Продолжительность ноль")
	}

	return stepsCount, duration, nil

}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		return ""
	}

	if steps <= 0 {
		return ""
	}

	distanceKm := float64(steps) * stepLength / mInKm
	spentCalories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceKm, spentCalories)
}
