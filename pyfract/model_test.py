import model as md
import rgba
import cfg

image = rgba.ReadImage("test/1.png")

model = md.Model(cfg.GetCFG())
ev = model.eval(image)
print(ev.points)
print(ev.len)
mfs = ev.mfs()
print(mfs.width)
print(mfs.height)
mfs_array = mfs.array
for i in range(0, mfs.width, 1):
    for j in range(0, mfs.height, 1):
        print(mfs_array[j][i])

for i in range(0, ev.len, 1):
    fd = ev.fdAt(i)
    #print(fd.len)
    print("i:", i, " fd:", fd.fd, " len:", fd.len)
    print(fd.logsizes)
    print(fd.logmeasures)
    
for i in range(0, ev.len, 1):
    um = ev.umbralAt(i)
    print("min:", um.min, " max:", um.max)
    
for i in range(0, ev.len, 1):
    loglog = ev.loglogAt(i)
    rgba.WriteImage(str(i)+"-loglog.png", loglog)