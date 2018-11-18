package markov

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// SubDictionary of follow up words and the cordict3ponding factor.
type SubDictionary map[string]int

// Dictionary of start words with follow ups and the cordict3ponding factor.
type Dictionary map[string]SubDictionary

// FitnessFunc signature that is allowed to be passed
type FitnessFunc func(Dictionary, []string) int

// Train from string of words
func Train(text string, factor int) Dictionary {
	dict := make(Dictionary)

	words := strings.Fields(text)

	words = cleanUpStrings(words)

	for i := 0; i < len(words)-1; i++ {
		if _, prefixAvail := dict[words[i]]; !prefixAvail {
			dict[words[i]] = make(SubDictionary)
		}
		if _, suffixAvail := dict[words[i]][words[i+1]]; !suffixAvail {
			dict[words[i]][words[i+1]] = factor
		} else {
			dict[words[i]][words[i+1]] = dict[words[i]][words[i+1]] + dict[words[i]][words[i+1]]*factor
		}
	}
	return dict
}

// TrainFromFile takes a file path, reads the file and passes the string to Train
func TrainFromFile(path string, factor int) Dictionary {
	buf, _ := ioutil.ReadFile(path)
	return Train(string(buf), factor)
}

// TrainFromFolder takes a path, and passes every text file it finds to TrainFromFile
func TrainFromFolder(path string, factor int) Dictionary {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(dir + string(os.PathSeparator) + path)
	if err != nil {
		log.Fatal(err)
	}

	dict := make(Dictionary)

	for _, file := range files {
		match, _ := regexp.MatchString(".*\\.txt$", file.Name())
		if match {
			fmt.Println(file.Name())
			fmt.Println(dir + string(os.PathSeparator) + path + string(os.PathSeparator) + file.Name())
			tmp := TrainFromFile(path+string(os.PathSeparator)+file.Name(), factor)
			dict = mergeDict(dict, tmp)
		}
	}

	return dict
}

// Generate takes a dictionary, a maximum length and a startword to generate a text based on the inputs
func Generate(dict Dictionary, maxLength int, startWord string) string {
	var word = ""
	if startWord == "" {
		word = pickRandom(dict.keys())
	} else {
		word = startWord
	}

	sentence := strings.Title(word)
	i := 0
	for maxLength == 0 || i < maxLength-1 {
		i++

		tmp := word

		for val := range dict[word] {
			if _, undefined := dict[word][val]; undefined {
				word = pickRandom(dict[word].keys())
			}
		}
		if word == tmp || word == "" {
			return sentence
		}
		tmp = word
		sentence = sentence + " " + word
	}
	return sentence
}

// AdjustFactors takes a dictionary and applies the fitness func on the dictionary
func AdjustFactors(dict Dictionary, maxLength int, f FitnessFunc) Dictionary {
	extract := strings.Fields(Generate(dict, maxLength, ""))

	var pairs [][]string
	i := 0
	for i < len(extract)-1 {
		if i >= len(extract) {
			i++
			continue
		}
		pairs = append(pairs, []string{extract[i], extract[i+1]})
		i++
	}
	i = 0
	for p := range pairs {
		fact := int((float64(f(dict, pairs[p])) - 0.5) * 2.0)
		dict = mergeDict(Train(pairs[p][0]+" "+pairs[p][1], fact), dict)
	}
	return dict
}

// BulkAdjustFactors takes a dictionary and runs the number of specified iterations applying the fitness function
func BulkAdjustFactors(dict Dictionary, iterations int, f []FitnessFunc) Dictionary {
	if len(f) < 1 {
		return dict
	}
	i := 0
	for i < iterations {
		i++
		j := 0
		for j < len(f) {
			dict = AdjustFactors(dict, 10, f[j])
			j++
		}
	}

	return dict
}

// FitnessFunction apply fitness to dictionary
func FitnessFunction(dict Dictionary, pair []string) int {
	if pair[1] == "" {
		return -1
	}
	if _, undefined := dict[pair[0]]; !undefined {
		return -1
	}
	return dict[pair[0]][pair[1]]
}

// mergeDict, given two dictionaries merges them into one
func mergeDict(dict1, dict2 Dictionary) Dictionary {
	dict3 := dict1
	for val := range dict2 {
		if _, undefined := dict3[val]; !undefined {
			dict3[val] = dict2[val]
		} else {
			for sval := range dict2[val] {
				if _, undefined := dict3[val][sval]; !undefined {
					dict3[val] = make(SubDictionary)
					dict3[val][sval] = dict2[val][sval]
				} else {
					dict3[val] = make(SubDictionary)
					dict3[val][sval] = dict3[val][sval] + dict2[val][sval]
				}
			}
		}
	}
	return dict3
}

func (m Dictionary) keys() []string {
	keys := make([]string, len(m))
	i := 0
	for val := range m {
		keys[i] = val
		i++
	}
	return keys
}

func (m SubDictionary) keys() []string {
	keys := make([]string, len(m))
	i := 0
	for val := range m {
		keys[i] = val
		i++
	}
	return keys
}

// pickRandom takes a dictionary and selects a random key
func pickRandom(keys []string) string {
	return keys[rand.Intn(len(keys))]
}

func cleanUpStrings(words []string) []string {
	for val := range words {
		// remove all non-alphanumeric characters from input
		reg, err := regexp.Compile("[^\\s\\w]")
		if err != nil {
			log.Fatal(err)
		}
		words[val] = reg.ReplaceAllString(words[val], "")

		// everything lowercase
		words[val] = strings.ToLower(words[val])
	}
	return words
}
