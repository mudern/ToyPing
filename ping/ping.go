package ping

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	ICMPTypeEchoRequest = 8 // ICMP Echo Request 类型
	ICMPTypeEchoReply   = 0 // ICMP Echo Reply 类型
)

// ICMP 结构体定义
type ICMP struct {
	Type     uint8  // ICMP 类型
	Code     uint8  // ICMP 代码
	Checksum uint16 // 校验和
	ID       uint16 // 请求 ID
	Seq      uint16 // 序列号
}

// Serialize 方法将 ICMP 结构体序列化为字节流
func (icmp *ICMP) Serialize() []byte {
	buf := make([]byte, 8)
	buf[0] = icmp.Type
	buf[1] = icmp.Code
	binary.BigEndian.PutUint16(buf[2:4], icmp.Checksum)
	binary.BigEndian.PutUint16(buf[4:6], icmp.ID)
	binary.BigEndian.PutUint16(buf[6:8], icmp.Seq)
	return buf
}

// Deserialize 方法将字节流反序列化为 ICMP 结构体
func (icmp *ICMP) Deserialize(data []byte) {
	icmp.Type = data[0]
	icmp.Code = data[1]
	icmp.Checksum = binary.BigEndian.Uint16(data[2:4])
	icmp.ID = binary.BigEndian.Uint16(data[4:6])
	icmp.Seq = binary.BigEndian.Uint16(data[6:8])
}

// CalculateChecksum 计算 ICMP 校验和
func CalculateChecksum(data []byte) uint16 {
	var sum uint32
	for i := 0; i < len(data); i += 2 {
		if i+1 < len(data) {
			sum += uint32(data[i])<<8 + uint32(data[i+1])
		} else {
			sum += uint32(data[i]) << 8
		}
	}
	for sum>>16 > 0 {
		sum = (sum & 0xffff) + (sum >> 16)
	}
	return uint16(^sum)
}

// SendPing 发送 Ping 请求
func SendPing(address string, count int) {
	conn, err := net.Dial("ip4:icmp", address)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	fmt.Printf("正在 Ping %s 具有 32 字节的数据:\n", address)

	var totalSent, totalReceived int
	var minRTT, maxRTT, totalRTT time.Duration

	for i := 0; i < count; i++ {
		icmp := ICMP{
			Type: ICMPTypeEchoRequest,
			Code: 0,
			ID:   uint16(os.Getpid() & 0xffff),
			Seq:  uint16(i + 1),
		}
		icmpData := icmp.Serialize()
		icmp.Checksum = CalculateChecksum(icmpData)
		icmpData = icmp.Serialize()

		start := time.Now()
		_, err := conn.Write(icmpData)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		totalSent++

		buf := make([]byte, 1500)
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("请求超时。")
			continue
		}

		ipHeader := buf[:20] // IP 头部的长度是 20 字节
		ttl := ipHeader[8]   // TTL 字段在 IP 头部的第 9 字节

		reply := ICMP{}
		reply.Deserialize(buf[20:n]) // 跳过 IP 头部
		if reply.Type == ICMPTypeEchoReply && reply.ID == icmp.ID {
			rtt := time.Since(start)
			totalReceived++
			if totalReceived == 1 || rtt < minRTT {
				minRTT = rtt
			}
			if rtt > maxRTT {
				maxRTT = rtt
			}
			totalRTT += rtt

			fmt.Printf("来自 %s 的回复: 字节=32 时间=%vms TTL=%d\n", address, rtt.Milliseconds(), ttl)
		} else {
			fmt.Println("收到无效的回复。")
		}

		time.Sleep(1 * time.Second)
	}

	fmt.Printf("\n%s 的 Ping 统计信息:\n", address)
	fmt.Printf("    数据包: 已发送 = %d，已接收 = %d，丢失 = %d (%.1f%% 丢失)，\n",
		totalSent, totalReceived, totalSent-totalReceived, float64(totalSent-totalReceived)/float64(totalSent)*100)
	if totalReceived > 0 {
		avgRTT := totalRTT / time.Duration(totalReceived)
		fmt.Printf("往返行程的估计时间(以毫秒为单位):\n")
		fmt.Printf("    最短 = %vms，最长 = %vms，平均 = %vms\n", minRTT.Milliseconds(), maxRTT.Milliseconds(), avgRTT.Milliseconds())
	}
}
