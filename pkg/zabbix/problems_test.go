package zabbix_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/sgaunet/zabbix-cli/pkg/zabbix"
)

func TestGetProblemOptionEventIDs(t *testing.T) {
	req := &zabbix.GetProblemRequest{Params: zabbix.ProblemParams{}}
	ids := []string{"1", "2"}
	zabbix.GetProblemOptionEventIDs(ids)(req)
	if !reflect.DeepEqual(req.Params.EventIDs, ids) {
		t.Errorf("expected EventIDs %v, got %v", ids, req.Params.EventIDs)
	}
}

func TestGetProblemOptionGroupsIDs(t *testing.T) {
	req := &zabbix.GetProblemRequest{Params: zabbix.ProblemParams{}}
	ids := []string{"g1", "g2"}
	zabbix.GetProblemOptionGroupsIDs(ids)(req)
	if !reflect.DeepEqual(req.Params.GroupsIDs, ids) {
		t.Errorf("expected GroupsIDs %v, got %v", ids, req.Params.GroupsIDs)
	}
}

func TestGetProblemOptionHostsIDs(t *testing.T) {
	req := &zabbix.GetProblemRequest{Params: zabbix.ProblemParams{}}
	ids := []string{"h1", "h2"}
	zabbix.GetProblemOptionHostsIDs(ids)(req)
	if !reflect.DeepEqual(req.Params.HostsIDs, ids) {
		t.Errorf("expected HostsIDs %v, got %v", ids, req.Params.HostsIDs)
	}
}

func TestGetProblemOptionObjectIDs(t *testing.T) {
	req := &zabbix.GetProblemRequest{Params: zabbix.ProblemParams{}}
	ids := []string{"o1", "o2"}
	zabbix.GetProblemOptionObjectIDs(ids)(req)
	if !reflect.DeepEqual(req.Params.ObjectIDs, ids) {
		t.Errorf("expected ObjectIDs %v, got %v", ids, req.Params.ObjectIDs)
	}
}

func TestGetProblemOptionSource(t *testing.T) {
	req := &zabbix.GetProblemRequest{Params: zabbix.ProblemParams{}}
	source := 42
	zabbix.GetProblemOptionSource(source)(req)
	if req.Params.Source != source {
		t.Errorf("expected Source %v, got %v", source, req.Params.Source)
	}
}

func TestGetProblemOptionObject(t *testing.T) {
	req := &zabbix.GetProblemRequest{Params: zabbix.ProblemParams{}}
	obj := 7
	zabbix.GetProblemOptionObject(obj)(req)
	if req.Params.Object != obj {
		t.Errorf("expected Object %v, got %v", obj, req.Params.Object)
	}
}

func TestGetProblemOptionAcknowledged(t *testing.T) {
	req := &zabbix.GetProblemRequest{Params: zabbix.ProblemParams{}}
	zabbix.GetProblemOptionAcknowledged(true)(req)
	if req.Params.Acknowledged != true {
		t.Errorf("expected Acknowledged true, got %v", req.Params.Acknowledged)
	}
}

func TestGetProblemOptionSuppressed(t *testing.T) {
	req := &zabbix.GetProblemRequest{Params: zabbix.ProblemParams{}}
	zabbix.GetProblemOptionSuppressed(true)(req)
	if req.Params.Suppressed != true {
		t.Errorf("expected Suppressed true, got %v", req.Params.Suppressed)
	}
}

func TestGetProblemOptionSeverities(t *testing.T) {
	req := &zabbix.GetProblemRequest{Params: zabbix.ProblemParams{}}
	sevs := []string{"1", "2"}
	zabbix.GetProblemOptionSeverities(sevs)(req)
	if !reflect.DeepEqual(req.Params.Severities, sevs) {
		t.Errorf("expected Severities %v, got %v", sevs, req.Params.Severities)
	}
}

func TestGetProblemOptionEvalType(t *testing.T) {
	req := &zabbix.GetProblemRequest{Params: zabbix.ProblemParams{}}
	zabbix.GetProblemOptionEvalType(2)(req)
	if req.Params.EvalType != 2 {
		t.Errorf("expected EvalType 2, got %v", req.Params.EvalType)
	}
}

