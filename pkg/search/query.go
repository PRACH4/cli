package search

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/cli/cli/v2/pkg/text"
)

const (
	KindRepositories = "repositories"
	KindIssues       = "issues"
)

type Query struct {
	Keywords   []string
	Kind       string
	Limit      int
	Order      string
	Page       int
	Qualifiers Qualifiers
	Sort       string
}

type Qualifiers struct {
	Archived         *bool
	Assignee         string
	Author           string
	Base             string
	Closed           string
	Commenter        string
	Comments         string
	Created          string
	Draft            *bool
	Followers        string
	Fork             string
	Forks            string
	GoodFirstIssues  string
	Head             string
	HelpWantedIssues string
	In               []string
	Interactions     string
	Involves         string
	Is               []string
	Label            []string
	Language         string
	License          []string
	Mentions         string
	Merged           string
	Milestone        string
	No               []string
	Org              string
	Project          string
	Pushed           string
	Reactions        string
	Repo             []string
	Review           string
	ReviewRequested  string
	ReviewedBy       string
	Size             string
	Stars            string
	State            string
	Status           string
	Team             string
	Topic            []string
	Topics           string
	Type             string
	Updated          string
}

func (q Query) String() string {
	qualifiers := formatQualifiers(q.Qualifiers)
	keywords := formatKeywords(q.Keywords)
	all := append(keywords, qualifiers...)
	return strings.Join(all, " ")
}

func (q Qualifiers) Map() map[string][]string {
	m := map[string][]string{}
	v := reflect.ValueOf(q)
	t := reflect.TypeOf(q)
	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name
		key := text.CamelToKebab(fieldName)
		typ := v.FieldByName(fieldName).Kind()
		value := v.FieldByName(fieldName)
		switch typ {
		case reflect.Ptr:
			if value.IsNil() {
				continue
			}
			v := reflect.Indirect(value)
			m[key] = []string{fmt.Sprintf("%v", v)}
		case reflect.Slice:
			if value.IsNil() {
				continue
			}
			s := []string{}
			for i := 0; i < value.Len(); i++ {
				s = append(s, fmt.Sprintf("%v", value.Index(i)))
			}
			m[key] = s
		default:
			if value.IsZero() {
				continue
			}
			m[key] = []string{fmt.Sprintf("%v", value)}
		}
	}
	return m
}

func quote(s string) string {
	if strings.ContainsAny(s, " \"\t\r\n") {
		return fmt.Sprintf("%q", s)
	}
	return s
}

func formatQualifiers(qs Qualifiers) []string {
	var all []string
	for k, vs := range qs.Map() {
		for _, v := range vs {
			all = append(all, fmt.Sprintf("%s:%s", k, quote(v)))
		}
	}
	sort.Strings(all)
	return all
}

func formatKeywords(ks []string) []string {
	for i, k := range ks {
		ks[i] = quote(k)
	}
	return ks
}
