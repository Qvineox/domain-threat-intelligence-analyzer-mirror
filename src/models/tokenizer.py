from keras.preprocessing.sequence import pad_sequences

validTokens = {
    '.': 1,
    '2': 2,
    'y': 3,
    'a': 4,
    '6': 5,
    'p': 6,
    't': 7,
    'w': 8,
    'x': 9,
    '_': 10,
    'u': 11,
    '1': 12,
    'f': 13,
    'l': 14,
    'i': 15,
    's': 16,
    'c': 17,
    'b': 18,
    '8': 19,
    'n': 20,
    'k': 21,
    'g': 22,
    'r': 23,
    'm': 24,
    'q': 25,
    'o': 26,
    '9': 27,
    '3': 28,
    '7': 29,
    '0': 30,
    'j': 31,
    'h': 32,
    '4': 33,
    'v': 34,
    '-': 35,
    'e': 36,
    'd': 37,
    '5': 38,
    'z': 39
}


def tokenize_domain(domain):
    d_ = [[validTokens[char] for char in domain]]
    d_ = pad_sequences(d_, maxlen=73)

    return d_
