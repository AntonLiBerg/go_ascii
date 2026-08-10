package helpers

import "slices"

func Filter[T any](values []T, keep func(T) bool) []T {
	result := make([]T, 0, len(values))
	for _, value := range values {
		if keep(value) {
			result = append(result, value)
		}
	}
	return result
}

func Any[T any](values []T, predicate func(T) bool) bool {
	for _, value := range values {
		if predicate(value) {
			return true
		}
	}
	return false
}

func All[T any](values []T, predicate func(T) bool) bool {
	for _, value := range values {
		if !predicate(value) {
			return false
		}
	}
	return true
}

func Transform[T any](values []T, transform func(T) T) []T {
	result := make([]T, 0, len(values))
	for _, value := range values {
		result = append(result, transform(value))
	}
	return result
}

func IsAllS1InS2(slice1 []string, slice2 []string) bool {
	return All(slice1, func(value string) bool {
		return slices.Contains(slice2, value)
	})
}

func IsUnique(values []string) bool {
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if _, exists := seen[value]; exists {
			return false
		}
		seen[value] = struct{}{}
	}
	return true
}
