package piper

type Piper struct {
}

func (p *Piper) Run(batchSize int, loader Loader, filter Filter, mutator Mutator, resultHandler ResultHandler) error {
	shallContinue := true
	for shallContinue {
		data, eof, err := loader.LoadBatch(batchSize)

		if err != nil {
			return err
		}

		data = filterFn(data, func(x interface{}) bool {
			return filter.Filter(x)
		})
		data, err = mutator.Mutate(data)

		if err != nil {
			return err
		}

		resultHandler(toMultiResult(data), eof)

		if eof {
			shallContinue = false
		}
	}

	return nil
}

func filterFn(vs []interface{}, f func(interface{}) bool) []interface{} {
	vsf := make([]interface{}, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func toMultiResult(vs []interface{}) MultiResult {
	vsm := make([]*Result, len(vs))
	for i, v := range vs {
		vsm[i] = &Result{
			Data: v,
			Err: nil,
		}
	}

	return MultiResult{
		Results: vsm,
		Err: nil,
	}
}