package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Query struct {
	TickerName   string
	PeriodName   string
	PeriodNumber int
	Periods      int
	StartDate    time.Time
	EndDate      time.Time
}

func (q *Query) Limit() bson.M {
	return bson.M{"$limit": q.Periods}
}

func (q *Query) MatchGroupSort() (bson.M, bson.M, bson.M) {
	group := bson.M{}
	sort := bson.M{}
	match := bson.M{}

	var startDate time.Time
	endDate := time.Now()
	now := time.Now()

	switch q.PeriodName {
	case "minute":
		startDate = now.Add(-1 * (time.Duration(q.Periods*q.PeriodNumber) * time.Minute))

		group = bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"year":      bson.M{"$year": "$date"},
					"dayOfYear": bson.M{"$dayOfYear": "$date"},
					"hour":      bson.M{"$hour": "$date"},
					"interval": bson.M{
						"$subtract": []bson.M{
							{"$minute": "$date"},
							{"$mod": []interface{}{
								bson.M{"$minute": "$date"},
								q.PeriodNumber,
							}},
						},
					},
				},
				"high":   bson.M{"$max": "$high"},
				"low":    bson.M{"$min": "$low"},
				"open":   bson.M{"$first": "$open"},
				"close":  bson.M{"$last": "$close"},
				"volume": bson.M{"$sum": "$volume"},
			},
		}

		sort = bson.M{"$sort": bson.M{
			"_id.interval":  1,
			"_id.hour":      1,
			"_id.dayOfYear": 1,
			"_id.year":      1,
		}}
	case "hour":
		startDate = now.Add(-1 * (time.Duration(q.Periods*q.PeriodNumber) * time.Hour))

		group = bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"year":      bson.M{"$year": "$date"},
					"dayOfYear": bson.M{"$dayOfYear": "$date"},
					"interval": bson.M{
						"$subtract": []bson.M{
							{"$hour": "$date"},
							{"$mod": []interface{}{
								bson.M{"$hour": "$date"},
								q.PeriodNumber,
							}},
						},
					},
				},
				"high":   bson.M{"$max": "$high"},
				"low":    bson.M{"$min": "$low"},
				"open":   bson.M{"$first": "$open"},
				"close":  bson.M{"$last": "$close"},
				"volume": bson.M{"$sum": "$volume"},
			},
		}

		sort = bson.M{"$sort": bson.M{
			"_id.interval":  1,
			"_id.dayOfYear": 1,
			"_id.year":      1,
		}}
	case "day":
		startDate = now.Add(-1 * (time.Duration(q.Periods*q.PeriodNumber*24) * time.Hour))

		group = bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"year":      bson.M{"$year": "$date"},
					"dayOfYear": bson.M{"$dayOfYear": "$date"},
				},
				"high":   bson.M{"$max": "$high"},
				"low":    bson.M{"$min": "$low"},
				"open":   bson.M{"$first": "$open"},
				"close":  bson.M{"$last": "$close"},
				"volume": bson.M{"$sum": "$volume"},
			},
		}

		sort = bson.M{"$sort": bson.M{
			"_id.dayOfYear": 1,
			"_id.year":      1,
		}}
	case "week":
		startDate = now.Add(-1 * (time.Duration(q.Periods*q.PeriodNumber*24*7) * time.Hour))

		group = bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"year": bson.M{"$year": "$date"},
					"week": bson.M{"$week": "$date"},
				},
				"high":   bson.M{"$max": "$high"},
				"low":    bson.M{"$min": "$low"},
				"open":   bson.M{"$first": "$open"},
				"close":  bson.M{"$last": "$close"},
				"volume": bson.M{"$sum": "$volume"},
			},
		}

		sort = bson.M{"$sort": bson.M{
			"_id.week": 1,
			"_id.year": 1,
		}}
	case "month":
		startDate = now.Add(-1 * (time.Duration(q.Periods*q.PeriodNumber*24*31) * time.Hour))

		group = bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"year":  bson.M{"$year": "$date"},
					"month": bson.M{"$month": "$date"},
				},
				"high":   bson.M{"$max": "$high"},
				"low":    bson.M{"$min": "$low"},
				"open":   bson.M{"$first": "$open"},
				"close":  bson.M{"$last": "$close"},
				"volume": bson.M{"$sum": "$volume"},
			},
		}

		sort = bson.M{"$sort": bson.M{
			"_id.month": 1,
			"_id.year":  1,
		}}
	}

	if false == q.StartDate.IsZero() {
		startDate = q.StartDate
	}

	if false == q.EndDate.IsZero() {
		endDate = q.EndDate
	}

	match = bson.M{"$match": bson.M{
		"ticker": q.TickerName,
		"date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}}

	return match, group, sort
}

type TickAggregator struct {
	db *mgo.Database
}

func (t *TickAggregator) setDB(db *mgo.Database) {
	t.db = db
}
func (t *TickAggregator) query(q Query) {
	match, group, sort := q.MatchGroupSort()

	c := t.db.C("ticks")
	pipe := c.Pipe([]bson.M{
		match,
		{"$sort": bson.M{"date": 1}},
		group,
		sort,
		{"$project": bson.M{
			"_id":      false,
			"interval": "$_id.interval",
			"open":     "$open",
			"close":    "$close",
			"high":     "$high",
			"low":      "$low",
			"ticker":   "$ticker",
			"volume":   "$volume",
		}},
	})

	results := []bson.M{}
	err := pipe.All(&results)
	if err != nil {
		fmt.Println("TODO: fault tolerance needed; ", err)
	}
	fmt.Println(results)
}
