package helper

func Diff[T comparable](listA, listB []T) (removed, added []T) {
	mapA := make(map[T]bool)
	mapB := make(map[T]bool)

	for _, item := range listA {
		mapA[item] = true
	}

	for _, item := range listB {
		mapB[item] = true
	}

	for item := range mapA {
		if !mapB[item] {
			removed = append(removed, item)
		}
	}

	for item := range mapB {
		if !mapA[item] {
			added = append(added, item)
		}
	}

	return removed, added
}
