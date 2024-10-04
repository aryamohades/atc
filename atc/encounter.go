package atc

import (
	"atc/atc/sqlhelper"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	defaultEncountersLimit = 100
	maxEncountersLimit     = 500
	maxEncountersOffset    = 1000
)

var defaultEncountersSort = fmt.Sprintf(`
	CASE
		WHEN status='waiting' THEN 1
		WHEN status='triaged' THEN 2
		WHEN status='completed' THEN 3
		ELSE 5
	END,
	e.alert_level %s,
	et.severity %s,
	e.started_at %s`,
	sqlhelper.SortDirDesc,
	sqlhelper.SortDirDesc,
	sqlhelper.SortDirAsc,
)

var encountersSortMap = map[string]string{
	"status": "e.status",
	"time":   "NOW() - e.status_changed_at",
}

type EncounterType struct {
	ID          string
	Description string
	Severity    uint8
}

type Encounter struct {
	ID              uuid.UUID
	PatientID       uuid.UUID
	EncounterTypeID string
	AlertLevel      int
	Severity        uint8
	Status          string
	StatusChangedAt time.Time
	StartedAt       time.Time
	CompletedAt     time.Time
	EncounterType   EncounterType
}

func (e *Encounter) TimeInStatus() time.Duration {
	return time.Since(e.StatusChangedAt).Truncate(time.Second)
}

func getEncountersSort(sort string) string {
	if sort == "" {
		return defaultEncountersSort
	}
	sortParts := strings.Split(sort, ",")
	if len(sortParts) != 2 {
		return defaultEncountersSort
	}

	sortDir, ok := sqlhelper.GetSortDir(sortParts[1])
	if !ok {
		return defaultEncountersSort
	}

	sortBy := sortParts[0]
	if sortBy == "type" {
		return fmt.Sprintf("et.id %s, e.started_at %s", sortDir, sqlhelper.SortDirAsc)
	}
	if sortBy == "severity" {
		return fmt.Sprintf("et.severity %s, e.alert_level %s", sortDir, sqlhelper.SortDirDesc)
	}
	if sortBy == "alert" {
		return fmt.Sprintf("e.alert_level %s, et.severity %s, e.started_at %s", sortDir, sqlhelper.SortDirDesc, sqlhelper.SortDirAsc)
	}
	if sortBy == "status" {
		return fmt.Sprintf("e.status %s, e.alert_level %s", sortDir, sqlhelper.SortDirDesc)
	}
	if sortBy == "time" {
		return fmt.Sprintf("NOW() - e.status_changed_at %s", sortDir)
	}
	if sortBy == "created" {
		return fmt.Sprintf("e.started_at %s", sortDir)
	}
	return defaultEncountersSort
}

func getEncountersLimit(limit int) int {
	if limit <= 0 || limit > maxEncountersLimit {
		return defaultEncountersLimit
	}
	return limit
}

func getEncountersOffset(offset int) int {
	if offset < 0 || offset > maxEncountersOffset {
		return 0
	}
	return offset
}

type QueryEncountersParams struct {
	Status     *string
	Severity   *int
	AlertLevel *int
	Sort       string
	Limit      int
	Offset     int
}

func (p *Platform) GetEncounters(ctx context.Context, params QueryEncountersParams) ([]Encounter, error) {
	params.Sort = getEncountersSort(params.Sort)
	params.Limit = getEncountersLimit(params.Limit)
	params.Offset = getEncountersOffset(params.Offset)

	q := fmt.Sprintf(`
	SELECT
		e.id,
		e.patient_id,
		e.encounter_type_id,
		e.alert_level,
		et.severity,
		et.description,
		e.status,
		e.status_changed_at,
		e.started_at
	FROM atc.encounters e
	JOIN atc.encounter_types et ON e.encounter_type_id = et.id
	WHERE ($1::INT IS NULL OR et.severity = $1)
	AND ($2::INT IS NULL OR e.alert_level = $2)
	AND ($3::TEXT IS NULL OR e.status = $3)
	ORDER BY %s
	LIMIT %d OFFSET %d;`,
		params.Sort,
		params.Limit,
		params.Offset,
	)

	rows, err := p.db.Query(ctx, q,
		params.Severity,
		params.AlertLevel,
		params.Status,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Encounter
	for rows.Next() {
		var e Encounter
		if err := rows.Scan(
			&e.ID,
			&e.PatientID,
			&e.EncounterTypeID,
			&e.AlertLevel,
			&e.EncounterType.Severity,
			&e.EncounterType.Description,
			&e.Status,
			&e.StatusChangedAt,
			&e.StartedAt,
		); err != nil {
			return nil, err
		}
		results = append(results, e)
	}
	return results, nil
}

func (p *Platform) ResetEncounters(ctx context.Context) error {
	q := `DELETE FROM atc.encounters;`
	_, err := p.db.Exec(ctx, q)

	p.mu.Lock()
	encountersCounter = 0
	p.mu.Unlock()
	return err
}
