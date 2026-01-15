package gcr

func PlaceHolderRange(s SqlSyntax, startAt, stopBefore int) []any {
	phs := make([]any, 0, stopBefore-startAt)
	for i := startAt; i < stopBefore; i++ {
		phs = append(phs, s.PlaceHolder(i))
	}
	return phs
}
