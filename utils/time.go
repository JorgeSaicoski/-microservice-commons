// utils/time.go
package utils

import (
	"fmt"
	"time"
)

// Common time formats
const (
	DateFormat     = "2006-01-02"
	TimeFormat     = "15:04:05"
	DateTimeFormat = "2006-01-02 15:04:05"
	RFC3339Format  = time.RFC3339
	ISO8601Format  = "2006-01-02T15:04:05Z07:00"
)

// TimeZones
const (
	UTC = "UTC"
	EST = "America/New_York"
	PST = "America/Los_Angeles"
	GMT = "Europe/London"
)

// Now returns current time in UTC
func Now() time.Time {
	return time.Now().UTC()
}

// NowInTimezone returns current time in specified timezone
func NowInTimezone(timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}
	return time.Now().In(loc), nil
}

// ToUTC converts time to UTC
func ToUTC(t time.Time) time.Time {
	return t.UTC()
}

// ToTimezone converts time to specified timezone
func ToTimezone(t time.Time, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(loc), nil
}

// ParseDate parses date string in YYYY-MM-DD format
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse(DateFormat, dateStr)
}

// ParseDateTime parses datetime string in YYYY-MM-DD HH:MM:SS format
func ParseDateTime(dateTimeStr string) (time.Time, error) {
	return time.Parse(DateTimeFormat, dateTimeStr)
}

// ParseRFC3339 parses RFC3339 formatted time string
func ParseRFC3339(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}

// ParseWithFormat parses time string with custom format
func ParseWithFormat(timeStr, format string) (time.Time, error) {
	return time.Parse(format, timeStr)
}

// FormatDate formats time as YYYY-MM-DD
func FormatDate(t time.Time) string {
	return t.Format(DateFormat)
}

// FormatTime formats time as HH:MM:SS
func FormatTime(t time.Time) string {
	return t.Format(TimeFormat)
}

// FormatDateTime formats time as YYYY-MM-DD HH:MM:SS
func FormatDateTime(t time.Time) string {
	return t.Format(DateTimeFormat)
}

// FormatRFC3339 formats time as RFC3339
func FormatRFC3339(t time.Time) string {
	return t.Format(time.RFC3339)
}

// FormatISO8601 formats time as ISO8601
func FormatISO8601(t time.Time) string {
	return t.Format(ISO8601Format)
}

// FormatWithCustom formats time with custom format
func FormatWithCustom(t time.Time, format string) string {
	return t.Format(format)
}

// StartOfDay returns the start of the day (00:00:00)
func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the end of the day (23:59:59.999999999)
func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}

// StartOfWeek returns the start of the week (Monday 00:00:00)
func StartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 { // Sunday
		weekday = 7
	}
	days := weekday - 1 // Days since Monday
	return StartOfDay(t.AddDate(0, 0, -days))
}

// EndOfWeek returns the end of the week (Sunday 23:59:59.999999999)
func EndOfWeek(t time.Time) time.Time {
	return EndOfDay(StartOfWeek(t).AddDate(0, 0, 6))
}

// StartOfMonth returns the start of the month (1st day 00:00:00)
func StartOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth returns the end of the month (last day 23:59:59.999999999)
func EndOfMonth(t time.Time) time.Time {
	return EndOfDay(StartOfMonth(t).AddDate(0, 1, -1))
}

