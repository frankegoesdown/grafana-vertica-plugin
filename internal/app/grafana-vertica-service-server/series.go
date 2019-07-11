package grafana_vertica_service

import (
	"context"
	"database/sql"
	"log"
	"regexp"

	_ "github.com/vertica/vertica-sql-go"
)

func SeriesResponse(connString string, qr QueryRequest) ([]QuerySeriesResponse, error) {
	connDB, err := sql.Open("vertica", connString)
	defer connDB.Close()
	ctx := context.Background()
	if err != nil {
		return nil, err

	}
	err = connDB.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	vals := make(map[string]*interface{})

	var resp []QuerySeriesResponse

	// iterate on targets (query blocks in grafana)
	for _, target := range qr.Targets {
		//if target.Type == "timeserie" {
		var valuesToScan []interface{}
		reGroup := regexp.MustCompile(`(?i)group by`)
		query := reGroup.ReplaceAllString(target.Target, "GROUP BY")

		// execute sql query
		rows, err := connDB.QueryContext(ctx, query)
		defer rows.Close()
		groupBys := myExp.FindGroupBy(query)

		if err != nil {
			log.Info("ERROR: ", err)
			return nil, err
		}

		cols, err := rows.Columns()
		if err != nil {
			log.Info("ERROR: ", err)
			return nil, err
		}

		// get all columns from query
		for _, column := range cols {
			vals[column] = new(interface{})
			valuesToScan = append(valuesToScan, vals[column])

			if len(groupBys) == 0 {
				if column == "time" {
					continue
				}
				resp = append(resp, QuerySeriesResponse{Target: column})

			}
		}

		var mapGroupBys = make(map[string]bool)
		for rows.Next() {
			err = rows.Scan(valuesToScan...)
			if err != nil {
				log.Info("ERROR: ", err)
				return nil, err
			}
			if len(groupBys) > 0 {
				for _, column := range cols {

					if stringInSlice(column, groupBys) && column != "time" {

						s, err := getString(*vals[column])
						if err != nil {
							log.Info("ERROR: ", err)
							return nil, err
						}
						if _, ok := mapGroupBys[s]; !ok {
							mapGroupBys[s] = true
							resp = append(resp, QuerySeriesResponse{Target: s})
						}
					} else {
						continue
					}
				}
			}

			for _, column := range cols {
				if len(groupBys) > 0 {
					if column == "time" {
						continue
					}
					for _, c := range cols {
						if !stringInSlice(c, groupBys) {
							val, err := getFloat(*vals[c], false)
							if err != nil {
								log.Info("ERROR: ", err)
								return nil, err
							}
							t, err := getFloat(*vals["time"], true)
							if err != nil {
								log.Info("ERROR: ", err)
								return nil, err
							}
							dp := [2]float64{val, t}
							for i := range resp {
								for k := range mapGroupBys {
									if resp[i].Target == k {
										resp[i].Datapoints = append(resp[i].Datapoints, dp)
									}
								}

							}
						}
					}
				} else {
					if column == "time" {
						continue
					}
					val, err := getFloat(*vals[column], false)
					if err != nil {
						log.Info("ERROR: ", err)
						return nil, err
					}
					t, err := getFloat(*vals["time"], true)
					if err != nil {
						log.Info("ERROR: ", err)
						return nil, err
					}

					dp := [2]float64{val, t}

					for i := range resp {
						if resp[i].Target == column {
							resp[i].Datapoints = append(resp[i].Datapoints, dp)
						}
					}
				}

			}
		}
	}

	return resp, nil
}