func TestGetProblemOptionTags(t *testing.T) {
	req := &zabbix.GetProblemRequest{Params: zabbix.ProblemParams{}}
	tags := []zabbix.FilterProblemTags{{Tag: "t", Value: "v", Operator: 1}}
	zabbix.GetProblemOptionTags(tags)(req)
	if !reflect.DeepEqual(req.Params.Tags, tags) {
		t.Errorf("expected Tags %v, got %v", tags, req.Params.Tags)
	}
}

func TestGetProblemOptionRecent(t *testing.T) {
	req := &zabbix.GetProblemRequest{Params: zabbix.ProblemParams{}}
	zabbix.GetProblemOptionRecent(true)(req)
	if req.Params.Recent != true {
		t.Errorf("expected Recent true, got %v", req.Params.Recent)
	}
}

func TestGetProblemOptionEventidFrom(t *testing.T) {
	req := &zabbix.GetProblemRequest{Params: zabbix.ProblemParams{}}
	zabbix.GetProblemOptionEventidFrom("123")(req)
	if req.Params.EventidFrom != "123" {
		t.Errorf("expected EventidFrom '123', got %v", req.Params.EventidFrom)
	}
}

func TestGetProblemOptionEventidTill(t *testing.T) {
	req := &zabbix.GetProblemRequest{Params: zabbix.ProblemParams{}}
	zabbix.GetProblemOptionEventidTill("456")(req)
	if req.Params.EventidTill != "456" {
		t.Errorf("expected EventidTill '456', got %v", req.Params.EventidTill)
	}
}

func TestGetProblemOptionTimeFrom(t *testing.T) {
	req := &zabbix.GetProblemRequest{Params: zabbix.ProblemParams{}}
	zabbix.GetProblemOptionTimeFrom(123456)(req)
	if req.Params.TimeFrom != 123456 {
		t.Errorf("expected TimeFrom 123456, got %v", req.Params.TimeFrom)
	}
}

func TestGetProblemOptionTimeTill(t *testing.T) {
	req := &zabbix.GetProblemRequest{Params: zabbix.ProblemParams{}}
	zabbix.GetProblemOptionTimeTill(654321)(req)
	if req.Params.TimeTill != 654321 {
		t.Errorf("expected TimeTill 654321, got %v", req.Params.TimeTill)
	}
}

func TestGetProblemOptionCountOutput(t *testing.T) {
	req := &zabbix.GetProblemRequest{Params: zabbix.ProblemParams{}}
	zabbix.GetProblemOptionCountOutput(true)(req)
	if req.Params.CountOutput != true {
		t.Errorf("expected CountOutput true, got %v", req.Params.CountOutput)
	}
}

func TestGetProblemOptionEditable(t *testing.T) {
	req := &zabbix.GetProblemRequest{Params: zabbix.ProblemParams{}}
	zabbix.GetProblemOptionEditable(true)(req)
	if req.Params.Editable != true {
		t.Errorf("expected Editable true, got %v", req.Params.Editable)
	}
}

func TestGetProblemOptionExcludeSearch(t *testing.T) {
	req := &zabbix.GetProblemRequest{Params: zabbix.ProblemParams{}}
	zabbix.GetProblemOptionExcludeSearch(true)(req)
	if req.Params.ExcludeSearch != true {
		t.Errorf("expected ExcludeSearch true, got %v", req.Params.ExcludeSearch)
	}
}

func TestGetProblemOptionLimit(t *testing.T) {
	req := &zabbix.GetProblemRequest{Params: zabbix.ProblemParams{}}
	zabbix.GetProblemOptionLimit(10)(req)
	if req.Params.Limit != 10 {
		t.Errorf("expected Limit 10, got %v", req.Params.Limit)
	}
}

func TestGetProblemOptionPreservekeys(t *testing.T) {
	req := &zabbix.GetProblemRequest{Params: zabbix.ProblemParams{}}
	zabbix.GetProblemOptionPreservekeys(true)(req)
	if req.Params.PreserveKeys != true { // Corrected field name
		t.Errorf("expected PreserveKeys true, got %v", req.Params.PreserveKeys) // Corrected field name
	}
}

func TestGetProblemOptionSearchByAny(t *testing.T) {
	req := &zabbix.GetProblemRequest{Params: zabbix.ProblemParams{}}
	zabbix.GetProblemOptionSearchByAny(true)(req)
	if req.Params.SearchByAny != true {
		t.Errorf("expected SearchByAny true, got %v", req.Params.SearchByAny)
	}
}

