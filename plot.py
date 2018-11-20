import matplotlib.pyplot as plt
import numpy as np
from IPython import embed

csv = np.loadtxt('result.csv', delimiter=',')

_, ax1 = plt.subplots()
ax2 = ax1.twinx()

ax1.plot(csv[:, 0], 'b')
ax1.plot(csv[:, 2], 'r') # Target
ax2.plot(csv[:, 1], 'g')

ax1.set_yticks(np.linspace(0, 100, 5))
ax1.set_ylim(-0.6, 100.6)
ax2.set_yticks(np.linspace(0, 60, 5))
ax2.set_ylim(-0.6, 60.6)

ax1.set_ylabel('CPU Usage',   color='b')
ax2.set_ylabel('Num Workers', color='g')
# plt.show()
plt.savefig("result.png")
