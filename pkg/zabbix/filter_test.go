package zabbix_test

import (
	"reflect"
	"testing"
	zabbix "github.com/sgaunet/zabbix-cli/pkg/zabbix"
)

func TestNewFilterGetMethod_Empty(t *testing.T) {
	f := zabbix.NewFilterGetMethod()
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
	if len(f.GetFilter()) != 0 {
		t.Errorf("expected empty filter map, got %v", f.GetFilter())
	}
}

func TestFilter_SingleKeyValue(t *testing.T) {
	f := zabbix.NewFilterGetMethod(zabbix.Filter("hostid", "12345"))
	got := f.GetFilter()
	want := map[string]interface{}{"hostid": "12345"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestFilter_OverwriteKey(t *testing.T) {
	f := zabbix.NewFilterGetMethod(zabbix.Filter("hostid", "12345"), zabbix.Filter("hostid", "67890"))
	got := f.GetFilter()
	want := map[string]interface{}{"hostid": "67890"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected overwrite to %v, got %v", want, got)
	}
}

func TestFilterByName(t *testing.T) {
	f := zabbix.NewFilterGetMethod(zabbix.FilterByName("myhost"))
	got := f.GetFilter()
	want := map[string]interface{}{"name": "myhost"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestFilterByNames(t *testing.T) {
	names := []string{"host1", "host2"}
	f := zabbix.NewFilterGetMethod(zabbix.FilterByNames(names))
	got := f.GetFilter()
	want := map[string]interface{}{"name": names}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestMultipleFilters(t *testing.T) {
	f := zabbix.NewFilterGetMethod(
		zabbix.Filter("hostid", "123"),
		zabbix.Filter("name", "abc"),
	)
	got := f.GetFilter()
	want := map[string]interface{}{"hostid": "123", "name": "abc"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}
