package file

func SplitIntoChunks(b []byte, chunkSize int) [][]byte {
	var chunks [][]byte

	for i := 0; i < len(b); i += chunkSize {
		end := i + chunkSize

		if end > len(b) {
			end = len(b)
		}

		chunks = append(chunks, b[i:end])
	}

	return chunks
}
