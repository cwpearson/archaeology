package levenshtein

import "fmt"
import "github.com/cwpearson/archaeology/archaeology"

type EditScript []EditOperation

type EditOperation interface {
	fmt.Stringer
	Apply(source, target *archaeology.BlockView, matrix [][]int, i, j int) (int, bool)
	Backtrack(matrix [][]int, i, j int) (int, int)
}

var _ EditOperation = &Insertion{}

var _ EditOperation = &Match{}

type Match struct {
}

func (o Match) Apply(source, target *archaeology.BlockView, matrix [][]int, i, j int) (int, bool) {
	if i > 0 && j > 0 {
		src, err := source.Get(i - 1)
		if err != nil {
			panic(err)
		}
		tgt, err := target.Get(j - 1)
		if err != nil {
			panic(err)
		}
		if src.Equals(tgt) {
			return matrix[i-1][j-1], true
		}
	}

	return 0, false
}

func (o Match) Backtrack(matrix [][]int, i, j int) (int, int) {
	return i - 1, j - 1
}

func (o Match) String() string {
	return "match"
}

type Insertion struct {
	Cost int
}

func (o Insertion) Apply(source, target *archaeology.BlockView, matrix [][]int, i, j int) (int, bool) {
	if j > 0 {
		return matrix[i][j-1] + o.Cost, true
	}

	return 0, false
}

func (o Insertion) Backtrack(matrix [][]int, i, j int) (int, int) {
	return i, j - 1
}

func (o Insertion) String() string {
	return "ins"
}

var _ EditOperation = &Deletion{}

type Deletion struct {
	Cost int
}

func (o Deletion) Apply(source, target *archaeology.BlockView, matrix [][]int, i, j int) (int, bool) {
	if i > 0 {
		return matrix[i-1][j] + o.Cost, true
	}

	return 0, false
}

func (o Deletion) Backtrack(matrix [][]int, i, j int) (int, int) {
	return i - 1, j
}

func (o Deletion) String() string {
	return "del"
}

var _ EditOperation = &Substitution{}

type Substitution struct {
	Cost int
}

func (o Substitution) Apply(source, target *archaeology.BlockView, matrix [][]int, i, j int) (int, bool) {
	if i > 0 && j > 0 {
		return matrix[i-1][j-1] + o.Cost, true
	}

	return 0, false
}

func (o Substitution) Backtrack(matrix [][]int, i, j int) (int, int) {
	return i - 1, j - 1
}

func (o Substitution) String() string {
	return "del"
}

var _ EditOperation = &Transposition{}

type Transposition struct {
	Cost int
}

func (o Transposition) Apply(source, target *archaeology.BlockView, matrix [][]int, i, j int) (int, bool) {
	if i > 1 && j > 1 {
		if source == nil && target == nil {
			return matrix[i-2][j-2] + o.Cost, true
		}

		srcm1, err := source.Get(i - 1)
		if err != nil {
			panic(err)
		}
		tgtm1, err := target.Get(j - 1)
		if err != nil {
			panic(err)
		}
		srcm2, err := source.Get(i - 2)
		if err != nil {
			panic(err)
		}
		tgtm2, err := target.Get(j - 2)
		if err != nil {
			panic(err)
		}

		if srcm1.Equals(tgtm2) && srcm2.Equals(tgtm1) {
			return matrix[i-2][j-2] + o.Cost, true
		}
	}

	return 0, false
}

func (o Transposition) Backtrack(matrix [][]int, i, j int) (int, int) {
	return i - 2, j - 2
}

func (o Transposition) String() string {
	return "trp"
}
