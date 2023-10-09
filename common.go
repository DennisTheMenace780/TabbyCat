package main

import (
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-billy/v5/util"
	fixtures "github.com/go-git/go-git-fixtures/v4"
	. "gopkg.in/check.v1"
)

type BaseSuite struct {
	fixtures.Suite
	Repository *git.Repository

	cache map[string]*git.Repository
}

func (s *BaseSuite) SetUpSuite(c *C) {
	s.buildBasicRepository(c)

	s.cache = make(map[string]*git.Repository)
}

func (s *BaseSuite) TearDownSuite(c *C) {
	s.Suite.TearDownSuite(c)
}

func (s *BaseSuite) buildBasicRepository(_ *C) {
	f := fixtures.Basic().One()
	s.Repository = s.NewRepository(f)
}

// NewRepository returns a new repository using the .git folder, if the fixture
// is tagged as worktree the filesystem from fixture is used, otherwise a new
// memfs filesystem is used as worktree.
func (s *BaseSuite) NewRepository(f *fixtures.Fixture) *git.Repository {
	var worktree, dotgit billy.Filesystem
	if f.Is("worktree") {
		r, err := git.PlainOpen(f.Worktree().Root())
		if err != nil {
			panic(err)
		}

		return r
	}

	dotgit = f.DotGit()
	worktree = memfs.New()

	st := filesystem.NewStorage(dotgit, cache.NewObjectLRUDefault())

	r, err := git.Open(st, worktree)
	if err != nil {
		panic(err)
	}

	return r
}

// NewRepositoryWithEmptyWorktree returns a new repository using the .git folder
// from the fixture but without a empty memfs worktree, the index and the
// modules are deleted from the .git folder.
func (s *BaseSuite) NewRepositoryWithEmptyWorktree(f *fixtures.Fixture) *git.Repository {
	dotgit := f.DotGit()
	err := dotgit.Remove("index")
	if err != nil {
		panic(err)
	}

	err = util.RemoveAll(dotgit, "modules")
	if err != nil {
		panic(err)
	}

	worktree := memfs.New()

	st := filesystem.NewStorage(dotgit, cache.NewObjectLRUDefault())

	r, err := git.Open(st, worktree)
	if err != nil {
		panic(err)
	}

	return r

}

func (s *BaseSuite) GetLocalRepositoryURL(f *fixtures.Fixture) string {
	return f.DotGit().Root()
}

func (s *BaseSuite) TemporalDir() (path string, clean func()) {
	fs := osfs.New(os.TempDir())
	relPath, err := util.TempDir(fs, "", "")
	if err != nil {
		panic(err)
	}

	path = fs.Join(fs.Root(), relPath)
	clean = func() {
		_ = util.RemoveAll(fs, relPath)
	}

	return
}

func (s *BaseSuite) TemporalHomeDir() (path string, clean func()) {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	fs := osfs.New(home)
	relPath, err := util.TempDir(fs, "", "")
	if err != nil {
		panic(err)
	}

	path = fs.Join(fs.Root(), relPath)
	clean = func() {
		_ = util.RemoveAll(fs, relPath)
	}

	return
}

func (s *BaseSuite) TemporalFilesystem() (fs billy.Filesystem, clean func()) {
	fs = osfs.New(os.TempDir())
	path, err := util.TempDir(fs, "", "")
	if err != nil {
		panic(err)
	}

	fs, err = fs.Chroot(path)
	if err != nil {
		panic(err)
	}

	clean = func() {
		_ = util.RemoveAll(fs, path)
	}

	return
}

type SuiteCommon struct{}

var _ = Suite(&SuiteCommon{})

var countLinesTests = [...]struct {
	i string // the string we want to count lines from
	e int    // the expected number of lines in i
}{
	{"", 0},
	{"a", 1},
	{"a\n", 1},
	{"a\nb", 2},
	{"a\nb\n", 2},
	{"a\nb\nc", 3},
	{"a\nb\nc\n", 3},
	{"a\n\n\nb\n", 4},
	{"first line\n\tsecond line\nthird line\n", 3},
}

func AssertReferences(c *C, r *git.Repository, expected map[string]string) {
	for name, target := range expected {
		expected := plumbing.NewReferenceFromStrings(name, target)

		obtained, err := r.Reference(expected.Name(), true)
		c.Assert(err, IsNil)

		c.Assert(obtained, DeepEquals, expected)
	}
}

func AssertReferencesMissing(c *C, r *git.Repository, expected []string) {
	for _, name := range expected {
		_, err := r.Reference(plumbing.ReferenceName(name), false)
		c.Assert(err, NotNil)
		c.Assert(err, Equals, plumbing.ErrReferenceNotFound)
	}
}

func CommitNewFile(c *C, repo *git.Repository, fileName string) plumbing.Hash {
	wt, err := repo.Worktree()
	c.Assert(err, IsNil)

	fd, err := wt.Filesystem.Create(fileName)
	c.Assert(err, IsNil)

	_, err = fd.Write([]byte("# test file"))
	c.Assert(err, IsNil)

	err = fd.Close()
	c.Assert(err, IsNil)

	_, err = wt.Add(fileName)
	c.Assert(err, IsNil)

	sha, err := wt.Commit("test commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "test",
			Email: "test@example.com",
			When:  time.Now(),
		},
		Committer: &object.Signature{
			Name:  "test",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	c.Assert(err, IsNil)

	return sha
}
