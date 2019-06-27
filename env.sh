lib_dir='lib/'$(uname -p)''
export CGO_LDFLAGS="-Wl,-rpath,\$ORIGIN/$lib_dir -lstdc++ -L$(pwd)/$lib_dir -lhsengine -lhswrapper -lhsdecoder"
export CGO_CFLAGS="-std=gnu99 -w"
if [ ! -f "./$lib_dir/libhswrapper.so" ];then echo "./$lib_dir/libhswrapper.so not exist. Please use cmake && make to build and cp to ./$lib_dir to continue"; fi
