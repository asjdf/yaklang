package synscanx

import (
	"github.com/google/gopacket/layers"
	"github.com/yaklang/yaklang/common/fp"
	"github.com/yaklang/yaklang/common/pcapx"
	"github.com/yaklang/yaklang/common/utils"
	"math/rand"
	"net"
	"runtime"
)

func (s *Scannerx) assembleSynPacket(host string, port int) ([]byte, error) {
	isLoopback := utils.IsLoopback(host)

	var packetBytes []byte
	var err error
	var opts []any

	dstMac := s.config.RemoteMac
	srcMac := s.config.SourceMac
	// 内网扫描时，这一步应该能够获取到目标机器的 MAC 地址
	if mac, ok := s.macCacheTable.Load(host); ok {
		dstMac = mac.(net.HardwareAddr)
	}

	if dstMac == nil {
		if !isLoopback {
			// 外网扫描时，目标机器的 MAC 地址就是网关的 MAC 地址
			dstMac, err = s.getGatewayMac()
			if err != nil {
				return nil, utils.Errorf("get gateway mac failed: %s", err)
			}
			// Ethernet
			opts = append(opts, pcapx.WithEthernet_SrcMac(srcMac))
			opts = append(opts, pcapx.WithEthernet_DstMac(dstMac))
		} else {
			// Loopback
			if runtime.GOOS == "windows" {
				opts = append(opts, pcapx.WithLoopback(isLoopback))
			} else {
				opts = append(opts, pcapx.WithEthernet_SrcMac(net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}))
				opts = append(opts, pcapx.WithEthernet_DstMac(net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}))
			}
		}
	} else {
		opts = append(opts,
			pcapx.WithEthernet_SrcMac(srcMac),
			pcapx.WithEthernet_DstMac(dstMac),
		)
	}

	var ipSrc string
	if isLoopback {
		//ipSrc = net.ParseIP("127.0.0.1").String()
		ipSrc = net.ParseIP("192.168.3.3").String()
		host = ipSrc
	} else {
		ipSrc = s.config.SourceIP.String()
	}
	srcPort := rand.Intn(65534) + 1
	// wireshark filter port
	//srcPort := 52873
	// IPv4
	opts = append(opts, pcapx.WithIPv4_Flags(layers.IPv4DontFragment))
	opts = append(opts, pcapx.WithIPv4_Version(4))
	opts = append(opts, pcapx.WithIPv4_NextProtocol(layers.IPProtocolTCP))
	opts = append(opts, pcapx.WithIPv4_TTL(64))
	opts = append(opts, pcapx.WithIPv4_ID(40000+rand.Intn(10000)))
	opts = append(opts, pcapx.WithIPv4_SrcIP(ipSrc))
	opts = append(opts, pcapx.WithIPv4_DstIP(host))
	opts = append(opts, pcapx.WithIPv4_Option(nil, nil))

	// TCP
	opts = append(opts,
		pcapx.WithTCP_SrcPort(srcPort),
		pcapx.WithTCP_DstPort(port),
		pcapx.WithTCP_Flags(pcapx.TCP_FLAG_SYN),
		pcapx.WithTCP_Window(1024),
		pcapx.WithTCP_Seq(500000+rand.Intn(10000)),
	)

	packetBytes, err = pcapx.PacketBuilder(opts...)
	if err != nil {
		return nil, utils.Wrapf(err, "assembleSynPacket failed")
	}
	//log.Infof("assembleSynPacket: %s", hex.EncodeToString(packetBytes))
	//0200000045000030baa640004006821f7f0000017f000001206900160007b64f0000000070020400a5470000020405b403030700
	return packetBytes, nil
}

func (s *Scannerx) assembleUdpPacket(host string, port int) ([]byte, error) {
	isLoopback := utils.IsLoopback(host)

	var packetBytes []byte
	var err error
	var opts []any

	dstMac := s.config.RemoteMac
	srcMac := s.config.SourceMac
	if mac, ok := s.macCacheTable.Load(host); ok {
		dstMac = mac.(net.HardwareAddr)
	}

	if dstMac == nil {
		if !isLoopback {
			dstMac, err = s.getGatewayMac()
			if err != nil {
				return nil, utils.Errorf("get gateway mac failed: %s", err)
			}
			// Ethernet
			opts = append(opts, pcapx.WithEthernet_SrcMac(srcMac))
			opts = append(opts, pcapx.WithEthernet_DstMac(dstMac))
		} else {
			// Loopback
			opts = append(opts, pcapx.WithLoopback(isLoopback))
		}
	} else {
		opts = append(opts,
			pcapx.WithEthernet_SrcMac(srcMac),
			pcapx.WithEthernet_DstMac(dstMac),
		)
	}

	var ipSrc string
	if isLoopback {
		ipSrc = net.ParseIP("127.0.0.1").String()
		host = ipSrc
	} else {
		ipSrc = s.config.SourceIP.String()
	}
	srcPort := rand.Intn(65534) + 1
	// wireshark filter port
	//srcPort := 52873

	// IPv4
	opts = append(opts, pcapx.WithIPv4_Flags(layers.IPv4DontFragment))
	opts = append(opts, pcapx.WithIPv4_Version(4))
	opts = append(opts, pcapx.WithIPv4_NextProtocol(layers.IPProtocolUDP))
	opts = append(opts, pcapx.WithIPv4_TTL(64))
	opts = append(opts, pcapx.WithIPv4_ID(40000+rand.Intn(10000)))
	opts = append(opts, pcapx.WithIPv4_SrcIP(ipSrc))
	opts = append(opts, pcapx.WithIPv4_DstIP(host))
	opts = append(opts, pcapx.WithIPv4_Option(nil, nil))

	// UDP
	opts = append(opts, pcapx.WithUDP_SrcPort(srcPort))
	opts = append(opts, pcapx.WithUDP_DstPort(port))
	var payload []byte
	nmapRuleConfig := fp.NewConfig(
		fp.WithActiveMode(true),
		fp.WithTransportProtos(fp.UDP),
		fp.WithProbesMax(3),
	)
	_, blocks, bestMode := fp.GetRuleBlockByConfig(port, nmapRuleConfig)
	if len(blocks) > 0 && bestMode {
		payload = []byte(blocks[0].Probe.Payload)
	}
	//payload := getUDPPayloadByPort(port)

	opts = append(opts, pcapx.WithPayload(payload))

	packetBytes, err = pcapx.PacketBuilder(opts...)
	if err != nil {
		return nil, utils.Wrapf(err, "assembleUdpPacket failed")
	}
	return packetBytes, nil
}

func (s *Scannerx) assembleArpPacket(host string) ([]byte, error) {
	var opts []any
	srcMac := s.config.SourceMac.String()
	srcIP := s.config.SourceIP.String()
	opts = append(opts, pcapx.WithEthernet_SrcMac(srcMac))
	opts = append(opts, pcapx.WithEthernet_DstMac(net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}))
	opts = append(opts, pcapx.WithEthernet_NextLayerType(layers.EthernetTypeARP))

	opts = append(opts, pcapx.WithArp_RequestEx(srcIP, srcMac, host))

	packetBytes, err := pcapx.PacketBuilder(opts...)
	if err != nil {
		return nil, err
	}
	return packetBytes, nil
}
