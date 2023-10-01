package crusado

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

var (
	errNoTemplateFoundForName = errors.New("no template found for name")
)

type Service struct {
	TemplatesDirectory string

	templates []Template
}

func (s *Service) GetAll() ([]Template, error) {
	if err := s.loadTemplatesFromDir(); err != nil {
		return nil, err
	}

	return s.templates, nil
}

func (s *Service) GetByName(name string) (*Template, error) {
	if err := s.loadTemplatesFromDir(); err != nil {
		return nil, err
	}

	for i := range s.templates {
		if s.templates[i].Meta.Name == name {
			return &s.templates[i], nil
		}
	}

	return nil, errNoTemplateFoundForName
}

func (s *Service) loadTemplatesFromDir() error {
	if s.templates != nil && len(s.templates) > 0 {
		return nil
	}

	s.templates = []Template{}

	entries, err := os.ReadDir(s.TemplatesDirectory)
	if err != nil {
		return err
	}

	for _, e := range entries {
		// TODO support recursively reading templates from subdirectories
		if e.IsDir() {
			continue
		}

		filePath := path.Join(s.TemplatesDirectory, e.Name())
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("Could not open %s: %q", filePath, err)
			continue
		}

		defer func() error {
			if err = file.Close(); err != nil {
				return err
			}
			return nil
		}()

		b, err := io.ReadAll(file)
		if err != nil {
			log.Printf("Could not read %s: %q", filePath, err)
			continue
		}

		tmpl, err := parse(string(b))
		if err != nil {
			log.Printf("Could not parse %s: %q", filePath, err)
			continue
		}

		s.templates = append(s.templates, *tmpl)
	}

	return nil
}

func parse(markdown string) (*Template, error) {
	md := goldmark.New(goldmark.WithExtensions(&frontmatter.Extender{}))

	var buf bytes.Buffer
	ctx := parser.NewContext()
	if err := md.Convert([]byte(markdown), &buf, parser.WithContext(ctx)); err != nil {
		return nil, err
	}

	meta := &Meta{}
	metaRaw := frontmatter.Get(ctx)

	if metaRaw == nil {
		return nil, fmt.Errorf("markdown doesn't seem to have expected frontmatter")
	}

	if err := metaRaw.Decode(meta); err != nil {
		return nil, err
	}

	return &Template{
		Meta:        *meta,
		Description: buf.String(),
	}, nil
}
