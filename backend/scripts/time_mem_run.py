import subprocess
import psutil
import time
import argparse


class Tester:
    INPUT_STOP_SYMBOL = "$$"

    def __init__(self, time_limit, memory_limit):
        self.time_limit = time_limit
        self.memory_limit = memory_limit

        self.inp = None
        self.cor_output = None
        self.err_output = None

        self.exit_code = None
        self.total_time = 0
        self.max_memory = 0
        self.output = None

        self.test_result = None


    def start(self, command):
        self.inp = self.read_from_console()
        self.cor_output = self.read_from_console()

        self.run_command(command)

        if self.total_time > self.time_limit:
            self.test_result = "TL"
        elif self.max_memory > self.memory_limit:
            self.test_result = "ML"
        elif self.exit_code != 0:
            self.test_result = "RE"
        elif not self.check_output():
            self.test_result = "WA"
        else:
            self.test_result = "OK"

            
    def check_output(self):
        return self.output == self.cor_output
    
    
    def print_test_result(self):
        print(self.test_result)
        print(self.total_time, self.max_memory)

        print("---Exit code:", self.exit_code)

        print("---Correct output:")
        print(self.cor_output)

        print("---Your output:")
        print(self.output)

        print("---Error output:")
        print(self.err_output)


    def run_command(self, command):
        process = subprocess.Popen(command, stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        process.stdin.write(self.inp.encode())
        process.stdin.close()

        start_time = time.time()
        self.exit_code = process.poll()
        while self.exit_code is None:
            time.sleep(0.01)

            self.total_time = round((time.time() - start_time) * 1000)
            self.max_memory = max(self.max_memory, psutil.Process(process.pid).memory_info().rss // 1024)
            if self.memory_limit < self.max_memory or self.time_limit < self.total_time:
                process.kill()

            self.exit_code = process.poll()

        self.output = list(map(bytes.decode, process.stdout))
        if self.output and self.output[-1] == "":self.output = self.output[:-1]
        self.output = "\n".join(map(str.strip, self.output))

        self.err_output = process.stderr.read().decode()

        process.stdout.close()
        process.stderr.close()


    def print_result(self):
        print(self.test_result)
        print(self.total_time, self.max_memory)


    @staticmethod
    def read_from_console():
        inp = []
        s = input().strip()
        while s != Tester.INPUT_STOP_SYMBOL:
            inp.append(s)
            s = input().strip()
        return "\n".join(inp)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.print_help = lambda: print(
            "Usage: python3 time_mem_run.py -t <time limit> -m <memory limit> [-v] <command>",
            "Enter input",
            Tester.INPUT_STOP_SYMBOL,
            "Enter correct output",
            Tester.INPUT_STOP_SYMBOL,
            "You will get",
            "<verdict>",
            "<total spent time(ms)> <max spent memory(kb)>",
            "---Exit code: <exit code>",
            "---Correct output:",
            "<Correct output>",
            "---Your output:",
            "<Your output>",
            "---Error output:",
            "<Error output>",
            sep="\n"
        )
    parser.add_argument('-t', '--time', type=int, required=True)
    parser.add_argument('-m', '--memory', type=int, required=True)
    parser.add_argument('-v', '--verbose')
    parser.add_argument('command', nargs=argparse.REMAINDER)
    args = parser.parse_args()

    time_limit = args.time
    memory_limit = args.memory
    verbose = args.verbose
    command = args.command

    if (verbose):
        print("Will run:", command)
        print("Time limit: " + str(time_limit))
        print("Memory limit: " + str(memory_limit))

    tester = Tester(time_limit, memory_limit)
    tester.start(command)
    tester.print_test_result()