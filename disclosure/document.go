package disclosure

import "context"

// DownloadDocument 은 접수번호(rcept_no)로 공시서류 원본 ZIP 을 그대로 반환한다.
// 압축 해제·파싱은 호출자 몫 (임의 공시 원본은 형태가 다양한 바이너리).
func (c *Client) DownloadDocument(ctx context.Context, rceptNo string) ([]byte, error) {
	return c.http.GetBytes(ctx, "/api/document.xml", map[string]string{"rcept_no": rceptNo})
}
