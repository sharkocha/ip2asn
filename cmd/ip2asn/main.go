package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/netip"
	"os"
	"runtime"

	"github.com/rdleal/intervalst/interval"
)

func main() {
	// open file
	f, err := os.Open("ip2asn-combined.tsv")
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	csvReader.Comma = '\t'
	// fmt.Printf("start	end	asn	country_code	as_organization\n")

	cmpFn := func(a, b netip.Addr) int {
		return a.Compare(b)
	}
	st := interval.NewSearchTree[string](cmpFn)

	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		start := netip.MustParseAddr(rec[0])
		// if start.Compare(next) != 0 {
		// 	fmt.Printf("found gap @ %s, %s\n", start, next)
		// }
		end := netip.MustParseAddr(rec[1])

		err = st.Insert(start, end, "ASN:"+rec[2]+" Country:\""+rec[3]+"\" Info:\""+rec[4]+"\"")
		if err != nil {
			fmt.Println(err)
			continue
		}

		// start, _ := strconv.Atoi(rec[0])
		// end, _ := strconv.Atoi(rec[1])
		// do something with read line
		// startBytes, endBytes := start.As16(), end.As16()
		// class := 0
		// for i := 15; i >= 0; i-- {
		// 	if startBytes[bit/8] & endBytes[bit/8] {
		// 		class = i - 1
		// 		break
		// 	}
		// }
		// fmt.Printf("%s\t%s\t%s\t%s\t%s\n", start, end, rec[2], rec[3], rec[4])
	}
	// addr := netip.MustParseAddr("1.24.220.1")
	// vs, found := st.AllIntersections(addr, addr)
	// if found {
	// 	fmt.Println(len(vs))
	// 	fmt.Println(vs)
	// }

	file, err := os.Open(os.Args[1])

	//handle errors while opening
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}

	defer file.Close()

	fileScanner := bufio.NewScanner(file)

	// read line by line
	for fileScanner.Scan() {
		// fmt.Println(fileScanner.Text())
		if len(fileScanner.Text()) <= 0 {
			continue
		}

		addr := netip.MustParseAddr(fileScanner.Text())
		fmt.Printf("%s=", addr)
		vs, found := st.AllIntersections(addr, addr)
		if found {
			// fmt.Println(len(vs))
			fmt.Println(vs)
		} else {
			fmt.Println("not found")
		}

		// fmt.Println(st.Size())
	}
	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	// PrintMemUsage()
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// func int2ip(n int) net.IP {
// 	nn := uint32(n)
// 	ip := make(net.IP, 4)
// 	binary.BigEndian.PutUint32(ip, nn)
// 	return ip
// }

// calculateCIDRs calculates the CIDR ranges that cover the given range of IP addresses.
// func calculateCIDRs(firstIP, lastIP netip.Addr) []string {
// 	cidrs := []string{}

// 	for firstIP.Less(lastIP) {
// 		// Calculate the maximum size of the CIDR range (the prefix length).
// 		fBts, lBts := firstIP.As4(), lastIP.As4()
// 		size := 32 - bits.LeadingZeros8(uint8(fBts[3]^lBts[3]))

// 		// Ensure that the CIDR range does not extend beyond the range we're covering.
// 		for firstIP.Add((1 << size) - 1).To4().GreaterThan(lastIP.To4()) {
// 			size--
// 		}

// 		// Add the CIDR range to the list.
// 		cidrs = append(cidrs, fmt.Sprintf("%s/%d", firstIP.String(), size))

// 		// Move on to the next IP not covered by the CIDR range.
// 		firstIP = firstIP.Add(1 << size)
// 	}

// 	return cidrs
// }

func calcCIDR(a, b [16]byte) uint {
	mask := [16]byte{}
	fmt.Println(a)
	fmt.Println(b)
	for i := 0; i < 16; i++ {
		mask[i] = (a[i] ^ b[i]) & 0xff
	}
	fmt.Println(mask)
	var cidrlen uint = 0
outer:
	for i := 0; i < 16; i++ {
		for j := 0; j < 8; j++ {

			set := (mask[16-i-1] & 1 << j)
			fmt.Println(i, j, set)
			if set == 0 {
				break outer
			}
			cidrlen++
		}
		// if (mask[i/8] ^ byte(i%8)) == 1 {
		// 	break
		// }
		// cidrlen++
	}
	return cidrlen
}
