package spotify

func StringSliceToIDs(ids []string) []ID {
	var result []ID
	for _, id := range ids {
		result = append(result, ID(id))
	}

	return result
}
