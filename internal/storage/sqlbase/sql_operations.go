package sqlbase

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/90amper/metmon/internal/logger"
	"github.com/90amper/metmon/internal/models"
)

func (sb *SQLBase) createMetric(name string, mtype bool) error {
	var err error = nil
	sqlQuery := loadSnippet("snippets/insert_metric.sql")
	res, err := sb.db.Exec(sqlQuery, name, mtype)
	if err != nil {
		return fmt.Errorf("create metric failed: %w", err)
	}
	aff, _ := res.RowsAffected()
	logger.Trace("Rows affected: %v", aff)
	return nil
}

type sqlmetric struct {
	time   time.Time
	name   string
	mtype  bool
	cvalue int
	gvalue float64
}

const (
	sgauge   = true
	scounter = false
)

func (sb *SQLBase) CleanGauges() error {
	return fmt.Errorf("not implemented yet")
}

func (sb *SQLBase) ResetCounters() error {
	return fmt.Errorf("not implemented yet")
}

func (sb *SQLBase) GetGauges() (models.GaugeStore, error) {
	return models.GaugeStore{}, fmt.Errorf("not implemented yet")
}

func (sb *SQLBase) SaveToFile() error {
	return fmt.Errorf("not implemented yet")
}

func (sb *SQLBase) LoadFromFile() error {
	return fmt.Errorf("not implemented yet")
}

func (sb *SQLBase) Dumper() error {
	return fmt.Errorf("not implemented yet")
}

func (sb *SQLBase) BatchAdd(ms []models.Metric) error {
	var errs []error
	for _, metric := range ms {
		var err error
		switch metric.MType {
		case "gauge":
			err = sb.AddGauge(metric.ID, models.Gauge(*metric.Value))
			if err != nil {
				errs = append(errs, err)
			}
		case "counter":
			err = sb.AddCounter(metric.ID, models.Counter(*metric.Delta))
			if err != nil {
				errs = append(errs, err)
			}
		default:
			err := fmt.Errorf("unsupported metric type")
			errs = append(errs, err)
			logger.Error(err)
			continue
		}
	}
	return errors.Join(errs...)
}

func (sb *SQLBase) AddGauge(name string, value models.Gauge) error {
	var err error = nil

	err = sb.createMetric(name, sgauge)
	if err != nil {
		return err
	}

	sqlQuery := loadSnippet("snippets/insert_gauge.sql")
	res, err := sb.db.Exec(sqlQuery, name, value)
	// sql.NamedArg("name", name),
	// sql.Named("type", true),
	// sql.Named("value", value),
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	logger.Trace("Rows affected: %v", aff)
	return nil
}

func (sb *SQLBase) AddCounter(name string, value models.Counter) error {
	var err error = nil

	err = sb.createMetric(name, scounter)
	if err != nil {
		return err
	}

	sqlQuery := loadSnippet("snippets/insert_counter.sql")
	res, err := sb.db.Exec(sqlQuery, name, value)

	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	logger.Trace("Rows affected: %v", aff)
	return nil
}

func (sb *SQLBase) TickCounter(name string) error {
	var err error = nil

	sqlQuery := loadSnippet("snippets/tick_counter.sql")
	res, err := sb.db.Exec(sqlQuery, name)

	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	logger.Trace("Rows affected: %v", aff)
	return nil
}
func (sb *SQLBase) GetCurrentGauge(mName string) (models.Gauge, error) {
	var err error = nil
	sqlQuery := loadSnippet("snippets/select_metric.sql")
	row := sb.db.QueryRow(sqlQuery, mName, sgauge)
	var m sqlmetric
	err = row.Scan(&m.time, &m.name, &m.mtype, &m.cvalue, &m.gvalue)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("metric not found: %v", mName)
		}
		return 0, err
	}
	gauge := models.Gauge(m.gvalue)
	return gauge, nil
}

func (sb *SQLBase) GetCounter(mName string) (models.Counter, error) {
	var err error = nil
	sqlQuery := loadSnippet("snippets/select_metric.sql")
	row := sb.db.QueryRow(sqlQuery, mName, scounter)
	var m sqlmetric
	err = row.Scan(&m.time, &m.name, &m.mtype, &m.cvalue, &m.gvalue)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("metric not found: %v", mName)
		}
		return 0, err
	}
	counter := models.Counter(m.cvalue)
	return counter, nil
}

func (sb *SQLBase) GetCounters() (models.CounterStore, error) {
	var err error = nil
	cs := make(models.CounterStore)
	sqlQuery := loadSnippet("snippets/select_all_metrics.sql")
	rows, err := sb.db.Query(sqlQuery, scounter)
	if err != nil {
		return cs, err
	}
	defer rows.Close()
	for rows.Next() {
		var m sqlmetric
		err = rows.Scan(&m.time, &m.mtype, &m.time, &m.cvalue, &m.gvalue)
		if err != nil {
			return cs, err
		}
		cs[m.name] = models.Counter(m.cvalue)
	}
	rows.Close()
	err = rows.Err()
	if err != nil {
		return cs, err
	}
	return cs, nil
}

func (sb *SQLBase) GetCurrentGauges() (models.GaugeList, error) {
	var err error = nil
	gl := make(models.GaugeList)
	sqlQuery := loadSnippet("snippets/select_all_metrics.sql")
	rows, err := sb.db.Query(sqlQuery, sgauge)
	if err != nil {
		return gl, err
	}

	defer rows.Close()
	for rows.Next() {
		var m sqlmetric
		err = rows.Scan(&m.time, &m.mtype, &m.time, &m.cvalue, &m.gvalue)
		if err != nil {
			return gl, err
		}
		gl[m.name] = models.Gauge(m.gvalue)
	}
	rows.Close()
	err = rows.Err()
	if err != nil {
		return gl, err
	}
	return gl, nil
}
