package ssdb

import (
	"bytes"
	"net"
	"fmt"
	"strconv"
)

type SSDBConnect struct {
	host string
	port int
	sock     *net.TCPConn
	recv_buf bytes.Buffer
}

func NewSSDBConnect(host string,port int,sock *net.TCPConn) *SSDBConnect{
	item:=new(SSDBConnect)
	item.host=host
	item.port=port
	item.sock=sock
	return item
}

func (this *SSDBConnect) close() error {
	if this.sock!=nil{
		return this.sock.Close()
	}

	return nil
}

func (this *SSDBConnect) do(args ...interface{}) ([]string, error) {
	err := this.send(args)
	if err != nil {
		return nil, err
	}
	resp, err := this.recv()
	return resp, err
}

func (this *SSDBConnect) set(key string, val string) error {
	resp, err := this.do("set", key, val)
	if err != nil {
		return err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return nil
	}
	return fmt.Errorf("bad response")
}

func (this *SSDBConnect) get(key string) (string, error) {
	resp, err := this.do("get", key)
	if err != nil {
		return "", err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return resp[1], nil
	}
	if resp[0] == "not_found" {
		return "", fmt.Errorf("not_found")
	}
	return "", fmt.Errorf("bad response")
}

func (this *SSDBConnect) del(key string) (interface{}, error) {
	resp, err := this.do("del", key)
	if err != nil {
		return nil, err
	}

	if len(resp) > 0 && resp[0] == "ok" {
		return true, nil
	}
	return nil, fmt.Errorf("bad response:resp:%v:", resp)
}

func (this *SSDBConnect) send(args []interface{}) error {
	var buf bytes.Buffer
	for _, arg := range args {
		var s string
		switch arg := arg.(type) {
		case string:
			s = arg
		case []byte:
			s = string(arg)
		case []string:
			for _, s := range arg {
				buf.WriteString(fmt.Sprintf("%d", len(s)))
				buf.WriteByte('\n')
				buf.WriteString(s)
				buf.WriteByte('\n')
			}
			continue
		case int:
			s = fmt.Sprintf("%d", arg)
		case int64:
			s = fmt.Sprintf("%d", arg)
		case float64:
			s = fmt.Sprintf("%f", arg)
		case bool:
			if arg {
				s = "1"
			} else {
				s = "0"
			}
		case nil:
			s = ""
		default:
			return fmt.Errorf("bad arguments")
		}
		buf.WriteString(fmt.Sprintf("%d", len(s)))
		buf.WriteByte('\n')
		buf.WriteString(s)
		buf.WriteByte('\n')
	}
	buf.WriteByte('\n')
	_, err := this.sock.Write(buf.Bytes())
	return err
}

func (this *SSDBConnect) recv() ([]string, error) {
	var tmp [8192]byte
	for {
		resp := this.parse()
		if resp == nil || len(resp) > 0 {
			return resp, nil
		}
		n, err := this.sock.Read(tmp[0:])
		if err != nil {
			return nil, err
		}
		this.recv_buf.Write(tmp[0:n])
	}
}

func (this *SSDBConnect) parse() []string {
	resp := []string{}
	buf := this.recv_buf.Bytes()
	var idx, offset int
	idx = 0
	offset = 0

	for {
		idx = bytes.IndexByte(buf[offset:], '\n')
		if idx == -1 {
			break
		}
		p := buf[offset : offset+idx]
		offset += idx + 1
		//fmt.Printf("> [%s]\n", p);
		if len(p) == 0 || (len(p) == 1 && p[0] == '\r') {
			if len(resp) == 0 {
				continue
			} else {
				var new_buf bytes.Buffer
				new_buf.Write(buf[offset:])
				this.recv_buf = new_buf
				return resp
			}
		}

		size, err := strconv.Atoi(string(p))
		if err != nil || size < 0 {
			return nil
		}
		if offset+size >= this.recv_buf.Len() {
			break
		}

		v := buf[offset : offset+size]
		resp = append(resp, string(v))
		offset += size + 1
	}

	return []string{}
}
