package scene

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	m "github.com/schmizzel/go-graphics/pkg/math"
)

// Parse a triangle mesh from an .obj file.
func ParseFromPath(path string) (*TriangleMesh, error) {
	objFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer objFile.Close()

	return ParseObj(objFile)
}

func ParseObj(objFile *os.File) (*TriangleMesh, error) {
	scanner := bufio.NewScanner(objFile)

	vertecies := make([]m.Vector3, 1, 1024)
	normals := make([]m.Vector3, 1, 1024)
	triangles := make([]*Triangle, 0, 1024)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		key := fields[0]
		values := fields[1:]

		switch key {
		case "v":
			if numbers, err := parseFloat(values); err != nil {
				return nil, err
			} else {
				vertecies = append(vertecies, m.NewVector3(numbers[0], numbers[1], numbers[2]))
			}
		case "vt":
		case "vn":
			if numbers, err := parseFloat(values); err != nil {
				return nil, err
			} else {
				normals = append(normals, m.NewVector3(numbers[0], numbers[1], numbers[2]))
			}
		case "f":
			if face, err := parseFace(values, vertecies, normals); err != nil {
				return nil, err
			} else {
				triangles = append(triangles, face...)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return NewTriangleMesh(triangles), nil
}

func parseFace(args []string, vertecies []m.Vector3, normals []m.Vector3) ([]*Triangle, error) {
	vIndeces := make([]int, 0, 4)
	nIndeces := make([]int, 0, 4)
	for _, arg := range args {
		indeces := strings.Split(arg, "/")
		vIndex, err := strconv.Atoi(indeces[0])
		if err != nil {
			return nil, err
		}
		if vIndex < 0 {
			vIndex = len(vertecies) + vIndex
		}
		vIndeces = append(vIndeces, vIndex)

		if len(indeces) >= 3 {
			nIndex, err := strconv.Atoi(indeces[2])
			if err != nil {
				return nil, err
			}
			if nIndex < 0 {
				nIndex = len(normals) + nIndex
			}
			nIndeces = append(nIndeces, nIndex)
		}
	}

	if len(nIndeces) == len(args) {
		triangles := make([]*Triangle, 0)
		for i := 1; i+2 <= len(vIndeces); i++ {
			triangles = append(triangles, triangleForIndeces(append(vIndeces[0:1], vIndeces[i:i+2]...), append(nIndeces[0:1], nIndeces[i:i+2]...), vertecies, normals))
		}
		return triangles, nil
	} else {
		triangles := make([]*Triangle, 0)
		for i := 1; i+2 <= len(vIndeces); i++ {
			triangles = append(triangles, triangleWithoutNormals(append(vIndeces[0:1], vIndeces[i:i+2]...), vertecies))
		}
		return triangles, nil
	}
}

func triangleWithoutNormals(vIndeces []int, vertecies []m.Vector3) *Triangle {
	return NewTriangleWithoutNormals(vertecies[vIndeces[0]], vertecies[vIndeces[1]], vertecies[vIndeces[2]])
}

func triangleForIndeces(vIndeces []int, nIndeces []int, vertecies []m.Vector3, normals []m.Vector3) *Triangle {
	var v [3]Vertex
	v[0] = Vertex{
		Position: vertecies[vIndeces[0]],
		Normal:   normals[nIndeces[0]],
	}
	v[1] = Vertex{
		Position: vertecies[vIndeces[1]],
		Normal:   normals[nIndeces[1]],
	}
	v[2] = Vertex{
		Position: vertecies[vIndeces[2]],
		Normal:   normals[nIndeces[2]],
	}
	return NewTriangle(v)
}

func parseFloat(args []string) ([]float64, error) {
	result := make([]float64, 0, len(args))
	for _, arg := range args {
		if num, err := strconv.ParseFloat(arg, 64); err != nil {
			return nil, err
		} else {
			result = append(result, num)
		}
	}
	return result, nil
}
