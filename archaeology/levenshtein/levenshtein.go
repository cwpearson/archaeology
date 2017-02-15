package levenshtein

import "math"
import "github.com/cwpearson/archaeology/archaeology"

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a int, b int) int {
	if b > a {
		return b
	}
	return a
}

func Matrix(src, tgt *archaeology.BlockView, ops []EditOperation) [][]int {

	if src.Length() > math.MaxInt32 {
		panic("files have too many blocks")
	}
	if tgt.Length() > math.MaxInt32 {
		panic("files have too many blocks")
	}

	rows := int(src.Length()) + 1
	cols := int(tgt.Length()) + 1
	mat := make([][]int, rows)

	for r := 0; r < rows; r++ {
		mat[r] = make([]int, cols)
		mat[r][0] = r
	}
	for c := 1; c < cols; c++ {
		mat[0][c] = c
	}

	for i := 1; i < rows; i++ {
		for j := 1; j < cols; j++ {
			lowestCost := math.MaxInt32
			for _, op := range ops {
				if cost, ok := op.Apply(src, tgt, mat, i, j); ok && cost < lowestCost {
					lowestCost = cost
				}
			}

			mat[i][j] = lowestCost
		}
	}
	return mat
}

func backtrace(i int, j int, matrix [][]int, ops []EditOperation) EditScript {
	for _, op := range ops {
		ib, jb := op.Backtrack(matrix, i, j)

		if ib < 0 || jb < 0 {
			continue
		}

		if cost, ok := op.Apply(nil, nil, matrix, ib, jb); ok && cost == matrix[i][j] {
			return append(backtrace(ib, jb, matrix, ops), op)
		}
	}

	return EditScript{}
}

func EditScriptForStrings(source, target *archaeology.BlockView, ops []EditOperation) EditScript {
	if source.Length() > math.MaxInt32 {
		panic("files have too many blocks")
	}
	if target.Length() > math.MaxInt32 {
		panic("files have too many blocks")
	}

	mat := Matrix(source, target, ops)

	return backtrace(int(source.Length()), int(target.Length()), mat, ops)
}

func EditScriptForMatrix(matrix [][]int, ops []EditOperation) EditScript {
	return backtrace(len(matrix[0])-1, len(matrix)-1, matrix, ops)
}
