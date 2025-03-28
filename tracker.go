package change_tracker

import (
	"fmt"
	"reflect"
)

// ChangeTracker - хелпер для отслеживания изменений полей структуры.
type ChangeTracker struct {
	originalValues map[string]any
	trackedFields  map[string]bool
}

func NewChangeTracker(source any, fieldsToTrack []string) (*ChangeTracker, error) {
	ct := &ChangeTracker{}
	err := ct.init(source, fieldsToTrack)
	if err != nil {
		return nil, err
	}
	return ct, nil
}

// init инициализирует трекер, сохраняя исходные значения указанных полей.
func (ct *ChangeTracker) init(source any, fieldsToTrack []string) error {
	ct.originalValues = make(map[string]any)
	ct.trackedFields = make(map[string]bool)

	// Запоминаем, какие поля нужно отслеживать
	for _, field := range fieldsToTrack {
		ct.trackedFields[field] = true
	}

	// Сохраняем исходные значения
	sourceValue := reflect.ValueOf(source)
	if sourceValue.Kind() == reflect.Ptr {
		sourceValue = sourceValue.Elem()
	}

	if sourceValue.Kind() != reflect.Struct {
		return fmt.Errorf("source must be a struct or pointer to struct")
	}

	sourceType := sourceValue.Type()

	for i := 0; i < sourceValue.NumField(); i++ {
		fieldName := sourceType.Field(i).Name
		if _, ok := ct.trackedFields[fieldName]; ok {
			ct.originalValues[fieldName] = sourceValue.Field(i).Interface()
		}
	}

	return nil
}

// Changes возвращает map изменённых полей и их старых значений.
func (ct *ChangeTracker) Changes(current any) (map[string]any, error) {
	if ct.originalValues == nil {
		return nil, fmt.Errorf("tracker not initialized, call Init first")
	}

	currentValue := reflect.ValueOf(current)
	if currentValue.Kind() == reflect.Ptr {
		currentValue = currentValue.Elem()
	}

	if currentValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("current must be a struct or pointer to struct")
	}

	changes := make(map[string]any)
	currentType := currentValue.Type()

	for i := 0; i < currentValue.NumField(); i++ {
		fieldName := currentType.Field(i).Name
		if _, ok := ct.trackedFields[fieldName]; !ok {
			continue // Пропускаем поля, которые не нужно отслеживать
		}

		currentFieldValue := currentValue.Field(i).Interface()
		originalValue, exists := ct.originalValues[fieldName]

		if !exists || !reflect.DeepEqual(originalValue, currentFieldValue) {
			changes[fieldName] = originalValue
		}
	}

	return changes, nil
}
