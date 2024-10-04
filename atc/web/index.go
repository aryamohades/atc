package web

import (
	"atc/atc"
	"atc/atc/webapi"
	"log/slog"
	"net/http"
	"strings"
)

var indexTemplate = webapi.ParseTemplate("index")

var indexMeta = &webapi.Meta{
	Title:       "ATC",
	Description: "Some description of ATC.",
	Keywords:    "atc, some, keywords",
}

type EncountersTableColumn struct {
	Field    string
	Label    string
	Link     string
	Sortable bool
	SortDir  string
}

var encountersTableColumns = []EncountersTableColumn{
	{Field: "patient_id", Label: "Patient ID", Sortable: false},
	{Field: "type", Label: "Type", Sortable: true},
	{Field: "status", Label: "Status", Sortable: true},
	{Field: "severity", Label: "Severity", Sortable: true},
	{Field: "alert", Label: "Alert Level", Sortable: true},
	{Field: "time", Label: "Status Time", Sortable: true},
	{Field: "created", Label: "Start Time", Sortable: true},
}

func (h *handler) handleIndexPage() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{
			"Meta":    indexMeta,
			"Mainjs":  h.mainjs,
			"Maincss": h.maincss,
		}

		ctx := r.Context()
		encounters, err := h.plt.GetEncounters(ctx, atc.QueryEncountersParams{
			Sort: r.URL.Query().Get("sort"),
		})
		if err != nil {
			h.logger.Error("get encounters", slog.String("error", err.Error()))
			webapi.RenderTemplate(w, errorTemplate, "error.tmpl", data)
			return
		}

		data["Rows"] = encounters
		data["Columns"] = getEncountersTableColumns(r)

		if r.Header.Get("X-Sort") != "" {
			data["Sort"] = true
		}
		if webapi.IsHtmxRequest(r) && !webapi.IsHtmxBoosted(r) {
			data["Partial"] = true
			webapi.RenderTemplate(w, indexTemplate, "content", data)
			return
		}
		webapi.RenderTemplate(w, indexTemplate, "page.tmpl", data)
	})
}

func getEncountersTableColumns(r *http.Request) []EncountersTableColumn {
	sortBy, sortDir := getSortFromRequest(r)
	u, _ := r.URL.Parse(r.URL.String())

	columns := make([]EncountersTableColumn, 0, len(encountersTableColumns))
	for _, col := range encountersTableColumns {
		if col.Sortable {
			q := u.Query()
			q.Set("sort", col.Field+",desc")

			if col.Field == sortBy {
				col.SortDir = sortDir
				if sortDir == "desc" {
					q.Set("sort", col.Field+",asc")
				}
			}
			u.RawQuery = q.Encode()
			col.Link = u.Path + "?" + u.RawQuery
		}
		columns = append(columns, col)
	}
	return columns
}

func getSortFromRequest(r *http.Request) (string, string) {
	sort := r.URL.Query().Get("sort")
	sortParts := strings.Split(sort, ",")
	if sort == "" || len(sortParts) != 2 {
		return "", ""
	}
	return sortParts[0], sortParts[1]
}
