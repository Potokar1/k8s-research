package cli

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/Potokar1/k8s-research/entry5/internal/k8s"
	"github.com/spf13/cobra"
)

// NewWatchCmd creates the watch command
func NewWatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "watch is a CLI tool for watching k8s resources",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse flags
			kingdom, err := cmd.Flags().GetString("kingdom")
			if err != nil {
				return err
			}
			town, err := cmd.Flags().GetString("town")
			if err != nil {
				return err
			}

			// start watching k8s, feed into watcher
			ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
			defer cancel()

			updates, err := k8s.WatchPods(ctx, kingdom, town)
			if err != nil {
				return err
			}

			w := &watcher{}

			// fan-in updates to the watcher
			go func() {
				for update := range updates {
					w.update(update.PodName, update.Inventory)
				}
			}()

			// render every second
			tick := time.NewTicker(200 * time.Millisecond)
			defer tick.Stop()
			for {
				select {
				case <-tick.C:
					w.render()
					// w.clearDiffs()
				case <-ctx.Done():
					fmt.Println("\nExiting watch command...")
					return nil
				}
			}
		},
	}

	cmd.Flags().String("kingdom", "", "Kingdom of the town")
	cmd.Flags().String("town", "", "Name of the town")
	cmd.MarkFlagRequired("kingdom")
	cmd.RegisterFlagCompletionFunc("kingdom", KingdomsValidArgsFunction)
	cmd.RegisterFlagCompletionFunc("town", TownsValidArgsFunction)

	return cmd
}

type watcher struct {
	sync.RWMutex
	pod map[string]podHelper
}

type podHelper struct {
	inventory map[string]int // product → amount
	diff      map[string]int // product → delta since last update
	changedAt changedAt      // product → time of last change
}

func (w *watcher) update(podName string, inventory map[string]string) {
	w.Lock()
	defer w.Unlock()

	if w.pod == nil {
		w.pod = make(map[string]podHelper)
	}

	ph := w.pod[podName]
	if ph.inventory == nil {
		ph.inventory = make(map[string]int)
		ph.diff = make(map[string]int)
		ph.changedAt = make(changedAt)
	}

	now := time.Now()
	for product, amountStr := range inventory {
		newAmount, err := strconv.Atoi(amountStr)
		if err != nil {
			continue
		}
		oldAmount := ph.inventory[product]
		// only update if the amount has changed
		if newAmount != oldAmount {
			ph.diff[product] = newAmount - oldAmount
			ph.changedAt[product] = now
		}
		ph.inventory[product] = newAmount
	}
	w.pod[podName] = ph
}

func (w *watcher) render() {
	w.RLock()
	defer w.RUnlock()

	// clear screen
	fmt.Print("\033[H\033[2J")

	// sort pods for stable output
	pods := make([]string, 0, len(w.pod))
	for name := range w.pod {
		pods = append(pods, name)
	}
	sort.Strings(pods)

	// render each pod once, then its inventory
	for _, podName := range pods {
		ph := w.pod[podName]
		fmt.Println(podName)

		// sort products for stable output
		prods := make([]string, 0, len(ph.inventory))
		for prod := range ph.inventory {
			prods = append(prods, prod)
		}
		sort.Strings(prods)

		for _, prod := range prods {
			amt := ph.inventory[prod]
			delta := ph.diff[prod]
			age := time.Since(ph.changedAt[prod])
			cell := createCell(amt-delta, amt, age)
			fmt.Printf("  %-20s %s\n", prod, cell)
		}
		fmt.Println()
	}
}

// track when each field last changed
type changedAt map[string]time.Time

// uses 24-bit ANSI escapes to fade arrow+delta
func createCell(old, new int, age time.Duration) string {
	const maxFade = 5 * time.Second

	if age > maxFade {
		// if the age is too old, just show the new value
		return fmt.Sprintf("%3d   ", new)
	}

	delta := new - old
	num := fmt.Sprintf("%3d", new)
	if delta == 0 {
		return num + "   "
	}

	// pick symbol and base RGB
	sym, r, g, b := "↑", 0, 255, 0
	if delta < 0 {
		sym, r, g, b = "↓", 255, 0, 0
	}

	// compute fade fraction
	frac := min(float64(age)/float64(maxFade), 1)

	// linearly fade each channel toward 0
	r = int(float64(r) * (1 - frac))
	g = int(float64(g) * (1 - frac))
	b = int(float64(b) * (1 - frac))

	arrow := fmt.Sprintf("%s%+d", sym, delta)
	// \x1b[38;2;<r>;<g>;<b>m sets true-color fg
	return fmt.Sprintf(
		"%s \x1b[38;2;%d;%d;%dm%s\x1b[0m",
		num, r, g, b, arrow,
	)
}
