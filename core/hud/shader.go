package hud

var (
	lineVertexSource = `
#version 330 core

in vec3 pos;

uniform mat4 matrix;

void main() {
    gl_Position = matrix *  vec4(pos, 1.0);
}
`

	lineFragmentSource = `
#version 330 core

out vec4 color;

void main() {
    color = vec4(1.0, 1.0, 1.0, 1.0);
}
`
)
