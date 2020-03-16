package pb

// this code satisfies kafka requirements of the record interface
// which reqiores Length() and Encoded() methods

func (r *Record) Length() int {
	return len(r.Payload)
}

func (r *Record) Encode() ([]byte, error) {
	return r.Payload, nil
}
