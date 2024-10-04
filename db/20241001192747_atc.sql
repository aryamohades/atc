-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

CREATE SCHEMA atc;

CREATE TABLE atc.encounter_types (
  id text PRIMARY KEY,
  description text NOT NULL,
  severity smallint NOT NULL CHECK (severity IN (1, 2, 3))
);

INSERT INTO atc.encounter_types (id, description, severity) VALUES
  ('consultation', 'Consultation', 1),
  ('examination', 'Examination', 1),
  ('lab', 'Lab', 1),
  ('ultrasound', 'Ultrasound', 1),
  ('ct', 'CT', 2),
  ('mri', 'MRI', 2),
  ('xray', 'X-Ray', 2),
  ('surgery', 'Surgery', 3),
  ('allergic_reaction', 'Allergic Reaction', 3),
  ('cardiac_arrest', 'Cardiac Arrest', 3),
  ('severe_trauma', 'Trauma', 3),
  ('bleeding', 'Bleeding', 3),
  ('mental_health_crisis', 'Mental Health Crisis', 3);

CREATE TABLE atc.encounters (
  id uuid PRIMARY KEY,
  patient_id text NOT NULL,
  encounter_type_id text NOT NULL REFERENCES atc.encounter_types (id) ON DELETE CASCADE,
  status text NOT NULL DEFAULT 'waiting' CHECK (status IN ('waiting', 'triaged', 'completed')),
  status_changed_at timestamptz NOT NULL DEFAULT NOW(),
  alert_level smallint NOT NULL DEFAULT 0 CHECK (alert_level >= 0),
  severity smallint NOT NULL CHECK (severity IN (1, 2, 3)),
  started_at timestamptz NOT NULL DEFAULT NOW(),
  triaged_at timestamptz,
  completed_at timestamptz
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE atc.encounters CASCADE;
DROP TABLE atc.encounter_types CASCADE;

DROP SCHEMA atc CASCADE;

DROP EXTENSION pg_stat_statements;
-- +goose StatementEnd
