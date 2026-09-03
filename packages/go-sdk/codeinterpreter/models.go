package codeinterpreter

import (
	"encoding/json"
	"fmt"
	"slices"
	"time"
)

// OutputMessage is one stdout or stderr line.
type OutputMessage struct {
	Line      string
	Timestamp int64
	Error     bool
}

func (message OutputMessage) String() string { return message.Line }

// ExecutionError is a kernel error and traceback.
type ExecutionError struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Traceback string `json:"traceback"`
}

func (e ExecutionError) Error() string { return fmt.Sprintf("%s: %s", e.Name, e.Value) }

type Logs struct {
	Stdout []string `json:"stdout"`
	Stderr []string `json:"stderr"`
}
type Execution struct {
	Results        []Result        `json:"results"`
	Logs           Logs            `json:"logs"`
	Error          *ExecutionError `json:"error,omitempty"`
	ExecutionCount int             `json:"execution_count,omitempty"`
}

func (execution Execution) Text() string {
	for _, result := range execution.Results {
		if result.IsMainResult {
			return result.Text
		}
	}
	return ""
}

// RawData preserves every MIME representation returned by the kernel.
type RawData map[string]json.RawMessage
type Result struct {
	Text, HTML, Markdown, SVG, PNG, JPEG, PDF, LaTeX, JSON, JavaScript string
	Data                                                               map[string]any
	Chart                                                              *Chart
	Extra                                                              map[string]json.RawMessage
	Raw                                                                RawData
	IsMainResult                                                       bool
}

func (result Result) Formats() []string {
	formats := make([]string, 0, len(result.Raw))
	for key := range result.Raw {
		formats = append(formats, key)
	}
	slices.Sort(formats)
	return formats
}

type ChartType string

const (
	ChartLine          ChartType = "line"
	ChartScatter       ChartType = "scatter"
	ChartBar           ChartType = "bar"
	ChartPie           ChartType = "pie"
	ChartBoxAndWhisker ChartType = "box_and_whisker"
	ChartSuper         ChartType = "superchart"
	ChartUnknown       ChartType = "unknown"
)

type ScaleType string

const (
	ScaleLinear      ScaleType = "linear"
	ScaleDatetime    ScaleType = "datetime"
	ScaleCategorical ScaleType = "categorical"
	ScaleLog         ScaleType = "log"
	ScaleSymlog      ScaleType = "symlog"
	ScaleLogit       ScaleType = "logit"
	ScaleFunction    ScaleType = "function"
	ScaleFunctionLog ScaleType = "functionlog"
	ScaleAsinh       ScaleType = "asinh"
)

// Chart retains typed common chart properties and unknown fields in Extra.
type Chart struct {
	Type     ChartType                  `json:"type"`
	Title    string                     `json:"title"`
	Elements []json.RawMessage          `json:"elements"`
	XLabel   string                     `json:"x_label,omitempty"`
	YLabel   string                     `json:"y_label,omitempty"`
	XUnit    string                     `json:"x_unit,omitempty"`
	YUnit    string                     `json:"y_unit,omitempty"`
	XScale   ScaleType                  `json:"x_scale,omitempty"`
	YScale   ScaleType                  `json:"y_scale,omitempty"`
	Extra    map[string]json.RawMessage `json:"-"`
}

func (chart *Chart) UnmarshalJSON(data []byte) error {
	type wire Chart
	var value wire
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	raw := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	known := map[string]bool{"type": true, "title": true, "elements": true, "x_label": true, "y_label": true, "x_unit": true, "y_unit": true, "x_scale": true, "y_scale": true}
	for key := range known {
		delete(raw, key)
	}
	*chart = Chart(value)
	chart.Extra = raw
	switch chart.Type {
	case ChartLine, ChartScatter, ChartBar, ChartPie, ChartBoxAndWhisker, ChartSuper:
	default:
		chart.Type = ChartUnknown
	}
	return nil
}

func parseOutput(line []byte, execution *Execution, options *RunCodeOptions) error {
	var message map[string]json.RawMessage
	if err := json.Unmarshal(line, &message); err != nil {
		return fmt.Errorf("codeinterpreter: invalid output: %w", err)
	}
	var kind string
	_ = json.Unmarshal(message["type"], &kind)
	switch kind {
	case "stdout", "stderr":
		var text string
		_ = json.Unmarshal(message["text"], &text)
		output := OutputMessage{Line: text, Timestamp: time.Now().UnixNano(), Error: kind == "stderr"}
		if output.Error {
			execution.Logs.Stderr = append(execution.Logs.Stderr, text)
			if options.OnStderr != nil {
				options.OnStderr(output)
			}
		} else {
			execution.Logs.Stdout = append(execution.Logs.Stdout, text)
			if options.OnStdout != nil {
				options.OnStdout(output)
			}
		}
	case "error":
		var value ExecutionError
		_ = json.Unmarshal(line, &value)
		execution.Error = &value
		if options.OnError != nil {
			options.OnError(value)
		}
	case "number_of_executions":
		_ = json.Unmarshal(message["execution_count"], &execution.ExecutionCount)
	case "result":
		result := resultFromRaw(message)
		execution.Results = append(execution.Results, result)
		if options.OnResult != nil {
			options.OnResult(result)
		}
	}
	return nil
}

func resultFromRaw(message map[string]json.RawMessage) Result {
	result := Result{Raw: RawData{}, Extra: map[string]json.RawMessage{}}
	_ = json.Unmarshal(message["is_main_result"], &result.IsMainResult)
	knownStrings := map[string]*string{"text": &result.Text, "html": &result.HTML, "markdown": &result.Markdown, "svg": &result.SVG, "png": &result.PNG, "jpeg": &result.JPEG, "pdf": &result.PDF, "latex": &result.LaTeX, "json": &result.JSON, "javascript": &result.JavaScript}
	for key, raw := range message {
		if key == "type" || key == "is_main_result" {
			continue
		}
		result.Raw[key] = append(json.RawMessage(nil), raw...)
		if target := knownStrings[key]; target != nil {
			_ = json.Unmarshal(raw, target)
			continue
		}
		switch key {
		case "data":
			_ = json.Unmarshal(raw, &result.Data)
		case "chart":
			var chart Chart
			if json.Unmarshal(raw, &chart) == nil {
				result.Chart = &chart
			}
		default:
			result.Extra[key] = append(json.RawMessage(nil), raw...)
		}
	}
	return result
}
