package eventloop

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"noelzubin/redis-go/protocol"
	"noelzubin/redis-go/store"
	"noelzubin/redis-go/utils"
	"strconv"
	"strings"
	"time"
)

var OK = "OK"

type Eventloop struct {
	reqChan chan ReqCommand
	store   store.Store
}

func InitEventloop(store store.Store) *Eventloop {
	return &Eventloop{
		store:   store,
		reqChan: make(chan ReqCommand),
	}
}

func (e *Eventloop) RunLoop() {
	fmt.Println("loop started")
	for cmd := range e.reqChan {
		var resp protocol.Value
		switch strings.ToLower(cmd.command[0]) {
		case "ping":
			r := e.store.Ping()
			resp = protocol.NewSimpleStringValue(r)
		case "set":

			if len(cmd.command) < 3 {
				resp = protocol.NewErrorValue("ERR wrong number of arguments for 'set' command")
				break
			}

			var exp *time.Time = nil
			if len(cmd.command) > 3 {
				if seconds, err := strconv.Atoi(cmd.command[3]); err == nil {
					expiry := time.Now().Add(time.Duration(seconds) * time.Second)
					exp = &expiry
				}
			}

			e.store.Set(cmd.command[1], cmd.command[2], exp)
			resp = protocol.NewSimpleStringValue(&OK)

		case "get":
			if (len(cmd.command)) < 2 {
				resp = protocol.NewErrorValue("ERR wrong number of arguments for 'set' command")
				break
			}
			r := e.store.Get(cmd.command[1])
			if r == nil {
				resp = protocol.NewNilValue()
				break
			}
			resp = protocol.NewSimpleStringValue(r)
		case "del":
			if len(cmd.command) < 2 {
				resp = protocol.NewErrorValue("ERR wrong number of arguments for 'del' command")
				break
			}

			r := e.store.Del(cmd.command[1:]...)
			resp = protocol.NewSimpleIntValue(int64(r))
		case "expire":
			if len(cmd.command) < 3 {
				resp = protocol.NewErrorValue("ERR wrong number of arguments for 'del' command")
				break
			}

			seconds, err := strconv.Atoi(cmd.command[2])
			if err != nil {
				resp = protocol.NewErrorValue("ERR value is not an integer or out of range")
				break
			}

			r := e.store.Expire(cmd.command[1], seconds)
			resp = protocol.NewSimpleIntValue(int64(r))
		case "keys":
			r := e.store.Keys("")
			resp = protocol.NewArrayStringValue(r)
		case "zadd":

			if (len(cmd.command))%2 != 0 || len(cmd.command) < 4 {
				resp = protocol.NewErrorValue("ERR wrong number of arguments for 'zadd' command")
				break
			}
			scoreMembers, err := utils.GetScoreMemberPairs(cmd.command[2:])
			if err != nil {
				resp = protocol.NewErrorValue("ERR syntax error")
				break
			}
			r := e.store.ZAdd(cmd.command[1], scoreMembers)
			resp = protocol.NewSimpleIntValue(int64(r))
		case "zrange":
			if len(cmd.command) < 4 {
				resp = protocol.NewErrorValue("ERR wrong number of arguments for 'zrange' command")
				break
			}

			fmt.Println(cmd.command)

			start, err := strconv.Atoi(cmd.command[2])
			if err != nil {
				resp = protocol.NewErrorValue("ERR value is not an integer or out of range")
				break
			}

			end, err := strconv.Atoi(cmd.command[3])
			if err != nil {
				resp = protocol.NewErrorValue("ERR value is not an integer or out of range")
				break
			}

			withScores := false
			for _, v := range cmd.command {
				if strings.ToLower(v) == "withscores" {
					withScores = true
					break
				}
			}

			r := e.store.ZRange(cmd.command[1], start, end, withScores)
			resp = protocol.NewArrayStringValue(r)
		default:
			resp = protocol.NewErrorValue("unknown command '" + cmd.command[0] + "'")
		}
		cmd.respChan <- resp
	}
}

type ReqCommand struct {
	command  []string
	respChan chan protocol.Value
}

func (e *Eventloop) HandleConnection(conn io.ReadWriteCloser) {
	defer conn.Close()

	for {
		value, err := protocol.DecodeRESP(bufio.NewReader(conn))
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Println("error decoding RESP: ", err.Error())
			return
		}

		respChan := make(chan protocol.Value)

		valuesArr := value.Array()
		strValues := make([]string, 0)
		for _, v := range valuesArr {
			strValues = append(strValues, v.String())
		}

		reqCmd := ReqCommand{
			command:  strValues,
			respChan: respChan,
		}

		e.reqChan <- reqCmd

		cmdRes := <-respChan
		conn.Write(cmdRes.Encode())
		fmt.Println("wrote to connetion", string(cmdRes.Encode()))
	}
}
