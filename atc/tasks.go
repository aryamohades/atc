package atc

import "context"

func (p *Platform) StartTasks(ctx context.Context) error {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.startAlerter(ctx)
	}()

	return nil
}
