#version 330

// Maximum Colors
const int colors = 256;

// Input frag attributes
in vec2 fragTexCoord;
in vec4 fragColor;

// Init uniform values
uniform sampler2D texture0;
uniform sampler2D palette;

out vec4 finalColor;

void main() {
  vec4 index = texture(texture0, fragTexCoord);
  finalColor = texture(palette, index.xy);
}
