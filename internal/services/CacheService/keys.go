package CacheService

import "context"

func (s *CacheService) GetAllKeys() ([]string, error) {

	var cursor uint64
	var n int
	var keys []string

	for {
		var err error
		var ks []string
		ks, cursor, err = s.redi.Scan(context.Background(), cursor, "", 1000).Result()
		if err != nil {
			return keys, err
		}
		keys = append(keys, ks...)
		n += len(ks)

		if cursor == 0 {
			break
		}
	}

	return keys, nil
}
