package grafana_vertica_service

type QueryRequest struct {
	Timezone      string    `json:"timezone,omitempty"`
	Format        string    `json:"format,omitempty"`
	PanelId       int32     `json:"panel_id,omitempty"`
	DashboardId   int32     `json:"dashboard_id,omitempty"`
	Range         *Range    `json:"range,omitempty"`
	RangeRaw      *RangeRaw `json:"rangeRaw,omitempty"`
	Interval      string    `json:"interval,omitempty"`
	IntervalMs    int64     `json:"interval_ms,omitempty"`
	Targets       []*Target `json:"targets,omitempty"`
	MaxDataPoints int64     `json:"maxDataPoints,omitempty"`
}

type Range struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
	Raw  *Raw   `json:"raw,omitempty"`
}

type Raw struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

type RangeRaw struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}
type Target struct {
	Target string `json:"target,omitempty"`
	RefId  string `json:"refId,omitempty"`
	Type   string `json:"type,omitempty"`
}

type QuerySeriesResponse struct {
	Target     string       `json:"target"`
	Datapoints [][2]float64 `json:"datapoints"`
}

type Column struct {
	Text string `json:"text"`
	Type string `json:"type"`
	//Sort bool   `json:"sort"`
	//Desc bool   `json:"desc"`
}

func (c *Column) setType(format string) {
	c.Type = format
}

type QueryTableResponse struct {
	Columns []Column        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
	Type    string          `json:"type"`
}
