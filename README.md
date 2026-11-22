install modules using go mod tidy

### post a task

curl http://localhost:8080/tasks \
--include \
--header "Content-Type: application/json" \
--request "POST" \
--data '{"task_id": "4", "title": "Chapter 8", "description": "Read Chapter 8 about generics", "is_completed": false}'

Concurrency: Gin handles requests in separate goroutines, so concurrent POSTs to a shared slice is a data race. For the in-memory stage, guard it with a sync.Mutex:
