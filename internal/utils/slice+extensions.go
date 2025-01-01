package utils

func Top[T any](items []T, top int) []T {
	if top > len(items) {
		top = len(items)
	}

	return items[:top]
}
