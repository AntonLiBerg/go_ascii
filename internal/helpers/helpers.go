package helpers

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
