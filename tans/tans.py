#!/usr/bin/python3
# date: 2022/08
# author: kozz

import math

print("Simple implementation of tANS.")
print("tANS is a table based ANS(asymmetic numeral systems) encoding algorithm.")
print("It's said ANS performs much better than the classic Huffman encoding in data compression, and it will be the replacement in future.")

class TANS:
    charStat = {}
    count = 0
    maxVal = 0
    collisionSet = set()
    lookupTable = []
    def encode(self, filepath, out, magicExtraBits = 5):
        p = {}
        with open(filepath) as f:
            while True:
                c = f.read(1)
                if not c:
                    break
                self.count += 1
                p[c] = p[c] + 1 if c in p else 1
        numPrecisionBits = math.ceil(math.log(len(p), 2)) + magicExtraBits
        self.maxVal = int(math.pow(2, numPrecisionBits)) - 1
        self.charStat = dict(map(lambda kv: (kv[0], (kv[1] / self.count, math.ceil((kv[1] * self.maxVal) / self.count))), p.items()))
        self.__info()
        self.__createLookupTable()
        return

    def __createLookupTable(self):
        p = sorted(self.charStat.items(), key=lambda kv: -kv[1][0])
        self.lookupTable = list(map(lambda l: tuple(filter(lambda n:n > 0, map(lambda x: -1 if x[1][1] < l else self.__resolveCollision(round(l/x[1][0])), p))), range(1, self.maxVal)))
        return

    def __resolveCollision(self, n):
        if n in self.collisionSet:
            return self.__resolveCollision(n + 1)
        else:
            self.collisionSet.add(n)
            return n
        
    def __info(self):
        print("original file size: " + str(self.count) + " bytes")
        print("input dict is:")
        print(self.charStat)
        print("maxVal is: " + str(self.maxVal))
        
    
    def decode(self, out):
        return

    def reset(self):
        self.p = {}
        self.count = 0
        self.maxVal = 0
