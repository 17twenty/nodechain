package nodechain

// SplitText splits a long text into chunks with overlap.
// If you chunk by token count later, we can plug in a tokenizer.
func SplitText(text string, chunkSize, overlap int) []string {
	runes := []rune(text)
	n := len(runes)

	if chunkSize <= 0 || overlap < 0 || overlap >= chunkSize {
		panic("invalid chunk or overlap size")
	}

	var chunks []string
	start := 0

	for start < n {
		end := start + chunkSize
		if end > n {
			end = n
		}
		chunks = append(chunks, string(runes[start:end]))

		if end == n {
			break
		}

		start = end - overlap
	}

	return chunks
}
