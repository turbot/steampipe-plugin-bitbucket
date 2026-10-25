package bitbucket

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

// Response shapes below are trimmed from real api.bitbucket.org src responses.

const testHash = "7b8b171c0dbfa8af665c24a63d91f2b185b8701b"

type srcServer struct {
	*httptest.Server
	requests atomic.Int32
	lastURL  atomic.Value
}

func newSrcServer(t *testing.T, handler http.HandlerFunc) *srcServer {
	s := &srcServer{}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.requests.Add(1)
		s.lastURL.Store(r.URL.String())
		handler(w, r)
	}))
	t.Cleanup(s.Close)
	return s
}

func (s *srcServer) get(urlStr string) (*http.Response, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("%s", resp.Status)
	}
	return resp, nil
}

func collect(t *testing.T, s *srcServer, req srcRequest) ([]SrcEntry, error) {
	t.Helper()
	var rows []SrcEntry
	err := readSrc(s.get, s.URL+"/2.0", req, func(e SrcEntry) bool {
		rows = append(rows, e)
		return true
	})
	return rows, err
}

func listingEntry(path, typ string, size int) string {
	sizeField := ""
	if typ == "commit_file" {
		sizeField = fmt.Sprintf(`"size":%d,"mimetype":null,`, size)
	}
	return fmt.Sprintf(`{"path":%q,"type":%q,%s"attributes":[],"commit":{"hash":%q,"type":"commit"}}`, path, typ, sizeField, testHash)
}

func TestListingFollowsNextAndKeepsTheFilterAndDepth(t *testing.T) {
	var s *srcServer
	s = newSrcServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if r.URL.Query().Get("page") == "" {
			if got := r.URL.Query().Get("pagelen"); got != "100" {
				t.Errorf("pagelen = %q, want 100", got)
			}
			if got := r.URL.Query().Get("max_depth"); got != "20" {
				t.Errorf("max_depth = %q, want 20", got)
			}
			if got := r.URL.Query().Get("q"); got != `path~"pom.xml"` {
				t.Errorf("q = %q", got)
			}
			_, _ = fmt.Fprintf(w, `{"pagelen":100,"page":1,"values":[%s,%s],"next":"%s/2.0/repositories/acme/widget/src/HEAD/?page=E8Ux"}`,
				listingEntry("src", "commit_directory", 0), listingEntry("pom.xml", "commit_file", 21845), s.URL)
			return
		}
		_, _ = fmt.Fprintf(w, `{"pagelen":100,"page":2,"values":[%s]}`, listingEntry("src/pom.xml", "commit_file", 512))
	})

	rows, err := collect(t, s, srcRequest{Workspace: "acme", RepoSlug: "widget", MaxDepth: 20, Query: `path~"pom.xml"`})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 3 || s.requests.Load() != 2 {
		t.Fatalf("got %d rows from %d requests, want 3 from 2", len(rows), s.requests.Load())
	}
	if rows[0].Type != "commit_directory" || rows[0].Size != nil {
		t.Errorf("directory row = %+v", rows[0])
	}
	if rows[1].Path != "pom.xml" || *rows[1].Size != 21845 || rows[1].Commit.Hash != testHash || rows[1].Content != nil {
		t.Errorf("file row = %+v", rows[1])
	}
	if rows[2].Path != "src/pom.xml" {
		t.Errorf("second page row = %+v", rows[2])
	}
}

// Stopping at a LIMIT must not spend another request of the hourly budget.
func TestListingStopsFetchingOnceEmitDeclines(t *testing.T) {
	var s *srcServer
	s = newSrcServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"values":[%s,%s],"next":"%s/next"}`,
			listingEntry("a", "commit_file", 1), listingEntry("b", "commit_file", 1), s.URL)
	})

	n := 0
	err := readSrc(s.get, s.URL+"/2.0", srcRequest{Workspace: "acme", RepoSlug: "widget", PageLen: 1}, func(SrcEntry) bool {
		n++
		return false
	})
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 || s.requests.Load() != 1 {
		t.Errorf("emitted %d rows over %d requests, want 1 over 1", n, s.requests.Load())
	}
	if got := s.lastURL.Load().(string); !strings.Contains(got, "pagelen=1") {
		t.Errorf("request %q did not carry the limit as pagelen", got)
	}
}

func TestListingAScopePathTargetsThatDirectory(t *testing.T) {
	s := newSrcServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"values":[]}`)
	})
	if _, err := collect(t, s, srcRequest{Workspace: "acme", RepoSlug: "widget", Commit: testHash, ScopePath: "/services/api/"}); err != nil {
		t.Fatal(err)
	}
	want := "/2.0/repositories/acme/widget/src/" + testHash + "/services/api/?pagelen=100"
	if got := s.lastURL.Load().(string); got != want {
		t.Errorf("requested %q, want %q", got, want)
	}
}

