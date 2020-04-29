import subprocess
import re
import copy
from matplotlib import rcParams as config
import matplotlib.pyplot as plt

config['font.family'] = 'menlo'

def collect():
    data = {}
    ps = subprocess.Popen('go test -bench=AverageWind', stdout=subprocess.PIPE, shell=True)

    while True:
        line = ps.stdout.readline().rstrip()

        if line == "":
            return data

        if line[:9] != "Benchmark":
            continue

        feat, impl, label = line[9:].split('/')[:3]

        label = float(label.split('-')[0])
        value = float(re.search("(\d+|[\d\.]+) ns/op", line).group(1))

        if not feat in data.keys():
            data[feat] = {}

        if not impl in data[feat].keys():
            data[feat][impl] = { 'labels': [], 'values': [] }

        data[feat][impl]['labels'].append(label)
        data[feat][impl]['values'].append(value)

        print(feat, impl, label, value)


def plot(data):
    for feat in data.keys():
        plt.figure(figsize=(10, 5))

        for key, impl in data[feat].items():
            # colors = {
            #     'gonum': 'red',
            #     'vector_(inline)': 'green',
            #     'vector_(gonum_style)': 'blue',
            #     'vector_(immutable)': 'orange',
            # }
            plt.plot(
                impl['labels'],
                impl['values'],
                label=key.replace('_', ' '),
                # color=colors[key],
                linewidth=2)

        plt.legend(loc="upper left")
        plt.title(feat)
        plt.ylabel('ns/op')
        plt.xlabel('winds')
        # plt.yscale('log')
        # plt.xscale('log')
        plt.savefig('plots/'+feat+'.svg', transparent=True)


data = collect()
plot(data)
