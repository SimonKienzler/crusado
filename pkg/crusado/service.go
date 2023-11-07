package crusado

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
	"gopkg.in/yaml.v2"
)

var (
	errNoTemplateFoundForName = errors.New("no template found for name")
)

type Service struct {
	TemplatesDirectory string
	templates          []Template
}

type FileType string

const (
	MarkdownFileType = FileType("md")
	YAMLFileType     = FileType("yaml")
)

func (s *Service) GetAll() ([]Template, error) {
	if err := s.loadTemplatesFromDir(); err != nil {
		return nil, err
	}

	sort.Slice(s.templates, func(i, j int) bool {
		return s.templates[i].Name < s.templates[j].Name
	})

	return s.templates, nil
}

func (s *Service) GetByName(name string) (*Template, error) {
	if err := s.loadTemplatesFromDir(); err != nil {
		return nil, err
	}

	for i := range s.templates {
		if s.templates[i].Name == name {
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

		supported, fileType := hasSupportedFileType(e.Name())

		if !supported {
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

		var tpls []Template

		switch fileType {
		case MarkdownFileType:
			tpls, err = parseMarkdown(b)
		case YAMLFileType:
			tpls, err = parseYAML(b)
		default:
			log.Printf("Could not parse %s: No matching parser found", filePath)
			continue
		}

		if err != nil {
			log.Printf("Could not parse %s: %q", filePath, err)
			continue
		}

		s.templates = append(s.templates, tpls...)
	}

	return ValidateTemplateList(s.templates)
}

func parseMarkdown(content []byte) ([]Template, error) {
	md := goldmark.New(goldmark.WithExtensions(&frontmatter.Extender{}))

	var buf bytes.Buffer
	ctx := parser.NewContext()
	if err := md.Convert(content, &buf, parser.WithContext(ctx)); err != nil {
		return nil, err
	}

	meta := Meta{}
	metaRaw := frontmatter.Get(ctx)

	if metaRaw == nil {
		return nil, fmt.Errorf("markdown doesn't seem to have expected frontmatter")
	}

	if err := metaRaw.Decode(&meta); err != nil {
		return nil, err
	}

	return []Template{
		{
			Meta:        meta,
			Description: buf.String(),
		},
	}, nil
}

// parseYAML can return multiple templates, due to how the profile YAML was
// designed in the first iteration of crusado
func parseYAML(content []byte) ([]Template, error) {
	templateMapList := map[string][]Template{}
	err := yaml.Unmarshal(content, &templateMapList)
	if err != nil {
		return nil, err
	}

	return templateMapList["templates"], nil
}

func hasSupportedFileType(fileName string) (bool, FileType) {
	if isMarkdownFile(fileName) {
		return true, MarkdownFileType
	}

	if isYAMLFile(fileName) {
		return true, YAMLFileType
	}

	return false, ""
}

func isMarkdownFile(fileName string) bool {
	if strings.HasSuffix(fileName, ".md") || strings.HasSuffix(fileName, ".markdown") {
		return true
	}

	return false
}

func isYAMLFile(fileName string) bool {
	if strings.HasSuffix(fileName, ".yaml") || strings.HasSuffix(fileName, ".yml") {
		return true
	}

	return false
}
