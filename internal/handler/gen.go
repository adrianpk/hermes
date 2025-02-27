package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/yosssi/gohtml"

	"github.com/adrianpk/gohermes/internal/hermes"
)

const debug = false

// GenHTML generates the HTML files from the markdown files.
func GenHTML() error {
	err := hermes.CheckHermes()
	if err != nil {
		return err
	}

	pp, err := startPreProcessor(hermes.ContentDir)
	if err != nil {
		return err
	}

	err = filepath.Walk(hermes.ContentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".md" {
			relativePath, err := filepath.Rel(hermes.ContentDir, path)
			if err != nil {
				return nil
			}

			fileData, exists := pp.FindFileData(relativePath)
			if !exists {
				return nil
			}

			if !fileData.Published {
				return nil
			}

			outputPath := determineOutputPath(relativePath)

			if shouldRender(path, outputPath) {
				fileContent, err := os.ReadFile(path)
				if err != nil {
					return nil
				}

				content, err := hermes.Parse(fileContent, path)
				if err != nil {
					return nil
				}

				layoutPath := findLayout(path)
				if layoutPath != "" {
					tmpl, err := template.New("webpage").Funcs(template.FuncMap{
						"safeHTML": safeHTML,
					}).ParseFiles(layoutPath)

					if err != nil {
						log.Printf("error parsing template files for %s: %v", path, err)
						return nil
					}

					var tmplBuf bytes.Buffer

					err = tmpl.Execute(&tmplBuf, content)
					if err != nil {
						log.Printf("error executing template for %s: %v", path, err)
						return nil
					}

					gohtml.Condense = true
					formattedHTML := gohtml.Format(tmplBuf.String())

					err = os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
					if err != nil {
						log.Printf("error creating directories for %s: %v", outputPath, err)
						return nil
					}

					outputFile, err := os.Create(outputPath)
					if err != nil {
						log.Printf("error creating output file %s: %v", outputPath, err)
						return nil
					}
					defer outputFile.Close()

					_, err = outputFile.WriteString(formattedHTML)
					if err != nil {
						log.Printf("error writing to output file %s: %v", outputPath, err)
						return nil
					}

					err = copyImages(path, outputPath)
					if err != nil {
						//log.Printf("error copying images for %s: %v", path, err)
						return nil
					}

				} else {
					log.Printf("no layout found for %s", path)
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	err = processRootSection(hermes.ContentDir, hermes.OutputDir, "assets/layout", pp)
	if err != nil {
		return fmt.Errorf("error processing root section: %w", err)
	}

	cfg, err := hermes.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	sections := cfg.Sections

	err = processSections(hermes.ContentDir, hermes.OutputDir, "assets/layout", pp, sections)
	if err != nil {
		return fmt.Errorf("error processing sections: %w", err)
	}

	err = addNoJekyll()
	if err != nil {
		return fmt.Errorf("error adding .nojekyll file: %w", err)
	}

	err = copyCSS("assets/css", filepath.Join(hermes.OutputDir, "css"))
	if err != nil {
		return fmt.Errorf("error copying CSS directory: %w", err)
	}

	log.Println("content generated!")
	return nil
}

// copyCSS copies the assets/css directory to the output/css directory.
func copyCSS(srcDir, destDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(destDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, os.ModePerm)
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		return err
	})
}

// genDefaultIndex generates the default index.html file.
// If no index.md is provided for the section then this one is used as default. It renders a list of all published
// content for this section.
// This is a WIP, a lot of logging is present to help debug the process. It will be removed as soon the logic is
// stable. Also, there are a lot of hardcoded values that will be replaced by dynamic ones.
// Finally, this is generating the default index for the root section, it includes references to all the content in the
// site.
// We will also need a similar logic to generate the section index when content is not provided for it.
// This will show all the content for the specific section.
// Worth mentioning that the partial used to render the content should also be improved to show a nice presentation of
// the content (image, title, excerpt, etc).
func genDefaultIndex(pp *hermes.PreProcessor, rootIndexPath, outputPath string, fd []hermes.FileData) error {
	// TODO : Replace this hardoded value by const based generated one
	partial := filepath.Join("assets", "layout", "default", "partial", "_index.html")
	log.Printf("using partial template: %s\n", partial)

	partialTmpl, err := template.New("_index.html").ParseFiles(partial)
	if err != nil {
		return err
	}

	var partialBuf bytes.Buffer

	err = partialTmpl.Execute(&partialBuf, fd)
	if err != nil {
		return err
	}

	content := map[string]interface{}{
		"HTML": partialBuf.String(),
	}

	layoutPath := findLayout(rootIndexPath)
	if layoutPath == "" {
		return fmt.Errorf("no layout found for %s", rootIndexPath)
	}

	layoutTmpl, err := template.New("webpage").Funcs(template.FuncMap{
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
	}).ParseFiles(layoutPath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	var finalBuf bytes.Buffer

	err = layoutTmpl.Execute(&finalBuf, content)
	if err != nil {
		return err
	}

	_, err = finalBuf.WriteTo(outputFile)
	if err != nil {
		log.Printf("error writing to output file: %v\n", err)
		return err
	}

	return nil
}

func processRootSection(contentDir, outputDir, layoutDir string, pp *hermes.PreProcessor) error {
	log.Println("processing root section")

	indexPath := filepath.Join(hermes.ContentDir, "root", "page", "index.md")
	outputPath := filepath.Join(hermes.OutputDir, "index.html")

	if !isValidIndex(indexPath, pp) {
		err := genDefaultIndex(pp, indexPath, outputPath, pp.GetAllPublished())
		if err != nil {
			return err
		}
	}

	return nil
}

func processSections(contentDir, outputDir, layoutDir string, pp *hermes.PreProcessor, sections []hermes.Section) error {
	for _, section := range sections {
		if section.Name == "root" {
			continue
		}

		indexPath := filepath.Join(hermes.ContentDir, section.Name, hermes.IndexMdFile)
		outputPath := filepath.Join(hermes.OutputDir, section.Name, hermes.IndexFile)
		layoutPath := filepath.Join(layoutDir, "index.html")

		log.Printf("layout path: %s\n", layoutPath)

		if !isValidIndex(indexPath, pp) {
			log.Printf("section index is not valid for section: %s, generating fallback index", section.Name)

			err := genDefaultIndex(pp, indexPath, outputPath, pp.GetPublishedBySection(section.Name))
			if err != nil {
				log.Printf("error generating fallback index for section: %s, error: %v\n", section.Name, err)
				continue
			}
		} else {
			log.Printf("section index is valid for section: %s, no need to generate fallback index", section.Name)
		}
	}

	log.Println("finished process sections")
	return nil
}

func isValidIndex(indexPath string, pp *hermes.PreProcessor) bool {
	relPath := strings.TrimPrefix(indexPath, "content/")

	fileData, _ := pp.FindFileData(relPath)

	return fileData.IsIndex()
}

func determineOutputPath(relativePath string) string {
	parts := strings.Split(relativePath, string(os.PathSeparator))

	if len(parts) < 2 {
		return outputPath(parts...)
	}

	section := parts[0]
	subdir := parts[1]
	needDir := needsCustomDir(subdir)

	switch section {
	case hermes.DefSection:
		if needDir {
			p := outputPath(append([]string{subdir}, parts[2:]...)...)
			return p
		} else {
			p := outputPath(append([]string{}, parts[2:]...)...)
			return p
		}

	case ct.Page, ct.Article:
		p := outputPath(parts[1:]...)
		return p

	default:
		if needDir {
			p := outputPath(append([]string{section, subdir}, parts[2:]...)...)
			return p
		} else {
			p := outputPath(append([]string{section}, parts[2:]...)...)
			return p
		}
	}
}

func needsCustomDir(dir string) bool {
	return dir != ct.Article && dir != ct.Page
}

func outputPath(parts ...string) string {
	trimmedPath := strings.TrimSuffix(strings.Join(parts, string(os.PathSeparator)), filepath.Ext(parts[len(parts)-1])) + ".html"
	return filepath.Join(hermes.OutputDir, trimmedPath)
}

// shouldRender checks if the markdown file is newer than the html file
// to determine if the html should be re-rendered.
func shouldRender(mdPath, htmlPath string) bool {
	htmlInfo, err := os.Stat(htmlPath)
	if os.IsNotExist(err) {
		return true
	}

	markdownInfo, err := os.Stat(mdPath)
	if err != nil {
		return false
	}

	return markdownInfo.ModTime().After(htmlInfo.ModTime())
}

func findLayout(path string) (layout string) {
	logDebug("")
	logDebug("=== findLayout start ===")
	defer func() {
		logDebug("=== findLayout end ===")
		logDebug("")
	}()

	path = strings.TrimPrefix(path, "content/")
	base := filepath.Base(path)
	base = strings.TrimSuffix(base, filepath.Ext(base))
	section, contentType := sectionTypeSegments(path)
	layoutDir := filepath.Join("assets", "layout")

	secLayoutDir := filepath.Join(layoutDir, section)
	secTypeLayoutDir := filepath.Join(layoutDir, section, contentType)
	defTypeLayoutDir := filepath.Join(layoutDir, hermes.DefLayoutDir, contentType)
	defLayoutDir := filepath.Join(layoutDir, hermes.DefLayoutDir)
	logDebug("section layout directory: " + defLayoutDir)

	layoutPaths := []string{
		filepath.Join(secTypeLayoutDir, base+".html"),
		filepath.Join(secTypeLayoutDir, hermes.DefLayout),
		filepath.Join(secLayoutDir, base+".html"),
		filepath.Join(secLayoutDir, hermes.DefLayout),
		filepath.Join(defTypeLayoutDir, base+".html"),
		filepath.Join(defTypeLayoutDir, hermes.DefLayout),
		filepath.Join(defLayoutDir, base+".html"),
		filepath.Join(defLayoutDir, hermes.DefLayout),
	}

	for _, layoutPath := range layoutPaths {
		if _, err := os.Stat(layoutPath); err == nil {
			logDebug("found layout path: " + layoutPath)
			return layoutPath
		}
	}

	logDebug("no layout path found")
	return layout
}

func sectionTypeSegments(path string) (sectionSegment string, typeSegment string) {
	dir := filepath.Dir(path)
	segments := strings.Split(dir, osFileSep)

	if len(segments) > 0 {
		sectionSegment = segments[0]
	}
	if len(segments) > 1 {
		typeSegment = segments[1]
	}

	return sectionSegment, typeSegment
}

// copyImages copies the images from the markdown directory to the output directory.
func copyImages(mdPath, htmlPath string) error {
	rootPrefix := "root/"
	imageDir := strings.TrimSuffix(mdPath, filepath.Ext(mdPath))
	relativeImageDir := strings.TrimPrefix(imageDir, hermes.ContentDir+"/")

	if strings.HasPrefix(relativeImageDir, rootPrefix) {
		relativeImageDir = strings.TrimPrefix(relativeImageDir, rootPrefix)

		parts := strings.Split(relativeImageDir, string(osFileSep))

		if len(parts) == 1 {
			relativeImageDir = filepath.Join(hermes.ImgDir, parts[0])
		} else if len(parts) > 1 {
			switch parts[0] {
			case ct.Blog, ct.Series:
				relativeImageDir = filepath.Join(hermes.ImgDir, parts[0], parts[1])
			case ct.Article, ct.Page:
				relativeImageDir = filepath.Join(hermes.ImgDir, parts[1])
			default:
				relativeImageDir = filepath.Join(hermes.ImgDir, strings.Join(parts, string(osFileSep)))
			}
		}
	} else {
		parts := strings.Split(relativeImageDir, string(osFileSep))

		if len(parts) > 2 && (parts[1] == ct.Article || parts[1] == ct.Page) {
			relativeImageDir = filepath.Join(hermes.ImgDir, parts[0], parts[2])
		} else if len(parts) > 2 && (parts[1] == ct.Blog || parts[1] == ct.Series) {
			relativeImageDir = filepath.Join(hermes.ImgDir, parts[0], parts[1], parts[2])
		} else {
			relativeImageDir = filepath.Join(hermes.ImgDir, strings.Join(parts, string(os.PathSeparator)))
		}
	}

	outputImageDir := filepath.Join(hermes.OutputDir, relativeImageDir)

	err := filepath.Walk(imageDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relativePath, err := filepath.Rel(imageDir, path)
			if err != nil {
				return err
			}

			destPath := filepath.Join(outputImageDir, relativePath)

			err = os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
			if err != nil {
				return err
			}

			srcFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			destFile, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer destFile.Close()

			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// addNoJekyll adds a .nojekyll file to the output directory.
func addNoJekyll() error {
	noJekyllPath := filepath.Join(hermes.OutputDir, hermes.NoJekyllFile)
	if _, err := os.Stat(noJekyllPath); os.IsNotExist(err) {
		file, err := os.Create(noJekyllPath)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}

func startPreProcessor(root string) (*hermes.PreProcessor, error) {
	pp := hermes.NewPreProcessor(root)
	err := pp.Build()
	if err != nil {
		log.Printf("error building pp: %v", err)
		return nil, err
	}

	err = pp.Sync()
	if err != nil {
		log.Printf("error syncing pp: %v", err)
		return nil, err
	}

	return pp, nil
}

// safeHTML function to replace the anonymous function
func safeHTML(s string) template.HTML {
	return template.HTML(s)
}

func logDebug(msg string) {
	if debug {
		log.Println(msg)
	}
}
