package seed

import (
	"fmt"
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/houyi-tracing/houyi/pkg/routing"
	"github.com/houyi-tracing/houyi/ports"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

// ATTENTION: all tests below must run after starting registry!

const (
	configServerIp   = "192.168.31.77"
	configServerPort = ports.ConfigServerGrpcListenPort
)

var (
	registryEndpoint = &routing.Endpoint{
		Addr: configServerIp,
		Port: configServerPort,
	}
)

func TestSeedCanConnectToRegistry(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	s := NewSeed(logger,
		Options.ConfigServerEndpoint(&routing.Endpoint{
			Addr: "192.168.31.77",
			Port: ports.ConfigServerGrpcListenPort,
		}),
		Options.ListenPort(22879),
	)
	_ = s.Start()
	time.Sleep(time.Minute)
	_ = s.Stop()
}

func TestSeedNodeIdMustBeUnique(t *testing.T) {
	logger, _ := zap.NewProduction()
	n := 1024
	set := make(map[int]bool)
	for i := 0; i < n; i++ {
		go func(id int) {
			s := NewSeed(logger,
				Options.ConfigServerEndpoint(registryEndpoint),
				Options.ListenPort(28391+id))
			nodeId := s.(*seed).nodeId
			assert.GreaterOrEqual(t, nodeId, int64(0))
			_, has := set[nodeId]
			assert.False(t, has)
			set[nodeId] = true
		}(i)
	}
}

func TestMessageMongering(t *testing.T) {
	logger, _ := zap.NewProduction()
	runtime.GOMAXPROCS(16)

	lock := &sync.RWMutex{}

	// params
	seeds := 64
	maxWaitTime := time.Second * 5

	msgReceivedCnt := make(map[int]map[int]int)

	for i := 0; i < seeds; i++ {
		go func(id int, msgCnt map[int]map[int]int, l *sync.RWMutex) {
			relHandler := func(rel *api_v1.Relation) {
				l.Lock()
				if _, has := msgCnt[id]; !has {
					msgCnt[id] = make(map[int]int)
				}
				msgCnt[id][0]++
				l.Unlock()
			}

			tagsHandler := func(tags *api_v1.EvaluatingTags) {
				l.Lock()
				if _, has := msgCnt[id]; !has {
					msgCnt[id] = make(map[int]int)
				}
				msgCnt[id][1]++
				l.Unlock()
			}

			opHandler := func(op *api_v1.Operation) {
				l.Lock()
				if _, has := msgCnt[id]; !has {
					msgCnt[id] = make(map[int]int)
				}
				msgCnt[id][2]++
				l.Unlock()
			}

			s := NewSeed(logger,
				Options.ConfigServerEndpoint(registryEndpoint),
				Options.ListenPort(28991+id),
				Options.OnNewRelation(relHandler))
			s.OnEvaluatingTags(tagsHandler)
			s.OnExpiredOperation(opHandler)
			_ = s.Start()
		}(i, msgReceivedCnt, lock)
	}

	s := NewSeed(logger,
		Options.ConfigServerEndpoint(registryEndpoint),
		Options.ListenPort(28000))
	_ = s.Start()
	time.Sleep(time.Second * 5)

	startTime := time.Now()
	s.MongerNewRelation(&api_v1.Relation{
		From: &api_v1.Operation{
			Service:   "svcA",
			Operation: "op1",
		},
		To: &api_v1.Operation{
			Service:   "svcB",
			Operation: "op7",
		},
	})
	s.MongerEvaluatingTags(&api_v1.EvaluatingTags{
		Tags: []*api_v1.EvaluatingTag{
			{
				TagName:       "",
				OperationType: 0,
				ValueType:     0,
				Value:         nil,
			},
		}})
	s.MongerExpiredOperation(&api_v1.Operation{
		Service:   "s",
		Operation: "o",
	})

	isConv := false
	for {
		time.Sleep(time.Second)
		lock.RLock()

		receivedAll := 0
		for _, perTypeCnt := range msgReceivedCnt {
			cnt := 0
			for _, received := range perTypeCnt {
				if received >= 1 {
					cnt++
				}
			}
			if cnt == 3 {
				receivedAll++
			}
		}

		if !isConv && receivedAll == seeds {
			isConv = true
			fmt.Println("Reaching convergence of message mongering consume time: ", time.Now().Sub(startTime).String())
		}
		if startTime.Add(maxWaitTime).Before(time.Now()) {
			break
		}
		lock.RUnlock()
	}

	totalMsg := 0
	for _, i := range msgReceivedCnt {
		for _, j := range i {
			totalMsg += j
		}
	}
	if !isConv {
		fmt.Println("Did not converge.")
	}
	fmt.Println("Maximum wait time:", maxWaitTime.String())
	fmt.Println("Convergence message msgReceivedCnt:", seeds*3)
	fmt.Println("Actual message msgReceivedCnt: ", totalMsg)
	assert.GreaterOrEqual(t, totalMsg, seeds*3)
}

func TestMongerNewRelation(t *testing.T) {
	logger, _ := zap.NewProduction()
	runtime.GOMAXPROCS(16)

	lock := &sync.RWMutex{}

	// params
	nSeed := 80
	maxWaitTime := time.Second * 5
	received := make(map[int]int)

	for i := 0; i < nSeed; i++ {
		go func(id int, msgCnt map[int]int, l *sync.RWMutex) {
			f := func(_ *api_v1.Relation) {
				l.Lock()
				defer l.Unlock()
				msgCnt[id]++
			}
			s := NewSeed(logger,
				Options.ConfigServerEndpoint(registryEndpoint),
				Options.ListenPort(28991+id),
				Options.OnNewRelation(f))
			_ = s.Start()
		}(i, received, lock)
	}

	s := NewSeed(logger,
		Options.ConfigServerEndpoint(registryEndpoint),
		Options.ListenPort(28990))
	_ = s.Start()

	// wait goroutines to get ready
	time.Sleep(time.Second * 3)

	s.MongerNewRelation(&api_v1.Relation{
		From: &api_v1.Operation{
			Service:   "svcA",
			Operation: "opA",
		},
		To: &api_v1.Operation{
			Service:   "svcB",
			Operation: "opB",
		},
	})

	startTime := time.Now()
	flag := false
	for {
		time.Sleep(time.Millisecond)
		if startTime.Add(maxWaitTime).Before(time.Now()) {
			break
		}

		if !flag && len(received) == nSeed {
			flag = true
			fmt.Println("converge: ", time.Now().Sub(startTime).String())
		}
	}

	total := 0
	for _, msgCnt := range received {
		total += msgCnt
	}
	fmt.Println("total message:", total)
	assert.GreaterOrEqual(t, total, nSeed)
}

func TestRemoveDeadSeeds(t *testing.T) {
	seeds := 100
	logger, _ := zap.NewProduction()

	for i := 0; i < seeds; i++ {
		s := NewSeed(logger,
			Options.ConfigServerEndpoint(registryEndpoint),
			Options.ListenPort(28391+i))
		_ = s.Start()
	}
	time.Sleep(time.Hour)
}

func TestUniqueMessageID(t *testing.T) {
	logger, _ := zap.NewProduction()
	s := NewSeed(logger,
		Options.ConfigServerEndpoint(registryEndpoint),
		Options.ListenPort(28000))
	se := s.(*seed)
	unique := make(map[int64]bool)
	msgN := 10000000
	for i := 0; i < msgN; i++ {
		unique[se.msgIdGenerator.Generate().Int64()] = true
	}
	assert.Equal(t, msgN, len(unique))
}

func BenchmarkMsgIdGenerating(b *testing.B) {
	logger, _ := zap.NewProduction()
	s := NewSeed(logger,
		Options.ConfigServerEndpoint(registryEndpoint),
		Options.ListenPort(28000))
	se := s.(*seed)
	for i := 0; i < b.N; i++ {
		se.msgIdGenerator.Generate().Int64()
	}
}

func BenchmarkRandom(b *testing.B) {
	prob := 0.2
	cnt := 0.0

	for i := 0; i < b.N; i++ {
		rand.Seed(time.Now().UnixNano())
		if rand.Float64() < prob {
			cnt += 1.0
		}
	}

	fmt.Println(cnt / float64(b.N))
}

func TestGetIp(t *testing.T) {
	ip, err := getClientIp()
	fmt.Println(ip, err)
}
