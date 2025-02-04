//Brandon Allison
//Receipt-Processor
//2-3-2025
package main

import (
    "encoding/json"
    "fmt"
    "errors"
    "math"
    "net/http"
    "regexp"
    "strconv"
    "strings"
    "sync"
    "time"

    "github.com/google/uuid"
)

type Receipt struct {
    Retailer     string `json:"retailer"`
    PurchaseDate string `json:"purchaseDate"`
    PurchaseTime string `json:"purchaseTime"`
    Items        []Item `json:"items"`
    Total        string `json:"total"`
}

type Item struct {
    ShortDescription string `json:"shortDescription"`
    Price            string `json:"price"`
}

type ResponseID struct {
    ID string `json:"id"`
}

type ResponsePoints struct {
    Points int `json:"points"`
}

var (
    receipts = make(map[string]int)
    mutex    = &sync.Mutex{}
)

func processReceipt(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "The receipt is invalid", http.StatusMethodNotAllowed)
        return
    }

    var receipt Receipt
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&receipt); err != nil {
        http.Error(w, "The receipt is invalid", http.StatusBadRequest)
        return
    }

    if err := validateReceipt(receipt); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    id := uuid.New().String()
    points := calculatePoints(receipt)

    mutex.Lock()
    receipts[id] = points
    mutex.Unlock()

    response := ResponseID{ID: id}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func getPoints(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/receipts/")
    id = strings.TrimSuffix(id, "/points")

    mutex.Lock()
    points, exists := receipts[id]
    mutex.Unlock()

    if !exists {
        http.Error(w, "No receipt found for that ID", http.StatusNotFound)
        return
    }

    response := ResponsePoints{Points: points}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func calculatePoints(receipt Receipt) int {
    points := 0

    // One point for every alph char in name
    reg := regexp.MustCompile(`[a-zA-Z0-9]`)
    points += len(reg.FindAllString(receipt.Retailer, -1))

    // 50 points if the total is round dollar amount
    if strings.HasSuffix(receipt.Total, ".00") {
        points += 50
    }

    // 25 points if total is a multiple of 0.25
    total, err := strconv.ParseFloat(receipt.Total, 64)
    if err == nil && math.Mod(total, 0.25) == 0 {
        points += 25
    }

    // 5 points for every two
    points += (len(receipt.Items) / 2) * 5

    // Points for item description length
    for _, item := range receipt.Items {
        descLength := len(strings.TrimSpace(item.ShortDescription))
        if descLength%3 == 0 {
            price, err := strconv.ParseFloat(item.Price, 64)
            if err == nil {
                points += int(math.Ceil(price * 0.2))
            }
        }
    }

    // 6 points if day in purchase date is odd
    dateParts := strings.Split(receipt.PurchaseDate, "-")
    if len(dateParts) == 3 {
        day, err := strconv.Atoi(dateParts[2])
        if err == nil && day%2 == 1 {
            points += 6
        }
    }

    // 10 points if the time 2:00pm - 4:00pm
    purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
    if err == nil {
        hour, min := purchaseTime.Hour(), purchaseTime.Minute()
        if (hour == 14 && min >= 0) || (hour == 15 && min < 60) {
            points += 10
        }
    }

    return points
}
func validateReceipt(receipt Receipt) error {
    if receipt.Retailer == "" || receipt.PurchaseDate == "" || receipt.PurchaseTime == "" || receipt.Total == "" || len(receipt.Items) == 0 {
        return errors.New("Missing required fields")
    }
    if !isValidDate(receipt.PurchaseDate) {
        return errors.New("Invalid date format")
    }
    if !isValidTime(receipt.PurchaseTime) {
        return errors.New("Invalid time format")
    }
    if _, err := strconv.ParseFloat(receipt.Total, 64); err != nil {
        return errors.New("Invalid total format")
    }
    for _, item := range receipt.Items {
        if item.ShortDescription == "" || item.Price == "" {
            return errors.New("Invalid item format")
        }
        if _, err := strconv.ParseFloat(item.Price, 64); err != nil {
            return errors.New("Invalid item price format")
        }
    }
    return nil
}

func isValidDate(date string) bool {
    _, err := time.Parse("2006-01-02", date)
    return err == nil
}

func isValidTime(t string) bool {
    _, err := time.Parse("15:04", t)
    return err == nil
}

func main() {
    http.HandleFunc("/receipts/process", processReceipt)
    http.HandleFunc("/receipts/", getPoints)

    fmt.Println("Server started on :8080")
    http.ListenAndServe(":8080", nil)
}
