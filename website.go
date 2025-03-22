package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Booking struct {
	UID    string
	From   string
	To     string
	Mode   string
	Date   string
	Cost   float64
	QRData string
}

var baseModeCosts = map[string]float64{
	"Bus":    500.00,
	"Train":  1000.00,
	"Flight": 5000.00,
}

var destinationMultipliers = map[string]float64{
	"Delhi": 1.0, "Mumbai": 1.1, "Bangalore": 1.2, "Chennai": 1.1, "Kolkata": 1.3, "Hyderabad": 1.2, "Pune": 1.0,
	"Ahmedabad": 1.1, "Jaipur": 1.0, "Lucknow": 1.2, "Chandigarh": 1.1, "Bhopal": 1.0, "Indore": 1.1,
	"Guwahati": 1.3, "Kochi": 1.4, "Patna": 1.2, "Nagpur": 1.1, "Bhubaneswar": 1.3, "Visakhapatnam": 1.2, "Coimbatore": 1.1,
	"Goa": 1.5, "Thrissur": 1.3, "Srinagar": 1.7, "Manali": 1.6,
}

var destinations = []string{
	"Delhi", "Mumbai", "Bangalore", "Chennai", "Kolkata", "Hyderabad", "Pune",
	"Ahmedabad", "Jaipur", "Lucknow", "Chandigarh", "Bhopal", "Indore",
	"Guwahati", "Kochi", "Patna", "Nagpur", "Bhubaneswar", "Visakhapatnam", "Coimbatore",
	"Goa", "Thrissur", "Srinagar", "Manali",
}

var travelModes = []string{"Bus", "Train", "Flight"}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/book", bookHandler)
	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	data := struct {
		Title        string
		Message      string
		Destinations []string
		Modes        []string
	}{
		Title:        "Book Your Ticket",
		Message:      "Fast & Easy Ticket Booking",
		Destinations: destinations,
		Modes:        travelModes,
	}
	tmpl.Execute(w, data)
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Request", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form Parsing Error", http.StatusInternalServerError)
		return
	}

	from := r.FormValue("from")
	to := r.FormValue("to")
	date := r.FormValue("date")
	mode := r.FormValue("mode")

	parsedDate, _ := time.Parse("2006-01-02", date)
	formattedDate := parsedDate.Format("02-01-2006")
	bookingID := strconv.FormatInt(time.Now().UnixNano(), 10)

	baseCost := baseModeCosts[mode]
	multiplier := destinationMultipliers[to]
	cost := baseCost * multiplier

	qrData := fmt.Sprintf("BookingID: %s | From: %s | To: %s | Mode: %s | Date: %s | Cost: ₹%.2f", bookingID, from, to, mode, formattedDate, cost)
	ticket := Booking{
		UID:    bookingID,
		From:   from,
		To:     to,
		Mode:   mode,
		Date:   formattedDate,
		Cost:   cost,
		QRData: qrData,
	}
	tmpl := template.Must(template.ParseFiles("ticket.html"))
	tmpl.Execute(w, ticket)
}