func TestGetProblemOptionSearchWildcardsEnabled(t *testing.T) {
	req := &zabbix.GetProblemRequest{Params: zabbix.ProblemParams{}}
	zabbix.GetProblemOptionSearchWildcardsEnabled(true)(req)
	if req.Params.SearchWildcardsEnabled != true {
		t.Errorf("expected SearchWildcardsEnabled true, got %v", req.Params.SearchWildcardsEnabled)
	}
}

// --- Problem struct methods ---
func TestProblem_GetClock(t *testing.T) {
	p := &zabbix.Problem{Clock: zabbix.StringInt64(time.Now().Unix())}
	tm := p.GetClock()
	if tm.IsZero() {
		t.Errorf("expected valid time, got zero")
	}
}

func TestProblem_GetClock_Zero(t *testing.T) {
	p := &zabbix.Problem{Clock: zabbix.StringInt64(0)}
	tm := p.GetClock()
	// 0 is Unix epoch (Jan 1, 1970), not invalid
	if tm.Unix() != 0 {
		t.Errorf("expected Unix epoch for clock=0, got %v", tm)
	}
}

func TestProblem_GetRClock(t *testing.T) {
	now := time.Now().Unix()
	p := &zabbix.Problem{Rclock: zabbix.StringInt64(now)}
	tm := p.GetRClock()
	if tm.IsZero() {
		t.Errorf("expected valid time, got zero")
	}
}

func TestProblem_GetRClock_Empty(t *testing.T) {
	p := &zabbix.Problem{Rclock: zabbix.StringInt64(0)}
	tm := p.GetRClock()
	if !tm.IsZero() {
		t.Errorf("expected zero time for empty rclock")
	}
}

func TestProblem_GetDuration(t *testing.T) {
	now := time.Now().Unix()
	p := &zabbix.Problem{Clock: zabbix.StringInt64(now - 10), Rclock: zabbix.StringInt64(now)}
	dur := p.GetDuration()
	if dur < 10*time.Second {
		t.Errorf("expected duration >= 10s, got %v", dur)
	}
}

func TestProblem_GetDurationStr(t *testing.T) {
	now := time.Now().Unix()
	p := &zabbix.Problem{Clock: zabbix.StringInt64(now - 3661), Rclock: zabbix.StringInt64(now)}
	str := p.GetDurationStr()
	if str == "" {
		t.Errorf("expected duration string, got empty")
	}
}

func TestProblem_GetAcknowledge(t *testing.T) {
	p := &zabbix.Problem{Acknowledged: "1"}
	if !p.GetAcknowledge() {
		t.Errorf("expected acknowledged true")
	}
	p.Acknowledged = "0"
	if p.GetAcknowledge() {
		t.Errorf("expected acknowledged false")
	}
}

func TestProblem_GetSuppressed(t *testing.T) {
	p := &zabbix.Problem{Suppressed: "1"}
	if !p.GetSuppressed() {
		t.Errorf("expected suppressed true")
	}
	p.Suppressed = "0"
	if p.GetSuppressed() {
		t.Errorf("expected suppressed false")
	}
}

func TestProblem_GetAcknowledgeStr(t *testing.T) {
	p := &zabbix.Problem{Acknowledged: "1"}
	if p.GetAcknowledgeStr() != "Yes" {
		t.Errorf("expected Yes for acknowledged")
	}
	p.Acknowledged = "0"
	if p.GetAcknowledgeStr() != "No" {
		t.Errorf("expected No for not acknowledged")
	}
}

func TestProblem_GetSuppressedStr(t *testing.T) {
	p := &zabbix.Problem{Suppressed: "1"}
	if p.GetSuppressedStr() != "Yes" {
		t.Errorf("expected Yes for suppressed")
	}
	p.Suppressed = "0"
	if p.GetSuppressedStr() != "No" {
		t.Errorf("expected No for not suppressed")
	}
}

func TestProblem_GetSeverity(t *testing.T) {
	p := &zabbix.Problem{Severity: "3"}
	if p.GetSeverity() != "Average" {
		t.Errorf("expected Average, got %v", p.GetSeverity())
	}
	p.Severity = "invalid"
	if p.GetSeverity() != "Unknown" {
		t.Errorf("expected Unknown for invalid severity")
	}
}
