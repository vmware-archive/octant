package localcontent

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/heptio/developer-dash/internal/hcli"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/pkg/errors"
)

type LocalContent struct {
	root string
}

func New(root string) *LocalContent {
	return &LocalContent{
		root: root,
	}
}

func (l *LocalContent) Root() string {
	return l.root
}

func (l *LocalContent) Name() string {
	return "local"
}

func (l *LocalContent) Content(ctx context.Context, contentPath string, prefix string, namespace string) (component.ContentResponse, error) {
	if contentPath == "/" || contentPath == "" {
		return l.list()
	}

	fileName := fmt.Sprintf("%s.json", contentPath)
	return l.content(fileName)
}

func (l *LocalContent) list() (component.ContentResponse, error) {
	var out component.ContentResponse

	cols := component.NewTableCols("Title", "File")
	table := component.NewTable("Local Components", cols)

	err := l.walk(func(name, base string, content component.ContentResponse) error {
		title, err := l.titleToText(content.Title)
		if err != nil {
			return errors.Wrap(err, "convert title to text")
		}

		table.Add(component.TableRow{
			"Title": component.NewLink("", title, path.Join("/content/local", base)),
			"File":  component.NewText(name),
		})

		return nil
	})

	if err != nil {
		return out, nil
	}

	out.Title = component.Title(component.NewText("Local Contents"))

	out.ViewComponents = []component.ViewComponent{
		table,
	}

	return out, nil
}

func (l *LocalContent) ContentPath() string {
	return fmt.Sprintf("/%s", l.Name())
}

func (l *LocalContent) content(name string) (component.ContentResponse, error) {
	var out component.ContentResponse

	f, err := os.Open(filepath.Join(l.root, name))
	if err != nil {
		return out, errors.Wrap(err, "open local content")
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(&out); err != nil {
		return out, errors.Wrap(err, "read JSON")
	}

	return out, nil
}

type walkFn func(name, base string, content component.ContentResponse) error

func (l *LocalContent) walk(fn walkFn) error {
	if fn == nil {
		return errors.New("walkFn is nil")
	}

	fis, err := ioutil.ReadDir(l.root)
	if err != nil {
		return errors.Wrapf(err, "read %s", l.root)
	}

	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}

		ext := filepath.Ext(fi.Name())
		if ext != ".json" {
			continue
		}

		content, err := l.content(fi.Name())
		if err != nil {
			return err
		}

		base := strings.TrimSuffix(fi.Name(), ext)
		if err = fn(fi.Name(), base, content); err != nil {
			return err
		}
	}

	return nil
}

func (l *LocalContent) Navigation(namespace, root string) (*hcli.Navigation, error) {
	if !strings.HasSuffix(root, "/") {
		root = fmt.Sprintf("%s/", root)
	}
	nav := &hcli.Navigation{
		Title:    "Local Content",
		Path:     root,
		Children: []*hcli.Navigation{},
	}

	err := l.walk(func(name, base string, content component.ContentResponse) error {
		title, err := l.titleToText(content.Title)
		if err != nil {
			return errors.Wrap(err, "convert title to text")
		}

		nav.Children = append(nav.Children, &hcli.Navigation{
			Title: title,
			Path:  path.Join(root, base),
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return nav, nil
}

func (l *LocalContent) titleToText(title []component.TitleViewComponent) (string, error) {
	if len(title) < 1 {
		return "", errors.New("title is empty")
	}

	var parts []string
	for _, titlePart := range title {
		text, ok := titlePart.(*component.Text)
		if !ok {
			return "", errors.New("title has a not text component")
		}
		parts = append(parts, text.Config.Text)
	}

	return strings.Join(parts, " > "), nil
}

func (l *LocalContent) SetNamespace(namespace string) error {
	return nil
}

func (l *LocalContent) Start() error {
	return nil
}

func (l *LocalContent) Stop() {
}

func (l *LocalContent) Handlers() map[string]http.Handler {
	return make(map[string]http.Handler)
}
