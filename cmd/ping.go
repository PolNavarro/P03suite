/*
Copyright © 2023 NAME HERE POL NAVARRO
*/
package cmd

import (
	"fmt"
	probing "github.com/prometheus-community/pro-bing"
	"net"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping scan. Put params <ipRange> <mask>",
	Long: `Ping scan. Put params <ipRange> <mask>
	Example: 192.168.0.0 24`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ping called")
		ipRange := args[0]
		mask := args[1]

		fmt.Println("IP Range: ", ipRange)
		fmt.Println("Mask: ", mask)

		ip, ipNet, err := net.ParseCIDR(ipRange + "/" + mask)
		if err != nil {
			fmt.Println(`Error to analyze IP range:`, err)
			os.Exit(1)
		}

		fmt.Printf("Scaning IP range %s/%s:\n", ip, mask)

		var wg sync.WaitGroup

		for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
			target := ip.String()
			wg.Add(1)
			go performPing(target, &wg)
		}

		wg.Wait()

	},
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func performPing(target string, wg *sync.WaitGroup) {
	defer wg.Done()

	pinger, err := probing.NewPinger(target)

	if err != nil {
		fmt.Println("Error to create pinger", target, err)
		return
	}

	pinger.Timeout = 1 * time.Second

	err = pinger.Run()
	if err != nil {
		fmt.Printf("FAIL PINGER %s: %s\n", target, err)
		return
	}

	stats := pinger.Statistics()

	if stats.PacketsRecv == 0 {
		fmt.Printf("NONE PINGER %s\n", target)
	} else {
		fmt.Printf("TRUE PINGER %s\n", target)
	}

	//fmt.Printf("%s - Paquetes Sent: %d, Paquetes Recibidos: %d\n", target, stats.PacketsSent, stats.PacketsRecv)

}

func init() {
	rootCmd.AddCommand(pingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
