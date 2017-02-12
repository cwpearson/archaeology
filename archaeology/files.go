package archaeology

import "os"

func readBytes(path string, n int) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := make([]byte, n)

	_, err = f.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
