# rwgopack
Example Linux based packer for ELF binaries that uses ZLib to compress and then XOR cipher single byte key the payload while creating a self unpacking binary. In the example code you can use an GCC compiled edition of a C hello world source and have that executed in a "packed" version as a wrapper using the above mechanisms.

## Why another packer?
I wrote one just because I wanted to figure out what the best mechanism of self-execution would be to prepend the executable after which is the harder part. We can take inspiration from my last obfuscation pet project generating [python payload](https://github.com/dc401/py-obfuscation-payloadgen). The name is the abbreviated for "RW for [Robin Williams](https://www.imdb.com/name/nm0000245/)" the late actor because it's just clowning around on this project (his role in [Patch Adams](https://www.imdb.com/title/tt0129290/))

## Build Summary
Recreated similar template of my "py-obfuscation-payloadgen" project by letting it create a skeleton. Add in the requirement for compression and encryption with a static key to produce a Python edition payload (again). Then I cheated, and used Claude 3.5 Gen AI to port the Python 3.x to a Go 1.18 compatible equivalent because I hate Go syntax so that's why I used Claude. Go was used because it can operate as a script and as a compiled version of itself. I modified it to use a standard shell sub process call to "build" itself and then pump out the binary to disk.

## Python PoC Demo
Python code written for Claude to interpret and translate into Go stuff later:
![enter image description here](https://github.com/dc401/rwgopack/blob/main/rwpypack-demo-replit.gif?raw=true)

## Go Portable Edition
Claude 3.5 refactored version and the only difference really just subprocess calling itself to compile a binary instead of just a script-wrapper. Then I check the entropy level before and after:
![enter image description here](https://github.com/dc401/rwgopack/blob/main/rwgopack-replit-demo.gif?raw=true)

