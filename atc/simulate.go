package atc

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

const (
	numEncounters = 100
	reset         = false
)

var encountersCounter = 0

var encounterTypes = []EncounterType{}

func (p *Platform) SimulateEncounters(ctx context.Context) error {
	res, err := p.db.Query(ctx, "SELECT id, severity FROM atc.encounter_types")
	if err != nil {
		return fmt.Errorf("select encounter types: %w", err)
	}
	defer res.Close()

	for res.Next() {
		var id string
		var severity uint8
		if err := res.Scan(&id, &severity); err != nil {
			return fmt.Errorf("scan encounter type: %w", err)
		}
		encounterTypes = append(encounterTypes, EncounterType{
			ID:       id,
			Severity: severity,
		})
	}

	if reset {
		// Delete all encounters in database to start fresh simulation.
		_, err = p.db.Exec(ctx, "DELETE FROM atc.encounters")
		if err != nil {
			return fmt.Errorf("delete encounters: %w", err)
		}
	}

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.simulatePendingEncounters(ctx)
	}()

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.simulateTriageEncounters(ctx)
	}()

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.simulateFinalizeEncounters(ctx)
	}()

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.simulateAlertEncounters(ctx)
	}()

	return nil
}

// Create new random encounters until numEncounters is reached.
func (p *Platform) simulatePendingEncounters(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.mu.RLock()
			if encountersCounter == numEncounters {
				p.mu.RUnlock()
				continue
			}
			p.mu.RUnlock()

			idx := rand.Intn(len(encounterTypes))
			encounterType := encounterTypes[idx]

			encounter := Encounter{
				ID:              uuid.New(),
				PatientID:       uuid.New(),
				EncounterTypeID: encounterType.ID,
				Severity:        encounterType.Severity,
			}

			q := `
			INSERT INTO atc.encounters (id, patient_id, encounter_type_id, severity)
			VALUES ($1, $2, $3, $4);`

			_, err := p.db.Exec(ctx, q,
				encounter.ID,
				encounter.PatientID,
				encounter.EncounterTypeID,
				encounter.Severity,
			)
			if err != nil {
				p.logger.Error("insert encounter",
					slog.Any("encounter", encounter),
					slog.String("error", err.Error()),
				)
			}

			p.mu.Lock()
			encountersCounter += 1
			p.mu.Unlock()
		}
	}
}

// Randomly select encounters in status 'waiting' and update them to 'triaged'.
func (p *Platform) simulateTriageEncounters(ctx context.Context) {
	ticker := time.NewTicker(4 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			q := `
			UPDATE atc.encounters SET
				status = 'triaged',
				status_changed_at = NOW(),
				alert_level = 0
			WHERE id = (
				SELECT id FROM atc.encounters
				WHERE status = 'waiting'
				ORDER BY random() LIMIT 1
			);`

			_, err := p.db.Exec(ctx, q)
			if err != nil {
				p.logger.Error("triage encounter",
					slog.String("error", err.Error()),
				)
			}
		}
	}
}

// Randomly select encounters in status 'triaged' and update them to 'completed'.
func (p *Platform) simulateFinalizeEncounters(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			q := `
			UPDATE atc.encounters SET
				status = 'completed',
				status_changed_at = NOW(),
				completed_at = NOW()
			WHERE id = (
				SELECT id FROM atc.encounters
				WHERE status = 'triaged'
				ORDER BY random() LIMIT 1
			);`

			_, err := p.db.Exec(ctx, q)
			if err != nil {
				p.logger.Error("finalize encounter",
					slog.String("error", err.Error()),
				)
			}
		}
	}
}

// Select encounters that need to be alerted and update their statuses.
func (p *Platform) simulateAlertEncounters(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			q := `
			UPDATE atc.encounters SET
				alert_level = alert_level + 1
			WHERE status = 'waiting'
			AND alert_level < 3
			AND (
				CASE
					WHEN severity = 1 THEN NOW() - status_changed_at > (alert_level + 1) * INTERVAL '1 minute'
					WHEN severity = 2 THEN NOW() - status_changed_at > (alert_level + 1) * INTERVAL '30 seconds'
					WHEN severity = 3 THEN NOW() - status_changed_at > (alert_level + 1) * INTERVAL '15 seconds'
				END
			)
			RETURNING id, alert_level;`

			rows, err := p.db.Query(ctx, q)
			if err != nil {
				p.logger.Error("alert encounter query failed",
					slog.String("error", err.Error()),
				)
				continue
			}

			for rows.Next() {
				var id string
				var alertLevel uint8
				if err := rows.Scan(&id, &alertLevel); err != nil {
					p.logger.Error("scan alert encounter",
						slog.String("error", err.Error()),
					)
					continue
				}

				p.logger.Info("encounter alerted",
					slog.String("id", id),
					slog.Int("alert_level", int(alertLevel)),
				)
			}
			rows.Close()
		}
	}
}
