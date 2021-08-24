#version 330

// Maximum Colors
const int colors = 256;

// Input frag attributes
in vec2 fragTexCoord;
in vec4 fragColor;

// Init uniform values
uniform sampler2D texture0;
// uniform ivec3 palette[colors];

out vec4 finalColor;

void main() {
  vec4 texelColor = texture(texture0, fragTexCoord) * fragColor;

  int index = int(texelColor.r * 255.0);
  // ivec3 color = palette[index];

  // finalColor = vec4(color / 255.0, texelColor.a);
  finalColor = vec4(1.0, 1.0, 1.0, 1.0);
}
