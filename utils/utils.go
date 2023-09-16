package utils

import (
	"noelzubin/redis-go/store"
	"strconv"
)

func GetScoreMemberPairs(s []string) ([]store.ScoreMember, error) {
	var pairs []store.ScoreMember
	for i := 0; i < len(s); i += 2 {
		score, err := strconv.Atoi(s[i])
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, store.NewScoreMember(int64(score), s[i+1]))
	}
	return pairs, nil
}
