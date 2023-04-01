package jsonutils

type stack []int

func (s *stack) peek() int {
	return (*s)[len(*s)-1]
}
func (s *stack) isEmpty() bool {
	return len(*s) == 0
}
func (s *stack) pop() {
	*s = (*s)[:len(*s)-1]
}
func (s *stack) push(i int) {
	*s = append(*s, i)

}
