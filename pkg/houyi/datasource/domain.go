package datasource

import (
	"context"
	"crypto/tls"
	"net"
	"strconv"
	"time"

	"github.com/aide-family/moon/pkg/label"
	"github.com/aide-family/moon/pkg/watch"
)

// DomainEval 函数用于获取域名的证书信息，并返回一个 map[Indexer]*Point 结构
//
//	其中 Indexer 是一个标签，用于标识一个数据点，Point 是一个数据点，包含标签和值。
func DomainEval(_ context.Context, domain string, port uint32, timeout time.Duration) (map[watch.Indexer]*Point, error) {
	now := time.Now()
	points := make(map[watch.Indexer]*Point)

	dsn := domain
	if port != 0 {
		dsn = domain + ":" + strconv.FormatUint(uint64(port), 10)
	}
	// 创建 TCP 连接
	conn, err := net.DialTimeout("tcp", dsn, timeout)
	if err != nil {
		// 超时或者连接失败，返回空切片和错误信息
		labels := label.NewLabels(map[string]string{label.Domain: domain, label.DomainPort: strconv.FormatUint(uint64(port), 10)})
		points[label.NewLabels(map[string]string{label.Domain: domain})] = &Point{
			Labels: labels.Map(),
			Values: []*Value{
				{
					Value:     0,
					Timestamp: now.Unix(),
				},
			},
		}
		return points, nil
	}
	// 函数执行完后关闭连接
	defer conn.Close()

	// 配置 TLS 的参数，ServerName 为域名，也就是我们调用函数时传递的参数
	config := &tls.Config{
		ServerName: domain,
	}

	// 创建一个 TLS 的连接
	tlsConn := tls.Client(conn, config)
	// 函数执行完后关闭连接
	defer tlsConn.Close()

	// 创建一个 TLS 的握手
	if err = tlsConn.Handshake(); err != nil {
		return nil, err
	}

	// 获取证书信息，返回的是一个切片
	certs := tlsConn.ConnectionState().PeerCertificates
	for _, cert := range certs {
		labels := label.NewLabels(map[string]string{
			label.Domain:     domain,
			label.DomainPort: strconv.FormatUint(uint64(port), 10),
		})
		points[labels] = &Point{
			Labels: labels.Map(),
			Values: []*Value{
				{
					Value:     float64(int(cert.NotAfter.Sub(now).Hours() / 24)),
					Timestamp: now.Unix(),
					Ext: map[string]any{
						label.DomainSubject:   cert.Subject.CommonName,
						label.DomainExpiresOn: cert.NotAfter.Format("2006-01-02 15:04:05"),
					},
				},
			},
		}
		break // 只取第一个证书信息
	}

	return points, nil
}
