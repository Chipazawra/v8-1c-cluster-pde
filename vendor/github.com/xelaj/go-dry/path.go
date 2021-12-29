package dry

import (
	"path/filepath"
	"strings"
)

func PathWithoutExt(filename string) string {
	ext := filepath.Ext(filename)
	return strings.TrimSuffix(filename, ext)
}

// делить название файла на само название и расширение
// при этом на название не влияет путь, в котором расположен файл
func PathSplitExt(path string) (basepath, ext string) {
	filename := filepath.Base(path)
	if filename == "." {
		return "", ""
	}

	hidden := false
	if strings.HasPrefix(filename, ".") {
		hidden = true
		filename = strings.TrimPrefix(filename, ".")
	}

	ext = filepath.Ext(filename)
	basepath = filename[:len(filename)-len(ext)]
	if hidden {
		basepath = "." + basepath
	}
	ext = strings.TrimPrefix(ext, ".")

	return
}
