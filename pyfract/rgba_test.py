import rgba

image = rgba.ReadImage("test/1.png")
print(image.width)
print(image.height)
print(image.stride)
print(image.len)
print(image.buffer[:30])
rgba.WriteImage("test/2.png", image)
np_array = image.as_np
image = rgba.NewImage(np_array)
rgba.WriteImage("test/3.png", image)