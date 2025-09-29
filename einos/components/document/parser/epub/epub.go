// Package epub provides a document parser for EPUB files.
package epub

import (
	"archive/zip"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/document/parser"
	"github.com/cloudwego/eino/schema"
	"github.com/jaytaylor/html2text"
	"github.com/mszlu521/go-epub"
)

// Parser is a document parser for EPUB files.
type Parser struct {
	// StripHTML determines whether to strip HTML tags from the content
	StripHTML bool
}

// Option is a function that configures the EPUB parser.
type Option func(*Parser)

// WithStripHTML sets whether to strip HTML tags from the content.
func WithStripHTML(stripHTML bool) Option {
	return func(p *Parser) {
		p.StripHTML = stripHTML
	}
}

// NewParser creates a new EPUB parser.
func NewParser(ctx context.Context, config *Config, opts ...Option) (parser.Parser, error) {
	p := &Parser{}
	
	// Apply options
	for _, opt := range opts {
		opt(p)
	}
	
	// If config is provided, apply its settings
	if config != nil {
		p.StripHTML = config.StripHTML
	}
	
	return p, nil
}

// Config holds configuration for the EPUB parser.
type Config struct {
	// StripHTML determines whether to strip HTML tags from the content
	StripHTML bool
}

// Parse parses an EPUB file from a reader and returns a slice of documents.
// This implementation uses the go-epub library to parse the EPUB file.
func (p *Parser) Parse(ctx context.Context, reader io.Reader, opts ...parser.Option) ([]*schema.Document, error) {
	// Handle common options
	commonOpts := parser.GetCommonOptions(&parser.Options{}, opts...)
	// Use the go-epub library to open and parse the EPUB file
	epubBook, err := epub.NewReader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to open EPUB file: %w", err)
	}
	defer epubBook.Close()

	title := epubBook.GetTitle()
	author := epubBook.GetAuthor()
	// Get chapters from the EPUB file
	chapters, err := epubBook.GetChapters()
	if err != nil {
		return nil, fmt.Errorf("failed to get chapters: %w", err)
	}

	// Convert chapters to documents
	documents := make([]*schema.Document, 0, len(chapters))
	
	// Process each chapter
	for i, chapter := range chapters {
		// Skip empty chapters
		if strings.TrimSpace(chapter.Content) == "" {
			continue
		}

		// Convert HTML to text if needed
		contentStr := chapter.Content
		stripHTML := p.StripHTML
		if commonOpts.ExtraMeta != nil {
			if strip, ok := commonOpts.ExtraMeta["strip_html"].(bool); ok {
				stripHTML = strip
			}
		}
		
		if stripHTML {
			textContent, err := html2text.FromString(contentStr, html2text.Options{PrettyTables: true})
			if err == nil {
				contentStr = textContent
			}
		}

		// Create a document for this chapter
		doc := &schema.Document{
			Content: strings.TrimSpace(contentStr),
			MetaData: map[string]any{
				"book_title": title,
				"author": author,
				"order":      i,
				"chapter":      chapter.Title,
				"source_uri": commonOpts.URI,
				"epub_order": chapter.Order,
			},
		}

		// Merge extra metadata if provided
		if commonOpts.ExtraMeta != nil {
			for k, v := range commonOpts.ExtraMeta {
				doc.MetaData[k] = v
			}
		}

		documents = append(documents, doc)
	}

	return documents, nil
}

