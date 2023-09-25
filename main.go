package main

import (
    "encoding/csv"
    "log"
    "os"
    "sync"

    "github.com/lawzava/emailscraper"
)

func main() {
    // Read the list of domains from a CSV file
    domains, err := readDomainsFromCSV("domains.csv")
    if err != nil {
        log.Fatalf("Error reading CSV: %v", err)
    }

    // Create a wait group to wait for all scrapers to finish
    var wg sync.WaitGroup

    s := emailscraper.New(emailscraper.DefaultConfig())

    // Iterate over the domains and start a scraper for each
    for _, domain := range domains {
        // Increment the wait group counter
        wg.Add(1)

        // Start the scraper in a goroutine
        go func(domain string) {
            defer wg.Done()

            extractedEmails, err := s.Scrape(domain)
            if err != nil {
                panic(err)
            }
            appendEmailsToCSV("savedEmails.csv", domain, extractedEmails)
        }(domain)
    }

    // Wait for all scrapers to finish
    wg.Wait()
}

func readDomainsFromCSV(filename string) ([]string, error) {
    var domains []string

    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    for {
        record, err := reader.Read()
        if err != nil {
            break
        }
        domains = append(domains, record[0])
    }

    return domains, nil
}

func appendEmailsToCSV(filename string, domain string, emails []string) {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatalf("Error opening CSV file: %v", err)
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    for _, email := range emails {
        // Write the domain and email to the CSV file
        record := []string{domain, email}
        if err := writer.Write(record); err != nil {
            log.Fatalf("Error writing to CSV file: %v", err)
        }
    }

    writer.Flush()
    if err := writer.Error(); err != nil {
        log.Fatalf("Error flushing CSV writer: %v", err)
    }
}


