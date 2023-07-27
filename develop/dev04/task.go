package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func makeAnagrams(words []string) map[string][]string {
	sets := make(map[string]map[string]struct{})
	for _, word := range words {
		word = strings.ToLower(word)
		hash := getHash(word)
		if sets[hash] == nil {
			sets[hash] = make(map[string]struct{})
		}
		sets[hash][word] = struct{}{}
	}

	ret := make(map[string][]string, len(sets))
	for _, set := range sets {
		if len(set) <= 1 {
			continue
		}

		sliceSet := make([]string, 0, len(set))
		for word := range set {
			sliceSet = append(sliceSet, word)
		}

		sort.Slice(sliceSet, func(i, j int) bool {
			return sliceSet[i] < sliceSet[j]
		})

		ret[sliceSet[0]] = sliceSet
		//if we dont need key word in set, than:
		//ret[sliceSet[0]] = sliceSet[1:]
	}

	return ret
}

const (
	alpPower    = 33
	startSymbol = 'а'
	separator   = '_'
)

// must generate one hash for two words that are anagrams
func getHash(s string) string {
	counts := make([]int, alpPower)
	for _, r := range s {
		counts[int(r)-int(startSymbol)]++
	}
	sb := strings.Builder{}
	for i, count := range counts {
		if i != 0 {
			sb.WriteRune(separator)
		}
		if count != 0 {
			sb.WriteString(strconv.Itoa(count))
		}
	}
	return sb.String()
}

func main() {
	words := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик"}
	fmt.Println(makeAnagrams(words))
}
