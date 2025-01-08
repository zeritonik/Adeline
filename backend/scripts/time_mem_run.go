package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Путь к исполняемому файлу программы
	// fmt.Scan(&program)
	var program string
	var ans string
	var in string
	var out bytes.Buffer
	var m int
	var t int
	inp := bufio.NewReader(os.Stdin)
	fmt.Scan(&program, &m, &t)
	in, _ = inp.ReadString('#')
	ans, _ = inp.ReadString('#')

	in = strings.Trim(in, "#")
	ans = strings.Trim(ans, "\n#")
	memLimit := m * 1024 * 1024
	timeLimit := time.Second * time.Duration(t)

	cmd := exec.Command("go", "run", program)
	cmd.Stdin = bytes.NewBuffer([]byte(in))
	cmd.Stdout = &out

	start := time.Now()
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Ошибка запуска программы: %v\n", err)
		return
	}

	time.AfterFunc(timeLimit, func() {
		cmd.Process.Kill()
	})

	pid := cmd.Process.Pid

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	var maxmem uint64
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()
	for {
		select {
		case <-ticker.C:
			memoryUsage, err := getMemoryUsage(pid)
			if err != nil {
				log.Printf("Ошибка получения информации об использовании памяти для pid %d: %v", pid, err)
			} else {
				// fmt.Printf("Текущее использование памяти для pid %d: %d КБ \n", pid, memoryUsage.VmRss)
				maxmem = max(maxmem, memoryUsage.VmRss)
				if memLimit < int(maxmem)*1024 {

					duration := time.Since(start)
					fmt.Println("ML")
					fmt.Println("program output:")
					fmt.Println(" ")
					fmt.Println("right output:")
					fmt.Println(ans)
					fmt.Println("memory:")
					fmt.Println(maxmem / 1024)
					fmt.Println("time:")
					fmt.Println(duration.Milliseconds())
					cmd.Process.Kill()
					return
				}

			}

		case err := <-done:
			if err != nil {
				if err.Error() == "signal: killed" {
					duration := time.Since(start)
					fmt.Println("TL")
					fmt.Println("program output:")
					fmt.Println(" ")
					fmt.Println("right output:")
					fmt.Println(ans)
					fmt.Println("memory:")
					fmt.Println(maxmem / 1024)
					fmt.Println("time:")
					fmt.Println(duration.Milliseconds())
				} else if err.Error() == "exit status 1" {
					duration := time.Since(start)
					fmt.Println("CE")
					fmt.Println("program output:")
					fmt.Println(" ")
					fmt.Println("right output:")
					fmt.Println(ans)
					fmt.Println("memory:")
					fmt.Println(maxmem / 1024)
					fmt.Println("time:")
					fmt.Println(duration.Milliseconds())
				}

				log.Println(err.Error())

				return
			}
			duration := time.Since(start)
			output := strings.Trim(out.String(), "\n")
			if ans == output {
				fmt.Println("OK")
			} else {
				fmt.Println("WA")
			}
			fmt.Println("program output:")
			fmt.Println(output)
			fmt.Println("right output:")
			fmt.Println(ans)
			fmt.Println("memory:")
			fmt.Println(maxmem / 1024)
			fmt.Println("time:")
			fmt.Println(duration.Milliseconds())
			return
		}
	}

}

func getMemoryUsage(pid int) (MemoryUsage, error) {
	file, err := os.Open(fmt.Sprintf("/proc/%d/status", pid))
	if err != nil {
		return MemoryUsage{}, fmt.Errorf("ошибка открытия /proc/%d/status: %v", pid, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	usage := MemoryUsage{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "VmRSS:") {
			fields := strings.Fields(line)
			if len(fields) < 2 {
				return MemoryUsage{}, fmt.Errorf("неправильная строка VmRSS: %s", line)
			}
			vmrssStr := strings.TrimSpace(fields[1])
			vmrss, err := strconv.ParseUint(vmrssStr, 10, 64)
			if err != nil {
				return MemoryUsage{}, fmt.Errorf("неправильное значение VmRSS : %v, строка :%s", err, line)
			}
			usage.VmRss = vmrss
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return MemoryUsage{}, fmt.Errorf("ошибка при чтении файла: %v", err)
	}
	return usage, nil
}

type MemoryUsage struct {
	VmRss uint64
}
