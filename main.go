package main

import (
	"flag"
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"
)

const ITERATIONS = 100_000_000

// leibnizChunk вычисляет сумму членов ряда Лейбница
// для нечётных знаменателей от startTerm до endTerm (включительно).
// startTerm и endTerm - это индексы членов ряда (0, 1, 2, ...),
// отсюда знаменатель = 2*index + 1, знак = (-1)^index
func leibnizChunk(startTerm, endTerm int, wg *sync.WaitGroup, results []float64, slot int) {
	defer wg.Done()

	sum := 0.0
	for i := startTerm; i <= endTerm; i++ {
		denominator := float64(2*i + 1)

		if i%2 == 0 {
			sum += 1.0 / denominator
		} else {
			sum -= 1.0 / denominator
		}
	}

	results[slot] = sum
}

func main() {
	iterations := flag.Int("it", ITERATIONS, "Количество итераций уточнения числа π")
	flag.Parse()

	numWorkers := min(*iterations, runtime.GOMAXPROCS(0))

	results := make([]float64, numWorkers)
	chunkSize := *iterations / numWorkers

	var wg sync.WaitGroup

	start := time.Now()

	for i := range numWorkers {
		startTerm := i * chunkSize
		endTerm := startTerm + chunkSize - 1

		if i == numWorkers-1 {
			endTerm = *iterations - 1
		}

		wg.Add(1)
		go leibnizChunk(startTerm, endTerm, &wg, results, i)
	}

	wg.Wait()

	pi := 0.0
	for _, partial := range results {
		pi += partial * 4
	}

	elapsed := time.Since(start)

	fmt.Printf("Воркеров:     %d\n", numWorkers)
	fmt.Printf("Итераций:     %d\n", *iterations)
	fmt.Printf("Число π:      %.16f\n", pi)
	fmt.Printf("Эталон π:     %.16f\n", math.Pi)
	fmt.Printf("Погрешность:  %.16f\n", pi-math.Pi)
	fmt.Printf("Время работы: %.3f сек\n", elapsed.Seconds())
}
