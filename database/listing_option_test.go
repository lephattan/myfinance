package database

import (
	"log"
	"regexp"
	"testing"
)

func TestListOptionWhereGroup(t *testing.T) {
	opt := ListOptions{}
	var query string
	var args []interface{}
	var want_query *regexp.Regexp
	tablename := "holdings"
	opt.SetTableName(tablename)

	where_group_1 := WhereGroup{Operator: "and"}
	where_group_1.Where("portfolio_id", 1, "and")
	where_group_1.Where("symbol", "nvl", "and")
	opt.WhereGroup(&where_group_1)
	query, args = opt.BuildQuery()
	log.Printf("Query: %s", query)
	log.Printf("Args: %+v", args)

	want_query = regexp.MustCompile(`(?mi)^select +\* +from +` + tablename + ` +where +\( +portfolio_id += +\? +and +symbol += \? ?\) ?;$`)
	if !want_query.MatchString(query) {
		t.Fatalf(`BuildQuery = %q, want match for %#q`, query, want_query)
	}
	if args[0] != 1 {
		t.Fatalf("args[0] is %v, want %v", args[0], 1)
	}
	if args[1] != "nvl" {
		t.Fatalf("args[0] is %v, want %v", args[1], "nvl")
	}

	where_group_or := WhereGroup{Operator: "or"}
	where_group_or.Where("volume", 100, "and")
	opt.WhereGroup(&where_group_or)
	query, args = opt.BuildQuery()
	log.Printf("Query: %s", query)
	log.Printf("Args: %+v", args)

	want_query = regexp.MustCompile(`(?mi)^select +\* +from +` + tablename + ` +where +\( +portfolio_id += +\? +and +symbol += \? ?\) +or +volume += +\? ?;$`)
	if !want_query.MatchString(query) {
		t.Fatalf(`BuildQuery = %q, want match for %#q`, query, want_query)
	}
	if args[2] != 100 {
		t.Fatalf("args[2] is %v, want %v", args[2], 100)
	}

}
