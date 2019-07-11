package grafana_vertica_service

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/pkg/errors"
)

type myRegexp struct {
	*regexp.Regexp
}

const groupBy = "group_by"

var errUnexpectedType = errors.New("Non-numeric type could not be converted to float")

var replacer = strings.NewReplacer(";", "")
var myExp = myRegexp{regexp.MustCompile(`(?P<group_by>GROUP BY) (?P<fields>(?:\S*,\s)+\S*)`)}

func stripSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}

		return r
	}, str)
}

func (r *myRegexp) FindGroupBy(s string) []string {
	var res []string
	var groupBys string

	match := r.FindStringSubmatch(s)
	if match == nil {
		return res
	}

	for i, name := range r.SubexpNames() {
		if i == 0 || name == groupBy {
			continue
		}
		groupBys = replacer.Replace(stripSpaces(match[i]))
	}
	res = strings.Split(groupBys, ",")
	return res
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func getTableType(unk interface{}) string {
	switch i := unk.(type) {
	case float64, float32, int, int32, int64:
		return "Number"
	case string:
		if _, err := strconv.ParseFloat(i, 64); err == nil {
			return "Number"
		}
		if _, err := time.Parse("2006-01-02", i); err == nil {
			return "Date"
		}
		return "String"
	default:
		return "String"
	}
}

func getFloat(unk interface{}, isTime bool) (float64, error) {
	var result float64
	switch i := unk.(type) {
	case float64:
		result = i
	case float32:
		result = float64(i)
	case int64:
		result = float64(i)
	case int32:
		result = float64(i)
	case int:
		result = float64(i)
	case uint64:
		result = float64(i)
	case uint32:
		result = float64(i)
	case uint:
		result = float64(i)
	case string:
		if s, err := strconv.ParseFloat(i, 64); err == nil {
			result = s
		} else {
			return math.NaN(), errUnexpectedType
		}
	default:
		return math.NaN(), errUnexpectedType
	}
	if isTime {
		return result * 1000, nil
	}
	return result, nil
}

func getString(unk interface{}) (string, error) {
	switch i := unk.(type) {
	case string:
		return string(i), nil
	default:
		return "", errUnexpectedType
	}
}
