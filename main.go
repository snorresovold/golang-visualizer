package main

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type cell struct {
	switching bool
	height    float32
	start     float32
	end       float32
}

func (c *cell) makeSlice() []float32 {
	return []float32{
		c.start, c.height, 0,
		c.start, -1, 0,
		c.end, -1, 0,

		c.start, c.height, 0,
		c.end, c.height, 0,
		c.end, -1, 0,
	}
}

func (c *cell) draw(shape []float32) {
	/*
		if !c.alive {
			return
		}
	*/

	gl.BindVertexArray(makeVao(c.makeSlice()))
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(shape)/3))
}

type cellList struct {
	cells []cell
}

func (cl *cellList) InsertionSort(items []cell) []cell {
	var n = len(items)
	for i := 1; i < n; i++ {
		j := i
		for j > 0 {
			// if the thing before j
			if items[j-1].height > items[j].height {
				// switch items[j-1] and items[j]
				items[j-1].start, items[j].start = items[j].start, items[j-1].start
				items[j-1].end, items[j].end = items[j].end, items[j-1].end
			}
			// sends j one step back
			j = j - 1
		}
	}
	return items
}

func (cl *cellList) MakeCells(amount int) cellList {
	start := -1.0
	var middle float32 = 2 / float32(amount)
	for x := 1; x <= amount; x++ {
		end := start + float64(middle)
		c := cell{false, rand32(), float32(start), float32(end)}
		cl.AddCell(c)
		start = end
	}
	return *cl
}

func (cl *cellList) AddCell(c cell) []cell {
	cl.cells = append(cl.cells, c)
	return cl.cells
}

func (cl *cellList) SwitchCell(c1 cell, c2 cell) []cell {
	swapC := reflect.Swapper(cl.cells)
	swapC(0, 1)
	return cl.cells
}

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()
	cl := cellList{}
	cl.MakeCells(10)
	//fmt.Println(cl.cells)
	cl.SwitchCell(cl.cells[0], cl.cells[1])
	fmt.Println(cl)

	for !window.ShouldClose() {
		for i := 0; i < len(cl.cells); i++ {
			cl.InsertionSort(cl.cells)
			fmt.Println(cl.cells[i])
			draw(cl, window, program, cl.cells[i].makeSlice())
		}
	}
}

func rand32() float32 {
	rand.Seed(time.Now().UnixNano())
	min := -1.0
	max := 1.0
	x := min + rand.Float64()*(max-min)
	return float32(x)
}

const (
	width  = 500
	height = 500

	vertexShaderSource = `
		#version 410
		in vec3 vp;
		void main() {
			gl_Position = vec4(vp, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 410
		out vec4 frag_colour;
		void main() {
			frag_colour = vec4(1, 1, 1, 1.0);
		}
	` + "\x00"
)

func draw(cl cellList, window *glfw.Window, program uint32, shape []float32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	for _, c := range cl.cells {
		c.draw(shape)
	}

	glfw.PollEvents()
	window.SwapBuffers()
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "visualizer", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
