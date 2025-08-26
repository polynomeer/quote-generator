package generator

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/polynomeer/quote-generator/internal/config"
	"github.com/polynomeer/quote-generator/internal/model"
)

type Options struct {
	Seed    int64
	Symbols []string
	Hz      float64 // >= 1.0
}

type Generator struct {
	opts  Options
	rng   *rand.Rand
	mu    sync.RWMutex
	price map[string]float64
	bid   map[string]float64
	ask   map[string]float64
	vol   map[string]int64
	quit  chan struct{}
}

func New(opts Options) *Generator {
	if opts.Hz <= 0 {
		opts.Hz = 1
	}
	g := &Generator{
		opts:  opts,
		rng:   rand.New(rand.NewSource(opts.Seed)),
		price: make(map[string]float64, len(opts.Symbols)),
		bid:   make(map[string]float64, len(opts.Symbols)),
		ask:   make(map[string]float64, len(opts.Symbols)),
		vol:   make(map[string]int64, len(opts.Symbols)),
		quit:  make(chan struct{}),
	}
	// 기본가 초기화
	for _, s := range opts.Symbols {
		base := 100 + g.rng.Float64()*100 // 100~200
		g.price[s] = base
		g.bid[s] = base - 0.01
		g.ask[s] = base + 0.01
		g.vol[s] = 0
	}
	return g
}

func (g *Generator) Start() {
	interval := time.Duration(float64(time.Second) / g.opts.Hz)
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-g.quit:
				return
			case <-ticker.C:
				g.stepAll()
			}
		}
	}()
}

func (g *Generator) Stop() { close(g.quit) }

func (g *Generator) stepAll() {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, s := range g.opts.Symbols {
		p := g.price[s]
		// 간단 GBM 스텁: μ=0, σ=0.2 연/초 환산 간략화
		dt := 1.0 / 252.0 / (6.5 * 3600) // 거래시간 근사(초당)
		sigma := 0.2
		z := g.rng.NormFloat64()
		pnext := p * math.Exp((-0.5*sigma*sigma)*dt+sigma*math.Sqrt(dt)*z)
		spread := 0.01 + 0.02*(g.rng.Float64()-0.5)
		g.price[s] = pnext
		g.bid[s] = pnext - spread/2
		g.ask[s] = pnext + spread/2
		g.vol[s] += int64(10 + g.rng.Intn(50))
	}
}

// Snapshot: 최신 틱을 가져간다.
func (g *Generator) Snapshot(symbols []string) []model.Quote {
	g.mu.RLock()
	defer g.mu.RUnlock()
	res := make([]model.Quote, 0, len(symbols))
	now := config.NowMillis()
	for _, s := range symbols {
		p := g.price[s]
		if p == 0 {
			continue
		}
		res = append(res, model.Quote{
			Symbol: s, Price: p, Bid: g.bid[s], Ask: g.ask[s],
			Volume: g.vol[s], Currency: "USD", Ts: now,
		})
	}
	return res
}
