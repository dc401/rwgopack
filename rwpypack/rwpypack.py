#!/usr/bin/env python3

import zlib
import argparse
import os
import sys

XOR_KEY = 0x42  # You can change this XOR key as needed


#XOR lets us encrypt directly on byte data types
def xor_cipher(data):
    return bytes([b ^ XOR_KEY for b in data])  #pythons bitwise operations


#by nature a packer doesnt have to be encryption enabled compression will work alone
def packbin(filename):
    with open(filename, 'rb') as file_handle:
        blob = file_handle.read()
    compressed_blob = zlib.compress(blob)
    ciphered_blob = xor_cipher(compressed_blob)
    print(f'zlib compress and XOR cipher complete for: {filename}')
    return ciphered_blob


#create a skeleton wrapper for the payload that self executes the elf from payload as XOR cipher that then decrypts itself uising the same key and then decompresses itself before running with a tempfile library to self execute
def create_self_extracting_script(ciphered_data, output_filename):
    script = f"""#!/usr/bin/env python3
import zlib
import os
import sys
import tempfile

XOR_KEY = {XOR_KEY}
CIPHERED_DATA = {ciphered_data}

def xor_decipher(data):
    return bytes([b ^ XOR_KEY for b in data])

if __name__ == "__main__":
    deciphered_data = xor_decipher(CIPHERED_DATA)
    original_data = zlib.decompress(deciphered_data)

    with tempfile.NamedTemporaryFile(delete=False, mode='wb') as temp_file:
        temp_file.write(original_data)
        temp_filename = temp_file.name

    os.chmod(temp_filename, 0o755)  # Make the file executable
    os.execl(temp_filename, temp_filename, *sys.argv[1:])
"""

    with open(output_filename, 'w') as f:
        f.write(script)

    os.chmod(output_filename, 0o755)  # Make the output script executable
    print(f"Self-extracting script created: {output_filename}")
    print(f"Run it with: ./{output_filename}")


#main driver statement with dunder you can run the script directly from the command line using parsed arguments vs syaargv for positionals
if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        prog='rwpypack',
        description='Test Linux ELF binary packer in Python 3.x',
        epilog='github.com/dc401/')

    parser.add_argument('-f',
                        '--file',
                        type=str,
                        required=True,
                        help='Path to the file to be processed')
    parser.add_argument('-o',
                        '--output',
                        type=str,
                        required=True,
                        help='Output filename for the self-extracting script')

    args = parser.parse_args()

    ciphered_data = packbin(args.file)
    print(f"Ciphered data size: {len(ciphered_data)} bytes")

    create_self_extracting_script(ciphered_data, args.output)