// ParseFromPath parses an EPUB file from a file path.
func (p *Parser) ParseFromPath(ctx context.Context, path string, opts ...parser.Option) ([]*schema.Document, error) {
	// Handle common options
	commonOpts := parser.GetCommonOptions(&parser.Options{}, opts...)
	
	// Open the EPUB file as a zip archive

	// Open the EPUB file as a zip archive
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open EPUB file: %w", err)
	}
	defer r.Close()

	// Find the container.xml file to get the OPF file path
	containerPath := "META-INF/container.xml"
	containerFile, err := r.Open(containerPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open container.xml: %w", err)
	}
	defer containerFile.Close()

	// Parse container.xml to find the OPF file
	opfPath, err := parseContainerXML(containerFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse container.xml: %w", err)
	}

	// Open and parse the OPF file
	opfFile, err := r.Open(opfPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open OPF file: %w", err)
	}
	defer opfFile.Close()

	// Parse the OPF file to get the manifest and spine
	manifest, spine, err := parseOPF(opfFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OPF file: %w", err)
	}

	// Parse each item in the spine
	var documents []*schema.Document
	baseDir := filepath.Dir(opfPath)
	// Keep track of document order
	docOrder := 0

	for _, itemRef := range spine.ItemRefs {
		// Find the item in the manifest
		item, exists := manifest[itemRef.IDRef]
		if !exists {
			continue
		}

		// Check if the item is an XHTML document
		if !strings.Contains(item.MediaType, "application/xhtml") &&
			!strings.Contains(item.MediaType, "text/html") {
			continue
		}

		// Construct the full path to the file
		fullPath := filepath.Join(baseDir, item.Href)
		if filepath.IsAbs(item.Href) {
			fullPath = item.Href[1:] // Remove leading slash
		}

		// Open the file
		contentFile, err := r.Open(fullPath)

		if err != nil {
			continue
		}

		// Read the content
		content, err := io.ReadAll(contentFile)
		if err != nil {
			contentFile.Close()
			continue
		}
		contentFile.Close()

		// Convert HTML to text if needed
		contentStr := string(content)
		stripHTML := p.StripHTML
		if commonOpts.ExtraMeta != nil {
			if strip, ok := commonOpts.ExtraMeta["strip_html"].(bool); ok {
				stripHTML = strip
			}
		}
		
		if stripHTML {
			contentStr, err = html2text.FromString(contentStr, html2text.Options{PrettyTables: true})
			if err != nil {
				// If HTML to text conversion fails, keep the original HTML
				contentStr = string(content)
			}
		}

		// Create a document for this section
		doc := &schema.Document{
			Content: strings.TrimSpace(contentStr),
			MetaData: map[string]any{
				"order":       docOrder,
				"href":        item.Href,
				"media_type":  item.MediaType,
				"source_uri":  commonOpts.URI,
				"epub_opf":    opfPath,
				"id_ref":      itemRef.IDRef,
			},
		}
		docOrder++

		// Try to extract title from the document
		if title := extractTitle(contentStr); title != "" {
			doc.MetaData["title"] = title
		}

		// Merge extra metadata if provided
		if commonOpts.ExtraMeta != nil {
			for k, v := range commonOpts.ExtraMeta {
				doc.MetaData[k] = v
			}
		}

		documents = append(documents, doc)
	}

	return documents, nil
}

// containerXML represents the structure of container.xml
type containerXML struct {
	Rootfiles []rootfile `xml:"rootfiles>rootfile"`
}

type rootfile struct {
	FullPath string `xml:"full-path,attr"`
}

// metadata represents the metadata section of the OPF file
type metadata struct {
	Title string `xml:"title"`
}

// parseContainerXML parses container.xml to find the OPF file path
func parseContainerXML(r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	var container containerXML
	if err := xml.Unmarshal(data, &container); err != nil {
		return "", err
	}

	if len(container.Rootfiles) == 0 {
		return "", fmt.Errorf("no rootfiles found in container.xml")
	}

	return container.Rootfiles[0].FullPath, nil
}

// extractTitle tries to extract a title from HTML content
func extractTitle(htmlContent string) string {
	// Simple title extraction from <title> tag
	titleStart := strings.Index(htmlContent, "<title>")
	if titleStart == -1 {
		return ""
	}
	
	titleEnd := strings.Index(htmlContent[titleStart:], "</title>")
	if titleEnd == -1 {
		return ""
	}
	
	title := htmlContent[titleStart+7 : titleStart+titleEnd]
	return strings.TrimSpace(title)
}

// opf represents the structure of the OPF file
type opf struct {
	Manifest manifest `xml:"manifest"`
	Spine    spine    `xml:"spine"`
	Metadata metadata `xml:"metadata"`
}

type manifest struct {
	Items []item `xml:"item"`
}

type item struct {
	ID         string `xml:"id,attr"`
	Href       string `xml:"href,attr"`
	MediaType  string `xml:"media-type,attr"`
}

type spine struct {
	ItemRefs []itemRef `xml:"itemref"`
}

type itemRef struct {
	IDRef string `xml:"idref,attr"`
}

// parseOPF parses the OPF file and returns the manifest and spine
func parseOPF(r io.Reader) (map[string]item, spine, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, spine{}, err
	}

	var opfFile opf
	if err := xml.Unmarshal(data, &opfFile); err != nil {
		return nil, spine{}, err
	}

	// Create a map of items by ID for easy lookup
	manifestMap := make(map[string]item)
	for _, item := range opfFile.Manifest.Items {
		manifestMap[item.ID] = item
	}

	return manifestMap, opfFile.Spine, nil
}