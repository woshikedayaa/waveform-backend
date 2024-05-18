import numpy as np
# import matplotlib.pyplot as plt
from scipy import signal
 
# 采样率
sampling_rate = 10e6  # 10Msps
 
# 时长
duration = 10  # seconds
 
# 时间序列
t = np.arange(0, duration, 1/sampling_rate)
 
# 生成正弦波信号
frequency = 1e3  # Hz
amplitude = 1
sin_wave = amplitude * np.sin(2 * np.pi * frequency * t)
 
# ADC 采样深度
adc_bits = 8
adc_levels = 2 ** adc_bits
 
# 量化
quantized_wave = np.round((sin_wave + amplitude) * (adc_levels - 1) / (2 * amplitude))
 
# 保存量化后的信号到文本文件
np.savetxt("../test/sin_wave.txt", quantized_wave, fmt='%d', delimiter=',')
 
# 绘制量化后的信号
# plt.plot(t, quantized_wave, 'r')
# plt.title('Quantized Sin Wave (8-bit ADC)')
# plt.xlabel('Time [s]')
# plt.ylabel('Amplitude')
# plt.show()
