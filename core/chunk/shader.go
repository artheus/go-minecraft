package chunk

var (
	blockVertexSource = `
#version 330 core

in vec3 pos;
in vec2 tex;
in vec3 normal;

uniform mat4 matrix;
uniform vec3 camera;
uniform float fogdis;

out vec2 Tex;
out float diff;
out float fog_factor;

const vec3 lightdir = normalize(vec3(-1, 1, -1));

void main() {
    gl_Position = matrix *  vec4(pos, 1.0);

    float camera_distance = distance(pos, camera)/2;
    fog_factor = pow(clamp(camera_distance/fogdis, 0, 1), 4);
    Tex = tex;
    diff = max(0, dot(normal, lightdir));
}
`

	blockFragmentSource = `
#version 330 core

in vec2 Tex;
in float diff;
in float fog_factor;
uniform sampler2D tex;

out vec4 FragColor;

const vec3 sky_color = vec3(0.57, 0.71, 0.77);

void main() {
    vec3 color = vec3(texture(tex, vec2(Tex.x, 1-Tex.y)));
    if (color == vec3(1,0,1)) {
        discard;
    }
    float df = diff;
    if (color == vec3(1,1,1)) {
        df = 1- diff * 0.2;
    }
    vec3 ambient = 0.05 * vec3(1, 1, 1);
    vec3 diffcolor = df * 0.5 * vec3(1,1,1);
    color = (ambient * 8 + diffcolor) * color;
    color = mix(color, sky_color, fog_factor/2);
    FragColor = vec4(color, 1);
}
`
)