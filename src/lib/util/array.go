package util

// ArrayStringShuffle ... string配列をシャッフルする
func ArrayStringShuffle(arr []string) []string {
	n := len(arr)
	for i := n - 1; i >= 0; i-- {
		j := IntRand(0, i+1)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

// ArrayIntShuffle ... int配列をシャッフルする
func ArrayIntShuffle(arr []int) []int {
	n := len(arr)
	for i := n - 1; i >= 0; i-- {
		j := IntRand(0, i+1)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

// ArrayInt64Shuffle ... int64配列をシャッフルする
func ArrayInt64Shuffle(arr []int64) []int64 {
	n := len(arr)
	for i := n - 1; i >= 0; i-- {
		j := IntRand(0, i+1)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

// ArrayStringInsert ... string配列の任意の場所に挿入する
func ArrayStringInsert(arr []string, v string, i int) []string {
	return append(arr[:i], append([]string{v}, arr[i:]...)...)
}

// ArrayIntInsert ... int配列の任意の場所に挿入する
func ArrayIntInsert(arr []int, v int, i int) []int {
	return append(arr[:i], append([]int{v}, arr[i:]...)...)
}

// ArrayInt64Insert ... int64配列の任意の場所に挿入する
func ArrayInt64Insert(arr []int64, v int64, i int) []int64 {
	return append(arr[:i], append([]int64{v}, arr[i:]...)...)
}

// ArrayStringDelete ... string配列の任意の値を削除する
func ArrayStringDelete(arr []string, i int) []string {
	return append(arr[:i], arr[i+1:]...)
}

// ArrayIntDelete ... int配列の任意の値を削除する
func ArrayIntDelete(arr []int, i int) []int {
	return append(arr[:i], arr[i+1:]...)
}

// ArrayInt64Delete ... int64配列の任意の値を削除する
func ArrayInt64Delete(arr []int64, i int) []int64 {
	return append(arr[:i], arr[i+1:]...)
}

// ArrayStringShift ... string配列の先頭を切り取る
func ArrayStringShift(arr []string) (string, []string) {
	return arr[0], arr[1:]
}

// ArrayIntShift ... int配列の先頭を切り取る
func ArrayIntShift(arr []int) (int, []int) {
	return arr[0], arr[1:]
}

// ArrayInt64Shift ... int64配列の先頭を切り取る
func ArrayInt64Shift(arr []int64) (int64, []int64) {
	return arr[0], arr[1:]
}

// ArrayStringBack ... string配列の後尾を切り取る
func ArrayStringBack(arr []string) (string, []string) {
	return arr[len(arr)-1], arr[:len(arr)-1]
}

// ArrayIntBack ... int配列の後尾を切り取る
func ArrayIntBack(arr []int) (int, []int) {
	return arr[len(arr)-1], arr[:len(arr)-1]
}

// ArrayInt64Back ... int64配列の後尾を切り取る
func ArrayInt64Back(arr []int64) (int64, []int64) {
	return arr[len(arr)-1], arr[:len(arr)-1]
}

// ArrayStringFilter ... string配列をフィルタする
func ArrayStringFilter(arr []string, fn func(string) bool) []string {
	ret := []string{}
	for _, v := range arr {
		if fn(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// ArrayIntFilter ... int配列をフィルタする
func ArrayIntFilter(arr []int, fn func(int) bool) []int {
	ret := []int{}
	for _, v := range arr {
		if fn(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// ArrayInt64Filter ... int64配列をフィルタする
func ArrayInt64Filter(arr []int64, fn func(int64) bool) []int64 {
	ret := []int64{}
	for _, v := range arr {
		if fn(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// ArrayStringUniq ... string配列の重複を排除する
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

// ArrayIntUniq ... int配列の重複を排除する
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

// ArrayInt64Uniq ... int64配列の重複を排除する
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

// ArrayStringContains ... string配列の値の存在確認
func ArrayStringContains(arr []string, e string) bool {
	for _, v := range arr {
		if e == v {
			return true
		}
	}
	return false
}

// ArrayIntContains ... int配列の値の存在確認
func ArrayIntContains(arr []int, e int) bool {
	for _, v := range arr {
		if e == v {
			return true
		}
	}
	return false
}

// ArrayInt64Contains ... int64配列の値の存在確認
func ArrayInt64Contains(arr []int64, e int64) bool {
	for _, v := range arr {
		if e == v {
			return true
		}
	}
	return false
}

// ArrayStringChunk ... string配列の分割
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

// ArrayIntChunk ... int配列の分割
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

// ArrayInt64Chunk ... int64配列の分割
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
