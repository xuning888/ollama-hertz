package util

import (
	"fmt"
	"strings"
)

type TextSplitter struct {
	separators   []string
	chunkSize    int
	chunkOverlap int
}

func NewTextSplitterWithArgs(separators []string, chunkSize int, chunkOverlap int) *TextSplitter {
	if separators == nil {
		separators = []string{"\n\n", "\n", " ", ""}
	}
	return &TextSplitter{
		separators:   separators,
		chunkSize:    chunkSize,
		chunkOverlap: chunkOverlap,
	}
}

func NewTextSplitter() *TextSplitter {
	return NewTextSplitterWithArgs(nil, 500, 50)
}

// splitText 方法，递归分割文本
func (r *TextSplitter) splitText(text string) []string {
	var finalChunks []string
	separator := r.separators[len(r.separators)-1]

	for _, s := range r.separators {
		if strings.Contains(text, s) || s == "" {
			separator = s
			break
		}
	}

	splits := strings.Split(text, separator)
	var goodSplits []string

	for _, s := range splits {
		if len(s) < r.chunkSize {
			goodSplits = append(goodSplits, s)
		} else {
			if len(goodSplits) > 0 {
				mergedText := r.mergeSplits(goodSplits, separator)
				finalChunks = append(finalChunks, mergedText...)
				goodSplits = []string{}
			}
			otherInfo := r.splitText(s)
			finalChunks = append(finalChunks, otherInfo...)
		}
	}

	if len(goodSplits) > 0 {
		mergedText := r.mergeSplits(goodSplits, separator)
		finalChunks = append(finalChunks, mergedText...)
	}

	return finalChunks
}

// mergeSplits 方法，合并小块成文本
func (r *TextSplitter) mergeSplits(splits []string, separator string) []string {
	separatorLen := len(separator)
	var docs []string
	var currentDoc []string
	total := 0

	for _, d := range splits {
		lenD := len(d)
		if total+lenD+(func() int {
			if separatorLen > 0 && len(currentDoc) != 0 {
				return separatorLen
			} else {
				return 0
			}
		}()) > r.chunkSize {
			if total > r.chunkSize {
				fmt.Printf("Warning: Created a chunk of size %d, which is longer than the specified %d\n", total, r.chunkSize)
			}
			if len(currentDoc) > 0 {
				doc := joinDocs(currentDoc, separator)
				if doc != "" {
					docs = append(docs, doc)
				}
				for total > r.chunkOverlap || (total+lenD+(func() int {
					if separatorLen > 0 && len(currentDoc) != 0 {
						return separatorLen
					} else {
						return 0
					}
				}()) > r.chunkSize && total > 0) {
					total -= len(currentDoc[0]) + (func() int {
						if separatorLen > 0 && len(currentDoc) > 1 {
							return separatorLen
						} else {
							return 0
						}
					}())
					currentDoc = currentDoc[1:]
				}
			}
		}
		currentDoc = append(currentDoc, d)
		total += lenD + (func() int {
			if separatorLen > 0 && len(currentDoc) > 1 {
				return separatorLen
			} else {
				return 0
			}
		}())
	}

	doc := joinDocs(currentDoc, separator)
	if doc != "" {
		docs = append(docs, doc)
	}

	return docs
}

// joinDocs 方法，连接文档块
func joinDocs(docs []string, separator string) string {
	if len(docs) == 0 {
		return ""
	}
	return strings.Join(docs, separator)
}
