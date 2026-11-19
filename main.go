package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

// Job represents one image download task.
type Job struct {
	ID  int
	URL string
}

// This function downloads one image.
func downloadImage(client *http.Client, job Job) error {

	response, err := client.Get(job.URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %s", response.Status)
	}

	fileName := fmt.Sprintf("images/image_%03d.jpg", job.ID)

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, response.Body)
	closeErr := file.Close()

	if err != nil {
		os.Remove(fileName)
		return err
	}

	if closeErr != nil {
		os.Remove(fileName)
		return closeErr
	}

	return nil
}

// This function represents one worker.
func worker(workerID int, jobs <-chan Job, wg *sync.WaitGroup, client *http.Client) {

	defer wg.Done()

	for job := range jobs {

		fmt.Printf("Worker %d downloading image %d\n", workerID, job.ID)

		err := downloadImage(client, job)

		if err != nil {
			fmt.Printf("Worker %d: Image %d failed: %v\n", workerID, job.ID, err)
		} else {
			fmt.Printf("Worker %d: Image %d downloaded successfully\n", workerID, job.ID)
		}
	}
}

func main() {

	const totalWorkers = 10

	// Create images directory.
	err := os.MkdirAll("images", 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	// Open URLs file.
	file, err := os.Open("urls.txt")
	if err != nil {
		fmt.Println("Error opening urls.txt:", err)
		return
	}
	defer file.Close()

	// Read URLs from the file.
	var urls []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		url := scanner.Text()

		if url != "" {
			urls = append(urls, url)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading URLs:", err)
		return
	}

	fmt.Println("Total URLs:", len(urls))

	// Create HTTP client.
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create jobs channel.
	jobs := make(chan Job)

	var wg sync.WaitGroup

	// Start workers.
	for i := 1; i <= totalWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg, client)
	}

	// Send jobs to workers.
	for i, url := range urls {
		jobs <- Job{
			ID:  i + 1,
			URL: url,
		}
	}

	// Close channel after sending all jobs.
	close(jobs)

	// Wait for all workers to finish.
	wg.Wait()

	fmt.Println("All download jobs completed!")
}
