package util

// ArrayStringShuffle ... shuffle string array
func ArrayStringShuffle(arr []string) []string {
	n := len(arr)
	for i := n - 1; i >= 0; i-- {
		j := IntRand(0, i+1)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

// ArrayIntShuffle ... shuffle int array
func ArrayIntShuffle(arr []int) []int {
	n := len(arr)
	for i := n - 1; i >= 0; i-- {
		j := IntRand(0, i+1)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

// ArrayInt64Shuffle ... shuffle int64 array
func ArrayInt64Shuffle(arr []int64) []int64 {
	n := len(arr)
	for i := n - 1; i >= 0; i-- {
		j := IntRand(0, i+1)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

// ArrayStringInsert ... insert string value to array
func ArrayStringInsert(arr []string, v string, i int) []string {
	return append(arr[:i], append([]string{v}, arr[i:]...)...)
}

// ArrayIntInsert ... insert int value to array
func ArrayIntInsert(arr []int, v int, i int) []int {
	return append(arr[:i], append([]int{v}, arr[i:]...)...)
}

// ArrayInt64Insert ... insert int64 value to array
func ArrayInt64Insert(arr []int64, v int64, i int) []int64 {
	return append(arr[:i], append([]int64{v}, arr[i:]...)...)
}

// ArrayStringDelete ... delete string value from array
func ArrayStringDelete(arr []string, i int) []string {
	return append(arr[:i], arr[i+1:]...)
}

// ArrayIntDelete ... delete int value from array
func ArrayIntDelete(arr []int, i int) []int {
	return append(arr[:i], arr[i+1:]...)
}

// ArrayInt64Delete ... delete int64 value from array
func ArrayInt64Delete(arr []int64, i int) []int64 {
	return append(arr[:i], arr[i+1:]...)
}

// ArrayStringShift ... shift string in array
func ArrayStringShift(arr []string) (string, []string) {
	return arr[0], arr[1:]
}

// ArrayIntShift ... shift int value in array
func ArrayIntShift(arr []int) (int, []int) {
	return arr[0], arr[1:]
}

// ArrayInt64Shift ... shift value from array of int64
func ArrayInt64Shift(arr []int64) (int64, []int64) {
	return arr[0], arr[1:]
}

// ArrayStringBack ... divide value from back of array
func ArrayStringBack(arr []string) (string, []string) {
	return arr[len(arr)-1], arr[:len(arr)-1]
}

// ArrayIntBack ... divide value from back of array
func ArrayIntBack(arr []int) (int, []int) {
	return arr[len(arr)-1], arr[:len(arr)-1]
}

// ArrayInt64Back ... divide value from back of array
func ArrayInt64Back(arr []int64) (int64, []int64) {
	return arr[len(arr)-1], arr[:len(arr)-1]
}

// ArrayStringFilter ... filter array
func ArrayStringFilter(arr []string, fn func(string) bool) []string {
	ret := []string{}
	for _, v := range arr {
		if fn(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// ArrayIntFilter ... filter array
func ArrayIntFilter(arr []int, fn func(int) bool) []int {
	ret := []int{}
	for _, v := range arr {
		if fn(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// ArrayInt64Filter ... filter array
func ArrayInt64Filter(arr []int64, fn func(int64) bool) []int64 {
	ret := []int64{}
	for _, v := range arr {
		if fn(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// ArrayStringUniq ... make unique array
func ArrayStringUniq(arr []string) []string {
	m := make(map[string]bool)
	uniq := []string{}
	for _, v := range arr {
		if !m[v] {
			m[v] = true
			uniq = append(uniq, v)
		}
	}
	return uniq
}

// ArrayIntUniq ... make unique array
func ArrayIntUniq(arr []int) []int {
	m := make(map[int]bool)
	uniq := []int{}
	for _, v := range arr {
		if !m[v] {
			m[v] = true
			uniq = append(uniq, v)
		}
	}
	return uniq
}

// ArrayInt64Uniq ... make unique array
func ArrayInt64Uniq(arr []int64) []int64 {
	m := make(map[int64]bool)
	uniq := []int64{}
	for _, v := range arr {
		if !m[v] {
			m[v] = true
			uniq = append(uniq, v)
		}
	}
	return uniq
}

// ArrayStringContains ... exist value in array
func ArrayStringContains(arr []string, e string) bool {
	for _, v := range arr {
		if e == v {
			return true
		}
	}
	return false
}

// ArrayIntContains ... exist value in array
func ArrayIntContains(arr []int, e int) bool {
	for _, v := range arr {
		if e == v {
			return true
		}
	}
	return false
}

// ArrayInt64Contains ... exist value in array
func ArrayInt64Contains(arr []int64, e int64) bool {
	for _, v := range arr {
		if e == v {
			return true
		}
	}
	return false
}

// ArrayStringChunk ... chunk array
func ArrayStringChunk(arr []string, size int) [][]string {
	var chunks [][]string
	arrSize := len(arr)
	for i := 0; i < arrSize; i += size {
		end := i + size
		if arrSize < end {
			end = arrSize
		}
		chunks = append(chunks, arr[i:end])
	}
	return chunks
}

// ArrayIntChunk ... chunk array
func ArrayIntChunk(arr []int, size int) [][]int {
	var chunks [][]int
	arrSize := len(arr)
	for i := 0; i < arrSize; i += size {
		end := i + size
		if arrSize < end {
			end = arrSize
		}
		chunks = append(chunks, arr[i:end])
	}
	return chunks
}

// ArrayInt64Chunk ... chunk array
func ArrayInt64Chunk(arr []int64, size int) [][]int64 {
	var chunks [][]int64
	arrSize := len(arr)
	for i := 0; i < arrSize; i += size {
		end := i + size
		if arrSize < end {
			end = arrSize
		}
		chunks = append(chunks, arr[i:end])
	}
	return chunks
}
