# Hacker News Top Stories 

Simple program to email the top stories from Hacker News 

## Usage

```go
from := mail.Account{"First LastName", "from@example.com", "password", "smtp.example.com", 587}
hackernews.Email(from, "to@example.com", 10)
```