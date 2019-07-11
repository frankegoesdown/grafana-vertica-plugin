package grafana_vertica_service

import (
	"context"
	"database/sql"
	"log"
	"regexp"

	_ "github.com/vertica/vertica-sql-go"
)

func TableResponse(connString string, qr QueryRequest) ([]QueryTableResponse, error) {
	connDB, err := sql.Open("vertica", connString)
	defer connDB.Close()
	ctx := context.Background()
	if err != nil {
		log.Println("ERROR: ", err)
		return nil, err

	}
	err = connDB.PingContext(ctx)
	if err != nil {
		log.Println("ERROR: ", err)
		return nil, err
	}

	vals := make(map[string]*interface{})
	var resp []QueryTableResponse

	// iterate on targets (query blocks in grafana)
	for _, target := range qr.Targets {
		var valuesToScan []interface{}
		reGroup := regexp.MustCompile(`(?i)group by`)
		query := reGroup.ReplaceAllString(target.Target, "GROUP BY")

		rows, err := connDB.QueryContext(ctx, query)
		defer rows.Close()
		groupBys := myExp.FindGroupBy(query)

		if err != nil {
			return nil, err
		}

		cols, err := rows.Columns()
		if err != nil {
			log.Println("ERROR: ", err)
			return nil, err
		}
		columns := []Column{}

		// get all columns from query
		for _, column := range cols {
			vals[column] = new(interface{})
			valuesToScan = append(valuesToScan, vals[column])
			if column == "time" {
				timeColumn := Column{Text: "Time", Type: "time"} //, Sort: true, Desc: true}
				columns = append(columns, timeColumn)
				continue
			}
			if len(groupBys) == 0 {
				col := Column{Text: column} //, Type: "number"}
				columns = append(columns, col)
			}
		}

		outputRows := [][]interface{}{}

		qRT := QueryTableResponse{Columns: columns, Type: "table"}

		itr := 0
		var mapGroupBys = make(map[string]bool)
		for rows.Next() {
			row := []interface{}{}
			err = rows.Scan(valuesToScan...)
			if err != nil {
				log.Println("ERROR: ", err)
				return nil, err
			}

			if len(groupBys) > 0 {
				for _, column := range cols {

					if stringInSlice(column, groupBys) && column != "time" {

						s, _ := getString(*vals[column])
						if _, ok := mapGroupBys[s]; !ok {
							mapGroupBys[s] = true
							col := Column{Text: s} //, Type: "number"}
							columns = append(columns, col)
						}
					} else {
						continue
					}
				}
				qRT = QueryTableResponse{Columns: columns, Type: "table"}
			}

			for _, column := range cols {
				if itr == 0 {
					tableType := getTableType(*vals[column])
					for r := range qRT.Columns {
						if qRT.Columns[r].Text != "Time" {
							qRT.Columns[r].setType(tableType)
						}
					}

				}

				if len(groupBys) > 0 {
					if column == "time" {
						r := []interface{}{}
						t, err := getFloat(*vals["time"], true)
						if err != nil {
							log.Println("ERROR: ", err)
							return nil, err
						}
						r = append(r, t)
						row = append(r, row...)
						continue
					}
					row = append(row, *vals[column])

				} else {
					if column == "time" {
						r := []interface{}{}
						t, err := getFloat(*vals["time"], true)
						if err != nil {
							log.Println("ERROR: ", err)
							return nil, err
						}
						r = append(r, t)
						row = append(r, row...)
						continue
					}
					row = append(row, *vals[column])

				}
			}
			outputRows = append(outputRows, row)
			if itr < 2 {
				itr++
			}
		}

		qRT.Rows = outputRows
		resp = append(resp, qRT)
	}

	return resp, nil
}
