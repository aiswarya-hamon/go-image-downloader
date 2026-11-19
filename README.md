# Go Concurrent Image Downloader

A concurrent image downloader built using Go.

This project downloads multiple images from URLs using goroutines, channels, and a worker pool.

## Features

- Reads image URLs from a text file.
- Uses goroutines for concurrent execution.
- Implements a worker pool with 10 workers.
- Uses channels to distribute download jobs.
- Downloads images using HTTP GET requests.
- Saves downloaded images in the `images` directory.
- Handles HTTP, network, and file errors.
- Uses a 30-second HTTP client timeout.
- Uses WaitGroup to wait for all workers to finish.

## Technologies Used

- Go
- Goroutines
- Channels
- Worker Pool
- WaitGroup
- HTTP Client
- File Handling

## Project Structure

```text
concurrent-image-downloader/
├── main.go
├── go.mod
├── README.md
├── .gitignore
├── urls.example.txt
└── images/
```

## How to Run

1. Clone the repository.

```bash
git clone https://github.com/YOUR-USERNAME/concurrent-image-downloader.git
```

2. Navigate into the project.

```bash
cd concurrent-image-downloader
```

3. Create your URL file from the example.

```bash
cp urls.example.txt urls.txt
```

4. Add your image URLs to `urls.txt`, one URL per line.

5. Run the program.

```bash
go run .
```

Downloaded images will be saved in the `images` directory.

## Concurrency

The application uses 10 worker goroutines. Each worker receives a download job from a shared jobs channel and downloads the image independently.

A WaitGroup ensures that the main function waits until all workers have completed their jobs.