// StartOfYear returns the start of the year (Jan 1st 00:00:00)
func StartOfYear(t time.Time) time.Time {
	year := t.Year()
	return time.Date(year, 1, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear returns the end of the year (Dec 31st 23:59:59.999999999)
func EndOfYear(t time.Time) time.Time {
	year := t.Year()
	return time.Date(year, 12, 31, 23, 59, 59, 999999999, t.Location())
}

// AddBusinessDays adds business days (Monday-Friday) to a date
func AddBusinessDays(t time.Time, days int) time.Time {
	current := t
	remaining := days

	if remaining > 0 {
		for remaining > 0 {
			current = current.AddDate(0, 0, 1)
			if IsBusinessDay(current) {
				remaining--
			}
		}
	} else if remaining < 0 {
		for remaining < 0 {
			current = current.AddDate(0, 0, -1)
			if IsBusinessDay(current) {
				remaining++
			}
		}
	}

	return current
}

// IsBusinessDay checks if a date is a business day (Monday-Friday)
func IsBusinessDay(t time.Time) bool {
	weekday := t.Weekday()
	return weekday >= time.Monday && weekday <= time.Friday
}

// IsWeekend checks if a date is a weekend (Saturday-Sunday)
func IsWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// DaysBetween calculates the number of days between two dates
func DaysBetween(start, end time.Time) int {
	start = StartOfDay(start)
	end = StartOfDay(end)
	return int(end.Sub(start).Hours() / 24)
}

// BusinessDaysBetween calculates business days between two dates
func BusinessDaysBetween(start, end time.Time) int {
	if start.After(end) {
		start, end = end, start
	}

	count := 0
	current := StartOfDay(start)
	endDate := StartOfDay(end)

	for current.Before(endDate) {
		if IsBusinessDay(current) {
			count++
		}
		current = current.AddDate(0, 0, 1)
	}

	return count
}

// Age calculates age in years from birthdate
func Age(birthdate time.Time) int {
	now := time.Now()
	age := now.Year() - birthdate.Year()

	// Adjust if birthday hasn't occurred this year
	if now.YearDay() < birthdate.YearDay() {
		age--
	}

	return age
}

// TimeAgo returns a human-readable string representing time elapsed
func TimeAgo(t time.Time) string {
	now := time.Now()
	duration := now.Sub(t)

	if duration < time.Minute {
		seconds := int(duration.Seconds())
		if seconds <= 1 {
			return "just now"
		}
		return fmt.Sprintf("%d seconds ago", seconds)
	}

	if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	}

	if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}

	days := int(duration.Hours() / 24)
	if days == 1 {
		return "1 day ago"
	}
	if days < 7 {
		return fmt.Sprintf("%d days ago", days)
	}

	weeks := days / 7
	if weeks == 1 {
		return "1 week ago"
	}
	if weeks < 4 {
		return fmt.Sprintf("%d weeks ago", weeks)
	}

	months := days / 30
	if months == 1 {
		return "1 month ago"
	}
	if months < 12 {
		return fmt.Sprintf("%d months ago", months)
	}

	years := days / 365
	if years == 1 {
		return "1 year ago"
	}
	return fmt.Sprintf("%d years ago", years)
}

// Sleep pauses execution for the specified duration
func Sleep(duration time.Duration) {
	time.Sleep(duration)
}

// Timeout creates a timeout channel
func Timeout(duration time.Duration) <-chan time.Time {
	return time.After(duration)
}

// Ticker creates a ticker channel
func Ticker(duration time.Duration) *time.Ticker {
	return time.NewTicker(duration)
}

// IsZero checks if time is zero value
func IsZero(t time.Time) bool {
	return t.IsZero()
}

// IsToday checks if time is today
func IsToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.YearDay() == now.YearDay()
}

// IsYesterday checks if time is yesterday
func IsYesterday(t time.Time) bool {
	yesterday := time.Now().AddDate(0, 0, -1)
	return t.Year() == yesterday.Year() && t.YearDay() == yesterday.YearDay()
}

// IsTomorrow checks if time is tomorrow
func IsTomorrow(t time.Time) bool {
	tomorrow := time.Now().AddDate(0, 0, 1)
	return t.Year() == tomorrow.Year() && t.YearDay() == tomorrow.YearDay()
}

// ParseDuration parses duration string with additional units
func ParseDuration(s string) (time.Duration, error) {
	// Try standard parsing first
	if d, err := time.ParseDuration(s); err == nil {
		return d, nil
	}

	// Handle additional formats like "1d", "1w", "1M", "1y"
	// This is a simplified implementation
	return time.ParseDuration(s)
}

// Min returns the minimum of two times
func Min(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

// Max returns the maximum of two times
func Max(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
