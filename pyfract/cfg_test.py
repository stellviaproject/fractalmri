import cfg

cf = cfg.GetCFG()
# cf.knn.k = 20
# cf.knn.distance = cfg.DIST_MANHATTAN
# cf.knn.selector = cfg.SELECTOR_SMOOTH_INVERSE
# cf.knn.minkowski_ratio = 1000
# cf.knn.weight_param = 90
# cf.knn.smoothing_param = 102
# print(cf.knn.k)
# print(cf.knn.distance)
# print(cf.knn.selector)
# print(cf.knn.minkowski_ratio)
# print(cf.knn.weight_param)
# print(cf.knn.smoothing_param)

# cf.buffer = 100
# print(cf.buffer)
# cf.parallel = 90
# print(cf.parallel)
# cf.window_ratio = 110
# print(cf.window_ratio)
# cf.denoiser_sigma_color = 19.1
# print(cf.denoiser_sigma_color)
# cf.denoiser_sigma_space = 36.5
# print(cf.denoiser_sigma_space)
# cf.denoiser_diameter = 100
# print(cf.denoiser_diameter)
# cf.denoiser_umbral_color = 99.2
# print(cf.denoiser_umbral_color)
# cf.min_umbral = 1.1
# print(cf.min_umbral)
# cf.max_umbral = 1.2
# print(cf.max_umbral)
# cf.min_area = 89
# print(cf.min_area)
# cf.max_area = 98
# print(cf.max_area)
# cf.ratio = 4
# print(cf.ratio)
print(cf.box_sizes)
cf.box_sizes = [3, 5, 7]
print(cf.box_sizes)
print(cf.umbral)
cf.umbral = [0.1, 0.5, 1.5]
print(cf.umbral)
# print(cf.box_sizes)
# cf.loglog_height = 500
# print(cf.loglog_height)
# cf.loglog_width = 400
# print(cf.loglog_width)
# cf.umbral = [0.1, 0.3, 0.5]
# print(cf.umbral)
#b = cf.buffer
#print(b)