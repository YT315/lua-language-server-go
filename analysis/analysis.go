package analysis

//Analysis 语义分析器
type Analysis struct {
	previous *Analysis //对象包含依赖,分析过程通过单项链表连接
	FilePath string    //正在分析的文件路径包含文件名
}
