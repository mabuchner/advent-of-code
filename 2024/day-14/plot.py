import numpy as np
import matplotlib.pyplot as plt
data = np.loadtxt("./variance.txt", delimiter=",")
plt.figure()
plt.plot(range(len(data)), data)
plt.show()
