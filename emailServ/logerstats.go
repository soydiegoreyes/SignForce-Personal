package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"sync"
	"time"
)

type LogEntry struct {
	Timestamp time.Time
	Level     string
	Service   string
	Message   string
}

type Stats struct {
	TotalEntries   int
	CountByLevel   map[string]int
	CountByService map[string]int
	FirstEntry     time.Time
	LastEntry      time.Time
}

// Mantenemos tu variable global tal cual
var GlobalStats Stats

func LogGen() {
	// Configuración
	files := []string{"auth.log", "api.log", "db.log"}
	linesPerFile := 10000000 // Cambia esto a 1,000,000 para pruebas pesadas

	levels := []string{"INFO", "WARN", "ERROR"}
	services := []string{"auth", "api", "db", "storage", "cache"}
	messages := map[string][]string{
		"INFO":  {"request received", "user login", "connection established", "sync complete"},
		"WARN":  {"slow query detected", "high memory usage", "retry attempt 1", "deprecated api call"},
		"ERROR": {"invalid token", "database connection lost", "permission denied", "panic recovery"},
	}

	for _, fileName := range files {
		file, _ := os.Create(fileName)
		writer := bufio.NewWriter(file)

		fmt.Printf("Generando %s...\n", fileName)

		startTime := time.Now().Add(-24 * time.Hour) // Empezar hace 24 horas

		for i := 0; i < linesPerFile; i++ {
			// Simular que el tiempo avanza unos segundos por cada log
			startTime = startTime.Add(time.Duration(rand.Intn(10)) * time.Second)

			lvl := levels[rand.Intn(len(levels))]
			srv := services[rand.Intn(len(services))]
			msg := messages[lvl][rand.Intn(len(messages[lvl]))]

			line := fmt.Sprintf("%s %s service=%s message=\"%s\"\n",
				startTime.Format(time.RFC3339), lvl, srv, msg)

			writer.WriteString(line)
		}
		writer.Flush()
		file.Close()
	}
	fmt.Println("¡Listo! Archivos generados con éxito.")
}

func processMassive(source io.Reader, ch chan<- *LogEntry, wg *sync.WaitGroup) {
	defer wg.Done()

	scanner := bufio.NewScanner(source)
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 1024*1024)

	match := `(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}-\d{2}:\d{2})\s(ERROR|INFO|WARN)\sservice=(\w+)\s+message="([^"]+)"`
	reg := regexp.MustCompile(match)
	// VALIDACIÓN CRUCIAL:
	// Si matches es nil o no tiene los 5 elementos (el total + los 4 grupos), sáltalo.
	for scanner.Scan() {
		matches := reg.FindSubmatch(scanner.Bytes())
		// VALIDACIÓN CRUCIAL:
		// Si matches es nil o no tiene los 5 elementos (el total + los 4 grupos), sáltalo.
		if len(matches) < 5 {
			continue
		}
		t_parsed, err := time.Parse(time.RFC3339, string(matches[1]))
		if err != nil {
			fmt.Println("Error procesando log")
			continue
		}

		ch <- &LogEntry{
			Timestamp: t_parsed,
			Level:     string(matches[2]),
			Service:   string(matches[3]),
			Message:   string(matches[4]),
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error de lectura: %v\n", err)
	}
}

func Merge(sources []io.Reader) ([]*LogEntry, Stats, error) {
	var wg sync.WaitGroup
	var entries []*LogEntry
	ch := make(chan *LogEntry)

	// Inicializamos los mapas aquí porque en Go un mapa nil da error al escribir
	GlobalStats.CountByLevel = make(map[string]int)
	GlobalStats.CountByService = make(map[string]int)

	// 1. Lanzamos las gorutinas
	for _, source := range sources {
		wg.Add(1)
		go processMassive(source, ch, &wg)
	}

	// 2. LA CLAVE: Una gorutina "asistente" que cierra el canal cuando todos acaben
	go func() {
		wg.Wait()
		close(ch)
	}()

	// 3. El receptor central (Hilo principal)
	// Este bucle "range" lee del canal hasta que se cierre.
	// Al ser el único que escribe en GlobalStats, ¡no necesitamos Mutex!
	for entry := range ch {
		entries = append(entries, entry)

		// Actualizamos estadísticas de forma segura y ordenada
		GlobalStats.TotalEntries++
		GlobalStats.CountByLevel[entry.Level]++
		GlobalStats.CountByService[entry.Service]++

		if GlobalStats.FirstEntry.IsZero() || entry.Timestamp.Before(GlobalStats.FirstEntry) {
			GlobalStats.FirstEntry = entry.Timestamp
		}
		if entry.Timestamp.After(GlobalStats.LastEntry) {
			GlobalStats.LastEntry = entry.Timestamp
		}
	}

	// 4. Ordenar al final
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})

	return entries, GlobalStats, nil
}

func main() {
	//LogGen()
	// Abrir archivos (Simplificado para el ejemplo)
	f1, _ := os.Open("auth.log")
	f2, _ := os.Open("api.log")
	f3, _ := os.Open("db.log")

	// Solo pasamos los que se abrieron correctamente
	var readers []io.Reader
	for _, f := range []*os.File{f1, f2, f3} {
		if f != nil {
			defer f.Close()
			readers = append(readers, f)
		}
	}

	entries, stats, err := Merge(readers)
	if err != nil {
		log.Fatal(err)
	}

	for i, e := range entries {
		if i == 100 {
			break
		}
		fmt.Printf("[%s] %s %s: %s\n", e.Level, e.Timestamp.Format(time.RFC3339), e.Service, e.Message)
	}
	fmt.Printf("\nTotal: %d, Errors: %d\n", stats.TotalEntries, stats.CountByLevel["ERROR"])
}
