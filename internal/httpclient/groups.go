package httpclient

import (
	"context"
	"encoding/json"
)

// Group 은 증권신고서(DS006) 응답의 그룹 하나(title + raw list).
type Group struct {
	Title string          `json:"title"` // 그룹명칭(예: 일반사항)
	List  json.RawMessage `json:"list"`  // 그룹 항목 배열(타입은 호출측에서 결정)
}

type groupEnvelope struct {
	Envelope
	Group []Group `json:"group"`
}

// GetGroups 는 그룹형(DS006) 응답을 디코딩해 그룹 목록을 반환한다.
func GetGroups(ctx context.Context, c *Client, path string, params map[string]string) ([]Group, error) {
	var env groupEnvelope
	if err := c.GetJSON(ctx, path, params, &env); err != nil {
		return nil, err
	}
	return env.Group, nil
}
