package utils

import (
	"strconv"
	"strings"
	"time"
)

/**
 * TimeFormat
 * @param time.Time time
 * @param string format 'Y-m-d H:i:s' like php
 * @return string
 */
func TimeFormat(t time.Time, format string) string {

	//year
	if strings.ContainsAny(format, "Y") {

		year := strconv.Itoa(t.Year())

		if strings.Count(format, "YY") == 1 && strings.Count(format, "Y") == 2 {
			format = strings.Replace(format, "yy", year[2:], 1)
		} else if strings.Count(format, "Y") == 1 {
			format = strings.Replace(format, "Y", year, 1)
		} else {
			panic("format year error! please 'YY' or 'Y'")
		}
	}

	//month
	if strings.ContainsAny(format, "m") {

		var month string

		if int(t.Month()) < 10 {
			month = "0" + strconv.Itoa(int(t.Month()))
		} else {
			month = strconv.Itoa(int(t.Month()))
		}

		if strings.Count(format, "m") == 1 {
			format = strings.Replace(format, "m", month, 1)
		} else {
			panic("format month error! please 'm'")
		}
	}

	//day
	if strings.ContainsAny(format, "d") {

		var day string

		if t.Day() < 10 {
			day = "0" + strconv.Itoa(t.Day())
		} else {
			day = strconv.Itoa(t.Day())
		}

		if strings.Count(format, "d") == 1 {
			format = strings.Replace(format, "d", day, 1)
		} else {
			panic("format day error! please 'd'")
		}
	}

	//hour
	if strings.ContainsAny(format, "H") {

		var hour string

		if t.Hour() < 10 {
			hour = "0" + strconv.Itoa(t.Hour())
		} else {
			hour = strconv.Itoa(t.Hour())
		}

		if strings.Count(format, "H") == 1 {
			format = strings.Replace(format, "H", hour, 1)
		} else {
			panic("format hour error! please 'H'")
		}
	}

	//minute
	if strings.ContainsAny(format, "i") {

		var minute string

		if t.Minute() < 10 {
			minute = "0" + strconv.Itoa(t.Minute())
		} else {
			minute = strconv.Itoa(t.Minute())
		}
		if strings.Count(format, "i") == 1 {
			format = strings.Replace(format, "i", minute, 1)
		} else {
			panic("format minute error! please 'i'")
		}
	}

	//second
	if strings.ContainsAny(format, "s") {

		var second string

		if t.Second() < 10 {
			second = "0" + strconv.Itoa(t.Second())
		} else {
			second = strconv.Itoa(t.Second())
		}

		if strings.Count(format, "s") == 1 {
			format = strings.Replace(format, "s", second, 1)
		} else {
			panic("format second error! please 's'")
		}
	}

	return format
}