// Bitbucket serves every file as text/plain, JSON files included, so a lockfile is read as content
// rather than parsed as a listing.
func TestFileReadReturnsContentForAJSONFile(t *testing.T) {
	const body = `{"name":"widget","values":[1],"lockfileVersion":3}`
	s := newSrcServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = fmt.Fprint(w, body)
	})

	rows, err := collect(t, s, srcRequest{Workspace: "acme", RepoSlug: "widget", Commit: testHash, Path: "package-lock.json"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("got %d rows, want 1", len(rows))
	}
	row := rows[0]
	if row.Path != "package-lock.json" || row.Type != "commit_file" || row.Content == nil || *row.Content != body {
		t.Errorf("row = %+v", row)
	}
	if *row.Size != int64(len(body)) || row.Commit.Hash != testHash {
		t.Errorf("size/commit = %v/%q", *row.Size, row.Commit.Hash)
	}
	if got := s.lastURL.Load().(string); got != "/2.0/repositories/acme/widget/src/"+testHash+"/package-lock.json" {
		t.Errorf("requested %q", got)
	}
}

func TestFileReadAtABranchLeavesTheHashUnknown(t *testing.T) {
	s := newSrcServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = fmt.Fprint(w, "x")
	})
	rows, err := collect(t, s, srcRequest{Workspace: "acme", RepoSlug: "widget", Commit: "main", Path: "go.mod"})
	if err != nil {
		t.Fatal(err)
	}
	if rows[0].Commit.Hash != "" {
		t.Errorf("commit hash = %q, want empty for a branch read", rows[0].Commit.Hash)
	}
}

func TestFileReadOfBinaryKeepsSizeButNoContent(t *testing.T) {
	s := newSrcServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write([]byte{0x89, 'P', 'N', 'G', 0x00, 0xff})
	})
	rows, err := collect(t, s, srcRequest{Workspace: "acme", RepoSlug: "widget", Path: "logo.png"})
	if err != nil {
		t.Fatal(err)
	}
	if rows[0].Content != nil || *rows[0].Size != 6 {
		t.Errorf("row = %+v", rows[0])
	}
}

func TestFileReadOfADirectoryReturnsTheDirectoryItself(t *testing.T) {
	s := newSrcServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = fmt.Fprintf(w, `{"values":[%s]}`, listingEntry("src/main", "commit_directory", 0))
	})
	rows, err := collect(t, s, srcRequest{Workspace: "acme", RepoSlug: "widget", Path: "src"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Path != "src" || rows[0].Type != "commit_directory" || rows[0].Content != nil {
		t.Errorf("rows = %+v", rows)
	}
	if rows[0].Commit.Hash != testHash {
		t.Errorf("commit hash = %q, want it taken from the listing", rows[0].Commit.Hash)
	}
}

func TestMissingFileIsNoRowButForbiddenIsAnError(t *testing.T) {
	status := http.StatusNotFound
	s := newSrcServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	})
	rows, err := collect(t, s, srcRequest{Workspace: "acme", RepoSlug: "widget", Path: "nope.txt"})
	if err != nil || len(rows) != 0 {
		t.Errorf("404: rows=%v err=%v, want none and nil", rows, err)
	}

	status = http.StatusForbidden
	if _, err := collect(t, s, srcRequest{Workspace: "acme", RepoSlug: "widget", Path: "nope.txt"}); err == nil {
		t.Error("403 returned no error, so a token without repository scope would read as an empty repository")
	}
	if _, err := collect(t, s, srcRequest{Workspace: "acme", RepoSlug: "widget"}); err == nil {
		t.Error("403 on a listing returned no error")
	}
}

func TestPathAndScopePathTogetherIsAnError(t *testing.T) {
	_, err := collect(t, newSrcServer(t, func(http.ResponseWriter, *http.Request) {}),
		srcRequest{Workspace: "acme", RepoSlug: "widget", Path: "a", ScopePath: "b"})
	if err == nil {
		t.Error("expected an error")
	}
}

func TestSrcURLEscapesSegmentsButKeepsSlashes(t *testing.T) {
	got := srcURL("https://api.bitbucket.org/2.0", "acme", "widget", "HEAD", "/docs/My File#1.md", false)
	want := "https://api.bitbucket.org/2.0/repositories/acme/widget/src/HEAD/docs/My%20File%231.md"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	if got := srcURL("https://api.bitbucket.org/2.0", "acme", "widget", "HEAD", "", true); got != "https://api.bitbucket.org/2.0/repositories/acme/widget/src/HEAD/" {
		t.Errorf("root listing url = %q", got)
	}
}
