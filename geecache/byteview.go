package project_for_practice

/*
	抽象一个只读数据结构 byteview：用来表示缓存值
	是geecache主要的数据结构之一
 */

// byteview 只有一个字段
// b会储存真实的缓存值， []byte类型可以存任意数据类型 ex：string 图片等
type ByteView struct {
	b	[]byte
}




func (v ByteView) Len() int {
	return len(v.b)
}

// b是只读！！所以使用byteslice方法返回拷贝！！，避免缓存被修改

func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c

}


