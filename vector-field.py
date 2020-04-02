import math
import numpy as np
import matplotlib.pyplot as plt

def avg(winds):
  r = np.array([0, 0])

  for wind in winds:
    r += wind

  dir = np.arctan2(r[1], r[0])

  speed = 0.
  for wind in winds:
    speed += math.sqrt(math.pow(wind[0], 2) + math.pow(wind[1], 2))
  speed /= len(winds)

  s = np.array([speed, 0])
  cos = math.cos(dir)
  sin = math.sin(dir)

  return [
    cos * s[0] - sin * s[1],
    sin * s[0] + cos * s[1]
  ]

x, y = np.meshgrid(np.linspace(0,4.9, 7),np.linspace(0,5.2, 7))

print(x)

u = -y/np.sqrt(x**1.2 + y**1.9)
v = x/np.sqrt(x**1.2 + y**1.9)

plt.quiver(
  x,y,u,v,
  gid="winds")


winds = []
for i in range(len(u)):
  for j in range(len(u[i])):
    if math.isnan(u[i][j]):
      continue
    winds.append(np.array([u[i][j], v[i][j]]))

w = avg(winds)
l = math.sqrt(math.pow(w[0], 2) + math.pow(w[1], 2))

X = np.array((2.85))
Y= np.array((2.5))
U = np.array((w[0]))
V = np.array((w[1]))

# fig, ax = plt.subplots()
ax = plt.axes()
q = ax.quiver(
  X, Y, U, V,
  units='xy',
  gid="avg-winds")

plt.legend(loc="upper left")
# plt.title("Winds")
plt.axis('off')
plt.savefig('plots/vector-field.svg', transparent=True)
