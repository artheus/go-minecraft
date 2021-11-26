package player

var playerVertexSource = `
#version 330 core

in vec3 pos;
in vec2 tex;
in vec3 normal;

uniform mat4 matrix;

out vec2 Tex;

void main() {
    gl_Position = matrix *  vec4(pos, 1.0);
    Tex = tex;
}
`

var playerFragmentSource = `
#version 330 core

in vec2 Tex;
uniform sampler2D tex;

out vec4 FragColor;

void main() {
    vec3 color = vec3(texture(tex, vec2(Tex.x, 1-Tex.y)));
    if (color == vec3(1,0,1)) {
        discard;
    }
    FragColor = vec4(color, 1);
}
`
