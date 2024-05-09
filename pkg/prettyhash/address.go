package prettyhash

func ShortenHash(hash string) string {
	if len(hash) < 12 {
		return ""
	}
	return hash[0:6] + "..." + hash[len(hash)-8:]
}
