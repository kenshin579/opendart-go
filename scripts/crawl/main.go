package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	listURLFmt   = "https://opendart.fss.or.kr/guide/main.do?apiGrpCd=%s"
	detailURLFmt = "https://opendart.fss.or.kr/guide/detail.do?apiGrpCd=%s&apiId=%s"
	docsRoot     = "docs/api"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	client := newHTTPClient()
	var done []APIRef
	var failures []string

	for _, cat := range categories {
		listHTML, err := httpGet(client, fmt.Sprintf(listURLFmt, cat.Code))
		if err != nil {
			failures = append(failures, fmt.Sprintf("list %s: %v", cat.Code, err))
			continue
		}
		refs, err := parseList(listHTML)
		if err != nil {
			failures = append(failures, fmt.Sprintf("parse list %s: %v", cat.Code, err))
			continue
		}
		for i := range refs {
			refs[i].GrpCd = cat.Code
			refs[i].Category = cat.Name
		}
		time.Sleep(politeDelay)

		for _, ref := range refs {
			detailHTML, err := httpGet(client, fmt.Sprintf(detailURLFmt, ref.GrpCd, ref.APIID))
			if err != nil {
				failures = append(failures, fmt.Sprintf("%s/%s fetch: %v", ref.Category, ref.Name, err))
				continue
			}
			spec, err := parseDetail(detailHTML)
			if err != nil {
				failures = append(failures, fmt.Sprintf("%s/%s parse: %v", ref.Category, ref.Name, err))
				continue
			}
			if err := writeDoc(ref, spec); err != nil {
				failures = append(failures, fmt.Sprintf("%s/%s write: %v", ref.Category, ref.Name, err))
				continue
			}
			done = append(done, ref)
			fmt.Printf("✓ %s / %s\n", ref.Category, ref.Name)
			time.Sleep(politeDelay)
		}
	}

	if err := writeIndex(done); err != nil {
		return err
	}

	fmt.Printf("\n완료: %d개 API 문서 생성\n", len(done))
	if len(failures) > 0 {
		fmt.Printf("\n실패 %d건:\n", len(failures))
		for _, f := range failures {
			fmt.Println("  -", f)
		}
	}
	return nil
}

// writeDoc 은 한 API md 를 docs/api/{카테고리}/{이름}.md 로 쓴다.
// 경로는 renderIndex 의 링크와 동일한 docRelPath 규칙을 공유한다.
func writeDoc(ref APIRef, spec APISpec) error {
	full := filepath.Join(docsRoot, filepath.FromSlash(docRelPath(ref)))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return err
	}
	return os.WriteFile(full, []byte(renderMarkdown(ref, spec)), 0o644)
}

// writeIndex 는 docs/api/README.md 를 쓴다.
func writeIndex(refs []APIRef) error {
	if err := os.MkdirAll(docsRoot, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(docsRoot, "README.md"), []byte(renderIndex(refs)), 0o644)
}
