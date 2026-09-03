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

// String returns the output line.
func (message OutputMessage) String() string { return message.Line }

// ExecutionError is a kernel error and traceback.
type ExecutionError struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Traceback string `json:"traceback"`
}

// Error formats the kernel error name and value.
func (e ExecutionError) Error() string { return fmt.Sprintf("%s: %s", e.Name, e.Value) }

// Logs contains collected standard output and standard error lines.
type Logs struct {
	Stdout []string `json:"stdout"`
	Stderr []string `json:"stderr"`
}

// Execution contains the complete result of one code execution.
type Execution struct {
	Results        []Result        `json:"results"`
	Logs           Logs            `json:"logs"`
	Error          *ExecutionError `json:"error,omitempty"`
	ExecutionCount int             `json:"execution_count,omitempty"`
}

// Text returns the text representation of the main result, if present.
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

// Result contains the decoded and raw MIME representations of one result.
type Result struct {
	Text, HTML, Markdown, SVG, PNG, JPEG, PDF, LaTeX, JSON, JavaScript string
	Data                                                               map[string]any
	Chart                                                              *Chart
	Extra                                                              map[string]json.RawMessage
	Raw                                                                RawData
	IsMainResult                                                       bool
}

// Formats returns the available raw MIME keys in sorted order.
func (result Result) Formats() []string {
	formats := make([]string, 0, len(result.Raw))
	for key := range result.Raw {
		formats = append(formats, key)
	}
	slices.Sort(formats)
	return formats
}

// ChartType identifies a supported chart representation.
type ChartType string

const (
	// ChartLine identifies a line chart.
	ChartLine ChartType = "line"
	// ChartScatter identifies a scatter chart.
	ChartScatter ChartType = "scatter"
	// ChartBar identifies a bar chart.
	ChartBar ChartType = "bar"
	// ChartPie identifies a pie chart.
	ChartPie ChartType = "pie"
	// ChartBoxAndWhisker identifies a box-and-whisker chart.
	ChartBoxAndWhisker ChartType = "box_and_whisker"
	// ChartSuper identifies a composite chart.
	ChartSuper ChartType = "superchart"
	// ChartUnknown identifies a chart type unknown to this SDK version.
	ChartUnknown ChartType = "unknown"
)

// ScaleType identifies a chart axis scale.
type ScaleType string

const (
	// ScaleLinear identifies a linear scale.
	ScaleLinear ScaleType = "linear"
	// ScaleDatetime identifies a date and time scale.
	ScaleDatetime ScaleType = "datetime"
	// ScaleCategorical identifies a categorical scale.
	ScaleCategorical ScaleType = "categorical"
	// ScaleLog identifies a logarithmic scale.
	ScaleLog ScaleType = "log"
	// ScaleSymlog identifies a symmetric logarithmic scale.
	ScaleSymlog ScaleType = "symlog"
	// ScaleLogit identifies a logit scale.
	ScaleLogit ScaleType = "logit"
	// ScaleFunction identifies a custom function scale.
	ScaleFunction ScaleType = "function"
	// ScaleFunctionLog identifies a logarithmic custom function scale.
	ScaleFunctionLog ScaleType = "functionlog"
	// ScaleAsinh identifies an inverse hyperbolic sine scale.
	ScaleAsinh ScaleType = "asinh"
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

// UnmarshalJSON decodes known chart fields and preserves unknown fields in Extra.
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
