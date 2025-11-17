package common

// type BinaryNode[T expression.ExpressionNode] struct {
// 	Left     *T
// 	Operator *Token
// 	Right    *BinaryNode[T]
// }

// type OrNode = BinaryNode[AndNode]
// type AndNode = BinaryNode[EqualityNode]
// type EqualityNode = BinaryNode[ComparisonNode]
// type ComparisonNode = BinaryNode[TermNode]
// type TermNode = BinaryNode[FactorNode]
// type FactorNode = BinaryNode[Unary]
// type

// func (bn *BinaryNode[T]) HasLeft() bool {
// 	return bn.Left != nil
// }

// func (bn *BinaryNode[T]) HasOperator() bool {
// 	return bn.Operator != nil
// }

// func (bn *BinaryNode[T]) HasRight() bool {
// 	return bn.Right != nil
// }

// func (bn *BinaryNode[T]) SetLeft(node *T) {
// 	bn.Left = node
// }

// func (bn *BinaryNode[T]) SetOperator(token *Token) {
// 	bn.Operator = token
// }

// func (bn *BinaryNode[T]) SetRight(node *BinaryNode[T]) {
// 	bn.Right = node
// }

// func (bn *BinaryNode[T]) GetLeft() MergableNode {
// 	if bn.Left == nil {
// 		panic("GetLeft: Left node is nil")
// 	}
// 	return bn.Left
// }

// func (bn *BinaryNode[T]) GetOperator() *Token {
// 	if bn.Operator == nil {
// 		panic("GetOperator: Operator is nil")
// 	}

// 	return bn.Operator
// }

// func (bn *BinaryNode[T]) GetRight() MergableNode {
// 	if bn.Right == nil {
// 		panic("GetRight: Right node is nil")
// 	}
// 	return bn.Right
// }
