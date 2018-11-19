package markov

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// GenerateText returns a text generated from a folder with text files
func GenerateText(folder string, nbPrefix int, nbRuns int, nbWordsPerRun int) string {
	startOnCapital := true
	stopAtSentence := true

	rand.Seed(time.Now().UnixNano())

	// contains all the folders to parse
	var allTextFilesInsideRoot []string

	// recursively read the folder and concatenate the .txt files
	err := filepath.Walk(folder,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(path, ".txt") {
				allTextFilesInsideRoot = append(allTextFilesInsideRoot, path)
			}
			return nil
		})
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBuffer(nil)
	for _, file := range allTextFilesInsideRoot {
		f, _ := os.Open(file) // Error handling elided for brevity.
		io.Copy(buf, f)       // Error handling elided for brevity.
		f.Close()
	}
	allTextFilesContent := string(buf.Bytes())

	m, err := NewMarkovFromString(allTextFilesContent, nbPrefix)
	if err != nil {
		log.Fatal(err)
	}

	var outString string
	outBuffer := bytes.NewBufferString(outString)

	for i := 0; i < nbRuns; i++ {
		err = m.Output(outBuffer, nbWordsPerRun, startOnCapital, stopAtSentence)
		if err != nil {
			log.Fatal(err)
		}
	}
	return outBuffer.String() + " " + os.Getenv("HASHTAGS_MARKOV")
}

// We'd like to use a map of []string -> []string (i.e. list of prefix
// words -> list of possible next words) but Go doesn't allow slices to be
// map keys.
//
// We could use arrays, e.g. map of [2]string -> []string, but array lengths
// are fixed at compile time. To work around that we could set a maximum value
// for n, say 8 or 16, and waste the extra array slots for smaller n.
//
// Or we could make the suffix map key just be the full prefix string. Then
// to get the words within the prefix we could either have a separate map
// (i.e. map of string -> []string) for the full prefix string -> the list
// of the prefix words. Or we could use strings.Fields() and strings.Join() to
// go back and forth (trading more runtime for less memory use).

// Markov is a Markov chain text generator.
type Markov struct {
	n           int
	capitalized int // number of suffix keys that start capitalized
	suffix      map[string][]string
}

// NewMarkovFromString initializes the Markov text generator
// with window `n` from the contents of a string s.
func NewMarkovFromString(s string, n int) (*Markov, error) {
	m := &Markov{
		n:      n,
		suffix: make(map[string][]string),
	}

	var r = bytes.NewBufferString(s)

	sc := bufio.NewScanner(r)
	sc.Split(bufio.ScanWords)
	window := make([]string, 0, n)
	for sc.Scan() {
		word := sc.Text()
		if len(window) > 0 {
			prefix := strings.Join(window, " ")
			m.suffix[prefix] = append(m.suffix[prefix], word)
			//log.Printf("%20q -> %q", prefix, m.suffix[prefix])
			if isCapitalized(prefix) {
				m.capitalized++
			}
		}
		window = appendMax(n, window, word)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return m, nil
}

// Output writes generated text of approximately `n` words to `w`.
// If `startCapital` is true it picks a starting prefix that is capitalized.
// If `stopSentence` is true it continues after `n` words until it finds
// a suffix ending with sentence ending punctuation ('.', '?', or '!').
func (m *Markov) Output(w io.Writer, n int, startCapital, stopSentence bool) error {
	// Use a bufio.Writer both for buffering and for simplified
	// error handling (it remembers any error and turns all future
	// writes/flushes into NOPs returning the same error).
	bw := bufio.NewWriter(w)

	var i int
	if startCapital {
		i = rand.Intn(m.capitalized)
	} else {
		i = rand.Intn(len(m.suffix))
	}
	var prefix string
	for prefix = range m.suffix {
		if startCapital && !isCapitalized(prefix) {
			continue
		}
		if i == 0 {
			break
		}
		i--
	}

	bw.WriteString(prefix) // nolint: errcheck
	prefixWords := strings.Fields(prefix)
	n -= len(prefixWords)

	for {
		suffixChoices := m.suffix[prefix]
		if len(suffixChoices) == 0 {
			break
		}
		i = rand.Intn(len(suffixChoices))
		suffix := suffixChoices[i]
		//log.Printf("prefix: %q, suffix: %q (from %q)", prefixWords, suffix, suffixChoices)
		bw.WriteByte(' ') // nolint: errcheck
		if _, err := bw.WriteString(suffix); err != nil {
			break
		}
		n--
		if n < 0 && (!stopSentence || isSentenceEnd(suffix)) {
			break
		}

		prefixWords = appendMax(m.n, prefixWords, suffix)
		prefix = strings.Join(prefixWords, " ")
	}
	return bw.Flush()
}

func isCapitalized(s string) bool {
	// We can't just look at s[0], which is the first *byte*,
	// if we want to support arbitrary Unicode input.
	// This still doesn't support combining runes :(.
	r, _ := utf8.DecodeRuneInString(s)
	return unicode.IsUpper(r)
}

func isSentenceEnd(s string) bool {
	r, _ := utf8.DecodeLastRuneInString(s)
	// Unfortunately, Unicode doesn't seem to provide
	// a test for sentence ending punctution :(.
	//return unicode.IsPunct(r)
	return r == '.' || r == '?' || r == '!'
}

func appendMax(max int, slice []string, value string) []string {
	// Often FIFO queues in Go are implemented via:
	//     fifo = append(fifo, newValues...)
	// and:
	//     fifo = fifo[numberOfValuesToRemove:]
	//
	// However, the append will periodically reallocate and copy. Since
	// we're dealing with a small number (usually two) of strings and we
	// only need to append a single new string it's better to (almost)
	// never reallocate the slice and just copy n-1 strings (which only
	// copies n-1 pointers, not the entire string contents) every time.
	if len(slice)+1 > max {
		n := copy(slice, slice[1:])
		slice = slice[:n]
	}
	return append(slice, value)
}
