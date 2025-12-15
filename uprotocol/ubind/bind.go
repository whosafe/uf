package ubind

// Bind 从 Value 绑定数据到对象
func Bind(val *Value, v Binder) error {
	if val == nil {
		return nil
	}

	if val.Type == TypeObject && val.Object != nil {
		for key, value := range val.Object {
			if err := v.Bind(key, value); err != nil {
				return err
			}
		}
	}

	return nil
}
