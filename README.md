# huffmanFileEncoder
An implementation off Huffman encoding for simple text files (in huffman.go). 

On the "longRandom.txt" file of randomly generated normal English sentences this implementation was able to compress it to 7379 bytes compared to an original size of 13541 bytes, so almost 50%.

Along the way I decided to implement the neccessary functions of a heap as well (heap.go). Originally I was going to make it a web page that let people upload files, get them compressed and sent back to them and then decompress them as well. Ultimately I decided I wouldn't learn much from this, but the file upload worked and the code is still in main.go. 

I used channels and a map reduce framework to count the occurences of different bytes because who can pass up a semi-reasonable use case for map reduce.
