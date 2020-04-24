// Copyright 2020 pyspa developers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	// Covid19JapanDailyDataURL is the URL of the daily patients data.
	Covid19JapanDailyDataURL = "https://raw.githubusercontent.com/kaz-ogiwara/covid19/master/data/prefectures.csv"

	// CSVNumField is the expected number of fields in the CSV file above.
	CSVNumField = 8

	// TimeZoneTokyo is the default timezone in this program
	TimeZoneTokyo = "Asia/Tokyo"
)

var (
	// DefaultCalendarStartDay is the default calendar start day.
	DefaultCalendarStartDay = time.Monday

	// DefaultCalendarBeginDate is the default calendar begin date.
	DefaultCalendarBeginDate time.Time
)

func init() {
	loc, err := time.LoadLocation(TimeZoneTokyo)
	if err != nil {
		loc = time.FixedZone(TimeZoneTokyo, 9*60*60)
	}
	time.Local = loc
	DefaultCalendarBeginDate = time.Date(2020, time.March, 1, 0, 0, 0, 0, time.Local)
}

func newCSVReader(r io.Reader) *csv.Reader {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = CSVNumField
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true
	return reader
}

// Calendar is the struct to display patients numbers
type Calendar struct {
	StartDay  time.Weekday
	BeginDate time.Time
}

// NewCalendar returns a reference to Calendar with default values
func NewCalendar() *Calendar {
	return &Calendar{
		StartDay:  DefaultCalendarStartDay,
		BeginDate: DefaultCalendarBeginDate,
	}
}

// Print returns io.Reader that contains the rendered data of the calendar.
func (c *Calendar) Print(r io.Reader) (io.Reader, error) {
	reader := newCSVReader(r)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	records = records[1:]
	var sb strings.Builder
	for _, rec := range records {
		record, err := NewRecord(rec)
		if err != nil {
			return nil, err
		}
		if record.Prefecture == TOKYO {
			sb.WriteString(fmt.Sprintf("%v %d\n", record.Date, record.Infected))
		}
	}

	return strings.NewReader(sb.String()), nil
}

// Record is the struct of the record from Toyo Keisai CSV.
type Record struct {
	Date         time.Time
	Prefecture   Prefecture
	Infected     int
	Hospitalized int
	Discharged   int
	Dead         int
}

// NewRecord returns the reference of Record from a raw CSV record.
func NewRecord(record []string) (*Record, error) {
	if len(record) != CSVNumField {
		return nil, fmt.Errorf("Number of fields (%v) in the CSV record is wrong: %v", len(record), record)
	}
	year, err := strconv.Atoi(record[0])
	if err != nil {
		return nil, err
	}
	month, err := strconv.Atoi(record[1])
	if err != nil {
		return nil, err
	}
	day, err := strconv.Atoi(record[2])
	if err != nil {
		return nil, err
	}
	infected, err := strconv.Atoi(record[4])
	if err != nil {
		return nil, err
	}
	hospitalized, err := strconv.Atoi(record[5])
	if err != nil {
		return nil, err
	}
	discharged, err := strconv.Atoi(record[6])
	if err != nil {
		return nil, err
	}
	dead, err := strconv.Atoi(record[7])
	if err != nil {
		return nil, err
	}
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	pref := PrefFromString(record[3])
	return &Record{
		Date:         date,
		Prefecture:   pref,
		Infected:     infected,
		Hospitalized: hospitalized,
		Discharged:   discharged,
		Dead:         dead,
	}, nil
}

func main() {
	resp, err := http.Get(Covid19JapanDailyDataURL)
	if err != nil {
		log.Fatalf("Error on fetching CSV data: %v", err)
	}
	defer resp.Body.Close()

	c := NewCalendar()
	r, err := c.Print(resp.Body)
	if err != nil {
		log.Fatalf("Error while convering io.Reader to calendar output: %v", err)
	}
	io.Copy(os.Stdout, r)
}
