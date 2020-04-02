import numpy as np
import math
import timeit
import matplotlib.pyplot as plt

n = 10000

labels = [2**x for x in range(0, 17)]

plots = {
  'numpy': [],
  'gonum.org/v1/gonum/mat': [178,170,174,179,189,210,242,327,486,804,1383,2439,4307,8664,17502,35292,78014],
  'github.com/kvartborg/vector': [14.3,14.4,14.8,15.9,18.6,21.4,27.7,41.0,75.0,130,252,510,1523,3059,6469,14209,29127],
}

for i in range(0, 17):
  k = int(math.pow(2, i))

  result = timeit.timeit('np.array(['+', '.join(['1'] * k)+']) + np.array(['+', '.join(['1'] * k)+'])', setup='import numpy as np', number=n)
  plots['numpy'].append((result / n) * 1e+9)
  print(k, (result / n) * 1e+9)


id = 1
for (key, plot) in plots.items():
  plt.plot(labels, plot, label=key, gid='plot_'+str(id), linewidth=2)
  id += 1

plt.legend(loc="upper left")
plt.title("addition")
plt.ylabel('ns/op')
plt.xlabel('dimensions')
plt.yscale('log')
plt.xscale('log')
plt.savefig('plots/numpy-compare.svg', transparent=True)
