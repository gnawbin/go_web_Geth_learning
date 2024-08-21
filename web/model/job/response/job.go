package response

// @Author: morris
// NAME    COMPLETIONS   DURATION   AGE
type Job struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	//控制Job成功完成的实例数目的  当指定的实例数目达到Completions字段所设定的值时，Job将被标记为成功完成
	Completions int32 `json:"completions"`
	//就绪的Job个数
	Succeeded int32 `json:"succeeded"`
	//Job的持续时间
	Duration int64 `json:"duration"`
	Age      int64 `json:"age"`
}
