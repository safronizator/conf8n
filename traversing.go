package conf8n

func getValueWithCompositeKey(m map[string]interface{}, keyChunks []string, current int) interface{} {
	data, _ := m[keyChunks[current]]
	if current == len(keyChunks)-1 {
		return data
	}
	if data := toStrMap(data); data != nil {
		return getValueWithCompositeKey(data, keyChunks, current+1)
	}
	return nil
}

func toStrMap(value interface{}) map[string]interface{} {
	if alreadyStrMap, ok := value.(map[string]interface{}); ok {
		return alreadyStrMap
	}
	if asIMap, ok := value.(map[interface{}]interface{}); ok {
		strMap := make(map[string]interface{}, len(asIMap))
		for k, v := range asIMap {
			if asStr, ok := k.(string); ok {
				strMap[asStr] = v
			} else {
				return nil
			}
		}
		return strMap
	}
	return nil
}
