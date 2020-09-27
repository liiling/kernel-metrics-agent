package otelstats

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	// Keywords for parsing .schema file
	labelKeyword  = "LABEL"
	metricKeyword = "METRIC"
	nameKeyword   = "NAME"
	flagKeyword   = "FLAG"
	typeKeyword   = "TYPE"
	descKeyword   = "DESC"

	// Flags found in .schema file
	cumulative = "CUMULATIVE"
	gauge      = "GAUGE"

	// Types found in .schema file
	intType   = "INT"
	floatType = "FLOAT"
)

// metricSchema contains information about one metric from a statsfs file,
// Name = metric name
// Path = from where the metric is retrieved
// Label = a list of key-value pairs of labels applied to the metric
// Flag = cumulative or gauge
// Type = int or float
// Desc = description of the metric
type metricSchema struct {
	mname   string
	mlabels []metricLabel
	mflag   string
	mtype   string
	mdesc   string
}

type metricLabel struct {
	key   string
	value string
}

// ParseSchema parses a statsfs .schema file
// returns a map with key = metric name and
// value = metricSchema
func parseSchema(path string) ([]metricSchema, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open .schema file at %v", path)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	labels := []metricLabel{}
	metrics := []metricSchema{}

	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, " ")
		switch words[0] {
		case labelKeyword:
			if labels, err = parseLabels(scanner); err != nil {
				return nil, fmt.Errorf("failed to parse schema at %s: %s", path, err)
			}
		case metricKeyword:
			if m, err := parseMetric(scanner); err != nil {
				return nil, fmt.Errorf("failed to parse schema at %s: %s", path, err)
			} else {
				metrics = append(metrics, *m)
			}
		case "":
			continue
		}
	}

	fmt.Printf("-------------------------------------\n")
	fmt.Printf(".schema file at %v\n", path)
	fmt.Printf("Prior to label update: metric = %v\n\n", metrics)

	fmt.Printf("Update labels...\n")
	// apply labels to all metrics listed in the .schema file
	for _, m := range metrics {
		m.mlabels = labels
		fmt.Printf("m = %v, m.mlabels: %v\n", m, m.mlabels)
	}

	fmt.Printf("\nAfter label update:\n")
	for _, m := range metrics {
		fmt.Printf("m = %v, m.mlabels: %v\n", m, m.mlabels)
	}
	fmt.Printf("-------------------------------------\n")

	return metrics, nil
}

// parseLabel parses the <labels> section of .schema file
// <labels> ::= "LABEL\n" (<label_key> " " <label_value> "\n")* "\n"
// <label_key> ::= <strset>+
// <label_value> ::= <strset>+
// <strset> ::= [A-Za-z_0-9]
func parseLabels(scanner *bufio.Scanner) ([]metricLabel, error) {
	labels := []metricLabel{}

	for scanner.Scan() {
		if line := scanner.Text(); line == "" {
			// Empty line denotes the end of <labels>
			break
		} else {
			words := strings.Split(strings.TrimSpace(line), " ")
			if len(words) != 2 {
				return nil, fmt.Errorf("syntax error in <labels> section")
			}
			labels = append(labels, metricLabel{key: words[0], value: words[1]})
		}
	}
	return labels, nil
}

// parseMetric parses one <metric_schema> section of .schema file
// <metric_schema> ::= "METRIC\n" <metric_name> <flag> <type> <description> “\n”
// <metric_name> ::= “NAME ” <strset>+ “\n”
// <flag> ::= “FLAG “ (CUMULATIVE | GAUGE) “\n”
// <type> ::= “TYPE “ (INT | FLOAT) “\n”
// <description> ::= “DESC “ <escaped-ASCII> “\n”
// <strset> ::= [A-Za-z_0-9]
func parseMetric(scanner *bufio.Scanner) (*metricSchema, error) {
	m := metricSchema{}

	for scanner.Scan() {
		if line := scanner.Text(); line == "" {
			// Empty line denotes the end of <metric_schema>
			break
		} else {
			words := strings.Split(line, " ")
			if len(words) != 2 && words[0] != descKeyword {
				return nil, fmt.Errorf("syntax error in <metric_schema> section")
			}
			val := words[1]

			switch words[0] {
			case nameKeyword:
				m.mname = val
			case flagKeyword:
				if val == cumulative {
					m.mflag = cumulative
				} else if val == gauge {
					m.mflag = gauge
				} else {
					return nil, fmt.Errorf("<metric_schema> unknown FLAG %s", val)
				}
			case typeKeyword:
				if val == intType {
					m.mtype = intType
				} else if val == floatType {
					m.mtype = floatType
				} else {
					return nil, fmt.Errorf("<metric_schema> unknown TYPE %s", val)
				}
			case descKeyword:
				m.mdesc = strings.Join(words[1:], " ")
			default:
				return nil, fmt.Errorf("<metric_schema> unknown keyword %s", words[0])
			}
		}
	}
	return &m, nil
}
