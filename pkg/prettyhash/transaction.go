package prettyhash

func ShortenTransactionHash(hash string) string {
	if len(hash) < 12 {
		return ""
	}
	return hash[0:12] + "..."
}
