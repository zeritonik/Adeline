import sys

MOD = 128

characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}|;:,.<>/?`~"

def gen_pass(s: str):
    s = s[:MOD]
    for i in range(MOD):
        s = spinY(spinX(s * MOD * 2, i * i), i * i * i)[:MOD]
    return s
    
def spinX(s: str, n: int):
    n = n % len(s)
    return s[-n:] + s[:-n]

def spinY(s: str, n: int):
    res = ""
    for i, ch in enumerate(s):
        if ch in characters:
            res += characters[(characters.index(ch) + i + n) % len(characters)]
        else:
            res += ch
    return res


if __name__ == "__main__":
    #inp = sys.argv[1]
    inp = input()
    print(gen_pass(inp))


