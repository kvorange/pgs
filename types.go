package pgs

type Record map[fieldI]interface{}

func (r Record) toMap() map[string]interface{} {
	insertMap := make(map[string]interface{})
	for field, value := range r {
		model := field.getModel()
		if model.joiner != nil {
			insertMap[model.joiner.From] = value
		} else {
			insertMap[field.getField()] = value
		}
	}
	return insertMap
}
